package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Generate_Run(t *testing.T) {
	out, err := runCobraCmd(GenerateCmd)

	assert.Nil(t, err)
	assert.Contains(t, out, "generate")
}
