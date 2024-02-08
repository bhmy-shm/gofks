package barrier

import (
	"sync"
)

type (
	Comparable interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | uint8 | uint32 | uint64 | ~string | ~float32 | ~float64 | struct{}
	}

	Slice[T Comparable] struct {
		lock  sync.RWMutex
		value []T
	}
)

func NewSlice[T Comparable](slice []T) *Slice[T] {
	return &Slice[T]{
		value: slice,
	}
}

// Append 渐进式扩容
func Append[T Comparable](slice []T, values ...T) []T {
	if slice == nil || len(slice) == 0 {
		slice = make([]T, 0)
		slice = append(slice, values...)
		return slice
	}

	var (
		length, caps = len(slice), cap(slice)
		newCap       int
	)

	//如果旧切片得长度已经大于切片容量，则扩容需要进行扩容
	if length > caps {
		//如果大于1024则翻倍
		if caps > 1024 {
			newCap = caps * 2
		} else {
			newCap = caps * 5 / 4
		}
	}

	//如果当前长度 + 添加内容得长度也大于cap容量，则也需要扩容
	l := length + len(values)
	if l > caps {
		if caps > 1024 {
			newCap = caps * 2
		} else {
			//需要用新的长度 / 容量，获取是否需要翻倍
			n := l / caps
			if n > 0 {
				//如果大于0，则表示需要翻几倍，否则1.25
				newCap = caps * n * 5 / 4
			} else {
				newCap = caps * 5 / 4
			}
		}
	} else {
		//以上两个条件都不满足无需扩容直接返回
		slice = append(slice, values...)
		return slice
	}

	//渐进式扩容
	var newSlice = make([]T, length, newCap)
	copy(newSlice, slice)
	newSlice = append(newSlice, values...)
	return newSlice
}

func (s *Slice[T]) checkIndexLength(index int) bool {
	return index > -1 && index <= len(s.value)
}

// QueryValue 获取指定索引的值
func (s *Slice[T]) QueryValue(index int) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.value[index]
}

// QueryIndex 获取指定值的索引
func (s *Slice[T]) QueryIndex(value T) int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if len(s.value) == 0 {
		return -1
	}

	for index, v := range s.value {
		//Notice：这里使用泛型比较需要注意，Comparable 约束，不能是 []byte、interface、map 等无法比较的类型
		if v == value {
			return index
		}
	}
	return -1
}

// Prepend 头插
func (s *Slice[T]) Prepend(value T) []T {
	s.lock.Lock()
	defer s.lock.Unlock()

	newSlice := make([]T, len(s.value)+1)
	newSlice[0] = value
	copy(newSlice[1:], s.value)
	return newSlice
}

// Append 尾插
func (s *Slice[T]) Append(value T) []T {
	s.lock.Lock()
	defer s.lock.Unlock()

	return Append(s.value, value)
}

// IndexAdd 基于索引添加
func (s *Slice[T]) IndexAdd(index int, value T) []T {

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.checkIndexLength(index) {
		panic("Index Out of Range")
	}

	//如果是尾部加入
	if len(s.value) == index {
		s.value = Append(s.value, value)
	}

	//创建两个全新得底层数组空间，函数内栈空间生成，结束函数自动回收
	var (
		temp, last []T
	)
	temp = Append(temp, s.value[:index+1]...) //从索引位置取取开头
	last = Append(last, s.value[index+1:]...) //从索引位置取取末尾

	temp = Append(temp, value)
	temp = Append(temp, last...)
	s.value = temp

	return s.value
}

// IndexEdit 基于索引编辑
func (s *Slice[T]) IndexEdit(index int, value T) []T {

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.checkIndexLength(index) {
		panic("Index Out of Range")
	}
	s.value[index] = value
	return s.value
}

// IndexRemove 基于索引删除
func (s *Slice[T]) IndexRemove(index int) []T {

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.checkIndexLength(index) {
		panic("Index Out of Range")
	}
	var (
		temp []T
	)
	temp = s.value[:index]
	temp = Append(temp, s.value[index+1:]...)
	s.value = temp

	return s.value
}

func (s *Slice[T]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.value = nil
}
