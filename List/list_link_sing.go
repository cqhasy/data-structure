package List

import "errors"

type SinNode[T any] struct {
	Val  T
	Next *SinNode[T]
}

type SinHead[T any] struct {
	len  int // 存储链表的长度
	Next *SinNode[T]
}

func NewSinList[T any]() *SinHead[T] {
	return &SinHead[T]{
		len:  0,
		Next: nil,
	}
}

func (head *SinHead[T]) Len() int {
	return head.len
}

func (n *SinNode[T]) Value() T {
	return n.Val
}

func NewSinNode[T any](val T) *SinNode[T] {
	return &SinNode[T]{
		Val:  val,
		Next: nil,
	}
}

func (head *SinHead[T]) Append(node *SinNode[T]) {
	if head.len == 0 {
		head.Next = node
		head.len++
		return
	}

	tar := head.Next
	for i := 1; i < head.len; i++ {
		tar = tar.Next
	}
	tar.Next = node
	head.len++
	return
}

// Insert 插入数据到第 i 个节点后
func (head *SinHead[T]) Insert(node *SinNode[T], i int) error {
	if i > head.len {
		return errors.New("index out of range")
	}

	if i == 0 {
		node.Next = head.Next
		head.Next = node
		head.len++
		return nil
	}
	tar := head.Next
	for j := 1; j < i; j++ {
		tar = tar.Next
	}
	node.Next = tar.Next
	tar.Next = node
	head.len++
	return nil
}

// Delete 删除第i个节点
func (head *SinHead[T]) Delete(i int) error {
	if i <= 0 || i > head.len {
		return errors.New("invalid index")
	}
	if head.len == 0 {
		return errors.New("no data to delete")
	}

	if i == 1 {
		head.Next = head.Next.Next
		head.len--
		return nil
	}
	var prev *SinNode[T] = nil
	tar := head.Next
	for j := 1; j < i; j++ {
		prev = tar
		tar = tar.Next
	}
	prev.Next = tar.Next
	head.len--
	return nil
}

// GetElem 获取第i个节点
func (head *SinHead[T]) GetElem(i int) (*SinNode[T], error) {
	if i <= 0 || i > head.len {
		return nil, errors.New("invalid index")
	}
	tar := head.Next
	for j := 1; j < i; j++ {
		tar = tar.Next
	}
	return tar, nil
}
