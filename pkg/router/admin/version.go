package admin

import (
	routing "github.com/go-ozzo/ozzo-routing"
	"goweb/pkg/components"
)

type VersionHandler struct {
}

func NewVersionHandler() VersionHandler {
	return VersionHandler{}
}

func (h VersionHandler) getVersion(c *routing.Context) error {
	return c.Write(components.Version)
}
