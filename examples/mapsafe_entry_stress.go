package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	. "github.com/enetx/g"
)

type TestResult struct {
	name       string
	total      int64
	races      int64
	badResults int64
}

func main() {
	tests := []struct {
		name string
		fn   func(*MapSafe[string, Int]) Int
	}{
		{
			name: "OrInsert",
			fn: func(ms *MapSafe[string, Int]) Int {
				return ms.Entry("key").
					AndModify(func(v *Int) { *v += 10 }).
					OrInsert(1)
			},
		},
		{
			name: "OrInsertWith",
			fn: func(ms *MapSafe[string, Int]) Int {
				return ms.Entry("key").
					AndModify(func(v *Int) { *v += 10 }).
					OrInsertWith(func() Int { return 1 })
			},
		},
		{
			name: "OrInsertWithKey",
			fn: func(ms *MapSafe[string, Int]) Int {
				return ms.Entry("key").
					AndModify(func(v *Int) { *v += 10 }).
					OrInsertWithKey(func(string) Int { return 1 })
			},
		},
		{
			name: "OrDefault",
			fn: func(ms *MapSafe[string, Int]) Int {
				return ms.Entry("key").
					AndModify(func(v *Int) { *v += 10 }).
					OrDefault()
			},
		},
		{
			name: "VacantEntry.Insert",
			fn: func(ms *MapSafe[string, Int]) Int {
				entry := ms.Entry("key")
				entry = entry.AndModify(func(v *Int) { *v += 10 })
				if v, ok := entry.(VacantSafeEntry[string, Int]); ok {
					return v.Insert(1)
				}
				return entry.OrInsert(1)
			},
		},
	}

	results := make([]TestResult, len(tests))

	for i, test := range tests {
		results[i] = runStressTest(test.name, test.fn)
	}

	// Print results
	fmt.Println("=" + strings.Repeat("=", 70))
	fmt.Printf("%-20s %12s %12s %12s\n", "Method", "Total", "Races", "Bad (0)")
	fmt.Println(strings.Repeat("-", 70))

	allPassed := true
	for _, r := range results {
		status := "✓"
		if r.badResults > 0 {
			status = "✗ BUG!"
			allPassed = false
		}

		fmt.Printf("%-20s %12d %12d %12d %s\n", r.name, r.total, r.races, r.badResults, status)
	}

	fmt.Println(strings.Repeat("=", 71))
	if allPassed {
		fmt.Println("All tests passed!")
	} else {
		fmt.Println("BUGS DETECTED!")
	}
}

func runStressTest(name string, testFn func(*MapSafe[string, Int]) Int) TestResult {
	var badResults atomic.Int64
	var totalRaces atomic.Int64
	var wg sync.WaitGroup

	const workers = 50
	const iterations = 50_000

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for range iterations {
				ms := NewMapSafe[string, Int]()
				ready := make(chan struct{})
				var result Int
				var done sync.WaitGroup

				// Goroutine A: Entry + AndModify + Or* method
				done.Add(1)
				go func() {
					defer done.Done()
					<-ready
					result = testFn(ms)
				}()

				// Goroutine B: Insert value 100
				done.Add(1)
				go func() {
					defer done.Done()
					<-ready
					ms.Insert("key", 100)
				}()

				// Goroutines C-F: aggressively remove key
				for range 4 {
					done.Add(1)
					go func() {
						defer done.Done()
						<-ready
						for range 50 {
							ms.Remove("key")
						}
					}()
				}

				close(ready)
				done.Wait()

				// Check result
				// Valid values: 0 (for OrDefault), 1, 100, 110
				// For OrDefault: valid are 0, 10, 100, 110
				isOrDefault := name == "OrDefault"

				if isOrDefault {
					// OrDefault inserts zero, so valid: 0, 10, 100, 110
					if result != 0 && result != 10 && result != 100 && result != 110 {
						badResults.Add(1)
					}
				} else {
					// Others insert 1, so valid: 1, 100, 110
					// 0 is a BUG
					if result == 0 {
						badResults.Add(1)
					}
				}

				// Count races (B inserted before A)
				if result == 100 || result == 110 || (isOrDefault && result == 10) {
					totalRaces.Add(1)
				}
			}
		}()
	}

	wg.Wait()

	return TestResult{
		name:       name,
		total:      workers * iterations,
		races:      totalRaces.Load(),
		badResults: badResults.Load(),
	}
}
