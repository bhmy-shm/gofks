package mathx

import (
	"testing"
)

func TestCalcEntropy(t *testing.T) {
	const total = 1000
	const count = 100
	m := make(map[interface{}]int, total)
	for i := 0; i < total; i++ {
		m[i] = count
	}
	//assert.True(t, CalcEntropy(m) > .99)
	println(CalcEntropy(m))
}
