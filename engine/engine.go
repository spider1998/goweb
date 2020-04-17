package engine

import (
	"fmt"
	"sync"
	"time"

	"goweb/code"
	"goweb/log"
)

type Handler func(*Job) (code code.Code, err error)

type GlobalHandler func() (code code.Code, err error)

type GlobalEngine struct {
	handlers     map[string]Handler
	emptyHandler Handler
	l            sync.RWMutex
	shutDowns    []func()
	stop         bool
	Logger       *log.WebLogger
}

var Engine = NewEngine()

var GlobalHandlers map[string]GlobalHandler

func init() {
	GlobalHandlers = make(map[string]GlobalHandler)
}

func Register(name string, handler GlobalHandler) {
	_, exists := GlobalHandlers[name]
	if exists {
		_ = fmt.Errorf("Can't overwrite global handler for command %s ", name)
		return
	}
	GlobalHandlers[name] = handler
	return
}

func NewEngine() *GlobalEngine {
	eng := &GlobalEngine{
		handlers: make(map[string]Handler),
		Logger:   log.NewLogger(),
	}
	go eng.Logger.Watch()
	return eng
}

func (eng *GlobalEngine) Register(name string, handler Handler) error {
	_, exists := eng.handlers[name]
	if exists {
		return fmt.Errorf("Can't overwrite handler for command %s ", name)
	}
	eng.handlers[name] = handler
	return nil
}

func (eng *GlobalEngine) Job(name string, args ...string) *Job {
	job := &Job{
		Eng:  eng,
		Name: name,
		Args: args,
	}

	if handler, exists := eng.handlers[name]; exists {
		job.handler = handler
	} else if eng.emptyHandler != nil && name != "" {
		job.handler = eng.emptyHandler
	}
	return job
}

func (eng *GlobalEngine) Shutdown() {
	eng.l.Lock()
	if eng.stop {
		eng.l.Unlock()
		return
	}
	eng.stop = true
	eng.l.Unlock()

	var wg sync.WaitGroup
	for _, h := range eng.shutDowns {
		wg.Add(1)
		go func(h func()) {
			defer wg.Done()
			h()
		}(h)
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-time.After(time.Second * 10):
	case <-done:
	}
	return
}

func (eng *GlobalEngine) AddShutdown(s func()) {
	eng.l.Lock()
	eng.shutDowns = append(eng.shutDowns, s)
	eng.l.Unlock()
}

func (eng *GlobalEngine) IsShutdown() bool {
	eng.l.RLock()
	defer eng.l.RUnlock()
	return eng.stop
}
