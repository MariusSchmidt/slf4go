package slf4go_logrus_provider

import (
	"github.com/MariusSchmidt/slf4go/slf4go_api"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogging(t *testing.T) {
	testConfig := newTestingSetup()

	scenarios := []struct {
		name          string
		logFn         func(string, ...interface{})
		logFnWithTags func(slf4go_api.LogTags, string, ...interface{})
		logLevel      logrus.Level
	}{
		{"fatal", testConfig.slf4GoLogrusLogger.Fatalf, testConfig.slf4GoLogrusLogger.FatalWithTagsf, logrus.FatalLevel},
		{"panic", testConfig.slf4GoLogrusLogger.Panicf, testConfig.slf4GoLogrusLogger.PanicWithTagsf, logrus.PanicLevel},
		{"error", testConfig.slf4GoLogrusLogger.Errorf, testConfig.slf4GoLogrusLogger.ErrorWithTagsf, logrus.ErrorLevel},
		{"warn", testConfig.slf4GoLogrusLogger.Warnf, testConfig.slf4GoLogrusLogger.WarnWithTagsf, logrus.WarnLevel},
		{"warning", testConfig.slf4GoLogrusLogger.Warningf, testConfig.slf4GoLogrusLogger.WarningWithTagsf, logrus.WarnLevel},
		{"info", testConfig.slf4GoLogrusLogger.Infof, testConfig.slf4GoLogrusLogger.InfoWithTagsf, logrus.InfoLevel},
		{"debug", testConfig.slf4GoLogrusLogger.Debugf, testConfig.slf4GoLogrusLogger.DebugWithTagsf, logrus.DebugLevel},
		{"trace", testConfig.slf4GoLogrusLogger.Tracef, testConfig.slf4GoLogrusLogger.TraceWithTagsf, logrus.TraceLevel},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-base", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}, "test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message with name=%s and value=%d", "beeblebrox", 42)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(
					map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"},
					"test message with name=%s and value=%d",
					"beeblebrox",
					42,
				)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	t.Run("unknown", func(t *testing.T) {
		testConfig.hook.Reset()
		var unknownLevel slf4go_api.LogLevel = 666
		testConfig.slf4GoLogrusLogger.Logf(unknownLevel, "Some message not displayed")
		assertLog(t, testConfig.hook).
			hasLevel(logrus.ErrorLevel).
			hasTags(map[string]interface{}{}).
			hasMessage("Mapping error level 'unknown' onto Logrus error level failed. Not logging event")
	})

}

func TestLogging_ForComponent(t *testing.T) {
	testConfig := newTestingSetup().
		forComponent("test-service")

	scenarios := []struct {
		name          string
		logFn         func(string, ...interface{})
		logFnWithTags func(slf4go_api.LogTags, string, ...interface{})
		logLevel      logrus.Level
	}{
		{"panic", testConfig.slf4GoLogrusLogger.Panicf, testConfig.slf4GoLogrusLogger.PanicWithTagsf, logrus.PanicLevel},
		{"error", testConfig.slf4GoLogrusLogger.Errorf, testConfig.slf4GoLogrusLogger.ErrorWithTagsf, logrus.ErrorLevel},
		{"warn", testConfig.slf4GoLogrusLogger.Warnf, testConfig.slf4GoLogrusLogger.WarnWithTagsf, logrus.WarnLevel},
		{"warning", testConfig.slf4GoLogrusLogger.Warningf, testConfig.slf4GoLogrusLogger.WarningWithTagsf, logrus.WarnLevel},
		{"info", testConfig.slf4GoLogrusLogger.Infof, testConfig.slf4GoLogrusLogger.InfoWithTagsf, logrus.InfoLevel},
		{"debug", testConfig.slf4GoLogrusLogger.Debugf, testConfig.slf4GoLogrusLogger.DebugWithTagsf, logrus.DebugLevel},
		{"trace", testConfig.slf4GoLogrusLogger.Tracef, testConfig.slf4GoLogrusLogger.TraceWithTagsf, logrus.TraceLevel},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-base", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					slf4go_api.DefaultAppComponentTag: slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}, "test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"dyn_key1":                        "dyn_val1",
					"dyn_key2":                        "dyn_val2",
					slf4go_api.DefaultAppComponentTag: slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message with name=%s and value=%d", "beeblebrox", 42)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					slf4go_api.DefaultAppComponentTag: slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(
					map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"},
					"test message with name=%s and value=%d",
					"beeblebrox",
					42,
				)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"dyn_key1":                        "dyn_val1",
					"dyn_key2":                        "dyn_val2",
					slf4go_api.DefaultAppComponentTag: slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	t.Run("unknown", func(t *testing.T) {
		testConfig.hook.Reset()
		var unknownLevel slf4go_api.LogLevel = 666
		testConfig.slf4GoLogrusLogger.Logf(unknownLevel, "Some message not displayed")
		assertLog(t, testConfig.hook).
			hasLevel(logrus.ErrorLevel).
			hasTags(map[string]interface{}{}).
			hasMessage("Mapping error level 'unknown' onto Logrus error level failed. Not logging event")
	})
}

func TestLogging_WithComponentLabel(t *testing.T) {
	testConfig := newTestingSetup().
		withComponentLabel("appLabel")

	scenarios := []struct {
		name          string
		logFn         func(string, ...interface{})
		logFnWithTags func(slf4go_api.LogTags, string, ...interface{})
		logLevel      logrus.Level
	}{
		{"panic", testConfig.slf4GoLogrusLogger.Panicf, testConfig.slf4GoLogrusLogger.PanicWithTagsf, logrus.PanicLevel},
		{"error", testConfig.slf4GoLogrusLogger.Errorf, testConfig.slf4GoLogrusLogger.ErrorWithTagsf, logrus.ErrorLevel},
		{"warn", testConfig.slf4GoLogrusLogger.Warnf, testConfig.slf4GoLogrusLogger.WarnWithTagsf, logrus.WarnLevel},
		{"warning", testConfig.slf4GoLogrusLogger.Warningf, testConfig.slf4GoLogrusLogger.WarningWithTagsf, logrus.WarnLevel},
		{"info", testConfig.slf4GoLogrusLogger.Infof, testConfig.slf4GoLogrusLogger.InfoWithTagsf, logrus.InfoLevel},
		{"debug", testConfig.slf4GoLogrusLogger.Debugf, testConfig.slf4GoLogrusLogger.DebugWithTagsf, logrus.DebugLevel},
		{"trace", testConfig.slf4GoLogrusLogger.Tracef, testConfig.slf4GoLogrusLogger.TraceWithTagsf, logrus.TraceLevel},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-base", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}, "test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message with name=%s and value=%d", "beeblebrox", 42)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(
					map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"},
					"test message with name=%s and value=%d",
					"beeblebrox",
					42,
				)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	t.Run("unknown", func(t *testing.T) {
		testConfig.hook.Reset()
		var unknownLevel slf4go_api.LogLevel = 666
		testConfig.slf4GoLogrusLogger.Logf(unknownLevel, "Some message not displayed")
		assertLog(t, testConfig.hook).
			hasLevel(logrus.ErrorLevel).
			hasTags(map[string]interface{}{}).
			hasMessage("Mapping error level 'unknown' onto Logrus error level failed. Not logging event")
	})

}

func TestLogging_WithStaticTags(t *testing.T) {
	testConfig := newTestingSetup().
		withStaticTags(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		})

	scenarios := []struct {
		name          string
		logFn         func(string, ...interface{})
		logFnWithTags func(slf4go_api.LogTags, string, ...interface{})
		logLevel      logrus.Level
	}{
		{"panic", testConfig.slf4GoLogrusLogger.Panicf, testConfig.slf4GoLogrusLogger.PanicWithTagsf, logrus.PanicLevel},
		{"error", testConfig.slf4GoLogrusLogger.Errorf, testConfig.slf4GoLogrusLogger.ErrorWithTagsf, logrus.ErrorLevel},
		{"warn", testConfig.slf4GoLogrusLogger.Warnf, testConfig.slf4GoLogrusLogger.WarnWithTagsf, logrus.WarnLevel},
		{"warning", testConfig.slf4GoLogrusLogger.Warningf, testConfig.slf4GoLogrusLogger.WarningWithTagsf, logrus.WarnLevel},
		{"info", testConfig.slf4GoLogrusLogger.Infof, testConfig.slf4GoLogrusLogger.InfoWithTagsf, logrus.InfoLevel},
		{"debug", testConfig.slf4GoLogrusLogger.Debugf, testConfig.slf4GoLogrusLogger.DebugWithTagsf, logrus.DebugLevel},
		{"trace", testConfig.slf4GoLogrusLogger.Tracef, testConfig.slf4GoLogrusLogger.TraceWithTagsf, logrus.TraceLevel},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-base", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1": "val1",
					"key2": "val2",
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}, "test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":     "val1",
					"key2":     "val2",
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message with name=%s and value=%d", "beeblebrox", 42)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1": "val1",
					"key2": "val2",
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(
					map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"},
					"test message with name=%s and value=%d",
					"beeblebrox",
					42,
				)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":     "val1",
					"key2":     "val2",
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	t.Run("unknown", func(t *testing.T) {
		testConfig.hook.Reset()
		var unknownLevel slf4go_api.LogLevel = 666
		testConfig.slf4GoLogrusLogger.Logf(unknownLevel, "Some message not displayed")
		assertLog(t, testConfig.hook).
			hasLevel(logrus.ErrorLevel).
			hasTags(map[string]interface{}{}).
			hasMessage("Mapping error level 'unknown' onto Logrus error level failed. Not logging event")
	})
}

func TestLogging_WithComponentLabel_WithStaticTags(t *testing.T) {
	testConfig := newTestingSetup().
		withComponentLabel("appLabel").
		withStaticTags(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		})

	scenarios := []struct {
		name          string
		logFn         func(string, ...interface{})
		logFnWithTags func(slf4go_api.LogTags, string, ...interface{})
		logLevel      logrus.Level
	}{
		{"panic", testConfig.slf4GoLogrusLogger.Panicf, testConfig.slf4GoLogrusLogger.PanicWithTagsf, logrus.PanicLevel},
		{"error", testConfig.slf4GoLogrusLogger.Errorf, testConfig.slf4GoLogrusLogger.ErrorWithTagsf, logrus.ErrorLevel},
		{"warn", testConfig.slf4GoLogrusLogger.Warnf, testConfig.slf4GoLogrusLogger.WarnWithTagsf, logrus.WarnLevel},
		{"warning", testConfig.slf4GoLogrusLogger.Warningf, testConfig.slf4GoLogrusLogger.WarningWithTagsf, logrus.WarnLevel},
		{"info", testConfig.slf4GoLogrusLogger.Infof, testConfig.slf4GoLogrusLogger.InfoWithTagsf, logrus.InfoLevel},
		{"debug", testConfig.slf4GoLogrusLogger.Debugf, testConfig.slf4GoLogrusLogger.DebugWithTagsf, logrus.DebugLevel},
		{"trace", testConfig.slf4GoLogrusLogger.Tracef, testConfig.slf4GoLogrusLogger.TraceWithTagsf, logrus.TraceLevel},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-base", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1": "val1",
					"key2": "val2",
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}, "test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":     "val1",
					"key2":     "val2",
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message with name=%s and value=%d", "beeblebrox", 42)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1": "val1",
					"key2": "val2",
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(
					map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"},
					"test message with name=%s and value=%d",
					"beeblebrox",
					42,
				)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":     "val1",
					"key2":     "val2",
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	t.Run("unknown", func(t *testing.T) {
		testConfig.hook.Reset()
		var unknownLevel slf4go_api.LogLevel = 666
		testConfig.slf4GoLogrusLogger.Logf(unknownLevel, "Some message not displayed")
		assertLog(t, testConfig.hook).
			hasLevel(logrus.ErrorLevel).
			hasTags(map[string]interface{}{}).
			hasMessage("Mapping error level 'unknown' onto Logrus error level failed. Not logging event")
	})
}

func TestLogging_ForComponent_WithComponentLabel(t *testing.T) {
	testConfig := newTestingSetup().
		forComponent("test-service").
		withComponentLabel("appLabel")

	scenarios := []struct {
		name          string
		logFn         func(string, ...interface{})
		logFnWithTags func(slf4go_api.LogTags, string, ...interface{})
		logLevel      logrus.Level
	}{
		{"panic", testConfig.slf4GoLogrusLogger.Panicf, testConfig.slf4GoLogrusLogger.PanicWithTagsf, logrus.PanicLevel},
		{"error", testConfig.slf4GoLogrusLogger.Errorf, testConfig.slf4GoLogrusLogger.ErrorWithTagsf, logrus.ErrorLevel},
		{"warn", testConfig.slf4GoLogrusLogger.Warnf, testConfig.slf4GoLogrusLogger.WarnWithTagsf, logrus.WarnLevel},
		{"warning", testConfig.slf4GoLogrusLogger.Warningf, testConfig.slf4GoLogrusLogger.WarningWithTagsf, logrus.WarnLevel},
		{"info", testConfig.slf4GoLogrusLogger.Infof, testConfig.slf4GoLogrusLogger.InfoWithTagsf, logrus.InfoLevel},
		{"debug", testConfig.slf4GoLogrusLogger.Debugf, testConfig.slf4GoLogrusLogger.DebugWithTagsf, logrus.DebugLevel},
		{"trace", testConfig.slf4GoLogrusLogger.Tracef, testConfig.slf4GoLogrusLogger.TraceWithTagsf, logrus.TraceLevel},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-base", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"appLabel": slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}, "test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
					"appLabel": slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message with name=%s and value=%d", "beeblebrox", 42)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"appLabel": slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(
					map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"},
					"test message with name=%s and value=%d",
					"beeblebrox",
					42,
				)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
					"appLabel": slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	t.Run("unknown", func(t *testing.T) {
		testConfig.hook.Reset()
		var unknownLevel slf4go_api.LogLevel = 666
		testConfig.slf4GoLogrusLogger.Logf(unknownLevel, "Some message not displayed")
		assertLog(t, testConfig.hook).
			hasLevel(logrus.ErrorLevel).
			hasTags(map[string]interface{}{}).
			hasMessage("Mapping error level 'unknown' onto Logrus error level failed. Not logging event")
	})
}

func TestLogging_ForComponent_WithStaticTags(t *testing.T) {
	testConfig := newTestingSetup().
		forComponent("test-service").
		withStaticTags(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		})

	scenarios := []struct {
		name          string
		logFn         func(string, ...interface{})
		logFnWithTags func(slf4go_api.LogTags, string, ...interface{})
		logLevel      logrus.Level
	}{
		{"panic", testConfig.slf4GoLogrusLogger.Panicf, testConfig.slf4GoLogrusLogger.PanicWithTagsf, logrus.PanicLevel},
		{"error", testConfig.slf4GoLogrusLogger.Errorf, testConfig.slf4GoLogrusLogger.ErrorWithTagsf, logrus.ErrorLevel},
		{"warn", testConfig.slf4GoLogrusLogger.Warnf, testConfig.slf4GoLogrusLogger.WarnWithTagsf, logrus.WarnLevel},
		{"warning", testConfig.slf4GoLogrusLogger.Warningf, testConfig.slf4GoLogrusLogger.WarningWithTagsf, logrus.WarnLevel},
		{"info", testConfig.slf4GoLogrusLogger.Infof, testConfig.slf4GoLogrusLogger.InfoWithTagsf, logrus.InfoLevel},
		{"debug", testConfig.slf4GoLogrusLogger.Debugf, testConfig.slf4GoLogrusLogger.DebugWithTagsf, logrus.DebugLevel},
		{"trace", testConfig.slf4GoLogrusLogger.Tracef, testConfig.slf4GoLogrusLogger.TraceWithTagsf, logrus.TraceLevel},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-base", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":                            "val1",
					"key2":                            "val2",
					slf4go_api.DefaultAppComponentTag: slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}, "test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":                            "val1",
					"key2":                            "val2",
					"dyn_key1":                        "dyn_val1",
					"dyn_key2":                        "dyn_val2",
					slf4go_api.DefaultAppComponentTag: slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message with name=%s and value=%d", "beeblebrox", 42)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":                            "val1",
					"key2":                            "val2",
					slf4go_api.DefaultAppComponentTag: slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(
					map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"},
					"test message with name=%s and value=%d",
					"beeblebrox",
					42,
				)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":                            "val1",
					"key2":                            "val2",
					"dyn_key1":                        "dyn_val1",
					"dyn_key2":                        "dyn_val2",
					slf4go_api.DefaultAppComponentTag: slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	t.Run("unknown", func(t *testing.T) {
		testConfig.hook.Reset()
		var unknownLevel slf4go_api.LogLevel = 666
		testConfig.slf4GoLogrusLogger.Logf(unknownLevel, "Some message not displayed")
		assertLog(t, testConfig.hook).
			hasLevel(logrus.ErrorLevel).
			hasTags(map[string]interface{}{}).
			hasMessage("Mapping error level 'unknown' onto Logrus error level failed. Not logging event")
	})
}

func TestLogging_ForComponent_WithComponentLabel_WithStaticTags(t *testing.T) {
	testConfig := newTestingSetup().
		forComponent("test-service").
		withComponentLabel("appLabel").
		withStaticTags(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		})

	scenarios := []struct {
		name          string
		logFn         func(string, ...interface{})
		logFnWithTags func(slf4go_api.LogTags, string, ...interface{})
		logLevel      logrus.Level
	}{
		{"panic", testConfig.slf4GoLogrusLogger.Panicf, testConfig.slf4GoLogrusLogger.PanicWithTagsf, logrus.PanicLevel},
		{"error", testConfig.slf4GoLogrusLogger.Errorf, testConfig.slf4GoLogrusLogger.ErrorWithTagsf, logrus.ErrorLevel},
		{"warn", testConfig.slf4GoLogrusLogger.Warnf, testConfig.slf4GoLogrusLogger.WarnWithTagsf, logrus.WarnLevel},
		{"warning", testConfig.slf4GoLogrusLogger.Warningf, testConfig.slf4GoLogrusLogger.WarningWithTagsf, logrus.WarnLevel},
		{"info", testConfig.slf4GoLogrusLogger.Infof, testConfig.slf4GoLogrusLogger.InfoWithTagsf, logrus.InfoLevel},
		{"debug", testConfig.slf4GoLogrusLogger.Debugf, testConfig.slf4GoLogrusLogger.DebugWithTagsf, logrus.DebugLevel},
		{"trace", testConfig.slf4GoLogrusLogger.Tracef, testConfig.slf4GoLogrusLogger.TraceWithTagsf, logrus.TraceLevel},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-base", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":     "val1",
					"key2":     "val2",
					"appLabel": slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"}, "test message")
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":     "val1",
					"key2":     "val2",
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
					"appLabel": slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFn("test message with name=%s and value=%d", "beeblebrox", 42)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":     "val1",
					"key2":     "val2",
					"appLabel": slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name+"-dyntags-formatted", func(t *testing.T) {
			testConfig.hook.Reset()
			fatalSafe(t, testConfig, scenario.logLevel, func() {
				scenario.logFnWithTags(
					map[string]interface{}{"dyn_key1": "dyn_val1", "dyn_key2": "dyn_val2"},
					"test message with name=%s and value=%d",
					"beeblebrox",
					42,
				)
			})
			assertLog(t, testConfig.hook).
				hasLevel(scenario.logLevel).
				hasTags(map[string]interface{}{
					"key1":     "val1",
					"key2":     "val2",
					"dyn_key1": "dyn_val1",
					"dyn_key2": "dyn_val2",
					"appLabel": slf4go_api.AppComponent("test-service"),
				}).
				hasMessage("test message with name=beeblebrox and value=42")
		})
	}

	t.Run("unknown", func(t *testing.T) {
		testConfig.hook.Reset()
		var unknownLevel slf4go_api.LogLevel = 666
		testConfig.slf4GoLogrusLogger.Logf(unknownLevel, "Some message not displayed")
		assertLog(t, testConfig.hook).
			hasLevel(logrus.ErrorLevel).
			hasTags(map[string]interface{}{}).
			hasMessage("Mapping error level 'unknown' onto Logrus error level failed. Not logging event")
	})
}

type testingSetup struct {
	slf4GoLogrusLogger *Slf4GoLogrusLogger
	hook               *test.Hook
}

func newTestingSetup() *testingSetup {
	var logrusLogger *logrus.Logger
	var hook *test.Hook
	logrusLogger, hook = test.NewNullLogger()
	logrusLogger.ExitFunc = func(int) {}
	logrusLogger.SetLevel(logrus.TraceLevel)
	return &testingSetup{
		slf4GoLogrusLogger: New(logrusLogger),
		hook:               hook,
	}
}

func (setup *testingSetup) forComponent(component slf4go_api.AppComponent) *testingSetup {
	setup.slf4GoLogrusLogger = setup.slf4GoLogrusLogger.ForComponent(component).(*Slf4GoLogrusLogger)
	return setup
}

func (setup *testingSetup) withComponentLabel(componentLabel string) *testingSetup {
	setup.slf4GoLogrusLogger = setup.slf4GoLogrusLogger.WithAppComponentLabel(componentLabel).(*Slf4GoLogrusLogger)
	return setup
}

func (setup *testingSetup) withStaticTags(tags map[string]interface{}) *testingSetup {
	setup.slf4GoLogrusLogger = setup.slf4GoLogrusLogger.WithStaticTags(tags).(*Slf4GoLogrusLogger)
	return setup
}

func fatalSafe(t *testing.T, setup *testingSetup, logLevel logrus.Level, logFn func()) {
	if logLevel == logrus.FatalLevel {
		origExit := setup.slf4GoLogrusLogger.logger.ExitFunc
		called := false
		setup.slf4GoLogrusLogger.logger.ExitFunc = func(code int) {
			called = true
		}
		logFn()
		assert.True(t, called, "exitFunc should have been called")
		setup.slf4GoLogrusLogger.logger.ExitFunc = origExit // restore
	} else if logLevel == logrus.PanicLevel {
		assert.Panics(t, logFn)
	} else {
		logFn()
	}
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

func (a *logAssertions) hasTags(tags slf4go_api.LogTags) *logAssertions {
	entry := a.hook.LastEntry()
	assert.Equal(a.t, logrus.Fields(tags), entry.Data)
	return a
}

func (a *logAssertions) hasMessage(msg string) *logAssertions {
	assert.Equal(a.t, msg, a.hook.LastEntry().Message)
	return a
}
