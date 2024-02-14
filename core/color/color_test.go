package color

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithColor(t *testing.T) {
	output := WithColor("Hello", BgRed)
	log.Println(output)
	assert.Equal(t, "Hello", output)
}

func TestWithColorPadding(t *testing.T) {
	output := WithColorPadding("Hello", BgRed)
	assert.Equal(t, " Hello ", output)
}
