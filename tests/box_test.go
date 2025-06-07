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

func TestNewBox_Load_Store(t *testing.T) {
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

func TestBox_Update(t *testing.T) {
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

func TestBox_ConcurrentUpdate(t *testing.T) {
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
