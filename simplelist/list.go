package simplelist

import (
	"linkedlist/linkedlist"
	"strconv"
	"strings"
	"sync"
)

type IntList struct {
	head   *intNode
	length int64
	mu     sync.RWMutex
}

type intNode struct {
	value int
	next  *intNode
}

func newIntNode(value int) *intNode {
	return &intNode{value: value}
}

func NewIntList() linkedlist.Linkedlist {
	return &IntList{head: newIntNode(0), length: 0}
}

// Insert inserts a new intNode with value to the IntList
// return if the operation successed
func (l *IntList) Insert(value int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	a := l.head
	b := a.next
	for b != nil && b.value < value {
		a, b = b, b.next
	}
	// Check if the node exists
	if b != nil && b.value == value {
		// value exist do not insert
		return false
	}
	x := newIntNode(value)
	x.next = b
	a.next = x
	l.length++
	return true
}

// Delete deletes a intNode of the IntList with value,
// return if the operation successed
func (l *IntList) Delete(value int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	a := l.head
	b := a.next
	for b != nil && b.value < value {
		a, b = b, b.next
	}
	// Check if the node not exists
	if b == nil || b.value != value {
		// value exist do not insert
		return false
	}
	a.next = b.next
	l.length--
	return true
}

// Contains check if an intNode with value exists in the IntList
func (l *IntList) Contains(value int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	x := l.head.next
	for x != nil && x.value < value {
		x = x.next
	}
	if x == nil {
		return false
	}
	return x.value == value
}

// Range implements the Range operation
func (l *IntList) Range(f func(value int) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	x := l.head.next
	for x != nil {
		if !f(x.value) {
			break
		}
		x = x.next
	}
}

// Len returns the length of IntList
func (l *IntList) Len() int {
	return int(l.length)
}

func (l *IntList) String() string {
	if l == nil {
		return "nil"
	}
	var builder strings.Builder
	builder.WriteString("[]int{")
	x := l.head.next
	for x != nil {
		builder.WriteString(strconv.FormatInt(int64(x.value), 10))
		if x.next != nil {
			builder.WriteString(",")
		} else {
			builder.WriteString("}")
		}
		x = x.next
	}
	return builder.String()
}
