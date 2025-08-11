package g_test

import (
	"sync"
	"testing"

	"github.com/enetx/g/cell"
)

type Config struct {
	Name  string
	Level int
}

func TestNewCellLoadStore(t *testing.T) {
	c := cell.New(&Config{Name: "init", Level: 1})

	val := c.Get()
	if val.Name != "init" || val.Level != 1 {
		t.Fatalf("unexpected initial value: %+v", val)
	}

	c.Set(&Config{Name: "updated", Level: 2})
	val = c.Get()
	if val.Name != "updated" || val.Level != 2 {
		t.Fatalf("unexpected stored value: %+v", val)
	}
}

func TestCellUpdate(t *testing.T) {
	c := cell.New(&Config{Name: "x", Level: 10})

	c.Update(func(current *Config) *Config {
		cp := *current
		cp.Level += 5
		return &cp
	})

	val := c.Get()
	if val.Level != 15 {
		t.Fatalf("expected level 15, got %d", val.Level)
	}
}

func TestCellConcurrentUpdate(t *testing.T) {
	c := cell.New(&Config{Name: "counter", Level: 0})

	var wg sync.WaitGroup
	const workers = 100

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Update(func(current *Config) *Config {
				cp := *current
				cp.Level++
				return &cp
			})
		}()
	}

	wg.Wait()

	val := c.Get()
	if val.Level != workers {
		t.Fatalf("expected level %d, got %d", workers, val.Level)
	}
}

func TestCellSwap(t *testing.T) {
	c := cell.New(&Config{Name: "first", Level: 1})
	newVal := &Config{Name: "second", Level: 2}

	// Perform the replace (equivalent to swap)
	oldVal := c.Replace(newVal)

	// 1. Check if the returned old value is correct
	if oldVal.Name != "first" || oldVal.Level != 1 {
		t.Fatalf("expected Replace to return the old value ('first'), but got: %+v", oldVal)
	}

	// 2. Check if the new value is now stored in the cell
	currentVal := c.Get()
	if currentVal.Name != "second" || currentVal.Level != 2 {
		t.Fatalf("expected the new value ('second') to be in the cell, but got: %+v", currentVal)
	}

	// 3. Ensure the pointer is the same one we passed in
	if currentVal != newVal {
		t.Fatal("expected the pointer in the cell to be the same as the one passed to Swap")
	}
}
