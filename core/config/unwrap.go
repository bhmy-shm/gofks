package pkg

import (
	"github.com/bhmy-shm/gofks/core/errorx"
)

type wrapAble interface {
	int | int8 | int16 | int32 | int64 |
		string | byte | float32 | float64 | bool | interface{}
}

type (
	WrapInter[T wrapAble] interface {
		Unwrap() T
		UnwrapOr(T) T
		UnwrapFunc(fn func())
	}

	wrap[T wrapAble] struct {
		data T
		err  error
	}
)

func newWrap[T wrapAble](data T, err error) WrapInter[T] {
	return &wrap[T]{
		data: data,
		err:  err,
	}
}

func (w *wrap[T]) Unwrap() T {
	if w.err != nil {
		errorx.Fatal(w.err)
	}
	return w.data
}

func (w *wrap[T]) UnwrapOr(data T) T {
	if w.err != nil {
		return data
	}
	return w.data
}

func (w *wrap[T]) UnwrapFunc(fn func()) {
	if w.err != nil {
		fn()
	}
}
