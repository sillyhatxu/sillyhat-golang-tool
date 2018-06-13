package main

import (
	"sillyhat-golang-tool/sillyhat_log/logrus"
	"runtime"
)

const (
	format_time = "2006-01-02 15:04:05"
	field_key_msg   = "message"
	field_key_level = "level"
	field_key_time  = "timestamp"
	field_module_name  = "module_name"
	field_file  = "file"
	field_line  = "line"
	field_func  = "func"
)

func getLogrusEntry() *logrus.Entry {
	logNew := logrus.New()
	jsonFormatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  field_key_time,
			logrus.FieldKeyLevel: field_key_level,
			logrus.FieldKeyMsg:   field_key_msg,
		},
		TimestampFormat: format_time,
	}
	logNew.Formatter = jsonFormatter
	logNew.Level = logrus.DebugLevel

	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		return logNew.WithField(field_file, file).WithField(field_line, line).WithField(field_func, fName).WithField(field_module_name,"test-module")
	}
	return logNew.WithField(field_module_name,"test-module")
}

func main() {
	//logstash-* logstash-2018.04.25

	//log := logrus.New()
	//jsonFormatter := &logrus.JSONFormatter{
	//	FieldMap: logrus.FieldMap{
	//		logrus.FieldKeyTime:  field_key_time,
	//		logrus.FieldKeyLevel: field_key_level,
	//		logrus.FieldKeyMsg:   field_key_msg,
	//	},
	//	TimestampFormat: format_time,
	//}
	//log.Formatter = jsonFormatter
	//log.Level = logrus.DebugLevel
	//hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	//
	//if err == nil {
	//	log.Hooks.Add(hook)
	//}
	//
	//log.Info(hook)

	getLogrusEntry().Infof("test %v","haha")
	getLogrusEntry().Infof("test %v","hehe")
}