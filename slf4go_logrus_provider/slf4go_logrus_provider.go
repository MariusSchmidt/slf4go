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

// New creates a new logger with optional configurations
func New(logrusLogger *logrus.Logger) slf4go_api.Slf4GoLogger {
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

func (l *Slf4GoLogrusLogger) Log(level slf4go_api.LogLevel, message string) {
	l.LogWithTagsf(level, slf4go_api.LogTags{}, message)
}

func (l *Slf4GoLogrusLogger) Logf(level slf4go_api.LogLevel, messageTemplate string, args ...interface{}) {
	l.LogWithTagsf(level, slf4go_api.LogTags{}, messageTemplate, args...)
}

func (l *Slf4GoLogrusLogger) logrusLogf(level logrus.Level, messageTemplate string, args ...interface{}) {
	switch level {
	case logrus.PanicLevel:
		l.logger.Panicf(messageTemplate, args...)
	case logrus.FatalLevel:
		l.logger.Fatalf(messageTemplate, args...)
	case logrus.ErrorLevel:
		l.logger.Errorf(messageTemplate, args...)
	case logrus.WarnLevel:
		l.logger.Warnf(messageTemplate, args...)
	case logrus.InfoLevel:
		l.logger.Infof(messageTemplate, args...)
	case logrus.DebugLevel:
		l.logger.Debugf(messageTemplate, args...)
	case logrus.TraceLevel:
		l.logger.WithFields(logrus.Fields{}).Tracef(messageTemplate, args...)
	}
}

func (l *Slf4GoLogrusLogger) logrusLogWithTagsf(level logrus.Level, fields logrus.Fields, format string, args ...interface{}) {
	entry := l.logger.WithFields(fields)
	switch level {
	case logrus.PanicLevel:
		entry.Panicf(format, args...)
	case logrus.FatalLevel:
		entry.Fatalf(format, args...)
	case logrus.ErrorLevel:
		entry.Errorf(format, args...)
	case logrus.WarnLevel:
		entry.Warnf(format, args...)
	case logrus.InfoLevel:
		entry.Infof(format, args...)
	case logrus.DebugLevel:
		entry.Debugf(format, args...)
	case logrus.TraceLevel:
		entry.Tracef(format, args...)
	}
}

func (l *Slf4GoLogrusLogger) LogWithTags(level slf4go_api.LogLevel, fields slf4go_api.LogTags, message string) {
	l.LogWithTagsf(level, fields, message)
}

func (l *Slf4GoLogrusLogger) LogWithTagsf(level slf4go_api.LogLevel, tags slf4go_api.LogTags, messageTemplate string, args ...interface{}) {
	logrusLevel, err := logrus.ParseLevel(level.Stringer())
	if err != nil {
		l.logger.Errorf("Mapping error level '%s' onto Logrus error level failed. Not logging event", level.Stringer())
		return
	}
	combineTags(l.tags, tags)
	if len(l.appComponent) == 0 && len(l.tags) == 0 {
		l.logrusLogf(logrusLevel, messageTemplate, args...)
	}
	if len(l.appComponent) == 0 && len(l.tags) >= 1 {
		l.logrusLogWithTagsf(logrusLevel, logrus.Fields(l.tags), messageTemplate, args...)
	}
	if len(l.appComponent) >= 1 && len(l.tags) == 0 {
		l.logrusLogWithTagsf(logrusLevel, logrus.Fields{l.componentTagLabel: l.appComponent}, messageTemplate, args...)
	}
	if len(l.appComponent) >= 1 && len(l.tags) >= 1 {
		tagsAsLogrusFields := combineTags(l.tags, slf4go_api.LogTags{l.componentTagLabel: l.appComponent})
		l.logrusLogWithTagsf(logrusLevel, logrus.Fields(tagsAsLogrusFields), messageTemplate, args...)
	}
}

func combineTags(t1 slf4go_api.LogTags, t2 slf4go_api.LogTags) slf4go_api.LogTags {
	merged := t1
	for k, v := range t2 {
		merged[k] = v
	}
	return merged
}

func (l *Slf4GoLogrusLogger) Tracef(format string, args ...interface{}) {
	l.Logf(slf4go_api.Trace, format, args...)
}

func (l *Slf4GoLogrusLogger) Debugf(format string, args ...interface{}) {
	l.Logf(slf4go_api.Debug, format, args...)
}

func (l *Slf4GoLogrusLogger) Infof(format string, args ...interface{}) {
	l.Logf(slf4go_api.Info, format, args...)
}

func (l *Slf4GoLogrusLogger) Warnf(format string, args ...interface{}) {
	l.Warningf(format, args...)
}

func (l *Slf4GoLogrusLogger) Warningf(format string, args ...interface{}) {
	l.Logf(slf4go_api.Warn, format, args...)
}

func (l *Slf4GoLogrusLogger) Errorf(format string, args ...interface{}) {
	l.Logf(slf4go_api.Error, format, args...)
}

func (l *Slf4GoLogrusLogger) Panicf(format string, args ...interface{}) {
	l.Logf(slf4go_api.Panic, format, args...)
}

func (l *Slf4GoLogrusLogger) Fatalf(format string, args ...interface{}) {
	l.Logf(slf4go_api.Fatal, format, args...)
}

func (l *Slf4GoLogrusLogger) TraceWithTagsf(fields slf4go_api.LogTags, format string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Trace, fields, format, args...)
}

func (l *Slf4GoLogrusLogger) DebugWithTagsf(fields slf4go_api.LogTags, format string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Debug, fields, format, args...)
}

func (l *Slf4GoLogrusLogger) InfoWithTagsf(fields slf4go_api.LogTags, format string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Info, fields, format, args...)
}

func (l *Slf4GoLogrusLogger) WarnWithTagsf(fields slf4go_api.LogTags, format string, args ...interface{}) {
	l.WarningWithTagsf(fields, format, args...)
}

func (l *Slf4GoLogrusLogger) WarningWithTagsf(fields slf4go_api.LogTags, format string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Warn, fields, format, args...)
}

func (l *Slf4GoLogrusLogger) ErrorWithTagsf(fields slf4go_api.LogTags, format string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Error, fields, format, args...)
}

func (l *Slf4GoLogrusLogger) PanicWithTagsf(fields slf4go_api.LogTags, format string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Panic, fields, format, args...)
}

func (l *Slf4GoLogrusLogger) FatalWithTagsf(fields slf4go_api.LogTags, format string, args ...interface{}) {
	l.LogWithTagsf(slf4go_api.Fatal, fields, format, args...)
}
