package barrier

import (
	"sync"
)

type Set interface {
	Add(interface{}) bool      //添加元素
	Remove(interface{}) bool   //删除元素
	Contains(interface{}) bool //判断是否包含元素

	Size() int             //返回集合大小
	Clear()                //清空集合
	Values() []interface{} //返回集合中所有元素的切片
	KeysStr() []string

	Union(other Set) Set      //返回另一个集合的并集
	Intersect(other Set) Set  //返回另一个集合的交集
	Difference(other Set) Set //返回另一个集合的差集
}

type MapSet struct {
	data  map[interface{}]struct{} //只是用key 作为元素值
	mutex sync.RWMutex             //读写锁
}

func NewSet() Set {
	return &MapSet{
		data: make(map[interface{}]struct{}),
	}
}

func (this *MapSet) Add(key interface{}) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, ok := this.data[key]; ok {
		return false //已经存在
	}
	//如果不存在直接追加
	this.data[key] = struct{}{}
	return true
}

func (this *MapSet) Remove(key interface{}) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, ok := this.data[key]; ok {
		delete(this.data, key) //找到直接删除
		return true
	}

	return false
}
func (this *MapSet) Contains(key interface{}) bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	_, ok := this.data[key]
	return ok
}

func (this *MapSet) Size() int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	return len(this.data)
}

func (this *MapSet) Clear() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.data = make(map[interface{}]struct{})
}

func (this *MapSet) Values() []interface{} {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	slice := make([]interface{}, len(this.data))
	for k, _ := range this.data {
		slice = append(slice, k)
	}
	return slice
}

// KeysStr returns string keys in s.
func (this *MapSet) KeysStr() []string {
	var keys []string

	for key, _ := range this.data {
		if strKey, ok := key.(string); ok {
			keys = append(keys, strKey)
		}
	}

	return keys
}

func (this *MapSet) Union(other Set) Set {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	union := NewSet()

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
func (this *MapSet) Intersect(other Set) Set {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	intersect := NewSet()
	for k := range this.data {
		if other.Contains(k) {
			intersect.Add(k)
		}
	}
	return intersect
}

func (this *MapSet) Difference(other Set) Set {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	diff := NewSet()
	for k := range this.data {
		if !other.Contains(k) {
			diff.Add(k)
		}
	}
	return diff
}
