package g_test

import (
	"sync"
	"testing"

	. "github.com/enetx/g"
)

func TestNewMutex(t *testing.T) {
	mutex := NewMutex(42)

	guard := mutex.Lock()
	defer guard.Unlock()

	if guard.Get() != 42 {
		t.Errorf("Expected initial value 42, got %d", guard.Get())
	}
}

func TestMutex_Lock_Unlock(t *testing.T) {
	mutex := NewMutex(0)

	guard := mutex.Lock()
	guard.Set(100)
	guard.Unlock()

	guard2 := mutex.Lock()
	defer guard2.Unlock()

	if guard2.Get() != 100 {
		t.Errorf("Expected value 100 after set, got %d", guard2.Get())
	}
}

func TestMutex_Get_Set(t *testing.T) {
	mutex := NewMutex("initial")

	guard := mutex.Lock()
	defer guard.Unlock()

	if guard.Get() != "initial" {
		t.Errorf("Expected 'initial', got '%s'", guard.Get())
	}

	guard.Set("modified")

	if guard.Get() != "modified" {
		t.Errorf("Expected 'modified', got '%s'", guard.Get())
	}
}

func TestMutex_Deref(t *testing.T) {
	mutex := NewMutex(10)

	guard := mutex.Lock()
	defer guard.Unlock()

	ptr := guard.Deref()
	*ptr = 20

	if guard.Get() != 20 {
		t.Errorf("Expected 20 after Deref modification, got %d", guard.Get())
	}
}

func TestMutex_TryLock_Success(t *testing.T) {
	mutex := NewMutex(42)

	opt := mutex.TryLock()
	if opt.IsNone() {
		t.Error("TryLock should succeed when mutex is unlocked")
	}

	guard := opt.Unwrap()
	defer guard.Unlock()

	if guard.Get() != 42 {
		t.Errorf("Expected 42, got %d", guard.Get())
	}
}

func TestMutex_TryLock_Failure(t *testing.T) {
	mutex := NewMutex(42)

	guard := mutex.Lock()
	defer guard.Unlock()

	opt := mutex.TryLock()
	if opt.IsSome() {
		opt.Unwrap().Unlock()
		t.Error("TryLock should fail when mutex is already locked")
	}
}

func TestMutex_Concurrent(t *testing.T) {
	mutex := NewMutex(0)
	var wg sync.WaitGroup

	iterations := 1000
	wg.Add(iterations)

	for range iterations {
		go func() {
			defer wg.Done()
			guard := mutex.Lock()
			defer guard.Unlock()
			current := guard.Get()
			guard.Set(current + 1)
		}()
	}

	wg.Wait()

	guard := mutex.Lock()
	defer guard.Unlock()

	if guard.Get() != iterations {
		t.Errorf("Expected %d, got %d", iterations, guard.Get())
	}
}

func TestMutex_WithStruct(t *testing.T) {
	type User struct {
		Name  string
		Count int
	}

	mutex := NewMutex(User{Name: "Alice", Count: 0})

	guard := mutex.Lock()
	user := guard.Deref()
	user.Name = "Bob"
	user.Count = 5
	guard.Unlock()

	guard2 := mutex.Lock()
	defer guard2.Unlock()

	result := guard2.Get()
	if result.Name != "Bob" || result.Count != 5 {
		t.Errorf("Expected User{Bob, 5}, got %+v", result)
	}
}

func TestMutex_WithSlice(t *testing.T) {
	mutex := NewMutex(SliceOf(1, 2, 3))

	guard := mutex.Lock()
	guard.Deref().Push(4, 5)
	guard.Unlock()

	guard2 := mutex.Lock()
	defer guard2.Unlock()

	slice := guard2.Get()
	if slice.Len() != 5 {
		t.Errorf("Expected slice length 5, got %d", slice.Len())
	}
}

func TestMutex_WithMap(t *testing.T) {
	mutex := NewMutex(NewMap[string, int]())

	guard := mutex.Lock()
	guard.Deref().Entry("key").OrInsert(42)
	guard.Unlock()

	guard2 := mutex.Lock()
	defer guard2.Unlock()

	val := guard2.Deref().Get("key")
	if val.IsNone() || val.Unwrap() != 42 {
		t.Errorf("Expected Some(42), got %v", val)
	}
}

func TestMutex_MultipleUnlocksSafe(t *testing.T) {
	mutex := NewMutex(0)

	guard := mutex.Lock()
	guard.Set(10)
	guard.Unlock()

	// Should be able to lock again
	guard2 := mutex.Lock()
	defer guard2.Unlock()

	if guard2.Get() != 10 {
		t.Errorf("Expected 10, got %d", guard2.Get())
	}
}

func TestMutex_ZeroValue(t *testing.T) {
	mutex := NewMutex(0)

	guard := mutex.Lock()
	defer guard.Unlock()

	if guard.Get() != 0 {
		t.Errorf("Expected zero value 0, got %d", guard.Get())
	}
}

func TestMutex_NilPointer(t *testing.T) {
	var ptr *int
	mutex := NewMutex(ptr)

	guard := mutex.Lock()
	defer guard.Unlock()

	if guard.Get() != nil {
		t.Error("Expected nil pointer")
	}

	newVal := 42
	guard.Set(&newVal)

	if *guard.Get() != 42 {
		t.Errorf("Expected pointer to 42, got %d", *guard.Get())
	}
}
