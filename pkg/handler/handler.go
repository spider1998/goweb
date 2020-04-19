package handler

import (
	"fmt"
	"goweb/pkg/code"
)

type GlobalHandler func() (code code.Code, err error)

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
