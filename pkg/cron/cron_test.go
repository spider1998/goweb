package cron

import (
	"goweb/pkg/log"
	"goweb/pkg/register"
	"testing"
	"time"
)

func TestCronService_Run(t *testing.T) {
	logger := log.NewLogger()
	register.Init()
	Cron.Run(logger)
	time.Sleep(time.Second * 20)
	Cron.Shutdown()
}
