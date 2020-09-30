package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Dev_Escort_Init(t *testing.T) {
	e := getEscort()
	assert.Nil(t, e.init())
}

func getEscort() *escort {
	return &escort{
		root:  ".",
		path:  ".",
		delay: time.Second,
	}
}
