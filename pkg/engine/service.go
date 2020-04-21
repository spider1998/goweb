package engine

import (
	"goweb/pkg/components"
)

type InitialService interface {
	Run(...components.Logger)
}

type OnShutdownService interface {
	OnShutdown() func()
}
