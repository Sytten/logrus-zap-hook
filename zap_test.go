package logrus_zap_hook

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const someLogMessage = "I am batman"
const someFieldKey = "Name"
const someFieldValue = "James Bond"
const someErrorMessage = "my martini is shaken"

func TestLogEntry(t *testing.T) {
	log, recordedLogs := newLogger()

	log.Info(someLogMessage)

	assert.Equal(t, recordedLogs.Len(), 1)
	assert.Equal(t, recordedLogs.All()[0].Message, someLogMessage)
	assert.Equal(t, recordedLogs.All()[0].Level, zapcore.InfoLevel)
}

func TestLogEntryWithField(t *testing.T) {
	log, recordedLogs := newLogger()

	log.WithField(someFieldKey, someFieldValue).Info(someLogMessage)

	assert.Equal(t, recordedLogs.Len(), 1)
	assert.Equal(t, recordedLogs.All()[0].Message, someLogMessage)
	assert.Equal(t, recordedLogs.All()[0].Level, zapcore.InfoLevel)
	assert.Equal(t, len(recordedLogs.All()[0].Context), 1)
	assert.Equal(t, recordedLogs.All()[0].Context[0].Key, someFieldKey)
	assert.Equal(t, recordedLogs.All()[0].Context[0].String, someFieldValue)
}

func TestLogEntryWithError(t *testing.T) {
	log, recordedLogs := newLogger()

	log.WithError(errors.New(someErrorMessage)).Info(someLogMessage)

	assert.Equal(t, recordedLogs.Len(), 1)
	assert.Equal(t, recordedLogs.All()[0].Message, someLogMessage)
	assert.Equal(t, recordedLogs.All()[0].Level, zapcore.InfoLevel)
	assert.Equal(t, len(recordedLogs.All()[0].Context), 1)
	assert.Equal(t, recordedLogs.All()[0].Context[0].Key, "error")
	assert.Equal(t, recordedLogs.All()[0].Context[0].Interface.(error).Error(), someErrorMessage)
}

func newLogger() (*logrus.Logger, *observer.ObservedLogs) {
	log := logrus.New()
	log.SetOutput(ioutil.Discard)

	hook, recordedLogs := newTestHook()
	log.Hooks.Add(hook)

	return log, recordedLogs
}

func newTestHook() (*ZapHook, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	hook, _ := NewZapHook(logger)

	return hook, recorded
}
