package httpsrv

import (
	"embed"
	"github.com/gin-gonic/gin"
	"time"
)

type Config struct {
	BindAddress     string           `json:"bind-address" yaml:"bind-address" mapstructure:"bind-address"`
	ListenPort      int              `json:"port" yaml:"port" mapstructure:"port"`
	ShutdownTimeout time.Duration    `json:"shutdown-timeout" yaml:"shutdown-timeout" mapstructure:"shutdown-timeout"`
	PathsNotToLog   []string         `json:"paths-not-to-log" yaml:"paths-not-to-log" mapstructure:"paths-not-to-log"`
	ServerMode      string           `json:"server-mode" yaml:"server-mode" mapstructure:"server-mode"`
	ServerCtx       ServerContextCfg `json:"server-context" yaml:"server-context" mapstructure:"server-context"`

	Statics     []StaticContent `json:"static-content" yaml:"static-content" mapstructure:"static-content"`
	HtmlContent string          `json:"html-content" yaml:"html-content" mapstructure:"html-content"`

	mwHandlers []H
	MwUse      []string `json:"mw-use" yaml:"mw-use" mapstructure:"mw-use"`
}

const (
	ServerContextMetricsEndpointProperty = "sys-metrics-endpoint"
)

type ServerContextCfg struct {
	Path          string                 `json:"path" yaml:"path" mapstructure:"path"`
	ContextParams map[string]interface{} `json:"context" yaml:"context" mapstructure:"context"`
}

type StaticContent struct {
	UrlPrefix       string `json:"url-prefix" yaml:"url-prefix" mapstructure:"url-prefix"`
	DocumentRoot    string `json:"document-root" yaml:"document-root" mapstructure:"document-root"`
	Indexes         bool
	EmbedFileSystem embed.FS
}

const (
	DefaultBindAddress     = "localhost"
	DefaultListenPort      = 8080
	DefaultShutdownTimeout = 500 * time.Millisecond
	DefaultServerMode      = gin.DebugMode
	DefaultContextPath     = "/api"
)

var DefaultConfig = Config{
	BindAddress:     DefaultBindAddress,
	ListenPort:      DefaultListenPort,
	ShutdownTimeout: DefaultShutdownTimeout,
	ServerCtx: ServerContextCfg{
		Path: DefaultContextPath,
	},
	ServerMode: DefaultServerMode,
}

/*
func GetConfigDefaults() []configuration.VarDefinition {
	return []configuration.VarDefinition{
		{"config.http.bind-address", DefaultBindAddress, "host reference"},
		{"config.http.server-context.path", DefaultContextPath, "context-path"},
		{"config.http.port", DefaultListenPort, "port"},
		{"config.http.shutdown-timeout", DefaultShutdownTimeout, "shutdown timeout"},
		{"config.http.server-mode", DefaultServerMode, "modalita' di lavoro server gin"},
	}
}
*/

// ConfigBuilder
//   WithBindAddress(string)                   bind address for this httpsrv
//   WithListenPort(uint16)                    listen port for this httpsrv
//   WithMiddlewareHandlers([]gin.HandlerFunc) array of middlewares for this httpsrv
//   WithShutdownTimeout(time.Duration)        shutdown  Timeout

type CfgOption func(*Config)

func WithBindAddress(ba string) CfgOption {
	return func(c *Config) {
		c.BindAddress = ba
	}
}

func WithListenPort(p int) CfgOption {
	return func(c *Config) {
		c.ListenPort = p
	}

}

func WithMiddlewareHandlers(mws ...H) CfgOption {
	return func(c *Config) {
		c.mwHandlers = append(c.mwHandlers, mws...)
	}

}

func WithShutdownTimeout(to time.Duration) CfgOption {
	return func(c *Config) {
		c.ShutdownTimeout = to
	}

}

func WithContextPath(cp string) CfgOption {
	return func(c *Config) {
		c.ServerCtx.Path = cp
	}
}

func WithServerMode(ginMode string) CfgOption {
	return func(c *Config) {
		c.ServerMode = ginMode
	}
}

func WithDocumentRoot(basePath string, aPath string, indexes bool) CfgOption {
	return func(c *Config) {
		c.Statics = append(c.Statics, StaticContent{UrlPrefix: basePath, DocumentRoot: aPath, Indexes: indexes})
	}
}
