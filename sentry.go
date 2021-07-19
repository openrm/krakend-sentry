package sentry

import (
	"errors"
	"strconv"

	"github.com/luraproject/lura/config"
	"github.com/getsentry/sentry-go"
)

// namespace
const Namespace = "github_com/openrm/krakend-sentry"
const Prefix = "sentry:"

// errors
var (
    ErrNoConfig = errors.New("no config for sentry")
    ErrInvalidConfig = errors.New("invalid config for sentry")
)

// initializer
func New(cfg config.ServiceConfig) error {
	return NewWithClientOptions(cfg, sentry.ClientOptions{})
}

func NewWithClientOptions(cfg config.ServiceConfig, opts sentry.ClientOptions) error {
	v, ok := ConfigGetter(cfg.ExtraConfig).(map[string]interface{})
	if !ok {
		return ErrNoConfig
	}
	parseConfig(v, opts)
	opts.Release = strconv.Itoa(cfg.Version)
	return sentry.Init(opts)
}

func ConfigGetter(cfg config.ExtraConfig) interface{} {
	v, ok := cfg[Namespace]
	if !ok {
		return nil
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return tmp
}

func parseConfig(cfg map[string]interface{}, opts sentry.ClientOptions) sentry.ClientOptions {
	getString(cfg, "dsn", &opts.Dsn)
	getString(cfg, "environment", &opts.Dsn)
	getStrings(cfg, "ignore_errors", &opts.IgnoreErrors)
	getBool(cfg, "debug", &opts.Debug)
	getBool(cfg, "attach_stacktrace", &opts.AttachStacktrace)
	getFloat64(cfg, "sample_rate", &opts.SampleRate)
	return opts
}

func getStrings(data map[string]interface{}, key string, v *[]string) {
	if vs, ok := data[key]; ok {
		result := []string{}
		for _, v := range vs.([]interface{}) {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		*v = result
	}
}

func getString(data map[string]interface{}, key string, v *string) {
	if val, ok := data[key]; ok {
		if s, ok := val.(string); ok && len(s) > 0 {
			*v = s
		}
	}
}

func getBool(data map[string]interface{}, key string, v *bool) {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			*v = b
		}
	}
}

func getFloat64(data map[string]interface{}, key string, v *float64) {
	if val, ok := data[key]; ok {
		switch i := val.(type) {
		case float64:
			*v = i
		case int:
			*v = float64(i)
		case int64:
			*v = float64(i)
		}
	}
}
