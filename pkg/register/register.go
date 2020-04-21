package register

import (
	"goweb/pkg/components"
	"goweb/pkg/handler"
)

var Register = new(RegisterService)

type RegisterService struct{}

func (r *RegisterService) Run(logger ...components.Logger) {
	handler.Register("ping", ping)
	handler.Register("version", version)
}
