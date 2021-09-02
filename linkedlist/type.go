package linkedlist

type Linkedlist interface {
	Insert(int) bool
	Delete(int) bool
	Contains(int) bool
	Range(f func(int) bool)
	Len() int
}
