package router

import (
	"fmt"
	"goweb/pkg/components"
	"goweb/pkg/router/admin"
	"goweb/pkg/router/middleware"
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
)

func Register(logger components.Logger) http.Handler {
	router := routing.New()
	router.NotFound(middleware.NotFound)
	router.Use(
		systemHeader,
		middleware.RoutingLogger(logger),
		middleware.ErrorHandler(logger),
		content.TypeNegotiator(content.JSON),
		/*	cors.Handler(cors.Options{
			AllowOrigins:  "*",
			AllowHeaders:  "*",
			AllowMethods:  "*",
			ExposeHeaders: "X-Page-Total, X-Page, X-Page-Size,X-Access-Token",
		}),*/
	)

	api := router.Group("/" + components.System)
	{
		api.Get("/version", func(c *routing.Context) error {
			return c.Write(map[string]string{
				"version":    components.Version,
				"build_time": components.BuildTime,
			})
		})
	}

	admin.Register(api.Group("/admin/v1"))

	for _, route := range router.Routes() {
		logger.Debugf(fmt.Sprintf("register route: \"%-6s -> %s\".", route.Method(), route.Path()))
	}

	return router
}

func systemHeader(c *routing.Context) error {
	c.Response.Header()["X-SD-Module"] = []string{components.System}
	return c.Next()
}
