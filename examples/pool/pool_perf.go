package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/pool"
)

const (
	TASKS_NUM = 100_000
	ITEMS_NUM = Int(10_000)
)

type data struct {
	name String
	age  Int
}

var dataMapPool = sync.Pool{
	New: func() any {
		return NewMap[String, data](ITEMS_NUM)
	},
}

func main() {
	start := time.Now()

	p := pool.New[time.Duration]()
	p.Limit(runtime.NumCPU())

	for range TASKS_NUM {
		p.Go(func() Result[time.Duration] {
			start := time.Now()

			var sum uint64

			dataMap := dataMapPool.Get().(Map[String, data])
			defer dataMapPool.Put(dataMap)
			dataMap.Clear()

			for i := range ITEMS_NUM {
				name := i.String()

				dataMap.Set(name, data{name: name, age: i})
				if val := dataMap.Get(name); val.IsSome() && val.Some().name.Eq(name) {
					sum += uint64(val.Some().age)
				}
			}

			return Ok(time.Since(start))
		})
	}

	results := TransformSlice(p.Wait(), Result[time.Duration].Ok)

	taskSum := results.Iter().Fold(0, func(acc, val time.Duration) time.Duration { return acc + val })
	taskMin := results.MinBy(cmp.Cmp)
	taskMax := results.MaxBy(cmp.Cmp)
	taskAvg := taskSum / TASKS_NUM
	total := time.Since(start).Seconds()

	fmt.Printf(
		" - finished in %.4fs, task avg %.4fs, min %.4fs, max %.4fs\n",
		total,
		taskAvg.Seconds(),
		taskMin.Seconds(),
		taskMax.Seconds(),
	)
}
