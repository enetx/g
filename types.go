package g

import (
	"context"
	"iter"
	"os"
	"sync"
)

type (
	// Result is a generic struct for representing a result value along with an error.
	Result[T any] struct {
		value *T    // Pointer to the value.
		err   error // Associated error.
	}

	// Option is a generic struct for representing an optional value.
	Option[T any] struct {
		value *T // Pointer to the value.
	}

	// File is a struct that represents a file along with an iterator for reading lines.
	File struct {
		file  *os.File // Underlying os.File.
		name  String   // File name.
		guard bool     // Guard indicates whether the file is protected against concurrent access.
	}

	// Dir is a struct representing a directory path.
	Dir struct {
		path String // Directory path.
	}

	// String is an alias for the string type.
	String string

	// Int is an alias for the int type.
	Int int

	// Float is an alias for the float64 type.
	Float float64

	// Bytes is an alias for the []byte type.
	Bytes []byte

	// Slice is a generic alias for a slice.
	Slice[T any] []T

	// Map is a generic alias for a map.
	Map[K comparable, V any] map[K]V

	// Set is a generic alias for a set implemented using a map.
	Set[T comparable] map[T]struct{}

	// Pair is a struct representing a key-value Pair for MapOrd.
	Pair[K, V any] struct {
		Key   K // Key of the pair.
		Value V // Value associated with the key.
	}

	// MapOrd is a generic alias for a slice of ordered key-value pairs.
	MapOrd[K, V any] []Pair[K, V]

	// MapSafe is a thread-safe wrapper around a generic map.
	// It provides synchronized access to the underlying map to ensure
	// data consistency in concurrent environments.
	MapSafe[K comparable, V any] struct {
		mu   sync.RWMutex // Mutex to synchronize access to the map.
		data Map[K, V]    // The underlying map storing key-value pairs.
	}

	// SeqSet is an iterator over sequences of unique values.
	SeqSet[V comparable] iter.Seq[V]

	// SeqSlice is an iterator over sequences of individual values.
	SeqSlice[V any] iter.Seq[V]

	// SeqSlices is an iterator over slices of sequences of individual values.
	SeqSlices[V any] iter.Seq[[]V]

	// SeqMapOrd is an iterator over sequences of ordered pairs of values, most commonly ordered key-value pairs.
	SeqMapOrd[K, V any] iter.Seq2[K, V]

	// SeqMap is an iterator over sequences of pairs of values, most commonly key-value pairs.
	SeqMap[K comparable, V any] iter.Seq2[K, V]

	// Pool[T any] is a goroutine pool that allows parallel task execution.
	Pool[T any] struct {
		ctx           context.Context          // Context for controlling cancellation and timeouts
		cancel        context.CancelCauseFunc  // Function to cancel the context
		semaphore     chan struct{}            // Semaphore for limiting concurrency
		results       *MapSafe[int, Result[T]] // Stores task results
		wg            sync.WaitGroup           // Waits for all tasks to complete
		totalTasks    int32                    // Total number of tasks submitted
		activeTasks   int32                    // Number of currently active tasks
		failedTasks   int32                    // Number of failed tasks
		cancelOnError bool                     // Cancels remaining tasks if any task fails
	}
)
