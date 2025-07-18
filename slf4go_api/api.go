package slf4go_api

//go:generate go run github.com/golang/mock/mockgen@latest -package=test_mocks -destination=./test_mocks/slf4go_mock_logger.go github.com/MariusSchmidt/slf4go/slf4go_api Slf4GoLogger

// AllLevels contains all available log levels in descending order of severity.
var AllLevels = []LogLevel{
	Fatal,
	Panic,
	Error,
	Warn,
	Info,
	Debug,
	Trace,
}

// Constants for the different logging levels, sorted by descending severity.
const (
	// Fatal logs the message and terminates the program with Exit(1),
	// even if the logging level is set to Panic.
	Fatal LogLevel = iota

	// Panic is the highest severity level. The logger logs the message and then calls panic().
	Panic

	// Error is used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	Error

	// Warn is used for non-critical events that deserve attention.
	Warn

	// Info is used for general operational entries about what's going on inside the application.
	Info

	// Debug is usually only enabled during development and produces very verbose logging.
	Debug

	// Trace designates even finer-grained informational events than Debug.
	Trace
)

// LogLevel defines the different logging levels as uint32.
type LogLevel uint32

// Stringer converts a log level to its string representation.
func (level LogLevel) Stringer() string {
	switch level {
	case Trace:
		return "trace"
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warning"
	case Error:
		return "error"
	case Fatal:
		return "fatal"
	case Panic:
		return "panic"
	default:
		return "unknown"
	}
}

// DefaultAppComponentTag defines the default key for component tags in log entries
const DefaultAppComponentTag string = "appComponent"

// Slf4GoLogger defines an interface for structured logging.
// It supports various log levels and the ability to add additional tags to log entries.
type Slf4GoLogger interface {
	// ForComponent creates a new Slf4GoLogger instance that will include the specified component
	// in all log entries. The component will be logged under the DefaultAppComponentTag key.
	// The returned logger inherits all other properties from the original logger.
	ForComponent(component AppComponent) Slf4GoLogger

	// WithAppComponentLabel creates a new Slf4GoLogger instance with a custom tag label for the component.
	// The component will be logged under the specified tag label instead of DefaultAppComponentTag.
	// The returned logger inherits all other properties from the original logger.
	WithAppComponentLabel(appComponentLabel string) Slf4GoLogger

	// WithStaticTags creates a new Slf4GoLogger instance with predefined tags
	// that will be added to every log entry.
	WithStaticTags(tags LogTags) Slf4GoLogger

	// Log logs a message with the specified level and formatted text.
	Log(level LogLevel, message string)

	// Logf logs a message with the specified level and formatted text.
	Logf(level LogLevel, format string, args ...interface{})

	// LogWithTags logs a message with the specified level, additional tags, and formatted text.
	LogWithTags(level LogLevel, tags LogTags, message string)

	// LogWithTagsf logs a message with the specified level, additional tags, and formatted text.
	LogWithTagsf(level LogLevel, tags LogTags, messageTemplate string, args ...interface{})

	// Fatalf logs critical errors using the specified format and arguments, then terminates the program.
	Fatalf(format string, args ...interface{})
	// Panicf logs severe errors using the specified format and arguments, then panics.
	Panicf(format string, args ...interface{})
	// Errorf logs errors using the specified format and arguments.
	Errorf(format string, args ...interface{})
	// Warnf logs warnings using the specified format and arguments.
	Warnf(format string, args ...interface{})
	// Warningf is an alias for Warnf that logs warnings using the specified format and arguments.
	Warningf(format string, args ...interface{})
	// Infof logs general information using the specified format and arguments.
	Infof(format string, args ...interface{})
	// Debugf logs debug information using the specified format and arguments.
	Debugf(format string, args ...interface{})
	// Tracef logs very detailed debug information using the specified format and arguments.
	Tracef(format string, args ...interface{})

	// FatalWithTagsf logs a fatal message with additional tags using the specified format and arguments, then terminates
	FatalWithTagsf(tags LogTags, format string, args ...interface{})
	// PanicWithTagsf logs a panic message with additional tags using the specified format and arguments, then panics
	PanicWithTagsf(tags LogTags, format string, args ...interface{})
	// ErrorWithTagsf logs an error message with additional tags using the specified format and arguments
	ErrorWithTagsf(tags LogTags, format string, args ...interface{})
	// WarnWithTagsf logs a warning message with additional tags using the specified format and arguments
	WarnWithTagsf(tags LogTags, format string, args ...interface{})
	// WarningWithTagsf is an alias for WarnWithTagsf that logs a warning message with additional tags
	WarningWithTagsf(tags LogTags, format string, args ...interface{})
	// InfoWithTagsf logs an info message with additional tags using the specified format and arguments
	InfoWithTagsf(tags LogTags, format string, args ...interface{})
	// DebugWithTagsf logs a debug message with additional tags using the specified format and arguments
	DebugWithTagsf(tags LogTags, format string, args ...interface{})
	// TraceWithTagsf logs a trace message with additional tags using the specified format and arguments
	TraceWithTagsf(tags LogTags, format string, args ...interface{})
}

// AppComponent represents a significant component of the application to be mentioned in logs.
type AppComponent string

// LogTags is a type alias for a map that adds additional structured data
// to a log entry. The values can be of any type.
type LogTags map[string]interface{}
