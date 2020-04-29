package engine

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"goweb/pkg/components"
	"goweb/pkg/register"
	"goweb/pkg/router"

	"goweb/pkg/code"
)

type Handler func(*Job) (code code.Code, err error)

type GlobalEngine struct {
	httpServer   *http.Server
	wg           sync.WaitGroup
	handlers     map[string]Handler
	emptyHandler Handler
	l            sync.RWMutex
	shutdowns    []func()
	stop         bool
	close        chan struct{}
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

func (eng *GlobalEngine) serveHTTP() {
	defer eng.wg.Done()

	eng.httpServer.Handler = router.Register(components.Log)

	components.Log.Infof("listen and serve http service.", "addr", components.Conf.Addr)

	err := eng.httpServer.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			fmt.Println(err)
			fmt.Println("an error was returned while listen and serve engine.")
			return
		}
	}
}

func (eng *GlobalEngine) registerSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case sig := <-ch:
		signal.Ignore(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
		fmt.Println("received signal, try to shutdown engine.", sig)
		close(ch)
		close(eng.close)
		eng.Shutdown()

	}
}

func (eng *GlobalEngine) Run() {
	go eng.registerSignal()

	eng.wg.Add(1)
	go eng.serveHTTP()
	eng.wg.Wait()
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
