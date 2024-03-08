package log_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Zulzi/jettison/errors"
	"github.com/Zulzi/jettison/j"
	"github.com/Zulzi/jettison/log"
)

type testLogger struct {
	logs []log.Entry
}

func (tl *testLogger) Log(_ context.Context, l log.Entry) string {
	tl.logs = append(tl.logs, l)
	return ""
}

func TestAddLoggers(t *testing.T) {
	tl := new(testLogger)
	log.SetLoggerForTesting(t, tl)

	log.Info(nil, "message", j.KV("some", "param"))
	log.Error(nil, errors.New("errMsg"))

	assert.Equal(t, "message,info,some,param,", toStr(tl.logs[0]))
	assert.Equal(t, "errMsg,error,", toStr(tl.logs[1]))
}

func toStr(l log.Entry) string {
	str := l.Message + ","
	str += string(l.Level) + ","
	if len(l.Parameters) == 0 {
		return str
	}
	for _, kv := range l.Parameters {
		str += kv.Key + "," + kv.Value + ","
	}
	return str
}
