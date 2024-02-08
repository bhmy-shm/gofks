package barrier

import (
	"reflect"
	"testing"
)

// 测试 int 类型的 Slice
func TestIntSlice(t *testing.T) {
	s := NewSlice([]int{1, 2, 3})
	if !reflect.DeepEqual(s.value, []int{1, 2, 3}) {
		t.Errorf("Expected slice to be [1, 2, 3], got %v", s.value)
	}

	s.Append(4)
	if !reflect.DeepEqual(s.value, []int{1, 2, 3, 4}) {
		t.Errorf("Expected slice after append to be [1, 2, 3, 4], got %v", s.value)
	}

	index := s.QueryIndex(3)
	if index != 2 {
		t.Errorf("Expected index of 3 to be 2, got %d", index)
	}

	v := s.QueryValue(2)
	t.Log("Expected querying value out of range", v)

	s.IndexAdd(2, 10)
	if !reflect.DeepEqual(s.value, []int{1, 2, 10, 3, 4}) {
		t.Errorf("Expected slice after IndexAdd to be [1, 2, 10, 3, 4], got %v", s.value)
	}
}

// 测试 string 类型的 Slice
func TestStringSlice(t *testing.T) {
	s := NewSlice([]string{"a", "b", "c"})
	if !reflect.DeepEqual(s.value, []string{"a", "b", "c"}) {
		t.Errorf("Expected slice to be [\"a\", \"b\", \"c\"], got %v", s.value)
	}

	s.Append("d")
	if !reflect.DeepEqual(s.value, []string{"a", "b", "c", "d"}) {
		t.Errorf("Expected slice after append to be [\"a\", \"b\", \"c\", \"d\"], got %v", s.value)
	}

	index := s.QueryIndex("c")
	if index != 2 {
		t.Errorf("Expected index of \"c\" to be 2, got %d", index)
	}

	v := s.QueryValue(2)
	t.Log("Expected querying value out of range", v)

}

// 自定义的 Comparable 类型
type MyComparable struct {
	value int
}

// 使 MyComparable 类型满足 Comparable 约束
func (c MyComparable) Equal(other MyComparable) bool {
	return c.value == other.value
}

// 测试自定义 Comparable 类型的 Slice
func TestMyComparableSlice(t *testing.T) {
	//注释部分会报错
	//s := NewSlice([]MyComparable{{1}, {2}, {3}})
	//if !reflect.DeepEqual(s.value, []MyComparable{{1}, {2}, {3}}) {
	//	t.Errorf("Expected slice to be [{1}, {2}, {3}], got %v", s.value)
	//}
}
