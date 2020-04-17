package cron

import (
	"fmt"
	"time"

	cr "github.com/robfig/cron"

	"goweb/engine"
)

type WebCron struct {
	Tasks []*Task
	Add   chan *Task
	Stop  chan struct{}
}

var Cron *WebCron

func Init() {
	Cron = NewCron()
	Cron.Watch()
}

func NewCron() *WebCron {
	return &WebCron{
		Add:  make(chan *Task, 1),
		Stop: make(chan struct{}),
	}
}

type Task struct {
	JobName      string
	RunTime      time.Time
	IntervalTime time.Duration
}

type TaskSource struct {
	JobName    string
	Expression string
}

func (c *WebCron) AddFunc(t TaskSource) error {
	schedule, err := cr.Parse(t.Expression)
	if err != nil {
		return err
	}
	var task Task
	task.JobName = t.JobName
	task.RunTime = schedule.Next(time.Now())
	task.IntervalTime = schedule.Next(task.RunTime).Sub(task.RunTime)
	c.Tasks = append(c.Tasks, &task)
	c.Add <- &task
	return nil
}

func (c *WebCron) Watch() {
	for {
		select {
		case t := <-c.Add:
			go func(task *Task) {
				for {
					code, err := engine.GlobalHandlers[task.JobName]()
					if err != nil {
						_ = fmt.Errorf("cron %s exceute error : %d : %s ", task.JobName, code, err)
					}
					now := time.Now()
					next := now.Add(task.IntervalTime)
					t := time.NewTimer(next.Sub(now))
					<-t.C
				}
			}(t)
		case <-c.Stop:
			close(c.Add)
			return

		}
	}
}
