package lexing

import (
	"errors"
	"testing"
)

func TestRingBufferPushAndPop(t *testing.T) {
	buf := NewRingBuffer[int](10)
	ExpectEqual(t, "capacity", 10, buf.Capacity())
	ExpectEqual(t, "size", 0, buf.Size())

	buf.Push(123)
	ExpectEqual(t, "size", 1, buf.Size())

	item, err := buf.Pop()
	if err != nil {
		t.Error(err)
	}

	ExpectEqual(t, "popped value", 123, item)
	ExpectEqual(t, "size", 0, buf.Size())
	ExpectEqual(t, "front", 1, buf.front)
	ExpectEqual(t, "back", 1, buf.back)

	_, err = buf.Pop()
	if !errors.Is(ErrBufferEmpty, err) {
		t.Error("buf is empty, should have errored")
	}
}

func TestRingBufferPushOverInitialCapacity(t *testing.T) {
	buf := NewRingBuffer[int](3)
	ExpectEqual(t, "capacity", 3, buf.Capacity())
	ExpectEqual(t, "size", 0, buf.Size())

	buf.Push(123)
	buf.Push(456)
	buf.Push(789)
	buf.Push(111)
	buf.Push(222)
	buf.Push(333)

	ExpectEqual(t, "size", 6, buf.Size())
	ExpectEqual(t, "capacity", 12, buf.Capacity())
	ExpectEqual(t, "front", 0, buf.front)
	ExpectEqual(t, "back", 6, buf.back)
}

func TestRingBufferPushAndPopOverCapacity(t *testing.T) {
	buf := NewRingBuffer[int](4)
	ExpectEqual(t, "capacity", 4, buf.Capacity())
	ExpectEqual(t, "size", 0, buf.Size())

	buf.Push(123)
	buf.Push(456)
	buf.Push(789)
	buf.Push(111)
	buf.Push(222)

	val, err := buf.Pop()
	if err != nil {
		t.Error(err)
	}
	ExpectEqual(t, "popped value", 123, val)
	ExpectEqual(t, "front", 1, buf.front)
	ExpectEqual(t, "back", 5, buf.back)

	val, err = buf.Pop()
	if err != nil {
		t.Error(err)
	}
	ExpectEqual(t, "popped value", 456, val)
	ExpectEqual(t, "front", 2, buf.front)
	ExpectEqual(t, "back", 5, buf.back)

	buf.Push(666)
	ExpectEqual(t, "front", 2, buf.front)
	ExpectEqual(t, "back", 6, buf.back)
	ExpectEqual(t, "size", 4, buf.Size())
}

func TestRingBufferRandomAccess(t *testing.T) {
	buf := NewRingBuffer[int](10)
	ExpectEqual(t, "capacity", 10, buf.Capacity())
	ExpectEqual(t, "size", 0, buf.Size())

	buf.Push(123)
	buf.Push(456)
	buf.Push(789)
	buf.Push(111)
	buf.Push(222)
	buf.Push(333)

	ExpectEqual(t, "size", 6, buf.Size())
	ExpectEqual(t, "capacity", 10, buf.Capacity())
	ExpectEqual(t, "front", 0, buf.front)
	ExpectEqual(t, "back", 6, buf.back)

	val, ok := buf.At(0)
	ExpectEqual(t, "found", true, ok)
	ExpectEqual(t, "value", 123, val)

	val, ok = buf.At(2)
	ExpectEqual(t, "found", true, ok)
	ExpectEqual(t, "value", 789, val)

	_, ok = buf.At(6)
	ExpectEqual(t, "found", false, ok)

	_, ok = buf.At(11)
	ExpectEqual(t, "found", false, ok)
}

func ExpectEqual[T comparable](t *testing.T, name string, expected T, actual T) {
	if actual != expected {
		t.Errorf("expected %s to be %v, but it was %v", name, expected, actual)
	}
}
