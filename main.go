package main

import (
	"fmt"
	"runtime"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("(%s:%d)", f.File, f.Line)
		},
	})

	log.SetReportCaller(true)
	log.Infof("started")
}

func main() {
	ct := &Crontab{
		cc: cron.New(),
	}

	ct.Setup()

	ct.cc.Start()

	select {}
}
