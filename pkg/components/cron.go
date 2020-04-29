package components

import (
	"fmt"

	"github.com/robfig/cron"

	"goweb/pkg/handler"
)

var Cron = new(CronService)

type CronService struct {
	cron   *cron.Cron
	logger Logger
}

type CronMsg struct {
	HandlerName string
	Exp         string
}

func (c *CronService) Run(logger ...Logger) {
	c.logger = logger[0]
	msgs := c.getCronMsg()
	c.cron = cron.New()
	for _, msg := range msgs {
		handlerFunc := c.getCronHandler(msg.HandlerName)
		if handlerFunc == nil {
			continue
		}
		err := c.cron.AddFunc(msg.Exp, handlerFunc)
		if err != nil {
			c.logger.Errorf(fmt.Sprintf("execute cron %s error", msg.HandlerName), err)
		}
	}
	c.cron.Start()
}

func (c *CronService) OnShutdown() func() {
	return func() {
		c.cron.Stop()
	}
}

func (c *CronService) getCronMsg() (m []CronMsg) {
	cronMap := Conf.Cron
	for k, v := range cronMap {
		m = append(m, CronMsg{
			HandlerName: k,
			Exp:         v,
		})
	}
	return
}

func (c *CronService) getCronHandler(name string) (f func()) {
	if h, ok := handler.GlobalHandlers[name]; ok {
		f = c.makeCronHandler(h)
	} else {
		c.logger.Errorf(fmt.Sprintf("cron %s not found", name))
	}
	return
}

func (c *CronService) makeCronHandler(globalHandler handler.GlobalHandler) func() {
	return func() {
		code, err := globalHandler()
		if err != nil {
			c.logger.Errorf("cron err", err, code)
		}
	}
}
