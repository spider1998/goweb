package engine

import (
	"fmt"
	"strings"

	"goweb/pkg/code"
)

type Job struct {
	Eng     *GlobalEngine
	Name    string
	Args    []string
	handler Handler
	code    code.Code
}

func (job *Job) Run() error {
	if job.Eng.IsShutdown() {
		return fmt.Errorf("engine is shutdown")
	}
	log := job.Eng.Logger
	log.Infof("+job %s", job.Call())
	defer func() {
		log.Infof("-job %s%s", job.Call(), job.CodeString())
	}()
	if job.handler == nil {
		log.Errorf("%s: command not found", job.Name)
	} else {
		var err error
		job.code, err = job.handler(job)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}
	return nil
}

func (job *Job) Call() string {
	return fmt.Sprintf("%s(%s)", job.Name, strings.Join(job.Args, ", "))
}

func (job *Job) CodeString() string {
	// If the job hasn't completed, status string is empty
	var err string
	if job.code == code.StatusOk {
		err = "OK"
	} else {
		err = "ERR"
	}
	return fmt.Sprintf(" = %s (%d)", err, job.code)
}
