package gofk

type IClass interface {
	Build(*Gofk)
	Injector() string
}

type Bean interface {
	Injector() string
}
