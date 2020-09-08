package log

import (
	"flag"
	"fmt"
	"io"
	"os"

	"k8s.io/klog/v2"
)

var defaultLogDir = "logs"

// InitFlags is for explicitly initializing the flags.
// Default to use logs as log dir.
func InitFlags() {
	klog.InitFlags(nil)

	if f := flag.Lookup("log_dir"); f == nil || f.Value.String() == "" {
		if err := flag.Set("log_dir", defaultLogDir); err != nil {
			panic(fmt.Errorf("log: failed to set log dir: %w", err))
		}
	}

	logDir := flag.Lookup("log_dir").Value.String()

	// Directory permissions should be 0750 or less
	// for security issue
	if err := os.MkdirAll(logDir, 0750); err != nil {
		panic(fmt.Errorf("log: failed to mkdir with %s: %w", logDir, err))
	}

	if err := flag.Set("logtostderr", "false"); err != nil {
		panic(fmt.Errorf("log: failed to set logtostderr: %w", err))
	}

	if err := flag.Set("alsologtostderr", "true"); err != nil {
		panic(fmt.Errorf("log: failed to set alsologtostderr: %w", err))
	}
}

// SetOutput sets the output destination for all severities
func SetOutput(w io.Writer) {
	klog.SetOutput(w)
}

// Flush flushes all pending log I/O.
func Flush() {
	klog.Flush()
}

// Warningln logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is always appended.
func Warningln(args ...interface{}) {
	klog.Warningln(args...)
}

// Warningf logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	klog.Warningf(format, args...)
}

// Errorln logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is always appended.
func Errorln(args ...interface{}) {
	klog.Errorln(args...)
}

// Errorf logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	klog.Errorf(format, args...)
}

// Infoln is equivalent to the global Infoln function, guarded by the value of v.
func Infoln(level int, args ...interface{}) {
	l := klog.Level(level)
	if klog.V(l).Enabled() {
		klog.V(l).Infoln(args...)
	}
}

// Infof is equivalent to the global Infof function, guarded by the value of v.
func Infof(level int, format string, args ...interface{}) {
	l := klog.Level(level)
	if klog.V(l).Enabled() {
		klog.V(l).Infof(format, args...)
	}
}
