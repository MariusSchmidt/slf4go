package slf4go_logrus_provider

import (
	"github.com/MariusSchmidt/slf4go/slf4go_api"
	"github.com/sirupsen/logrus"
)

type Slf4GoLogrusLogger struct {
	logger            *logrus.Logger
	appComponent      slf4go_api.AppComponent
	tags              slf4go_api.LogTags
	componentTagLabel string
}

// New creates a new slf4GoLogrusLogger with optional configurations
func New(logrusLogger *logrus.Logger) *Slf4GoLogrusLogger {
	return &Slf4GoLogrusLogger{
		logger:            logrusLogger,
		appComponent:      "",
		tags:              make(slf4go_api.LogTags),
		componentTagLabel: slf4go_api.DefaultAppComponentTag,
	}
}

func (l *Slf4GoLogrusLogger) ForComponent(component slf4go_api.AppComponent) slf4go_api.Slf4GoLogger {
	return &Slf4GoLogrusLogger{
		logger:            l.logger,
		appComponent:      component,
		tags:              l.tags,
		componentTagLabel: l.componentTagLabel,
	}
}

func (l *Slf4GoLogrusLogger) WithAppComponentLabel(componentTagLabel string) slf4go_api.Slf4GoLogger {
	return &Slf4GoLogrusLogger{
		logger:            l.logger,
		appComponent:      l.appComponent,
		tags:              l.tags,
		componentTagLabel: componentTagLabel,
	}
}

func (l *Slf4GoLogrusLogger) WithStaticTags(tags slf4go_api.LogTags) slf4go_api.Slf4GoLogger {
	return &Slf4GoLogrusLogger{
		logger:            l.logger,
		appComponent:      l.appComponent,
		tags:              tags,
		componentTagLabel: l.componentTagLabel,
	}
}

func (l *Slf4GoLogrusLogger) Logf(level slf4go_api.LogLevel, msgTemplate string, args ...interface{}) {
	l.LogWithTagsf(level, slf4go_api.LogTags{}, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) logrusLogf(level logrus.Level, msgTemplate string, args ...interface{}) {
	switch level {
	case logrus.FatalLevel:
		l.logger.Fatalf(msgTemplate, args...)
	case logrus.PanicLevel:
		l.logger.Panicf(msgTemplate, args...)
	case logrus.ErrorLevel:
		l.logger.Errorf(msgTemplate, args...)
	case logrus.WarnLevel:
		l.logger.Warnf(msgTemplate, args...)
	case logrus.InfoLevel:
		l.logger.Infof(msgTemplate, args...)
	case logrus.DebugLevel:
		l.logger.Debugf(msgTemplate, args...)
	case logrus.TraceLevel:
		l.logger.WithFields(logrus.Fields{}).Tracef(msgTemplate, args...)
	}
}

func (l *Slf4GoLogrusLogger) logrusLogWithTagsf(level logrus.Level, fields logrus.Fields, msgTemplate string, args ...interface{}) {
	entry := l.logger.WithFields(fields)
	switch level {
	case logrus.FatalLevel:
		entry.Fatalf(msgTemplate, args...)
	case logrus.PanicLevel:
		entry.Panicf(msgTemplate, args...)
	case logrus.ErrorLevel:
		entry.Errorf(msgTemplate, args...)
	case logrus.WarnLevel:
		entry.Warnf(msgTemplate, args...)
	case logrus.InfoLevel:
		entry.Infof(msgTemplate, args...)
	case logrus.DebugLevel:
		entry.Debugf(msgTemplate, args...)
	case logrus.TraceLevel:
		entry.Tracef(msgTemplate, args...)
	}
}

func (l *Slf4GoLogrusLogger) LogWithTagsf(level slf4go_api.LogLevel, tags slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	logrusLevel, err := logrus.ParseLevel(level.Stringer())
	if err != nil {
		l.logger.Errorf("Mapping error level '%s' onto Logrus error level failed. Not logging event", level.Stringer())
		return
	}
	tags = combineTags(l.tags, tags)
	if len(l.appComponent) == 0 && len(tags) == 0 {
		l.logrusLogf(logrusLevel, msgTemplate, args...)
	}
	if len(l.appComponent) == 0 && len(tags) >= 1 {
		l.logrusLogWithTagsf(logrusLevel, logrus.Fields(tags), msgTemplate, args...)
	}
	if len(l.appComponent) >= 1 && len(tags) == 0 {
		l.logrusLogWithTagsf(logrusLevel, logrus.Fields{l.componentTagLabel: l.appComponent}, msgTemplate, args...)
	}
	if len(l.appComponent) >= 1 && len(tags) >= 1 {
		tagsAsLogrusFields := combineTags(tags, slf4go_api.LogTags{l.componentTagLabel: l.appComponent})
		l.logrusLogWithTagsf(logrusLevel, logrus.Fields(tagsAsLogrusFields), msgTemplate, args...)
	}
}

func combineTags(tags ...slf4go_api.LogTags) slf4go_api.LogTags {
	merged := make(slf4go_api.LogTags)
	for _, m := range tags {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func (l *Slf4GoLogrusLogger) Tracef(msgTemplate string, args ...interface{}) {
	l.Logf(slf4go_api.Trace, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) Debugf(msgTemplate string, args ...interface{}) {
	l.Logf(slf4go_api.Debug, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) Infof(msgTemplate string, args ...interface{}) {
	l.Logf(slf4go_api.Info, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) Warnf(msgTemplate string, args ...interface{}) {
	l.Warningf(msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) Warningf(msgTemplate string, args ...interface{}) {
	l.Logf(slf4go_api.Warn, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) Errorf(msgTemplate string, args ...interface{}) {
	l.Logf(slf4go_api.Error, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) Panicf(msgTemplate string, args ...interface{}) {
	l.Logf(slf4go_api.Panic, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) Fatalf(msgTemplate string, args ...interface{}) {
	l.Logf(slf4go_api.Fatal, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) TraceWithTagsf(fields slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Trace, fields, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) DebugWithTagsf(fields slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Debug, fields, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) InfoWithTagsf(fields slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Info, fields, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) WarnWithTagsf(fields slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	l.WarningWithTagsf(fields, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) WarningWithTagsf(fields slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Warn, fields, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) ErrorWithTagsf(fields slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Error, fields, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) PanicWithTagsf(fields slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Panic, fields, msgTemplate, args...)
}

func (l *Slf4GoLogrusLogger) FatalWithTagsf(fields slf4go_api.LogTags, msgTemplate string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Fatal, fields, msgTemplate, args...)
}
