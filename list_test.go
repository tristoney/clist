package linkedlist

import (
	"fmt"
	"linkedlist/linkedlist"
	"linkedlist/simplelist"
	"sync"
	"sync/atomic"
	"testing"
	_ "unsafe"
)

//go:linkname fastrand runtime.fastrand
func fastrand() uint32

//go:nosplit
func fastrandn(n uint32) uint32 {
	return uint32(uint64(fastrand()) * uint64(n) >> 32)
}

func testLinkedList(t *testing.T, f func() linkedlist.Linkedlist) {
	// correctness test
	l := f()
	if l.Len() != 0 {
		t.Fatal("invalid length")
	}
	if l.Contains(0) {
		t.Fatal("invalid contains")
	}
	if !l.Insert(0) || l.Len() != 1 {
		t.Fatal("invalid insert")
	}
	if !l.Contains(0) {
		t.Fatal("invalid contains")
	}

	if !l.Insert(20) || l.Len() != 2 {
		t.Fatal("invalid insert")
	}
	if !l.Insert(22) || l.Len() != 3 {
		t.Fatal("invalid insert")
	}
	if !l.Insert(21) || l.Len() != 4 {
		t.Fatal("invalid insert")
	}
	// fmt.Println(l)

	l = simplelist.NewIntList()
	const num = 1000
	// Make rand shuffle array.
	// The testArray contains [1,num]
	testArray := make([]int, num)
	testArray[0] = num + 1
	for i := 1; i < num; i++ {
		// We left 0, because it is the default score for head and tail.
		// If we check the skiplist contains 0, there must be something wrong.
		testArray[i] = int(i)
	}
	for i := len(testArray) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := fastrandn(uint32(i + 1))
		testArray[i], testArray[j] = testArray[j], testArray[i]
	}

	// Concurrent insert.
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			l.Insert(testArray[i])
			wg.Done()
		}()
	}
	wg.Wait()

	if l.Len() != num {
		t.Fatalf("invalid length expected %d, got %d", num, l.Len())
	}

	// Don't contains 0 after concurrent insertion.
	if l.Contains(0) {
		t.Fatal("contains 0 after concurrent insertion")
	}

	// Concurrent contains.
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			if !l.Contains(testArray[i]) {
				wg.Done()
				panic(fmt.Sprintf("insert doesn't contains %d", i))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Test all methods.
	var tmp uint64
	var smallZone uint64 = 100
	l = simplelist.NewIntList()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			r := fastrandn(2)
			if r == 0 {
				l.Insert(int(atomic.AddUint64(&tmp, 1) % smallZone))
			} else {
				l.Contains(int(fastrandn(uint32(smallZone))))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if l.Len() != int(smallZone) {
		t.Fatal("invalid length")
	}

	// Concurrent Insert and Delete in small zone.
	x := simplelist.NewIntList()
	var (
		insertcount uint64 = 0
		deletecount uint64 = 0
	)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				if fastrandn(2) == 0 {
					if x.Delete(int(fastrandn(10))) {
						atomic.AddUint64(&deletecount, 1)
					}
				} else {
					if x.Insert(int(fastrandn(10))) {
						atomic.AddUint64(&insertcount, 1)
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(insertcount, deletecount)
	fmt.Println(x)
}

func TestLinkedList(t *testing.T) {
	t.Run("simplelist", func(t *testing.T) {
		testLinkedList(t, simplelist.NewIntList)
	})
	t.Run("clist", func(t *testing.T) {
		testLinkedList(t, NewIntList)
	})
}
