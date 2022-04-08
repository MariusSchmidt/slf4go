package logging

import "github.com/sirupsen/logrus"

type LogrusFieldLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func (l *LogrusFieldLogger) Trace(args ...interface{}) {
	l.entry.Trace(args)
}

func (l *LogrusFieldLogger) Debug(args ...interface{}) {
	l.entry.Debug(args)
}

func (l *LogrusFieldLogger) Print(args ...interface{}) {
	l.entry.Print(args)
}

func (l *LogrusFieldLogger) Info(args ...interface{}) {
	l.entry.Info(args)
}

func (l *LogrusFieldLogger) Warning(args ...interface{}) {
	l.entry.Warning(args)
}

func (l *LogrusFieldLogger) Error(args ...interface{}) {
	l.entry.Error(args)
}

func (l *LogrusFieldLogger) Panic(args ...interface{}) {
	l.entry.Panic(args)
}

func (l *LogrusFieldLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args)
}

func (l *LogrusFieldLogger) Tracef(format string, args ...interface{}) {
	l.entry.Tracef(format, args)
}

func (l *LogrusFieldLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args)
}

func (l *LogrusFieldLogger) Printf(format string, args ...interface{}) {
	l.entry.Printf(format, args)
}

func (l *LogrusFieldLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args)
}

func (l *LogrusFieldLogger) Warningf(format string, args ...interface{}) {
	l.entry.Warningf(format, args)
}

func (l *LogrusFieldLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args)
}

func (l *LogrusFieldLogger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args)
}

func (l *LogrusFieldLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args)
}

func (l *LogrusFieldLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusFieldLogger{entry: l.logger.WithFields(fields)}
}
