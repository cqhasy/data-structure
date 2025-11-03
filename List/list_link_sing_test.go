package List

import "testing"

func TestNewSinList(t *testing.T) {
	list := NewSinList[int]()
	if list.Len() != 0 {
		t.Errorf("expected length 0, got %d", list.Len())
	}
	if list.Next != nil {
		t.Errorf("expected Next nil, got %+v", list.Next)
	}
}

// 测试 Append 操作
func TestAppendOfSing(t *testing.T) {
	list := NewSinList[int]()
	list.Append(NewSinNode(10))
	list.Append(NewSinNode(20))
	list.Append(NewSinNode(30))

	if list.Len() != 3 {
		t.Errorf("expected length 3, got %d", list.Len())
	}

	node, _ := list.GetElem(3)
	if node.Val != 30 {
		t.Errorf("expected last value 30, got %d", node.Val)
	}
}

// 测试 Insert 操作
func TestInsert(t *testing.T) {
	list := NewSinList[int]()
	list.Append(NewSinNode(1))
	list.Append(NewSinNode(2))
	list.Append(NewSinNode(3))

	err := list.Insert(NewSinNode(99), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	node, _ := list.GetElem(3)
	if node.Val != 99 {
		t.Errorf("expected value 99 at position 3, got %d", node.Val)
	}

	if list.Len() != 4 {
		t.Errorf("expected length 4, got %d", list.Len())
	}
}

// 测试 Delete 操作
func TestDelete(t *testing.T) {
	list := NewSinList[int]()
	list.Append(NewSinNode(1))
	list.Append(NewSinNode(2))
	list.Append(NewSinNode(3))

	err := list.Delete(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	node, _ := list.GetElem(2)
	if node.Val != 3 {
		t.Errorf("expected 3 at position 2 after delete, got %d", node.Val)
	}

	if list.Len() != 2 {
		t.Errorf("expected length 2, got %d", list.Len())
	}
}

// 测试 GetElem 越界
func TestGetElemOutOfRange(t *testing.T) {
	list := NewSinList[int]()
	list.Append(NewSinNode(1))

	_, err := list.GetElem(5)
	if err == nil {
		t.Errorf("expected error for invalid index, got nil")
	}
}
