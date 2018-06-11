package main

import (
	"sillyhat-golang-tool/sillyhat_log/logrus"
	//lSyslog "sillyhat-golang-tool/sillyhat_log/logrus/hooks/syslog"
	//"log/syslog"
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

func main() {
	log := logrus.New()
	jsonFormatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  field_key_time,
			logrus.FieldKeyLevel: field_key_level,
			logrus.FieldKeyMsg:   field_key_msg,
		},
		TimestampFormat: format_time,
	}
	log.Formatter = jsonFormatter
	log.Level = logrus.DebugLevel
	//hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")

	//if err == nil {
	//	log.Hooks.Add(hook)
	//}


	log.Infof("test %v","haha")
	log.Infof("test %v","hehe")
}