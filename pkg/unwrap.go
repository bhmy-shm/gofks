package pkg

type InterfaceResult struct {
	err   error
	value Value
}

func Result(result interface{}, err error) *InterfaceResult {
	return &InterfaceResult{err: err, value: UnwrapValue(result)}
}

func (r *InterfaceResult) Unwraps() Value {
	if r.err != nil {
		panic(r.err.Error())
	}
	return r.value
}

func (r *InterfaceResult) Unwrap() interface{} {
	if r.err != nil {
		panic(r.err.Error())
	}
	return r.value.Interface()
}

func (r *InterfaceResult) UnwrapOr(data interface{}) interface{} {
	if r.err != nil {
		return data
	}
	return r.value.Interface()
}

func (r *InterfaceResult) UnwrapStr(format ...string) interface{} {
	if r.err != nil {
		return format
	}
	return r.value.Interface()
}

/***
 * @用于某一个事件没有拿到想要的结果，此时可以传入一个 f func() interface{} 来执行另一个事件
 * @属于一种双向策略选择，如果A没实现，旧用B实现。如：redis命中
**
*/
func (r *InterfaceResult) UnwrapFunc(f func() interface{}) interface{} {
	if r.err != nil {
		return f()
	}
	return r.value.Interface()
}
