package main

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	hookSyslog "sillyhat-golang-tool/sillyhat_log/logrus/hooks/syslog"
	"log/syslog"
)

const (
	format_time = "2006-01-02 15:04:05"
)

//func getLogrusEntry() *sillyhat_logrus.Entry {
//	logNew := sillyhat_logrus.New()
//	jsonFormatter := &sillyhat_logrus.JSONFormatter{
//		FieldMap: sillyhat_logrus.FieldMap{
//			sillyhat_logrus.FieldKeyTime:  field_key_time,
//			sillyhat_logrus.FieldKeyLevel: field_key_level,
//			sillyhat_logrus.FieldKeyMsg:   field_key_msg,
//		},
//		TimestampFormat: format_time,
//	}
//	logNew.Formatter = jsonFormatter
//	logNew.Level = sillyhat_logrus.DebugLevel
//
//	if pc, file, line, ok := runtime.Caller(1); ok {
//		fName := runtime.FuncForPC(pc).Name()
//		return logNew.WithField(field_file, file).WithField(field_line, line).WithField(field_func, fName).WithField(field_module_name,"test-module")
//	}
//	return logNew.WithField(field_module_name,"test-module")
//}

func init() {
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: format_time})
	//log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetModuleName("test-module")
	log.SetHookType(log.Elasticsearch)
	log.SetWriteLogProperties(log.WriteLogProperties{URL:"http://172.28.2.25:9200"})
	hook, err := hookSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err != nil {
		panic(err)
	}
	log.AddHook(hook)
}

func main() {
	//logstash-* logstash-2018.04.25

	//log := sillyhat_logrus.New()
	//jsonFormatter := &sillyhat_logrus.JSONFormatter{
	//	FieldMap: sillyhat_logrus.FieldMap{
	//		sillyhat_logrus.FieldKeyTime:  field_key_time,
	//		sillyhat_logrus.FieldKeyLevel: field_key_level,
	//		sillyhat_logrus.FieldKeyMsg:   field_key_msg,
	//	},
	//	TimestampFormat: format_time,
	//}
	//log.Formatter = jsonFormatter
	//log.Level = sillyhat_logrus.DebugLevel

	log.Infof("test %v","haha")
	log.Infof("test %v","12314")
	log.Infof("test %v","34623464")
	log.Infof("test %v","hedfgsdfghdhe")
	log.Infof("test %v","45734")
	log.Infof("test %v","asgasgasg")
	log.Infof("test %v","he357457428dvbxcvbhe")
}