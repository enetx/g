package g_test

import (
	"fmt"
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

func TestCellSwapTwoCells(t *testing.T) {
	// Test swapping two different cells
	cell1 := cell.New(&Config{Name: "cell1", Level: 10})
	cell2 := cell.New(&Config{Name: "cell2", Level: 20})

	cell1.Swap(cell2)

	val1 := cell1.Get()
	val2 := cell2.Get()

	if val1.Name != "cell2" || val1.Level != 20 {
		t.Errorf("Expected cell1 to have cell2's value, got %+v", val1)
	}
	if val2.Name != "cell1" || val2.Level != 10 {
		t.Errorf("Expected cell2 to have cell1's value, got %+v", val2)
	}

	// Test swapping the same cell (should be no-op)
	cell3 := cell.New(&Config{Name: "same", Level: 30})
	original := cell3.Get()
	cell3.Swap(cell3)

	after := cell3.Get()
	if after.Name != original.Name || after.Level != original.Level {
		t.Errorf("Swapping same cell should not change value: original %+v, after %+v", original, after)
	}

	// Test swapping multiple cells to trigger different pointer ordering scenarios
	cells := make([]*cell.Cell[*Config], 10)
	for i := 0; i < 10; i++ {
		cells[i] = cell.New(&Config{Name: fmt.Sprintf("cell%d", i), Level: i * 10})
	}

	// Perform swaps between different cells to ensure all ordering branches are hit
	for i := 0; i < len(cells)-1; i++ {
		for j := i + 1; j < len(cells); j++ {
			// Store original values
			origI := cells[i].Get()
			origJ := cells[j].Get()

			// Swap
			cells[i].Swap(cells[j])

			// Verify swap occurred
			newI := cells[i].Get()
			newJ := cells[j].Get()

			if newI.Name != origJ.Name || newI.Level != origJ.Level {
				t.Errorf("Swap failed: cells[%d] should have original cells[%d] value", i, j)
			}
			if newJ.Name != origI.Name || newJ.Level != origI.Level {
				t.Errorf("Swap failed: cells[%d] should have original cells[%d] value", j, i)
			}

			// Swap back to restore original state
			cells[i].Swap(cells[j])
		}
	}

	// Test swap with pointer ordering to cover unsafe.Pointer comparison branches
	// Create cells where second has lower memory address than first
	for k := 0; k < 100; k++ {
		cell1 := cell.New(&Config{Name: "lower", Level: 10})
		cell2 := cell.New(&Config{Name: "higher", Level: 20})

		// Try swapping both ways to hit different pointer ordering scenarios
		orig1 := cell1.Get()
		orig2 := cell2.Get()

		cell2.Swap(cell1) // second parameter has different address

		new1 := cell1.Get()
		new2 := cell2.Get()

		if new1.Name != orig2.Name || new1.Level != orig2.Level {
			t.Errorf("Pointer ordering swap failed: cell1 should have original cell2 value")
		}
		if new2.Name != orig1.Name || new2.Level != orig1.Level {
			t.Errorf("Pointer ordering swap failed: cell2 should have original cell1 value")
		}

		// Swap back
		cell1.Swap(cell2)
	}
}
