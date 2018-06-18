// +build !windows,!nacl,!plan9

package syslog

import (
	"fmt"
	"log/syslog"
	"os"
	"sillyhat-golang-tool/sillyhat_log/logrus"
	elastic "sillyhat-golang-tool/sillyhat_elasticsearch"
	"time"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	Writer        *syslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewSyslogHook(network, raddr string, priority syslog.Priority, tag string) (*SyslogHook, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &SyslogHook{w, network, raddr}, err
}

func (hook *SyslogHook) Fire(entry *sillyhat_logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case sillyhat_logrus.PanicLevel:
		write(entry,line)
		return hook.Writer.Crit(line)
	case sillyhat_logrus.FatalLevel:
		write(entry,line)
		return hook.Writer.Crit(line)
	case sillyhat_logrus.ErrorLevel:
		write(entry,line)
		return hook.Writer.Err(line)
	case sillyhat_logrus.WarnLevel:
		write(entry,line)
		return hook.Writer.Warning(line)
	case sillyhat_logrus.InfoLevel:
		write(entry,line)
		return hook.Writer.Info(line)
	case sillyhat_logrus.DebugLevel:
		write(entry,line)
		return hook.Writer.Debug(line)
	default:
		return nil
	}
}

func write(entry *sillyhat_logrus.Entry,line string) {
	switch entry.HookType {
	case sillyhat_logrus.Elasticsearch:
		elasticsearch(entry.WriteLogProperties,line)
	case sillyhat_logrus.KAFKA:

	default:

	}
}

const (
	elastic_type = "gostash"
)

func getElasticsearchIndex() string {
	return "gostash-" + time.Now().Format("2006.01.02")
}

func elasticsearch(writeLogProperties sillyhat_logrus.WriteLogProperties,line string)  {
	elastic.SetURL(writeLogProperties.URL)
	elastic.SetIndex(getElasticsearchIndex())
	elastic.SetType(elastic_type)
	elastic.Index(line)
}

func (hook *SyslogHook) Levels() []sillyhat_logrus.Level {
	return sillyhat_logrus.AllLevels
}
