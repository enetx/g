<p align="center">
  <img src="https://user-images.githubusercontent.com/65846651/229838021-741ff719-8c99-45f6-88d2-1a32927bd863.png">
</p>

# g - Functional programming framework for Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/enetx/g.svg)](https://pkg.go.dev/github.com/enetx/g)
[![Go Report Card](https://goreportcard.com/badge/github.com/enetx/g)](https://goreportcard.com/report/github.com/enetx/g)
[![Coverage Status](https://coveralls.io/repos/github/enetx/g/badge.svg?branch=main&service=github)](https://coveralls.io/github/enetx/g?branch=main)
[![Go](https://github.com/enetx/g/actions/workflows/go.yml/badge.svg)](https://github.com/enetx/g/actions/workflows/go.yml)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/enetx/g)

```bash
go get github.com/enetx/g
```

Requires **Go 1.24+**

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/enetx/g"
)

func main() {
    // Slice with functional operations
    g.SliceOf(1, 2, 3, 4, 5).
        Iter().
        Filter(func(x int) bool { return x%2 == 0 }).
        Map(func(x int) int { return x * x }).
        Collect().
        Println() // Slice[4, 16]

    // Safe map access with Option
    m := g.NewMap[string, int]()
    m.Set("key", 42)

    value := m.Get("key").UnwrapOr(0)
    fmt.Println(value) // 42

    // Result for error handling
    result := g.String("123").ToInt()
    if result.IsOk() {
        fmt.Println(result.Unwrap()) // 123
    }
}
```

---

## Navigation

| Core | Collections | Sync | Types |
|:----:|:-----------:|:----:|:-----:|
| [Option](#option) | [Slice](#slice) | [Mutex](#mutex) | [String](#string) |
| [Result](#result) | [Map](#map) | [RwLock](#rwlock) | [Int/Float](#int--float) |
| [Iterators](#iterators) | [Set](#set) | [Pool](#pool) | [Bytes](#bytes) |
| | [Heap](#heap) | | [File/Dir](#file--directory) |
| | [Deque](#deque) | | |

---

## Option

Safe nullable values: `Some(value)` or `None`.

```go
opt := g.Some(42)           // Some(42)
none := g.None[int]()       // None

opt.IsSome()                // true
opt.Unwrap()                // 42 (panics if None)
opt.UnwrapOr(0)             // 42 (returns default if None)
opt.UnwrapOrDefault()       // 42 (returns zero value if None)

// From map lookup
v, ok := myMap[key]
opt := g.OptionOf(v, ok)

// Chaining
g.Some(5).Then(func(x int) g.Option[int] {
    return g.Some(x * 2)
}) // Some(10)
```

<details>
<summary><b>All Option Methods</b></summary>

| Method | Description |
|--------|-------------|
| `IsSome()` | Returns true if contains value |
| `IsNone()` | Returns true if empty |
| `Some()` | Returns value (same as Unwrap) |
| `Unwrap()` | Returns value, panics if None |
| `UnwrapOr(default)` | Returns value or default |
| `UnwrapOrDefault()` | Returns value or zero value |
| `UnwrapOrElse(fn)` | Returns value or result of fn |
| `Expect(msg)` | Returns value, panics with msg if None |
| `Then(fn)` | Transforms value if Some |
| `Option()` | Returns (value, bool) |
| `Take()` | Takes value, leaves None |
| `Filter(fn)` | Returns None if predicate fails |

</details>

---

## Result

Success (`Ok`) or failure (`Err`) with error.

```go
ok := g.Ok(42)                              // Ok(42)
err := g.Err[int](errors.New("failed"))     // Err

ok.IsOk()                   // true
ok.Unwrap()                 // 42
err.UnwrapOr(0)             // 0

// From standard (value, error) pattern
result := g.ResultOf(strconv.Atoi("42"))    // Ok(42)
result := g.String("42").ToInt()            // Ok(42)

// Chaining
g.Ok(10).Then(func(x int) g.Result[int] {
    return g.Ok(x * 2)
}) // Ok(20)

// Convert to Option (discards error)
opt := result.Option()      // Some(42)
```

<details>
<summary><b>All Result Methods</b></summary>

| Method | Description |
|--------|-------------|
| `IsOk()` | Returns true if success |
| `IsErr()` | Returns true if error |
| `Ok()` | Returns value (same as Unwrap) |
| `Err()` | Returns error |
| `Unwrap()` | Returns value, panics if Err |
| `UnwrapOr(default)` | Returns value or default |
| `UnwrapOrDefault()` | Returns value or zero value |
| `UnwrapOrElse(fn)` | Returns value or result of fn |
| `Expect(msg)` | Returns value, panics with msg |
| `Then(fn)` | Chains operation if Ok |
| `ThenOf(fn)` | Chains (T, error) function |
| `MapErr(fn)` | Transforms error if Err |
| `Option()` | Converts to Option |
| `Result()` | Returns (value, error) |

</details>

---

## Slice

Extended slice with 90+ methods.

```go
s := g.SliceOf(1, 2, 3, 4, 5)
s := g.NewSlice[int](10)        // with capacity

s.Len()                         // Int(5)
s.Get(0)                        // Some(1)
s.Get(100)                      // None (safe!)
s.Last()                        // Some(5)
s.Contains(3)                   // true

s.Push(6, 7)                    // append in place
s.Pop()                         // Some(7)
s.Clone()                       // safe copy

// Sorting
s.SortBy(cmp.Cmp)               // ascending
s.Reverse()                     // in place
s.Shuffle()                     // random order
```

<details>
<summary><b>All Slice Methods</b></summary>

| Category | Methods |
|----------|---------|
| **Access** | `Get`, `Last`, `First`, `Random`, `RandomSample` |
| **Modify** | `Set`, `Push`, `Pop`, `Insert`, `Delete`, `Clear` |
| **Search** | `Contains`, `ContainsBy`, `Index`, `IndexBy` |
| **Transform** | `Clone`, `Reverse`, `Shuffle`, `SortBy`, `SubSlice` |
| **Convert** | `Iter`, `ToHeap`, `ToSet`, `Std` |
| **Info** | `Len`, `Cap`, `Empty`, `NotEmpty` |

</details>

---

## Map

Extended map with Entry API.

```go
m := g.NewMap[string, int]()

m.Set("a", 1)                   // set value
m.Get("a")                      // Some(1)
m.Get("x")                      // None
m.Contains("a")                 // true
m.Delete("a")                   // delete

m.Keys()                        // Slice of keys
m.Values()                      // Slice of values
```

### Entry API

Efficient update without multiple lookups:

```go
// Insert if absent
m.Entry("counter").OrInsert(0)

// Update if present, insert if absent
m.Entry("counter").AndModify(func(v *int) { *v++ }).OrInsert(1)

// Lazy initialization
m.Entry("key").OrInsertWith(func() int {
    return expensiveComputation()
})
```

**Word frequency example:**
```go
words := g.SliceOf("apple", "banana", "apple", "cherry", "banana", "apple")
freq := g.NewMap[string, int]()

for _, word := range words {
    freq.Entry(word).AndModify(func(v *int) { *v++ }).OrInsert(1)
}
// {"apple": 3, "banana": 2, "cherry": 1}
```

<details>
<summary><b>Entry Pattern Matching</b></summary>

```go
switch e := m.Entry("key").(type) {
case g.OccupiedEntry[string, int]:
    fmt.Println("Exists:", e.Get())
    e.Insert(newValue)      // replace
    e.Remove()              // delete
case g.VacantEntry[string, int]:
    e.Insert(defaultValue)  // insert
}
```

</details>

---

## Set

Collection of unique elements.

```go
s := g.SetOf(1, 2, 3)
s.Insert(4)                     // add
s.Remove(1)                     // remove
s.Contains(2)                   // true
```

### Set Operations

```go
a := g.SetOf(1, 2, 3, 4)
b := g.SetOf(3, 4, 5, 6)

a.Union(b).Collect()            // {1, 2, 3, 4, 5, 6}
a.Intersection(b).Collect()     // {3, 4}
a.Difference(b).Collect()       // {1, 2}
a.SymmetricDifference(b).Collect() // {1, 2, 5, 6}

a.Subset(b)                     // false
a.Superset(b)                   // false
```

---

## Heap

Priority queue (binary heap).

```go
import "github.com/enetx/g/cmp"

// Min-heap (smallest first)
h := g.NewHeap(cmp.Cmp[int])
h.Push(5, 3, 8, 1, 9)

h.Peek()                        // Some(1)
h.Pop()                         // Some(1)
h.Pop()                         // Some(3)

// Max-heap (largest first)
maxH := g.NewHeap(func(a, b int) cmp.Ordering {
    return cmp.Cmp(b, a)
})

// From slice
heap := g.SliceOf(5, 3, 8, 1).ToHeap(cmp.Cmp)
```

---

## Deque

Double-ended queue (ring buffer). O(1) at both ends.

```go
dq := g.NewDeque[int]()
dq := g.DequeOf(1, 2, 3)

dq.PushFront(0)                 // add to front
dq.PushBack(4)                  // add to back

dq.Front()                      // Some(0)
dq.Back()                       // Some(4)

dq.PopFront()                   // Some(0)
dq.PopBack()                    // Some(4)

dq.Get(1)                       // element at index
dq.Len()                        // length
dq.IsEmpty()                    // true/false
```

---

## Iterators

Functional operations on sequences.

```go
g.SliceOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
    Iter().
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Take(3).
    Collect()
// [4, 16, 36]
```

### Common Patterns

```go
// Chain iterators
g.SliceOf(1, 2).Iter().Chain(g.SliceOf(3, 4).Iter()).Collect()
// [1, 2, 3, 4]

// Sum with Fold
g.SliceOf(1, 2, 3, 4, 5).Iter().Fold(0, func(acc, x int) int { return acc + x })
// 15

// Enumerate
g.SliceOf("a", "b", "c").Iter().Enumerate().ForEach(func(i g.Int, v string) {
    fmt.Printf("%d: %s\n", i, v)
})

// Using f package predicates
import "github.com/enetx/g/f"

g.SliceOf(1, 2, 3, 4, 5).Iter().Filter(f.Ne(3)).Collect()       // [1, 2, 4, 5]
g.SliceOf("", "a", "").Iter().Exclude(f.IsZero).Collect()       // ["a"]
```

<details>
<summary><b>All Iterator Methods</b></summary>

| Transform | Slice | Combine | Aggregate | Search | Other |
|-----------|-------|---------|-----------|--------|-------|
| `Map` | `Take` | `Chain` | `Collect` | `Find` | `ForEach` |
| `Filter` | `Skip` | `Zip` | `Fold` | `Any` | `Inspect` |
| `Exclude` | `StepBy` | `Enumerate` | `Reduce` | `All` | `Range` |
| `FilterMap` | `First` | `Intersperse` | `Count` | `MaxBy` | `Scan` |
| `FlatMap` | `Last` | `Cycle` | `Counter` | `MinBy` | `Combinations` |
| `Flatten` | `Nth` | | `Partition` | | `Permutations` |
| `Dedup` | `Chunks` | | `GroupBy` | | `Context` |
| `Unique` | `Windows` | | | | `ToChan` |
| `SortBy` | | | | | `Parallel` |

</details>

---

## Mutex

Typed mutex — data bound to lock.

```go
counter := g.NewMutex(0)

// Lock and modify
guard := counter.Lock()
guard.Set(guard.Get() + 1)
guard.Unlock()

// With defer (recommended)
func increment(c *g.Mutex[int]) {
    guard := c.Lock()
    defer guard.Unlock()
    guard.Set(guard.Get() + 1)
}

// Direct pointer access
guard := counter.Lock()
defer guard.Unlock()
*guard.Deref() += 1

// Non-blocking
if opt := counter.TryLock(); opt.IsSome() {
    guard := opt.Unwrap()
    defer guard.Unlock()
    // got the lock
}
```

**Why typed mutex?**
```go
// Traditional — easy to forget locking
type Old struct {
    mu   sync.Mutex
    data map[string]int  // what protects this?
}

// Typed — impossible to access without lock
type New struct {
    data *g.Mutex[g.Map[string, int]]
}

func (s *New) Get(key string) g.Option[int] {
    guard := s.data.Lock()
    defer guard.Unlock()
    return guard.Deref().Get(key)
}
```

---

## RwLock

Multiple readers OR single writer.

```go
config := g.NewRwLock(Config{Port: 8080})

// Read (concurrent)
func getPort(c *g.RwLock[Config]) int {
    guard := c.Read()
    defer guard.Unlock()
    return guard.Get().Port
}

// Write (exclusive)
func setPort(c *g.RwLock[Config], port int) {
    guard := c.Write()
    defer guard.Unlock()
    guard.Deref().Port = port
}

// Non-blocking
if opt := config.TryRead(); opt.IsSome() { ... }
if opt := config.TryWrite(); opt.IsSome() { ... }
```

| Use Case | Choose |
|----------|--------|
| Read-heavy (config, cache) | `RwLock` |
| Write-heavy or balanced | `Mutex` |
| Simple counters | `Mutex` |

---

## Pool

Goroutine pool for parallel tasks.

```go
import "github.com/enetx/g/pool"

p := pool.New[int]().Limit(4)

for i := range 10 {
    p.Go(func() g.Result[int] {
        return g.Ok(i * 2)
    })
}

for result := range p.Wait() {
    if result.IsOk() {
        fmt.Println(result.Unwrap())
    }
}
```

<details>
<summary><b>With Context and Cancel on Error</b></summary>

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

p := pool.New[int]().
    Context(ctx).
    CancelOnError().
    Limit(4)

for i := range 100 {
    p.Go(func() g.Result[int] {
        if i == 50 {
            return g.Err[int](errors.New("error"))
        }
        return g.Ok(i)
    })
}

for result := range p.Wait() {
    // stops after first error
}
```

</details>

---

## String

Extended string with 80+ methods.

```go
s := g.String("Hello, World!")

s.Len()                         // 13
s.Upper()                       // "HELLO, WORLD!"
s.Lower()                       // "hello, world!"
s.Contains("World")             // true
s.Split(", ").Collect()         // ["Hello", "World!"]
s.Trim()                        // remove whitespace

// Conversions
s.Std()                         // Go string
s.Bytes()                       // Bytes type
g.String("42").ToInt()          // Result[Int]

// Encoding
s.Hash().MD5()                  // hash
s.Enc().Base64()                // encode
s.Comp().Gzip()                 // compress
```

---

## Int / Float

Numeric types with utility methods.

```go
i := g.Int(42)

i.Add(8)                        // 50
i.Mul(2)                        // 84
i.Abs()                         // absolute value
i.Min(10, 20)                   // minimum
i.Max(10, 20)                   // maximum

i.IsZero()                      // false
i.IsPositive()                  // true
i.IsNegative()                  // false

i.Binary()                      // binary string
i.Hex()                         // hex string

// Float
f := g.Float(3.14159)
f.Round()                       // Int(3)
f.RoundDecimal(2)               // Float(3.14)
f.Sqrt()                        // square root

// Compare
i.Cmp(g.Int(50))                // cmp.Less
```

---

## Bytes

Extended byte slice.

```go
b := g.Bytes([]byte("Hello"))

b.Len()
b.Upper()
b.Lower()
b.Contains([]byte("ell"))

b.String()                      // to String
b.Std()                         // to []byte
```

---

## File / Directory

### File

```go
// Read/Write
content := g.NewFile("config.json").Read()      // Result[String]
g.NewFile("output.txt").Write("Hello")
g.NewFile("log.txt").Append("line\n")

// Lazy line reading (memory efficient)
for line := range g.NewFile("large.txt").Lines() {
    if line.IsOk() {
        fmt.Println(line.Unwrap())
    }
}

// Properties
f := g.NewFile("test.txt")
f.Exist()                       // check existence
f.Stat()                        // file info
f.Copy("backup.txt")            // copy
f.Rename("new.txt")             // rename
f.Remove()                      // delete
```

### File Guard (Exclusive Lock)

```go
f := g.NewFile("data.txt").Guard()  // holds exclusive lock
f.Write("exclusive data")
f.Close()                           // release lock
```

### Directory

```go
g.NewDir("mydir").Create()
g.NewDir("path/to/deep").CreateAll()

// Read contents
for file := range g.NewDir(".").Read() {
    if file.IsOk() {
        file.Ok().Name().Println()
    }
}

// Recursive walk
for file := range g.NewDir(".").Walk() { ... }

// Glob
for file := range g.NewDir("*.txt").Glob() { ... }

g.NewDir("src").Copy("dst")         // copy directory
g.NewDir("mydir").Remove()          // remove (recursive)
```

---

## Other Maps

### MapOrd — Ordered Map

Maintains insertion order.

```go
m := g.NewMapOrd[string, int]()
m.Set("c", 3)
m.Set("a", 1)
m.Set("b", 2)

for k, v := range m.Iter() {
    fmt.Println(k, v)  // c, a, b order
}

m.SortByKey(cmp.Cmp)            // sort by key
m.SortByValue(cmp.Cmp)          // sort by value
```

### MapSafe — Concurrent Map

Thread-safe for concurrent access.

```go
m := g.NewMapSafe[string, int]()

// Safe from multiple goroutines
go func() { m.Set("a", 1) }()
go func() { m.Get("a") }()
```

---

## License

MIT License
