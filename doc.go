// Package g provides an ergonomic standard library extension for Go:
// monadic error handling, rich generic containers, lazy iterators with
// type-changing generic methods (Go 1.27+), and ergonomic wrappers for
// primitive types and the filesystem.
//
// # Monads
//
// [Option] represents an optional value (Some/None) and [Result] represents
// success or failure (Ok/Err). Both offer chainable combinators — Map, Then,
// UnwrapOr, MapOr, Inspect and friends — so error and nil handling become
// expressions instead of if-ladders.
//
// # Containers
//
//   - [Slice]: an extended slice with 90+ methods.
//   - [Map]: a map with ergonomic accessors and entry API.
//   - [MapOrd]: an insertion-ordered map.
//   - [MapSafe]: a concurrency-safe map.
//   - [Set]: a hash set with algebraic operations (union, intersection, ...).
//   - [Deque]: a double-ended queue backed by a ring buffer.
//   - [Heap]: a binary min/max heap driven by a comparison function.
//
// # Iterators
//
// Every container exposes Iter/IntoIter methods returning lazy sequences —
// [SeqSlice], [SeqSlices], [SeqSet], [SeqMap], [SeqMapOrd], [SeqDeque],
// [SeqHeap], [SeqResult] and [SeqPairs] — built on top of
// github.com/enetx/iter. With Go 1.27 generic methods, transformations can
// change the element type mid-chain:
//
//	g.SliceOf(1, 2, 3).Iter().Map[string](strconv.Itoa).Collect() // Slice[string]
//
// Iterators are lazy: nothing is computed until a consumer (Collect, ForEach,
// Fold, ...) runs the chain.
//
// # Primitive wrappers
//
// [String], [Int], [Float] and [Bytes] wrap the built-in types with fluent
// methods, including conversion pipelines via Encode/Decode (Base64, Hex,
// Octal, Binary, JSON, ...), Compress/Decompress (gzip, zstd, brotli, ...)
// and Hash (MD5, SHA1, SHA256, SHA512).
//
// # Filesystem
//
// [File] and [Dir] provide chainable, Result-based file and directory
// operations, including lazy line/chunk iterators over file contents.
//
// # Subpackages
//
//   - pool: a generic goroutine pool with limits, rate limiting and streaming.
//   - cmp: ordering primitives (cmp.Ordering, cmp.Cmp) used by sorts and heaps.
//   - f: predicate combinators for filters (f.Eq, f.Gt, f.Contains, ...).
//   - ref: pointer helpers (ref.Of).
//   - constraints: generic type constraints shared across the library.
//   - dbg: debugging helpers that print expressions with source locations.
//
// # Panics
//
// The library reports failures through [Result] and [Option]; panics are
// reserved for programmer errors. Only two families of API panic:
//
//   - the Unwrap/Expect family on [Option] and [Result] (Unwrap, UnwrapErr,
//     Expect) when called on the wrong variant;
//   - documented index-based operations (e.g. Slice.Swap/Insert/Replace/SubSlice,
//     Deque.Insert/Swap) when given out-of-range indices, and constructors with
//     documented preconditions (e.g. NewHeap with a nil comparison function).
//
// Everything else returns Option/Result instead of panicking.
//
// # Security
//
// [Format] templates can invoke arbitrary exported methods on their arguments
// via reflection ({key.Method(...)} placeholders). Treat templates as code:
// never pass untrusted or user-controlled input as the template string —
// untrusted data belongs only in argument values.
package g
