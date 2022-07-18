package gofk

type IClass interface {
	Build(*Gofk)
	Name() string
}

type Bean interface {
	Name() string
}
