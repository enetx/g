package main

import (
	"fmt"
	"sync"
	"time"

	. "github.com/enetx/g"
)

const (
	TASKS_NUM = 100_000
	ITEMS_NUM = Int(10_000)
)

type data struct {
	name String
	age  Int
}

func main() {
	start := time.Now()

	mapPool := sync.Pool{
		New: func() any { return NewMap[String, data](ITEMS_NUM).Ptr() },
	}

	pool := NewPool[time.Duration]().Limit(100)

	for range TASKS_NUM {
		pool.Go(func() Result[time.Duration] {
			start := time.Now()

			var sum uint64
			dataMap := mapPool.Get().(*Map[String, data])
			defer mapPool.Put(dataMap)

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

	var (
		taskSum time.Duration
		taskMin time.Duration = time.Hour
		taskMax time.Duration
	)

	pool.Wait().Iter().
		ForEach(func(v Result[time.Duration]) {
			val := v.Ok()
			if val < taskMin {
				taskMin = val
			}
			if val > taskMax {
				taskMax = val
			}
			taskSum += val
		})

	taskAvg := taskSum / time.Duration(TASKS_NUM)
	total := time.Since(start)

	fmt.Printf(
		" - finished in %.4fs, task avg %.4fs, min %.4fs, max %.4fs\n",
		total.Seconds(),
		taskAvg.Seconds(),
		taskMin.Seconds(),
		taskMax.Seconds(),
	)
}
