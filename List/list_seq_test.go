package List

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// 测试 NewSeqList：正常容量、非法容量（panic）
func TestNewSeqList(t *testing.T) {
	// 场景1：合法容量（cap=5）
	t.Run("valid_capacity", func(t *testing.T) {
		sl := NewSeqList[int](5)
		assert.NotNil(t, sl)
		assert.Equal(t, 5, sl.Cap)
		assert.Equal(t, 0, sl.Len)
		assert.Len(t, sl.Data, 0)        // 切片长度为0
		assert.Equal(t, 5, cap(sl.Data)) // 切片容量为5
	})

	// 场景2：非法容量（cap<=0）→ 触发panic
	t.Run("invalid_capacity_panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewSeqList[string](0) // cap=0
		}, "cap=0 应触发 panic")

		assert.Panics(t, func() {
			NewSeqList[float64](-3) // cap=-3
		}, "cap=-3 应触发 panic")
	})
}

// 测试 Append：正常追加、扩容场景、多次追加
func TestAppendOfSeq(t *testing.T) {
	// 场景1：正常追加（不触发扩容）
	t.Run("append_without_expand", func(t *testing.T) {
		sl := NewSeqList[int](3)
		sl.Append(10)
		sl.Append(20)
		sl.Append(30)

		assert.Equal(t, 3, sl.Cap) // 未扩容，容量保持3
		assert.Equal(t, 3, sl.Len) // 长度为3
		assert.Equal(t, []int{10, 20, 30}, sl.Data)
	})

	// 场景2：追加触发扩容（二倍扩容）
	t.Run("append_with_expand", func(t *testing.T) {
		sl := NewSeqList[string](2) // 初始容量2
		sl.Append("a")
		sl.Append("b")
		assert.Equal(t, 2, sl.Cap) // 前2个元素不扩容

		// 第3个元素触发扩容（cap=2→4）
		sl.Append("c")
		assert.Equal(t, 4, sl.Cap) // 扩容后容量4
		assert.Equal(t, 3, sl.Len) // 长度3
		assert.Equal(t, []string{"a", "b", "c"}, sl.Data)
		assert.Equal(t, 4, cap(sl.Data)) // 切片容量4

		// 第4个元素不扩容，第5个元素再次扩容（cap=4→8）
		sl.Append("d")
		sl.Append("e")
		assert.Equal(t, 8, sl.Cap) // 再次扩容后容量8
		assert.Equal(t, 5, sl.Len)
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, sl.Data)
	})

	// 场景3：自定义类型追加
	t.Run("append_custom_type", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}
		sl := NewSeqList[User](2)
		sl.Append(User{1, "张三"})
		sl.Append(User{2, "李四"})

		assert.Equal(t, 2, sl.Len)
		assert.Equal(t, User{1, "张三"}, sl.Data[0])
		assert.Equal(t, User{2, "李四"}, sl.Data[1])
	})
}

// 测试 Remove：正常删除、边界索引、越界删除、删除后清空冗余
func TestRemove(t *testing.T) {
	// 初始化测试数据：[10,20,30,40]，cap=4，len=4
	prepareList := func() *SeqList[int] {
		sl := NewSeqList[int](4)
		sl.Append(10)
		sl.Append(20)
		sl.Append(30)
		sl.Append(40)
		return sl
	}

	// 场景1：正常删除中间元素（索引1）
	t.Run("remove_middle_element", func(t *testing.T) {
		sl := prepareList()
		data, ok := sl.Remove(1)

		assert.True(t, ok)
		assert.Equal(t, 20, data)                            // 删除的元素是20
		assert.Equal(t, 3, sl.Len)                           // 长度变为3
		assert.Equal(t, 4, sl.Cap)                           // 容量不变
		assert.Equal(t, []int{10, 30, 40, 0}, sl.Data)       // 最后一个元素清空为0（int零值）
		assert.Equal(t, []int{10, 30, 40}, sl.Data[:sl.Len]) // 有效元素正确
	})

	// 场景2：删除首尾边界元素
	t.Run("remove_boundary_element", func(t *testing.T) {
		// 删除首元素（索引0）
		sl1 := prepareList()
		data1, ok1 := sl1.Remove(0)
		assert.True(t, ok1)
		assert.Equal(t, 10, data1)
		assert.Equal(t, []int{20, 30, 40, 0}, sl1.Data)

		// 删除尾元素（索引3→删除后len=3）
		sl2 := prepareList()
		data2, ok2 := sl2.Remove(3)
		assert.True(t, ok2)
		assert.Equal(t, 40, data2)
		assert.Equal(t, []int{10, 20, 30, 0}, sl2.Data)
	})

	// 场景3：越界删除（索引<0或>=len）
	t.Run("remove_out_of_range", func(t *testing.T) {
		sl := prepareList()

		// 索引=-1（负数）
		data1, ok1 := sl.Remove(-1)
		assert.False(t, ok1)
		assert.Equal(t, 0, data1)  // 返回int零值
		assert.Equal(t, 4, sl.Len) // 长度不变

		// 索引=4（>=len=4）
		data2, ok2 := sl.Remove(4)
		assert.False(t, ok2)
		assert.Equal(t, 0, data2)
		assert.Equal(t, 4, sl.Len) // 长度不变
	})

	// 场景4：删除引用类型（验证清空冗余，避免内存泄漏）
	t.Run("remove_reference_type", func(t *testing.T) {
		type RefType struct {
			Val string
		}
		sl := NewSeqList[*RefType](2)
		sl.Append(&RefType{"a"})
		sl.Append(&RefType{"b"})

		// 删除索引0
		data, ok := sl.Remove(0)
		assert.True(t, ok)
		assert.Equal(t, "a", data.Val)
		assert.Nil(t, sl.Data[1]) // 冗余位置清空为nil（引用类型零值）
		assert.Equal(t, 1, sl.Len)
	})
}

// 测试 Get：正常获取、边界索引、越界获取
func TestGet(t *testing.T) {
	// 初始化测试数据：[5,15,25,35]
	sl := NewSeqList[int](4)
	sl.Append(5)
	sl.Append(15)
	sl.Append(25)
	sl.Append(35)

	// 场景1：正常获取中间元素
	t.Run("get_valid_middle", func(t *testing.T) {
		data, ok := sl.Get(1)
		assert.True(t, ok)
		assert.Equal(t, 15, data)
	})

	// 场景2：获取边界元素（首、尾）
	t.Run("get_valid_boundary", func(t *testing.T) {
		// 获取首元素（索引0）
		data0, ok0 := sl.Get(0)
		assert.True(t, ok0)
		assert.Equal(t, 5, data0)

		// 获取尾元素（索引3）
		data3, ok3 := sl.Get(3)
		assert.True(t, ok3)
		assert.Equal(t, 35, data3)
	})

	// 场景3：越界获取
	t.Run("get_out_of_range", func(t *testing.T) {
		// 索引=-2
		data1, ok1 := sl.Get(-2)
		assert.False(t, ok1)
		assert.Equal(t, 0, data1)

		// 索引=4（>=len=4）
		data2, ok2 := sl.Get(4)
		assert.False(t, ok2)
		assert.Equal(t, 0, data2)
	})
}

// 测试 Locate：找到元素、未找到元素、自定义比较器、重复元素（返回第一个）
func TestLocate(t *testing.T) {
	// 场景1：基础类型（int），找到第一个匹配元素
	t.Run("locate_int_found", func(t *testing.T) {
		sl := NewSeqList[int](5)
		sl.Append(10)
		sl.Append(20)
		sl.Append(30)
		sl.Append(20) // 重复元素

		// 自定义比较器：相等判断
		comparator := func(a, b int) bool {
			return a == b
		}

		// 查找20，返回第一个匹配的索引1
		idx, ok := sl.Locate(20, comparator)
		assert.True(t, ok)
		assert.Equal(t, 1, idx)
	})

	// 场景2：未找到匹配元素
	t.Run("locate_not_found", func(t *testing.T) {
		sl := NewSeqList[string](3)
		sl.Append("apple")
		sl.Append("banana")

		comparator := func(a, b string) bool {
			return a == b
		}

		idx, ok := sl.Locate("orange", comparator)
		assert.False(t, ok)
		assert.Equal(t, -1, idx)
	})

	// 场景3：自定义类型（User），按字段匹配
	t.Run("locate_custom_type", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}
		sl := NewSeqList[User](3)
		sl.Append(User{1, "张三"})
		sl.Append(User{2, "李四"})
		sl.Append(User{3, "王五"})

		// 比较器：按ID匹配
		idComparator := func(a, b User) bool {
			return a.ID == b.ID
		}
		idx1, ok1 := sl.Locate(User{ID: 2}, idComparator)
		assert.True(t, ok1)
		assert.Equal(t, 1, idx1)

		// 比较器：按Name模糊匹配（包含"王"）
		nameComparator := func(a, b User) bool {
			return strings.HasPrefix(a.Name, "王")
		}
		idx2, ok2 := sl.Locate(User{}, nameComparator) // 目标参数仅用于触发比较器
		assert.True(t, ok2)
		assert.Equal(t, 2, idx2)
	})

	// 场景4：传入nil比较器（虽然代码没校验，但测试覆盖异常场景）
	t.Run("locate_nil_comparator", func(t *testing.T) {
		sl := NewSeqList[int](2)
		sl.Append(10)

		// 传入nil比较器→触发panic（代码未处理nil，需注意）
		assert.Panics(t, func() {
			sl.Locate(10, nil)
		}, "nil比较器应触发panic")
	})
}

// 测试 Clear：保留原容量、重置为默认容量
func TestClear(t *testing.T) {
	// 初始化测试数据：cap=8，len=3，data=[1,2,3]
	prepareExpandedList := func() *SeqList[int] {
		sl := NewSeqList[int](2)
		sl.Append(1)
		sl.Append(2)
		sl.Append(3) // 触发扩容→cap=4，再追加2个元素→cap=8
		sl.Append(4)
		sl.Append(5)
		return sl
	}

	// 场景1：ifBack=false → 保留原容量，仅清空元素
	t.Run("clear_keep_capacity", func(t *testing.T) {
		sl := prepareExpandedList()
		oldCap := sl.Cap // 原容量8

		sl.Clear(false)
		assert.Equal(t, 0, sl.Len)       // 长度重置为0
		assert.Equal(t, oldCap, sl.Cap)  // 容量保持8
		assert.Len(t, sl.Data, 0)        // 切片长度0
		assert.Equal(t, 8, cap(sl.Data)) // 切片容量8
	})

	// 场景2：ifBack=true → 重置为默认容量（DefaultCap=10）
	t.Run("clear_reset_to_default_cap", func(t *testing.T) {
		sl := prepareExpandedList()

		sl.Clear(true)
		assert.Equal(t, 0, sl.Len)                // 长度重置为0
		assert.Equal(t, DefaultCap, sl.Cap)       // 容量重置为10
		assert.Len(t, sl.Data, 0)                 // 切片长度0
		assert.Equal(t, DefaultCap, cap(sl.Data)) // 切片容量10
	})

	// 场景3：清空后追加元素（验证容量有效性）
	t.Run("append_after_clear", func(t *testing.T) {
		sl := prepareExpandedList()
		sl.Clear(false) // 保留容量8

		sl.Append(100)
		assert.Equal(t, 1, sl.Len)
		assert.Equal(t, 8, sl.Cap)
		assert.Equal(t, []int{100}, sl.Data)
	})
}

// 测试综合场景：多方法联动（Append→Remove→Get→Clear→Append）
func TestIntegration(t *testing.T) {
	sl := NewSeqList[string](3)

	// 1. 追加3个元素
	sl.Append("a")
	sl.Append("b")
	sl.Append("c")
	assert.Equal(t, 3, sl.Len)

	// 2. 删除索引1（元素"b"）
	_, ok := sl.Remove(1)
	assert.True(t, ok)
	assert.Equal(t, 2, sl.Len)

	// 3. 获取索引1（元素"c"）
	data, ok := sl.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "c", data)

	// 4. 定位元素"a"
	comparator := func(a, b string) bool {
		return a == b
	}
	idx, ok := sl.Locate("a", comparator)
	assert.True(t, ok)
	assert.Equal(t, 0, idx)

	// 5. 清空（保留容量）
	sl.Clear(false)
	assert.Equal(t, 0, sl.Len)
	assert.Equal(t, 3, sl.Cap)

	// 6. 清空后追加新元素
	sl.Append("x")
	sl.Append("y")
	assert.Equal(t, 2, sl.Len)
	assert.Equal(t, []string{"x", "y"}, sl.Data)
}
