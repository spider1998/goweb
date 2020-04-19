package register

import (
	"goweb/pkg/handler"
	"goweb/pkg/log"
)

var Register = new(RegisterService)

type RegisterService struct{}

func (r *RegisterService) Run(logger ...log.Logger) {
	handler.Register("ping", ping)
	handler.Register("version", version)
}
