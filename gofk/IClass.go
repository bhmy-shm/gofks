package gofk

type IClass interface {
	Build(*Gofk)
	Name() string
}
