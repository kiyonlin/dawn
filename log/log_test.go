package log

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitFlags(t *testing.T) {
	defaultLogFile = "test_log_file"
	defer os.RemoveAll(defaultLogFile)

	InitFlags()

	assert.Equal(t, defaultLogFile, flag.Lookup("log_file").Value.String())
	assert.Equal(t, "false", flag.Lookup("logtostderr").Value.String())
	assert.Equal(t, "true", flag.Lookup("alsologtostderr").Value.String())
	assert.Equal(t, "0", flag.Lookup("v").Value.String())
	assert.Equal(t, "2", flag.Lookup("stderrthreshold").Value.String())
}

func Test_All(t *testing.T) {
	SetOutput(new(bytes.Buffer))
	Errorln("errorln")
	Errorf("%s", "errorf")
	Warningln("warningln")
	Warningf("%s", "warningf")
	Infoln(0, "infoln level 0")
	Infof(0, "%s", "infof level 0")
	Infof(1, "%s", "infof level 1")
	Flush()
}
