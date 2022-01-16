package logs

// Logger defines what we want to be able to use to log information while handling a request
type Logger interface {
	Printf(format string, v ...interface{})
}
