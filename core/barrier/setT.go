package barrier

import (
	"sync"
)

type SetT[T Comparable] interface {
	Add(T) bool      //添加元素
	Remove(T) bool   //删除元素
	Contains(T) bool //判断是否包含元素

	Size() int   //返回集合大小
	Clear()      //清空集合
	Values() []T //返回集合中所有元素的切片
	ValueList() []T

	Union(other SetT[T]) SetT[T]      //返回另一个集合的并集
	Intersect(other SetT[T]) SetT[T]  //返回另一个集合的交集
	Difference(other SetT[T]) SetT[T] //返回另一个集合的差集
}

type mapSet[T Comparable] struct {
	data map[T]struct{} //只是用key 作为元素值

	values *Slice[T] //实时维护更新的value值

	lock sync.RWMutex //读写锁
}

func NewSetT[T Comparable]() SetT[T] {
	return &mapSet[T]{
		data:   make(map[T]struct{}),
		values: NewSlice(make([]T, 0)),
	}
}

func (this *mapSet[T]) Add(key T) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.data[key]; ok {
		return false //已经存在
	}

	//如果不存在直接追加
	this.data[key] = struct{}{}  //追加set
	s := this.values.Append(key) //追加值集合
	this.values.value = s
	return true
}

func (this *mapSet[T]) Remove(key T) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.data[key]; ok {
		delete(this.data, key)                               //删除结果集合
		this.values.IndexRemove(this.values.QueryIndex(key)) //删除值集合
		return true
	}

	return false
}

func (this *mapSet[T]) Contains(key T) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	_, ok := this.data[key]
	index := this.values.QueryIndex(key)

	return ok && index >= 0
}

func (this *mapSet[T]) Size() int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return len(this.data)
}

func (this *mapSet[T]) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data = make(map[T]struct{})
	this.values.Clear()
}

func (this *mapSet[T]) ValueList() []T {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.values.value
}

func (this *mapSet[T]) Values() []T {
	this.lock.RLock()
	defer this.lock.RUnlock()

	result := make([]T, 0)
	for key, _ := range this.data {
		result = append(result, key)
	}

	return result
}

func (this *mapSet[T]) Union(other SetT[T]) SetT[T] {
	this.lock.RLock()
	defer this.lock.RUnlock()

	union := NewSetT[T]()

	//将当前集合数据加入到临时集合中
	for k, _ := range this.data {
		union.Add(k)
	}

	//转换成切片，然后依次比较，如果相同则不会写入
	for _, v := range other.Values() {
		union.Add(v)
	}
	return union
}

// Intersect 求交集，当数据在两个集合内都存在的时候才会生成交集
func (this *mapSet[T]) Intersect(other SetT[T]) SetT[T] {

	intersect := NewSetT[T]()

	this.lock.RLock()
	defer this.lock.RUnlock()

	for k := range this.data {
		if other.Contains(k) {
			intersect.Add(k)
		}
	}
	return intersect
}

func (this *mapSet[T]) Difference(other SetT[T]) SetT[T] {
	this.lock.RLock()
	defer this.lock.RUnlock()

	diff := NewSetT[T]()
	for k := range this.data {
		if !other.Contains(k) {
			diff.Add(k)
		}
	}
	return diff
}
