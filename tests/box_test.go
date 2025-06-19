package g_test

import (
	"sync"
	"testing"

	"github.com/enetx/g/box"
)

type Config struct {
	Name  string
	Level int
}

func TestNewBoxLoadStore(t *testing.T) {
	b := box.New(&Config{Name: "init", Level: 1})

	val := b.Load()
	if val.Name != "init" || val.Level != 1 {
		t.Fatalf("unexpected initial value: %+v", val)
	}

	b.Store(&Config{Name: "updated", Level: 2})
	val = b.Load()
	if val.Name != "updated" || val.Level != 2 {
		t.Fatalf("unexpected stored value: %+v", val)
	}
}

func TestBoxUpdate(t *testing.T) {
	b := box.New(&Config{Name: "x", Level: 10})

	b.Update(func(c *Config) *Config {
		cp := *c
		cp.Level += 5
		return &cp
	})

	val := b.Load()
	if val.Level != 15 {
		t.Fatalf("expected level 15, got %d", val.Level)
	}
}

func TestBoxConcurrentUpdate(t *testing.T) {
	b := box.New(&Config{Name: "counter", Level: 0})

	var wg sync.WaitGroup
	const workers = 100

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Update(func(c *Config) *Config {
				cp := *c
				cp.Level++
				return &cp
			})
		}()
	}

	wg.Wait()

	val := b.Load()
	if val.Level != workers {
		t.Fatalf("expected level %d, got %d", workers, val.Level)
	}
}

func TestBoxSwap(t *testing.T) {
	b := box.New(&Config{Name: "first", Level: 1})
	newVal := &Config{Name: "second", Level: 2}

	// Perform the swap
	oldVal := b.Swap(newVal)

	// 1. Check if the returned old value is correct
	if oldVal.Name != "first" || oldVal.Level != 1 {
		t.Fatalf("expected Swap to return the old value ('first'), but got: %+v", oldVal)
	}

	// 2. Check if the new value is now stored in the box
	currentVal := b.Load()
	if currentVal.Name != "second" || currentVal.Level != 2 {
		t.Fatalf("expected the new value ('second') to be in the box, but got: %+v", currentVal)
	}

	// 3. Ensure the pointer is the same one we passed in
	if currentVal != newVal {
		t.Fatal("expected the pointer in the box to be the same as the one passed to Swap")
	}
}

func TestBoxCompareAndSwap(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		oldVal := &Config{Name: "old", Level: 1}
		newVal := &Config{Name: "new", Level: 2}
		b := box.New(oldVal)

		// This should succeed because the current value matches `oldVal`
		swapped := b.CompareAndSwap(oldVal, newVal)
		if !swapped {
			t.Fatal("expected CompareAndSwap to succeed, but it failed")
		}

		// Verify the new value is in place
		if b.Load() != newVal {
			t.Fatal("box does not contain the new value after successful CAS")
		}
	})

	t.Run("Failure", func(t *testing.T) {
		currentVal := &Config{Name: "current", Level: 1}
		staleOldVal := &Config{Name: "stale", Level: 0} // A different pointer
		newVal := &Config{Name: "new", Level: 2}
		b := box.New(currentVal)

		// This should fail because `staleOldVal` does not match `currentVal`
		swapped := b.CompareAndSwap(staleOldVal, newVal)
		if swapped {
			t.Fatal("expected CompareAndSwap to fail, but it succeeded")
		}

		// Verify the value in the box has not changed
		if b.Load() != currentVal {
			t.Fatal("box value was changed after a failed CAS")
		}
	})
}

func TestBoxTryUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		b := box.New(&Config{Level: 10})

		ok := b.TryUpdate(func(c *Config) *Config {
			cp := *c
			cp.Level = 20
			return &cp
		})

		if !ok {
			t.Fatal("expected TryUpdate to succeed, but it returned false")
		}
		if b.Load().Level != 20 {
			t.Fatalf("expected level 20 after TryUpdate, got %d", b.Load().Level)
		}
	})

	t.Run("Failure due to contention", func(t *testing.T) {
		b := box.New(&Config{Level: 0})

		// This TryUpdate will be forced to fail because another goroutine
		// will change the value before its CompareAndSwap can execute.
		ok := b.TryUpdate(func(c *Config) *Config {
			// Just as this function is about to return, another goroutine
			// will sneak in and change the value in the box.
			b.Store(&Config{Level: 999}) // Simulate contention

			// This new value will be computed based on the original `c`.
			cp := *c
			cp.Level = 1
			return &cp
		})

		if ok {
			t.Fatal("expected TryUpdate to fail due to contention, but it succeeded")
		}

		// The final value should be the one set by the 'contending' goroutine.
		if b.Load().Level != 999 {
			t.Fatalf("expected final level to be 999, got %d", b.Load().Level)
		}
	})
}

func TestBoxUpdateAndGet(t *testing.T) {
	b := box.New(&Config{Level: 5})

	// Atomically update the value and get the new state back
	newValue := b.UpdateAndGet(func(c *Config) *Config {
		cp := *c
		cp.Level *= 3 // 5 * 3 = 15
		return &cp
	})

	// 1. Check if the returned value is correct
	if newValue.Level != 15 {
		t.Fatalf("expected UpdateAndGet to return level 15, got %d", newValue.Level)
	}

	// 2. Check if the value in the box is also correct
	currentVal := b.Load()
	if currentVal.Level != 15 {
		t.Fatalf("expected box to contain level 15, got %d", currentVal.Level)
	}

	// 3. Ensure the returned pointer is the same as the one now in the box
	if currentVal != newValue {
		t.Fatal("pointer returned by UpdateAndGet differs from the one in the box")
	}
}

func TestBoxConcurrentUpdateAndGet(t *testing.T) {
	b := box.New(&Config{Level: 0})

	var wg sync.WaitGroup
	const workers = 100

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.UpdateAndGet(func(c *Config) *Config {
				cp := *c
				cp.Level++
				return &cp
			})
		}()
	}

	wg.Wait()

	finalVal := b.Load()
	if finalVal.Level != workers {
		t.Fatalf("expected final level to be %d, got %d", workers, finalVal.Level)
	}
}
