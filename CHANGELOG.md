# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Breaking:** the library now requires **Go 1.27** (currently `1.27.0-rc.1`): the
  whole chain API is built on generic methods.
- **Breaking:** package-level generic helpers are replaced by generic methods that
  can change the payload type mid-chain: `Option.Map[U]`/`Option.Then[U]`,
  `Result.Map[U]`/`Result.Then[U]`, `SeqSlice.Map[U]`/`SeqSlice.Fold[A]`,
  `SeqMap.Fold[A]`, `SeqSlices.Map[U]`, `Transform[U]` on `Bytes`, `Float`,
  `Deque`, and friends.
- **Breaking:** `String.Split`/`SplitAfter` and `Bytes.Split`/`SplitAfter` take a
  mandatory `sep` argument; the old variadic `Split(sep ...)` form is gone.
- **Breaking:** `Int.IsPositive` is strict (`> 0`); it previously returned `true`
  for `0`.
- **Breaking:** `Zip` on `SeqSlice`/`SeqDeque`/`SeqHeap` is generic —
  `Zip[U](two)` — and returns the new typed `SeqPairs[V, U]` instead of
  `SeqMapOrd[any, any]`.
- **Breaking:** `SeqMapOrd.Unzip` is eager: it returns `(Slice[K], Slice[V])`
  collected in a single pass, instead of lazy `(SeqSlice[K], SeqSlice[V])` that
  iterated the source sequence twice.
- **Breaking:** `GroupBy` on sequence types is renamed to `ChunkBy` (semantics
  unchanged: it chunks *consecutive* elements).
- **Breaking:** `Dir.Exist` and `File.Exist` are renamed to `Exists`.
- **Breaking:** `String.TryBigInt` returns `Result[*big.Int]` instead of
  `Option[*big.Int]`.
- **Breaking:** `Dir.Temp` and `Dir.CreateTemp` are now the package-level
  functions `TempDir` and `CreateTempDir`. (`String.Random` stays a method.)
- **Breaking:** `pool.Pool.Reset` panics if tasks are still running instead of
  returning an `error`.
- **Breaking:** JSON support is backed by `encoding/json/v2` + `jsontext`
  (`String`/`Bytes`/`File` encoders, `Option`, `Result`): duplicate object keys
  are rejected, nil slices/maps marshal as `[]`/`{}`, and invalid UTF-8 is an
  error instead of being silently replaced.

### Removed

- **Breaking:** package-level mapping helpers, superseded by generic methods:
  `TransformSlice`, `TransformSet`, `TransformOption`, `TransformResult`,
  `TransformResultOf`, `MapOption`, `MapSeqResult`, and the free
  `Flatten`/`FlattenResult` functions for `SeqSlice`/`SeqResult`.
- **Breaking:** `Slice.AsAny` — with typed generic methods there is no need to
  erase to `Slice[any]`.
- **Breaking:** eager `Slice.MaxBy`/`Slice.MinBy` — use
  `.Iter().MaxBy/MinBy` (they return `Option[V]`).
- **Breaking:** the any-typed `SeqSlice.Counter` (`SeqMapOrd[any, Int]`) — use
  the generic `CounterBy[K comparable]` instead.
- **Breaking:** `SeqHeap.Eq` — use the new `Heap.Eq`/`Heap.Ne`.
- **Breaking:** `Option.Result(err)` — use `Option.OkOr`/`Option.OkOrElse`.
- **Breaking:** `Int.IsNonNegative` — use `!i.IsNegative()`.

### Added

- `SeqPairs[K, V]`: typed lazy key/value sequence produced by `Zip`, with
  `Keys`, `Values`, `Unzip`, `Map[T]`, `Filter`, `Exclude`, `FilterByKey`,
  `FilterByValue`, `Take`, `TakeWhile`, `SkipWhile`, and more.
- `Int`: checked, saturating, and overflowing arithmetic —
  `CheckedAdd/Sub/Mul/Div/Rem/Neg/Abs/Pow` (returning `Option[Int]`),
  `SaturatingAdd/Sub/Mul`, `OverflowingAdd/Sub/Mul`, and `Clamp`.
- `Float`: math suite — `IsNaN`, `IsInf`, `IsFinite`, `IsNormal`, `Signum`,
  `IsSignPositive/Negative`, `Ceil`, `Floor`, `Trunc`, `Fract`, `Clamp`,
  `Recip`, `Copysign`, `MulAdd`, `Hypot`, `Cbrt`, `Exp/Exp2/ExpM1`,
  `Ln/Ln1p/Log2/Log10`, trigonometric and hyperbolic functions.
- `Bytes`: text-API parity with `String` — `IsASCII`, `IsDigit`,
  `ReplaceMulti`, `Remove`, `ReplaceNth`, `Chunks`, `Cut`, `SubBytes`,
  `Similarity`, `Truncate`, `LeftJustify`, `RightJustify`, `Center`.
- `Bytes`: encode/decode parity with `String` — `JSON`, `URL`, `HTML`,
  `Rot13`, `Octal` on `Encode()`/`Decode()`.
- `Result`: JSON support — `MarshalJSON`/`UnmarshalJSON` plus json/v2
  `MarshalJSONTo`/`UnmarshalJSONFrom` (`jsontext`), externally tagged as
  `{"ok": …}` / `{"err": …}`.
- New sequence adapters: `TakeWhile`/`SkipWhile` on all sequence types
  (`SeqSlice`, `SeqSet`, `SeqMap`, `SeqMapOrd`, `SeqDeque`, `SeqHeap`,
  `SeqPairs`); `CounterBy[K]` on `SeqSlice`/`SeqSet`/`SeqDeque`/`SeqHeap`;
  `FilterByKey`/`FilterByValue` on `SeqMap`/`SeqMapOrd`/`SeqPairs`;
  `SeqResult.TryCollect`.
- Constructors: `MapOf`, `MapOrdOf`, `HeapOf`, `PairOf`.
- `Set.Disjoint`.
- `Heap.Eq`/`Heap.Ne` (value equality, replacing `SeqHeap.Eq`).
- `f.Id` identity helper.

### Fixed

- `Slice.SubSlice` (and `Bytes.SubBytes`): index-out-of-range panic when a
  negative step started at the end of the slice; `s.SubSlice(100, 0, -1)` now
  follows Python slicing semantics and starts at the last element.
- `SeqMapOrd.Collect`: was O(n²) (a linear `Insert` scan per element); now
  builds an index map and runs in O(n).
- `Cycle` on sequences: no longer probes the source with an extra up-front
  pass (which consumed an element of single-use sources), and terminates once
  the source stops yielding elements instead of spinning forever.
- `Next` on sequence iterators: now advances via a pull iterator, so the
  source is walked exactly once — O(1) per element and correct for
  non-deterministic sources (map-backed sets, channels); previously each call
  re-ran the source from the start.
- `RPosition`: single forward pass with O(1) memory instead of materializing
  the entire sequence into a slice.
- SQL: `uint`/`uint64` values above `math.MaxInt64` are rejected when
  converting to `driver.Value` instead of silently wrapping into negative
  numbers.
