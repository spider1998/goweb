package admin

import (
	routing "github.com/go-ozzo/ozzo-routing"
	"goweb/pkg/router/middleware"
)

func Register(router *routing.RouteGroup) {
	router.Use(
		middleware.SessionChecker,
	)
	{
		handler := NewVersionHandler()
		router.Get("/version", handler.getVersion)
	}
}
