package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

/*
 * Provides a mock clock to ease testing of timestamps
 */
type constantClock time.Time

const timestamp string = "2077-01-23T10:15:13.000Z"

func (c constantClock) Now() time.Time { return time.Time(c) }
func (c constantClock) NewTicker(d time.Duration) *time.Ticker {
	return &time.Ticker{}
}

var clock constantClock

// Set up our mock clock
func init() {
	date := time.Date(2077, 1, 23, 10, 15, 13, 441, time.UTC)
	clock = constantClock(date)
}

var (
	rescueStdout, rescueStderr, rout, wout, rerr, werr *os.File
)

func redirectStdout() {
	rescueStdout = os.Stdout
	rout, wout, _ = os.Pipe()
	os.Stdout = wout
}

func redirectStderr() {
	rescueStderr = os.Stderr
	rerr, werr, _ = os.Pipe()
	os.Stderr = werr
}

func captureStdout() []byte {
	wout.Close()
	out, _ := ioutil.ReadAll(rout)
	os.Stdout = rescueStdout
	return out
}

func captureStderr() []byte {
	werr.Close()
	out, _ := ioutil.ReadAll(rerr)
	os.Stderr = rescueStderr
	return out
}

func TestServerLogger(t *testing.T) {
	os.Setenv("PROFILE", "server")
	redirectStderr()
	e := fmt.Sprintf(`{"level":"DEBUG","@timestamp":"%s","caller":"common-logger/logger_test.go:\d+","message":"TestServerLogger"}`, timestamp)
	l := New(zap.WithClock(clock))
	l.Debug("TestServerLogger")
	out := strings.Trim(string(captureStderr()), "\n")
	assert.Regexp(t, e, out)
}

func TestServerLoggerNotStdout(t *testing.T) {
	os.Setenv("PROFILE", "server")
	redirectStdout()
	l := New(zap.WithClock(clock))
	l.Debug("TestServerLogger")
	out := strings.Trim(string(captureStdout()), "\n")
	assert.Empty(t, out)
}

func TestConsoleLogger(t *testing.T) {
	os.Setenv("PROFILE", "console")
	redirectStderr()
	e := fmt.Sprintf("%s\tDEBUG\tcommon-logger/logger_test.go:\\d+\tTestConsoleLogger", timestamp)
	l := New(zap.WithClock(clock))
	l.Debugf("TestConsoleLogger")
	out := strings.Trim(string(captureStderr()), "\n")
	assert.Regexp(t, e, out)
}

func TestConsoleLoggerNotStdout(t *testing.T) {
	os.Setenv("PROFILE", "console")
	redirectStdout()
	l := New(zap.WithClock(clock))
	l.Debug("TestConsoleLogger")
	out := strings.Trim(string(captureStdout()), "\n")
	assert.Empty(t, out)
}

func TestDefaultLogger(t *testing.T) {
	os.Setenv("PROFILE", "")
	redirectStderr()
	e := fmt.Sprintf(`{"level":"DEBUG","@timestamp":"%s","caller":"common-logger/logger_test.go:\d+","message":"TestDefaultLogger"}`, timestamp)
	l := New(zap.WithClock(clock))
	l.Debugf("TestDefaultLogger")
	out := strings.Trim(string(captureStderr()), "\n")
	assert.Regexp(t, e, out)
}
