package sync

import (
	"runtime"
	"sync"
	"testing"
)

// A testData is the state shared between the threads in each of the tests below.
type testData struct {
	nThreads  int // Number of test threads; constant after init.
	loopCount int // Iteration count for each test thread; constant after init.

	mu *Mutex // Protects i, id, and finishedThreads.
	i  int    // Counter incremented by test loops.
	id int    // id of current lock-holding thread in some tests.

	mutex sync.Mutex // Protects i and id when in countingLoopMutex.

}

func TestMuNThread(t *testing.T) {
	wg := sync.WaitGroup{}

	td := testData{nThreads: 5, loopCount: 1000000, mu: NewMutex(1)}
	for i := 0; i != td.nThreads; i++ {
		wg.Add(1)
		go countingLoopMu(&td, i, &wg)
	}

	wg.Wait()

	if td.i != td.nThreads*td.loopCount {
		t.Fatalf("TestMuNThread final count inconsistent: want %d, got %d",
			td.nThreads*td.loopCount, td.i)
	}
}

func countingLoopMu(td *testData, id int, wg *sync.WaitGroup) {
	var n = td.loopCount
	for i := 0; i != n; i++ {
		td.mu.Lock()
		td.id = id
		td.i++
		if td.id != id {
			panic("td.id != id")
		}
		td.mu.Unlock()
	}

	wg.Done()
}

func TestTryMuNThread(t *testing.T) {
	wg := sync.WaitGroup{}

	td := testData{nThreads: 5, loopCount: 100000, mu: NewMutex(1)}
	for i := 0; i != td.nThreads; i++ {
		wg.Add(1)
		go countingLoopTryMu(&td, i, &wg)
	}

	wg.Wait()
	if td.i != td.nThreads*td.loopCount {
		t.Fatalf("TestTryMuNThread final count inconsistent: want %d, got %d",
			td.nThreads*td.loopCount, td.i)
	}
}

func countingLoopTryMu(td *testData, id int, wg *sync.WaitGroup) {
	var n = td.loopCount
	for i := 0; i != n; i++ {
		for !td.mu.TryLock() {
			runtime.Gosched()
		}
		td.id = id
		td.i++
		if td.id != id {
			panic("td.id != id")
		}
		td.mu.Unlock()
	}

	wg.Done()
}

func countingLoopMutex(td *testData, id int, wg *sync.WaitGroup) {
	var n = td.loopCount
	for i := 0; i != n; i++ {
		td.mutex.Lock()
		td.id = id
		td.i++
		if td.id != id {
			panic("td.id != id")
		}
		td.mutex.Unlock()
	}

	wg.Done()
}

func TestMutexNThread(t *testing.T) {
	wg := sync.WaitGroup{}

	td := testData{nThreads: 5, loopCount: 1000000, mu: NewMutex(1)}

	for i := 0; i != td.nThreads; i++ {
		wg.Add(1)
		go countingLoopMutex(&td, i, &wg)
	}

	wg.Wait()
	if td.i != td.nThreads*td.loopCount {
		t.Fatalf("TestMutexNThread final count inconsistent: want %d, got %d",
			td.nThreads*td.loopCount, td.i)
	}
}
