package engine

import "goweb/pkg/log"

type InitialService interface {
	Run(...log.Logger)
}

type OnShutdownService interface {
	OnShutdown() func()
}
