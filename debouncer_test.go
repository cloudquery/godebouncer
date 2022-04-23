package godebouncer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/vnteamopen/godebouncer"
)

func Example() {
	wait := 10 * time.Second
	debouncer := godebouncer.New(wait).WithTriggered(func() {
		fmt.Println("Trigger") // Triggered func will be called after 10 seconds from last SendSignal().
	})

	fmt.Println("Action 1")
	debouncer.SendSignal()

	time.Sleep(3 * time.Second)

	fmt.Println("Action 2")
	debouncer.SendSignal()
	// After 10 seconds, the trigger will be called.
	//Previous `SendSignal()` will be ignore to trigger the triggered function.

	time.Sleep(10 * time.Second)
}

func createIncrementCount(counter int) (*int, func()) {
	return &counter, func() {
		fmt.Println("Triggered")
		counter++
	}
}

func TestDebounceDoBeforeExpired(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(1)

	debouncer.Do(func() {
		fmt.Println("Action 1")
	})

	time.Sleep(50 * time.Millisecond)

	debouncer.Do(func() {
		fmt.Println("Action 2")
	})

	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceDoAfterExpired(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(2)

	debouncer.Do(func() {
		fmt.Println("Action 1")
	})

	time.Sleep(400 * time.Millisecond)

	debouncer.Do(func() {
		fmt.Println("Action 2")
	})

	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceMixed(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(2)

	debouncer.Do(func() {
		fmt.Println("Action 1")
	})

	debouncer.Do(func() {
		fmt.Println("Action 2")
	})

	time.Sleep(400 * time.Millisecond)

	debouncer.Do(func() {
		fmt.Println("Action 3")
	})

	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceWithoutTriggeredFunc(t *testing.T) {
	debouncer := godebouncer.New(200 * time.Millisecond)

	debouncer.Do(func() {
		fmt.Println("Action 1")
	})
	time.Sleep(400 * time.Millisecond)
	fmt.Println("debouncer.Do() finished successfully!")
}

func TestDebounceSendSignal(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(1)

	debouncer.SendSignal()
	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceUpdateTriggeredFuncBeforeDuration(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(2)

	debouncer.SendSignal()
	time.Sleep(100 * time.Millisecond)

	debouncer.UpdateTriggeredFunc(func() {
		*countPtr += 2
	})
	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceUpdateTriggeredFuncAfterDuration(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(3)

	debouncer.SendSignal()
	time.Sleep(400 * time.Millisecond)

	debouncer.UpdateTriggeredFunc(func() {
		*countPtr += 2
	})
	debouncer.SendSignal()
	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceCancel(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(0)

	debouncer.SendSignal()
	time.Sleep(100 * time.Millisecond)

	debouncer.Cancel()
	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceUpdateDuration(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(600 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(1)

	debouncer.UpdateTimeDuration(200 * time.Millisecond)
	debouncer.SendSignal()
	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceUpdateDurationAfterSendSignal(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	debouncer := godebouncer.New(400 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(1)

	debouncer.SendSignal()
	time.Sleep(200 * time.Millisecond)

	debouncer.UpdateTimeDuration(600 * time.Millisecond)
	time.Sleep(300 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}
