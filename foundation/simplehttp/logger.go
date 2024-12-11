package simplehttp

// RestyLogger should be refactored. It was implemented to avoid logs emitted by default by resty.
type RestyLogger struct{}

func (rl RestyLogger) Errorf(format string, v ...interface{}) {
	return
}
func (rl RestyLogger) Warnf(format string, v ...interface{}) {
	return
}
func (rl RestyLogger) Debugf(format string, v ...interface{}) {
	return
}
