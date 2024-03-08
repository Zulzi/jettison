package errors

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-stack/stack"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

	"github.com/Zulzi/jettison/internal"
	"github.com/Zulzi/jettison/trace"
)

//go:generate go test -update

func TestSetTraceConfig(t *testing.T) {
	type packageType struct{}
	cfg := trace.StackConfig{
		RemoveLambdas: true,
		PackagesShown: []string{trace.PackagePath(packageType{})},
		TrimRuntime:   true,
		FormatStack: func(call stack.Call) string {
			return fmt.Sprintf("%+k:%n", call, call)
		},
	}
	SetTraceConfig(cfg)
	_, st := getTrace(0)
	assert.Equal(t, []string{"github.com/Zulzi/jettison/errors:TestSetTraceConfig"}, st)

	assert.Panics(t, func() {
		SetTraceConfig(trace.StackConfig{})
	})
}

// TestStack tests the stack trace including line numbers.
// Adding anything to this file might break the test.
func TestStack(t *testing.T) {
	SetTraceConfigTesting(t, TestingConfig)
	err := stackCalls(5)
	tr := []byte(strings.Join(err.StackTrace, "\n") + "\n")
	goldie.New(t).Assert(t, t.Name(), tr)
}

func stackCalls(i int) *internal.Error {
	if i == 0 {
		return New("stack").(*internal.Error)
	}
	return stackCalls(i - 1)
}

func TestGetSourceCode(t *testing.T) {
	SetTraceConfigTesting(t, TestingConfig)
	assert.Equal(t, "trace_test.go TestGetSourceCode", getSourceCode(0))
}
