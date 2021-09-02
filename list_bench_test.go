package linkedlist

import (
	"linkedlist/simplelist"
	"math"
	"testing"
)

const randN = math.MaxUint32

func BenchmarkInsert(b *testing.B) {
	b.Run("simplelist", func(b *testing.B) {
		l := simplelist.NewIntList()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int(fastrandn(randN)))
			}
		})
	})
	b.Run("concurrentlist", func(b *testing.B) {
		l := NewIntList()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int(fastrandn(randN)))
			}
		})
	})
}

func BenchmarkInsertDupl(b *testing.B) {
	var randN = uint32(500)
	b.Run("simplelist", func(b *testing.B) {
		l := simplelist.NewIntList()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int(fastrandn(randN)))
			}
		})
	})
	b.Run("concurrentlist", func(b *testing.B) {
		l := NewIntList()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int(fastrandn(randN)))
			}
		})
	})
}

func Benchmark30Insert70Contains(b *testing.B) {
	b.Run("simplelist", func(b *testing.B) {
		l := simplelist.NewIntList()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Insert(int(fastrandn(randN)))
				} else {
					l.Contains(int(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("concurrentlist", func(b *testing.B) {
		l := NewIntList()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Insert(int(fastrandn(randN)))
				} else {
					l.Contains(int(fastrandn(randN)))
				}
			}
		})
	})
}

func Benchmark1Delete9Insert90Contains(b *testing.B) {
	const initsize = 1000
	b.Run("simplelist", func(b *testing.B) {
		l := simplelist.NewIntList()
		for i := 0; i < initsize; i++ {
			l.Insert(i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u < 9 {
					l.Insert(int(fastrandn(randN)))
				} else if u == 10 {
					l.Delete(int(fastrandn(randN)))
				} else {
					l.Contains(int(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("concurrentlist", func(b *testing.B) {
		l := NewIntList()
		for i := 0; i < initsize; i++ {
			l.Insert(i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u < 9 {
					l.Insert(int(fastrandn(randN)))
				} else if u == 10 {
					l.Delete(int(fastrandn(randN)))
				} else {
					l.Contains(int(fastrandn(randN)))
				}
			}
		})
	})
}
