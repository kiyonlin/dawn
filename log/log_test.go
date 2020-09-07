package log

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitFlags(t *testing.T) {
	defaultLogDir = "test_dir"
	defer os.RemoveAll(defaultLogDir)

	InitFlags()

	assert.Equal(t, defaultLogDir, flag.Lookup("log_dir").Value.String())
	assert.Equal(t, "false", flag.Lookup("logtostderr").Value.String())
	assert.Equal(t, "true", flag.Lookup("alsologtostderr").Value.String())
	assert.Equal(t, "0", flag.Lookup("v").Value.String())
	assert.Equal(t, "2", flag.Lookup("stderrthreshold").Value.String())
}

func Test_All(t *testing.T) {
	SetOutput(new(bytes.Buffer))
	Errorln("errorln")
	Errorf("%s", "errorf")
	Infoln(0, "infoln level 0")
	Infof(0, "%s", "infof level 0")
	Infof(1, "%s", "infof level 1")
	Flush()
}
