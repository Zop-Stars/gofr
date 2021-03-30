package logging

import (
	"encoding/json"
	"testing"

	"github.com/vikash/gofr/pkg/gofr/testutil"
)

const testLogStatement = "hello log!"

func TestLogger_Log(t *testing.T) {
	f := func() {
		logger := NewLogger(DEBUG)
		logger.Log(testLogStatement)
	}

	output := testutil.StdoutOutputForFunc(f)
	assertMessageInJSONLog(t, output, testLogStatement)
}

func TestLogger_Logf(t *testing.T) {
	f := func() {
		logger := NewLogger(DEBUG)
		logger.Logf("%s", testLogStatement)
	}

	output := testutil.StdoutOutputForFunc(f)

	assertMessageInJSONLog(t, output, testLogStatement)
}

func TestLogger_Error(t *testing.T) {
	f := func() {
		logger := NewLogger(DEBUG)
		logger.Error(testLogStatement)
	}

	output := testutil.StderrOutputForFunc(f)

	assertMessageInJSONLog(t, output, testLogStatement)
}

func TestLogger_Errorf(t *testing.T) {
	f := func() {
		logger := NewLogger(DEBUG)
		logger.Errorf("%s", testLogStatement)
	}

	output := testutil.StderrOutputForFunc(f)

	assertMessageInJSONLog(t, output, testLogStatement)
}

func assertMessageInJSONLog(t *testing.T, logLine, expectation string) {
	var l logEntry
	_ = json.Unmarshal([]byte(logLine), &l)

	if l.Message != expectation {
		t.Errorf("Log mismatch. Expected: %s Got: %s", expectation, l.Message)
	}
}
