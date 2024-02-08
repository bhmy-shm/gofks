package barrier

import (
	"log"
	"testing"
)

func TestSetOperations(t *testing.T) {
	s := NewSetT[int]()

	// 测试添加
	added := s.Add(1)
	if !added || s.Size() != 1 {
		t.Error("Failed to add item to the set")
	}

	// 测试重复添加
	added = s.Add(12)
	if added {
		t.Error("Should not be able to add the same item twice")
	}

	log.Println("value1", s.Values())
	log.Println("value2", s.ValueList())

	// 测试是否包含
	if !s.Contains(13) {
		t.Error("Set should contain the item after being added")
	}

	// 测试删除
	removed := s.Remove(12)
	if !removed || s.Size() != 0 {
		t.Error("Failed to remove item from the set")
	}

	// 测试清空
	s.Add(13)
	s.Add(24)

	if s.Size() != 0 {
		t.Error("Set should be empty after clear")
	}

	log.Println("value3", s.Values())
	log.Println("value4", s.ValueList())

	s.Clear()

	log.Println("value clear", s.Values())

	// 测试并集
	s1 := NewSet[int]()
	s2 := NewSet[int]()
	s1.Add(1)
	s2.Add(2)
	union := s1.Union(s2)
	if union.Size() != 2 {
		t.Error("Union set should contain two items")
	}

	// 测试交集
	s1.Add(2) // 现在 s1 包含 {1, 2}
	intersect := s1.Intersect(s2)
	if intersect.Size() != 1 || !intersect.Contains(2) {
		t.Error("Intersect set should contain only one item and be {2}")
	}

	// 测试差集
	diff := s1.Difference(s2)
	if diff.Size() != 1 || !diff.Contains(1) {
		t.Error("Difference set should contain only one item and be {1}")
	}
}
