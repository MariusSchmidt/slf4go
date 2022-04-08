package logging

import "github.com/sirupsen/logrus"

type LogrusLogger struct {
	logger *logrus.Logger
}

func (l *LogrusLogger) Trace(args ...interface{}) {
	l.logger.Trace(args)
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args)
}

func (l *LogrusLogger) Print(args ...interface{}) {
	l.logger.Print(args)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args)
}

func (l *LogrusLogger) Warning(args ...interface{}) {
	l.logger.Warning(args)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args)
}

func (l *LogrusLogger) Panic(args ...interface{}) {
	l.logger.Panic(args)
}

func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args)
}

func (l *LogrusLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args)
}

func (l *LogrusLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args)
}

func (l *LogrusLogger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, args)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args)
}

func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args)
}

func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusFieldLogger{entry: l.logger.WithFields(fields)}
}
