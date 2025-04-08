// https://github.com/curvednebula/perf-tests

package main

import (
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

const (
	TASKS_NUM               = 100_000
	ITEMS_NUM               = Int(10_000)
	TASKS_IN_BUNCH          = 10
	TIME_BETWEEN_BUNCHES_MS = 1
)

type data struct {
	name String
	age  Int
}

func main() {
	start := time.Now()

	pool := NewPool[float64]()
	pool.Limit(1000)

	for taskID := range TASKS_NUM {
		pool.Go(func() Result[float64] {
			start := time.Now()

			var sum uint64
			dataMap := NewMap[String, data](ITEMS_NUM)

			for i := range ITEMS_NUM {
				name := i.String()

				dataMap.Set(name, data{name: name, age: i})
				if val := dataMap.Get(name); val.IsSome() && val.Some().name.Eq(name) {
					sum += uint64(val.Some().age)
				}
			}

			return Ok(time.Since(start).Seconds())
		})

		if taskID%TASKS_IN_BUNCH == 0 {
			time.Sleep(TIME_BETWEEN_BUNCHES_MS * time.Millisecond)
		}
	}

	tasks := pool.Wait().Iter().Filter(Result[float64].IsOk).Collect()
	results := TransformSlice(tasks, Result[float64].Ok)

	taskSum := results.Iter().Fold(0,
		func(acc, val float64) float64 {
			return acc + val
		})

	taskMin := results.MinBy(cmp.Cmp)
	taskMax := results.MaxBy(cmp.Cmp)

	taskAvg := taskSum / TASKS_NUM

	total := time.Since(start).Seconds()
	fmt.Printf(" - finished in %.4fs, task avg %.4fs, min %.4fs, max %.4fs\n", total, taskAvg, taskMin, taskMax)
}
