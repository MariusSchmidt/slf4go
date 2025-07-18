package slf4go_logrus_provider

import (
	"github.com/MariusSchmidt/slf4go/slf4go_api"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructuredLogging(t *testing.T) {
	testConfig := newTestingSetup().
		withComponent("test-service").
		withComponentLabel("appLabel").
		withStaticTags(map[string]interface{}{"env": "test"})

	levels := []struct {
		name     string
		logFn    func(string, ...interface{})
		logLevel logrus.Level
	}{
		{"error", testConfig.logger.Errorf, logrus.ErrorLevel},
		{"warn", testConfig.logger.Warnf, logrus.WarnLevel},
		{"warning", testConfig.logger.Warningf, logrus.WarnLevel},
		{"info", testConfig.logger.Infof, logrus.InfoLevel},
		{"debug", testConfig.logger.Debugf, logrus.DebugLevel},
		{"trace", testConfig.logger.Tracef, logrus.TraceLevel},
	}

	for _, level := range levels {
		t.Run(level.name, func(t *testing.T) {
			testConfig.hook.Reset()
			level.logFn("test message")
			assertLog(t, testConfig.hook).
				hasField("env", "test").
				hasField("appLabel", slf4go_api.AppComponent("test-service")).
				hasLevel(level.logLevel).
				hasMessage("test message")
		})
	}
}

type testingSetup struct {
	logger slf4go_api.Slf4GoLogger
	hook   *test.Hook
}

func newTestingSetup() *testingSetup {
	var logrusLogger *logrus.Logger
	var hook *test.Hook
	logrusLogger, hook = test.NewNullLogger()
	logrusLogger.SetLevel(logrus.TraceLevel)
	return &testingSetup{
		logger: New(logrusLogger),
		hook:   hook,
	}
}

func (setup *testingSetup) withComponent(component slf4go_api.AppComponent) *testingSetup {
	setup.logger = setup.logger.ForComponent(component)
	return setup
}

func (setup *testingSetup) withComponentLabel(componentLabel string) *testingSetup {
	setup.logger = setup.logger.WithAppComponentLabel(componentLabel)
	return setup
}

func (setup *testingSetup) withStaticTags(tags map[string]interface{}) *testingSetup {
	setup.logger = setup.logger.WithStaticTags(tags)
	return setup
}

type logAssertions struct {
	t    *testing.T
	hook *test.Hook
}

func assertLog(t *testing.T, hook *test.Hook) *logAssertions {
	return &logAssertions{t, hook}
}

func (a *logAssertions) hasLevel(level logrus.Level) *logAssertions {
	assert.Equal(a.t, level, a.hook.LastEntry().Level)
	return a
}

func (a *logAssertions) hasField(expectedKey string, expectedValue any) *logAssertions {
	entry := a.hook.LastEntry()
	assert.Equal(a.t, entry.Data[expectedKey], expectedValue)
	return a
}

func (a *logAssertions) hasMessage(msg string) *logAssertions {
	assert.Equal(a.t, msg, a.hook.LastEntry().Message)
	return a
}
