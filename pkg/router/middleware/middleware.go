package middleware

import (
	"bytes"
	"fmt"
	"git.sdkeji.top/share/sqmslib/api"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"goweb/pkg/code"
	"goweb/pkg/components"
	"net/http"
	"os"
	"runtime/debug"
	"sort"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/pkg/errors"
)

func RoutingLogger(logger components.Logger) routing.Handler {
	return access.CustomLogger(func(req *http.Request, rw *access.LogResponseWriter, elapsed float64) {
		clientIP := access.GetClientIP(req)
		logger.Infof("access log.",
			"ip", clientIP,
			"proto", req.Proto,
			"method", req.Method,
			"url", req.URL.String(),
			"status", rw.Status,
			"size", rw.BytesWritten,
			"duration", elapsed,
		)
	})
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ErrorHandler(logger components.Logger) routing.Handler {
	return func(c *routing.Context) (err error) {
		defer func() {
			if e := recover(); e != nil {
				logger.Errorf("recovered from panic.", "v", e)
				fmt.Print(string(debug.Stack()))
				sendError(logger, c, code.CodeError(code.InternalServerError), http.StatusInternalServerError)
				c.Abort()
				err = nil
			}
		}()
		err = c.Next()
		if err != nil {
			c.Abort()
			if err, ok := err.(stackTracer); ok {
				buf := new(bytes.Buffer)
				buf.WriteString(fmt.Sprintf("error with stacktrace returned: %v\n", err))
				for _, f := range err.StackTrace() {
					buf.WriteString(fmt.Sprintf("%+v\n", f))
				}
				fmt.Fprint(os.Stderr, buf.String())
			}

			switch e := errors.Cause(err).(type) {
			case validation.Errors:
				type validationError struct {
					Field string `json:"field"`
					Error string `json:"error"`
				}
				result := make([]validationError, 0)
				fields := make([]string, 0)
				for field := range e {
					fields = append(fields, field)
				}
				sort.Strings(fields)
				for _, field := range fields {
					err := e[field]
					result = append(result, validationError{
						Field: field,
						Error: err.Error(),
					})
				}
				apiErr := code.CodeError(code.InvalidData).WithDetails(result)
				sendError(logger, c, apiErr, apiErr.StatusCode())
			case code.APIError:
				components.Log.Debugf("api error.", "response", e)
				sendError(logger, c, e, e.StatusCode())
			case routing.HTTPError:
				sendError(logger, c, code.CodeError(code.InvalidData).WithMessage(e.Error()), e.StatusCode())
			default:
				components.Log.Errorf("unknown error.", "error", err)
				sendError(logger, c, code.CodeError(code.InternalServerError), http.StatusInternalServerError)
				return nil
			}
		}
		return nil
	}
}

func NotFound(c *routing.Context) error {
	methods := c.Router().FindAllowedMethods(c.Request.URL.Path)
	if len(methods) == 0 {
		return code.CodeError(code.NotFound)
	}
	methods["OPTIONS"] = true
	ms := make([]string, len(methods))
	i := 0
	for method := range methods {
		ms[i] = method
		i++
	}
	sort.Strings(ms)
	c.Response.Header().Set("Allow", strings.Join(ms, ", "))
	if c.Request.Method != "OPTIONS" {
		return code.CodeError(code.MethodNotAllowed)
	}
	c.Abort()
	return nil
}

func sendError(logger components.Logger, c *routing.Context, err error, status int) {
	c.Response.WriteHeader(status)
	c.Response.Header().Set("X-Content-Type-Options", "nosniff")
	err = c.Write(err)
	if err != nil {
		logger.Errorf("fail to write error.", "error", err)
	}
}

const (
	sessionTokenHeaderKey = "X-Access-Token"
)

func SessionChecker(c *routing.Context) error {
	token := c.Request.Header.Get(sessionTokenHeaderKey)
	if token == "" {
		token = c.Query(strings.ToLower(sessionTokenHeaderKey))
		if token == "" {
			return code.CodeError(code.InvalidToken).WithStatus(http.StatusForbidden)
		}
	}
}
