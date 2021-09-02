package linkedlist

import (
	"linkedlist/linkedlist"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

type IntList struct {
	head   *intNode
	length int64
}

type intNode struct {
	value  int
	marked uint32
	next   *intNode
	mu     sync.Mutex
}

func newIntNode(value int) *intNode {
	return &intNode{value: value}
}

func (n *intNode) loadNext() *intNode {
	return (*intNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next))))
}

func (n *intNode) storeNext(node *intNode) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.next)), unsafe.Pointer(node))
}

func (n *intNode) setMarked() {
	atomic.StoreUint32(&n.marked, 1)
}

func (n *intNode) isMakred() bool {
	return atomic.LoadUint32(&n.marked) == 1
}

func NewIntList() linkedlist.Linkedlist {
	return &IntList{head: newIntNode(0)}
}

// Insert inserts a new intNode with value to the IntList
// return if the operation successed
func (l *IntList) Insert(value int) bool {
	for {
		// Step1 Find a and b
		a := l.head
		b := a.loadNext()
		for b != nil && b.value < value {
			a = b
			b = b.loadNext()
		}
		// check if node exists
		if b != nil && b.value == value {
			if b.isMakred() {
				// The node has been marked, try next time
				continue
			}
			return false
		}
		// Step 2. Lock a
		a.mu.Lock()
		if a.next != b || a.isMakred() {
			// check if a.next == b && !a.isMarked() is true
			a.mu.Unlock()
			continue
		}
		// Step 3. Create new node x
		x := newIntNode(value)
		// Step 4. x.next = b; a.next = x
		x.next = b
		a.storeNext(x)
		// Step 5. Unlock a
		a.mu.Unlock()
		atomic.AddInt64(&l.length, 1)
		return true
	}
}

// Delete deletes a intNode of the IntList with value,
// return if the operation successed
func (l *IntList) Delete(value int) bool {
	for {
		// Step1 Find a and b
		a := l.head
		b := a.loadNext()
		for b != nil && b.value < value {
			a = b
			b = b.loadNext()
		}
		// check if node not exists
		if b == nil || b.value != value {
			return false
		}
		// Step 2. Lock b
		b.mu.Lock()
		if b.isMakred() {
			// Check if b has been deleted or another goroutine has delete it.
			b.mu.Unlock()
			return false
		}
		// Step 3. Lock a
		if a.next != b || a.isMakred() {
			// check if a.next == b && !a.isMarked() is true
			a.mu.Unlock()
			b.mu.Unlock()
			continue
		}
		// Step 4. Mark node and delete
		b.setMarked()
		a.storeNext(b.next)
		atomic.AddInt64(&l.length, -1)
		// Step 5. Unlock a and b
		a.mu.Unlock()
		b.mu.Unlock()
		return true
	}
}

// Contains check if an intNode with value exists in the IntList
func (l *IntList) Contains(value int) bool {
	x := l.head.loadNext()
	for x != nil && x.value < value {
		x = x.loadNext()
	}
	if x == nil {
		return false
	}
	return x.value == value && !x.isMakred()
}

// Range implements the Range operation
func (l *IntList) Range(f func(value int) bool) {
	x := l.head.loadNext()
	for x != nil {
		if x.isMakred() {
			x = x.loadNext()
			continue
		}
		if !f(x.value) {
			break
		}
		x = x.loadNext()
	}
}

// Len returns the length of IntList
func (l *IntList) Len() int {
	return int(atomic.LoadInt64(&l.length))
}

func (l *IntList) String() string {
	if l == nil {
		return "nil"
	}
	var builder strings.Builder
	builder.WriteString("[]int{")
	x := l.head.loadNext()
	for x != nil {
		builder.WriteString(strconv.FormatInt(int64(x.value), 10))
		if x.next != nil {
			builder.WriteString(",")
		} else {
			builder.WriteString("}")
		}
		x = x.loadNext()
	}
	return builder.String()
}
