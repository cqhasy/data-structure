// sequential storage

package List

import (
	"log"
)

const (
	DefaultCap = 10
)

type SeqList[T any] struct {
	Cap  int
	Len  int
	Data []T
}

func NewSeqList[T any](cap int) *SeqList[T] {
	if cap <= 0 {
		panic("cap must be greater than zero")
	}
	return &SeqList[T]{
		Cap:  cap,
		Len:  0,
		Data: make([]T, 0, cap),
	}
}

// Append adds an element to the end of the sequential list
func (sl *SeqList[T]) Append(v T) {
	// 检查当前长度，是否需要扩容
	if sl.Len == sl.Cap {
		/*
			这里始终默认二倍扩容，但go的底层会帮助你扩容且有更好的扩容策略，
			这里只是想尝试自己实现
		*/
		sl.Cap += sl.Cap
		// 准备迁移数据
		newData := make([]T, sl.Len, sl.Cap)
		copy(newData, sl.Data)
		sl.Data = newData
	}

	sl.Data = append(sl.Data, v)
	sl.Len++
}

// Remove 删除对应索引的值
func (sl *SeqList[T]) Remove(i int) (T, bool) {
	if i < 0 || i >= sl.Len {
		var zero T
		log.Println("index out of range")
		return zero, false
	}

	var data = sl.Data[i]
	for j := i; j < sl.Len-1; j++ {
		sl.Data[j] = sl.Data[j+1]
	}

	var zero T
	sl.Data[sl.Len-1] = zero
	sl.Len--

	return data, true
}

func (sl *SeqList[T]) Get(i int) (T, bool) {
	if i < 0 || i >= sl.Len {
		var zero T
		return zero, false
	}

	var data = sl.Data[i]
	return data, true
}

// Locate 返回检索到的第一个相等的数据的下标
func (sl *SeqList[T]) Locate(data T, comparator func(a, b T) bool) (int, bool) {
	for i := 0; i < sl.Len; i++ {
		if comparator(sl.Data[i], data) {
			return i, true
		}
	}
	return -1, false
}

func (sl *SeqList[T]) Clear(ifBack bool) {
	if ifBack {
		sl.Data = make([]T, 0, DefaultCap)
		sl.Len = 0
		sl.Cap = DefaultCap
		return
	}

	sl.Data = make([]T, 0, sl.Cap)
	sl.Len = 0
}
