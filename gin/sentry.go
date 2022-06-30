package gin

import (
	"fmt"
	"net/http"

	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	ginlura "github.com/luraproject/lura/v2/router/gin"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	krakendsentry "github.com/openrm/krakend-sentry/v2"
)

var nopHandler = func(c *gin.Context) { c.Next() }

func Register(cfg config.ServiceConfig, logger logging.Logger, engine *gin.Engine) {
	RegisterWithClientOptions(cfg, logger, sentry.ClientOptions{}, engine)
}

func RegisterWithClientOptions(cfg config.ServiceConfig, logger logging.Logger, opts sentry.ClientOptions, engine *gin.Engine) {
	if err := krakendsentry.NewWithClientOptions(cfg, opts); err != nil {
		logger.Debug(krakendsentry.Prefix, "middleware disabled")
		return
	}
	engine.Use(sentrygin.New(sentrygin.Options{ Repanic: true }))
	logger.Debug(krakendsentry.Prefix, "middleware enabled")
}

func HandlerFactory(logger logging.Logger, hf ginlura.HandlerFactory) ginlura.HandlerFactory {
	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handler := hf(cfg, p)
		logger.Debug(krakendsentry.Prefix, "enabled for the endpoint", cfg.Endpoint)
		return func(c *gin.Context) {
			hub := sentrygin.GetHubFromContext(c)
			if hub == nil {
				handler(c)
				return
			}
			hub.Scope().SetTag("transaction", fmt.Sprintf("%s %s", cfg.Method, cfg.Endpoint))
			handler(c)
			if len(c.Errors) > 0 {
				for _, err := range c.Errors {
					hub.CaptureException(err)
				}
			} else if status := c.Writer.Status(); status >= http.StatusInternalServerError {
				// TODO report status >= 500
			}
		}
	}
}
