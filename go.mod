module github.com/openrm/krakend-sentry/v2

go 1.17

require (
	github.com/getsentry/sentry-go v0.11.0
	github.com/gin-gonic/gin v1.7.7
	github.com/luraproject/lura/v2 v2.0.5
)

retract [v2.0.0, v2.0.1]
