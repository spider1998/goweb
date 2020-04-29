package engine

import (
	"fmt"
	"sync"
	"time"

	"goweb/pkg/components"
	"goweb/pkg/register"

	"goweb/pkg/code"
)

type Handler func(*Job) (code code.Code, err error)

type GlobalEngine struct {
	handlers     map[string]Handler
	emptyHandler Handler
	l            sync.RWMutex
	shutdowns    []func()
	stop         bool
}

func NewEngine() *GlobalEngine {
	eng := &GlobalEngine{
		handlers: make(map[string]Handler),
	}
	return eng
}

func (eng *GlobalEngine) Initial() *GlobalEngine {
	if err := eng.initial(register.Register, components.Cron, components.Log); err != nil {
		panic(err)
	}
	eng.onShutdown(components.Cron)

	return eng
}

func (eng *GlobalEngine) initial(services ...InitialService) error {
	for _, s := range services {
		v, ok := s.(InitialService)
		if ok {
			go v.Run(components.Log)
		}
	}

	return nil
}

func (eng *GlobalEngine) onShutdown(services ...OnShutdownService) {
	for _, s := range services {
		v, ok := s.(OnShutdownService)
		if ok {
			eng.AddShutdown(v.OnShutdown())
		}
	}
	return
}

func (eng *GlobalEngine) Run() {
	fmt.Println("start engine")
	components.Log.Infof("test info", "test engine")
	time.Sleep(time.Second * 120)
	eng.Shutdown()
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
	for _, h := range eng.shutdowns {
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
	eng.shutdowns = append(eng.shutdowns, s)
	eng.l.Unlock()
}

func (eng *GlobalEngine) IsShutdown() bool {
	eng.l.RLock()
	defer eng.l.RUnlock()
	return eng.stop
}
