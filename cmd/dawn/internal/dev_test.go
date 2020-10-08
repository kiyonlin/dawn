package internal

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/stretchr/testify/assert"
)

func Test_Dev_Escort_New(t *testing.T) {
	ctx := cli.NewContext(
		cli.NewApp(),
		flag.NewFlagSet("test", 1),
		nil)
	assert.NotNil(t, newEscort(ctx))
}

func Test_Dev_Escort_Init(t *testing.T) {
	at := assert.New(t)

	e := getEscort()
	at.Nil(e.init())

	at.Contains(e.root, "internal")
	at.NotEmpty(e.binPath)
	at.Nil(os.Remove(e.binPath))
}

func getEscort() *escort {
	return &escort{
		root:   ".",
		target: ".",
		delay:  time.Second,
	}
}
