package g_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/enetx/g"
)

func TestNewRwLock(t *testing.T) {
	rwlock := NewRwLock(42)

	guard := rwlock.Read()
	defer guard.Unlock()

	if guard.Get() != 42 {
		t.Errorf("Expected initial value 42, got %d", guard.Get())
	}
}

func TestRwLock_Read(t *testing.T) {
	rwlock := NewRwLock("hello")

	guard := rwlock.Read()
	defer guard.Unlock()

	if guard.Get() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", guard.Get())
	}
}

func TestRwLock_Write(t *testing.T) {
	rwlock := NewRwLock(0)

	guard := rwlock.Write()
	guard.Set(100)
	guard.Unlock()

	readGuard := rwlock.Read()
	defer readGuard.Unlock()

	if readGuard.Get() != 100 {
		t.Errorf("Expected 100, got %d", readGuard.Get())
	}
}

func TestRwLock_WriteGuard_Get_Set(t *testing.T) {
	rwlock := NewRwLock(10)

	guard := rwlock.Write()
	defer guard.Unlock()

	if guard.Get() != 10 {
		t.Errorf("Expected 10, got %d", guard.Get())
	}

	guard.Set(20)

	if guard.Get() != 20 {
		t.Errorf("Expected 20, got %d", guard.Get())
	}
}

func TestRwLock_ReadGuard_Deref(t *testing.T) {
	rwlock := NewRwLock(42)

	guard := rwlock.Read()
	defer guard.Unlock()

	ptr := guard.Deref()
	if *ptr != 42 {
		t.Errorf("Expected Deref to return pointer to 42, got %d", *ptr)
	}
}

func TestRwLock_WriteGuard_Deref(t *testing.T) {
	rwlock := NewRwLock(10)

	guard := rwlock.Write()
	defer guard.Unlock()

	ptr := guard.Deref()
	*ptr = 30

	if guard.Get() != 30 {
		t.Errorf("Expected 30 after Deref modification, got %d", guard.Get())
	}
}

func TestRwLock_TryRead_Success(t *testing.T) {
	rwlock := NewRwLock(42)

	opt := rwlock.TryRead()
	if opt.IsNone() {
		t.Error("TryRead should succeed when rwlock is unlocked")
	}

	guard := opt.Unwrap()
	defer guard.Unlock()

	if guard.Get() != 42 {
		t.Errorf("Expected 42, got %d", guard.Get())
	}
}

func TestRwLock_TryWrite_Success(t *testing.T) {
	rwlock := NewRwLock(42)

	opt := rwlock.TryWrite()
	if opt.IsNone() {
		t.Error("TryWrite should succeed when rwlock is unlocked")
	}

	guard := opt.Unwrap()
	defer guard.Unlock()

	guard.Set(100)
	if guard.Get() != 100 {
		t.Errorf("Expected 100, got %d", guard.Get())
	}
}

func TestRwLock_TryWrite_FailsWhenReadLocked(t *testing.T) {
	rwlock := NewRwLock(42)

	readGuard := rwlock.Read()
	defer readGuard.Unlock()

	opt := rwlock.TryWrite()
	if opt.IsSome() {
		opt.Unwrap().Unlock()
		t.Error("TryWrite should fail when read lock is held")
	}
}

func TestRwLock_TryRead_FailsWhenWriteLocked(t *testing.T) {
	rwlock := NewRwLock(42)

	writeGuard := rwlock.Write()
	defer writeGuard.Unlock()

	opt := rwlock.TryRead()
	if opt.IsSome() {
		opt.Unwrap().Unlock()
		t.Error("TryRead should fail when write lock is held")
	}
}

func TestRwLock_MultipleReaders(t *testing.T) {
	rwlock := NewRwLock(42)
	var wg sync.WaitGroup

	readers := 10
	wg.Add(readers)

	for range readers {
		go func() {
			defer wg.Done()
			guard := rwlock.Read()
			defer guard.Unlock()

			if guard.Get() != 42 {
				t.Errorf("Expected 42, got %d", guard.Get())
			}

			// Hold the lock briefly
			time.Sleep(10 * time.Millisecond)
		}()
	}

	wg.Wait()
}

func TestRwLock_ConcurrentReadWrite(t *testing.T) {
	rwlock := NewRwLock(0)
	var wg sync.WaitGroup

	writers := 10
	readers := 100
	wg.Add(writers + readers)

	// Writers
	for range writers {
		go func() {
			defer wg.Done()
			guard := rwlock.Write()
			defer guard.Unlock()
			current := guard.Get()
			guard.Set(current + 1)
		}()
	}

	// Readers
	for range readers {
		go func() {
			defer wg.Done()
			guard := rwlock.Read()
			defer guard.Unlock()
			_ = guard.Get() // Just read
		}()
	}

	wg.Wait()

	guard := rwlock.Read()
	defer guard.Unlock()

	if guard.Get() != writers {
		t.Errorf("Expected %d, got %d", writers, guard.Get())
	}
}

func TestRwLock_WithStruct(t *testing.T) {
	type Config struct {
		Host string
		Port int
	}

	rwlock := NewRwLock(Config{Host: "localhost", Port: 8080})

	// Read
	readGuard := rwlock.Read()
	config := readGuard.Get()
	readGuard.Unlock()

	if config.Host != "localhost" || config.Port != 8080 {
		t.Errorf("Expected {localhost, 8080}, got %+v", config)
	}

	// Write
	writeGuard := rwlock.Write()
	writeGuard.Set(Config{Host: "0.0.0.0", Port: 9000})
	writeGuard.Unlock()

	// Read again
	readGuard2 := rwlock.Read()
	defer readGuard2.Unlock()
	config2 := readGuard2.Get()

	if config2.Host != "0.0.0.0" || config2.Port != 9000 {
		t.Errorf("Expected {0.0.0.0, 9000}, got %+v", config2)
	}
}

func TestRwLock_WithMap(t *testing.T) {
	rwlock := NewRwLock(NewMap[string, int]())

	// Write
	guard := rwlock.Write()
	guard.Deref().Entry("a").OrInsert(1)
	guard.Deref().Entry("b").OrInsert(2)
	guard.Unlock()

	// Read
	readGuard := rwlock.Read()
	defer readGuard.Unlock()

	m := readGuard.Deref()
	if m.Get("a").UnwrapOr(0) != 1 {
		t.Error("Expected a=1")
	}
	if m.Get("b").UnwrapOr(0) != 2 {
		t.Error("Expected b=2")
	}
}

func TestRwLock_WithSlice(t *testing.T) {
	rwlock := NewRwLock(SliceOf(1, 2, 3))

	// Write
	guard := rwlock.Write()
	guard.Deref().Push(4, 5)
	guard.Unlock()

	// Read
	readGuard := rwlock.Read()
	defer readGuard.Unlock()

	if readGuard.Deref().Len() != 5 {
		t.Errorf("Expected length 5, got %d", readGuard.Deref().Len())
	}
}

func TestRwLock_ReadersNotBlocked(t *testing.T) {
	rwlock := NewRwLock(42)

	var readersStarted atomic.Int32
	var readersFinished atomic.Int32

	// Start multiple readers
	for range 5 {
		go func() {
			guard := rwlock.Read()
			readersStarted.Add(1)
			time.Sleep(50 * time.Millisecond)
			_ = guard.Get()
			guard.Unlock()
			readersFinished.Add(1)
		}()
	}

	// Wait for readers to start
	time.Sleep(20 * time.Millisecond)

	// All readers should have started (not blocked by each other)
	if readersStarted.Load() != 5 {
		t.Errorf("Expected all 5 readers to start concurrently, got %d", readersStarted.Load())
	}

	// Wait for completion
	time.Sleep(100 * time.Millisecond)

	if readersFinished.Load() != 5 {
		t.Errorf("Expected all 5 readers to finish, got %d", readersFinished.Load())
	}
}

func TestRwLock_ZeroValue(t *testing.T) {
	rwlock := NewRwLock(0)

	guard := rwlock.Read()
	defer guard.Unlock()

	if guard.Get() != 0 {
		t.Errorf("Expected zero value 0, got %d", guard.Get())
	}
}

func TestRwLock_NilPointer(t *testing.T) {
	var ptr *int
	rwlock := NewRwLock(ptr)

	readGuard := rwlock.Read()
	if readGuard.Get() != nil {
		t.Error("Expected nil pointer")
	}
	readGuard.Unlock()

	writeGuard := rwlock.Write()
	newVal := 42
	writeGuard.Set(&newVal)
	writeGuard.Unlock()

	readGuard2 := rwlock.Read()
	defer readGuard2.Unlock()

	if *readGuard2.Get() != 42 {
		t.Errorf("Expected pointer to 42, got %d", *readGuard2.Get())
	}
}

func TestRwLock_TryReadMultiple(t *testing.T) {
	rwlock := NewRwLock(42)

	// Should be able to acquire multiple read locks via TryRead
	opt1 := rwlock.TryRead()
	opt2 := rwlock.TryRead()
	opt3 := rwlock.TryRead()

	if opt1.IsNone() || opt2.IsNone() || opt3.IsNone() {
		t.Error("Should be able to acquire multiple read locks")
	}

	opt1.Unwrap().Unlock()
	opt2.Unwrap().Unlock()
	opt3.Unwrap().Unlock()
}

func TestRwLock_With(t *testing.T) {
	rwlock := NewRwLock(10)
	rwlock.With(func(v *int) {
		*v = 42
	})
	guard := rwlock.Read()
	defer guard.Unlock()
	if guard.Get() != 42 {
		t.Errorf("Expected 42 after With, got %d", guard.Get())
	}
}

func TestRwLock_RWith(t *testing.T) {
	rwlock := NewRwLock(42)
	var got int
	rwlock.RWith(func(v int) {
		got = v
	})
	if got != 42 {
		t.Errorf("Expected 42 from RWith, got %d", got)
	}
}

func TestRwLock_RWith_ReceivesCopy(t *testing.T) {
	rwlock := NewRwLock(SliceOf(1, 2, 3))
	rwlock.RWith(func(sl Slice[int]) {
		if sl.Len() != 3 {
			t.Errorf("Expected slice length 3, got %d", sl.Len())
		}
	})
}

func TestRwLock_With_Struct(t *testing.T) {
	type Config struct {
		Host string
		Port int
	}
	rwlock := NewRwLock(Config{Host: "localhost", Port: 8080})
	rwlock.With(func(c *Config) {
		c.Host = "0.0.0.0"
		c.Port = 9000
	})
	guard := rwlock.Read()
	defer guard.Unlock()
	config := guard.Get()
	if config.Host != "0.0.0.0" || config.Port != 9000 {
		t.Errorf("Expected {0.0.0.0, 9000}, got %+v", config)
	}
}

func TestRwLock_With_Slice(t *testing.T) {
	rwlock := NewRwLock(SliceOf(1, 2, 3))
	rwlock.With(func(sl *Slice[int]) {
		sl.Push(4, 5)
	})
	guard := rwlock.Read()
	defer guard.Unlock()
	if guard.Deref().Len() != 5 {
		t.Errorf("Expected length 5, got %d", guard.Deref().Len())
	}
}

func TestRwLock_With_Map(t *testing.T) {
	rwlock := NewRwLock(NewMap[string, int]())
	rwlock.With(func(m *Map[string, int]) {
		m.Entry("a").OrInsert(1)
		m.Entry("b").OrInsert(2)
	})
	guard := rwlock.Read()
	defer guard.Unlock()
	if guard.Deref().Get("a").UnwrapOr(0) != 1 {
		t.Error("Expected a=1")
	}
	if guard.Deref().Get("b").UnwrapOr(0) != 2 {
		t.Error("Expected b=2")
	}
}

func TestRwLock_With_Concurrent(t *testing.T) {
	rwlock := NewRwLock(0)
	var wg sync.WaitGroup

	writers := 100
	wg.Add(writers)
	for range writers {
		go func() {
			defer wg.Done()
			rwlock.With(func(v *int) { *v++ })
		}()
	}
	wg.Wait()

	guard := rwlock.Read()
	defer guard.Unlock()
	if guard.Get() != writers {
		t.Errorf("Expected %d, got %d", writers, guard.Get())
	}
}

func TestRwLock_RWith_Concurrent(t *testing.T) {
	rwlock := NewRwLock(42)
	var wg sync.WaitGroup
	var readersActive atomic.Int32

	readers := 10
	wg.Add(readers)
	for range readers {
		go func() {
			defer wg.Done()
			rwlock.RWith(func(v int) {
				readersActive.Add(1)
				time.Sleep(20 * time.Millisecond)
				if v != 42 {
					t.Errorf("Expected 42, got %d", v)
				}
				readersActive.Add(-1)
			})
		}()
	}

	time.Sleep(10 * time.Millisecond)
	if readersActive.Load() < 2 {
		t.Error("Expected multiple concurrent readers in RWith")
	}

	wg.Wait()
}

func TestRwLock_With_PanicUnlocks(t *testing.T) {
	rwlock := NewRwLock(42)
	func() {
		defer func() { recover() }()
		rwlock.With(func(v *int) {
			*v = 100
			panic("test panic")
		})
	}()

	// Lock should be released after panic
	guard := rwlock.Read()
	defer guard.Unlock()
	if guard.Get() != 100 {
		t.Errorf("Expected 100 after panic in With, got %d", guard.Get())
	}
}

func TestRwLock_RWith_PanicUnlocks(t *testing.T) {
	rwlock := NewRwLock(42)
	func() {
		defer func() { recover() }()
		rwlock.RWith(func(v int) {
			_ = v
			panic("test panic")
		})
	}()

	// Lock should be released after panic
	guard := rwlock.Write()
	defer guard.Unlock()
	guard.Set(99)
	if guard.Get() != 99 {
		t.Errorf("Expected 99 after panic in RWith, got %d", guard.Get())
	}
}
