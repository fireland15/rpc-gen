package parser

import "errors"

var (
	ErrBufferEmpty = errors.New("ring buffer empty")
)

type RingBuffer[T any] struct {
	data  []T
	front int
	back  int
}

func NewRingBuffer[T any](cap int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:  make([]T, cap),
		front: 0,
		back:  0,
	}
}

func (b *RingBuffer[T]) Capacity() int {
	return len(b.data)
}

func (b *RingBuffer[T]) Size() int {
	if b.back < b.front {
		return len(b.data) - b.front + b.back
	}
	return b.back - b.front
}

func (b *RingBuffer[T]) IsEmpty() bool {
	return b.Size() == 0
}

func (b *RingBuffer[T]) Push(item T) {
	if (b.back+1)%len(b.data) == b.front {
		b.resize()
	}
	b.data[b.back] = item
	b.back = (b.back + 1) % len(b.data)
}

func (b *RingBuffer[T]) Pop() (T, error) {
	var t T
	if b.front == b.back {
		return t, ErrBufferEmpty
	}
	t = b.data[b.front]
	b.front = (b.front + 1) % len(b.data)
	return t, nil
}

func (q *RingBuffer[T]) At(n int) (T, bool) {
	var val T
	if n >= q.Size() {
		return val, false
	}

	idx := (q.front + n) % len(q.data)
	return q.data[idx], true
}

func (q *RingBuffer[T]) resize() {
	cap := len(q.data) * 2
	back := len(q.data) - 1
	newData := make([]T, cap)
	if q.front > q.back {
		copy(newData, q.data[q.front:q.back-1])
	} else if q.front <= q.back {
		copy(newData, q.data[q.front:])
		copy(newData, q.data[:q.back-1])
	}

	q.data = newData
	q.front = 0
	q.back = back
}
