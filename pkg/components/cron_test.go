package components

import (
	"testing"
	"time"

	"goweb/pkg/register"
)

func TestCronService_Run(t *testing.T) {
	logger := NewLogger()
	register.Init()
	Cron.Run(logger)
	time.Sleep(time.Second * 20)
	Cron.Shutdown()
}
