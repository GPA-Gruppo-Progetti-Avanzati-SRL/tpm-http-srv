package httpsrv

import (
	"errors"
	"github.com/rs/zerolog/log"
)

const (
	SRVCTX_SERVER_PORT = "http-server-port"
)

type ServerContext interface {
	GetContextPath() string
	GetConfig(n string) (interface{}, bool)
	Add(k string, v interface{}) error
	Get(n string) (interface{}, bool)
}

type serverContext struct {
	ServerContextCfg
	Context map[string]interface{}
}

func NewServerContext(cfg Config) ServerContext {

	sctx := &serverContext{ServerContextCfg: cfg.ServerCtx, Context: make(map[string]interface{})}

	if sctx.ServerContextCfg.ContextParams == nil {
		sctx.ServerContextCfg.ContextParams = make(map[string]interface{})
	}

	// Augmenting server context with config related stuff required by some components.
	sctx.ServerContextCfg.ContextParams[SRVCTX_SERVER_PORT] = cfg.ListenPort

	return sctx
}

func (sctx *serverContext) GetContextPath() string {
	return sctx.Path
}

func (sctx *serverContext) GetConfig(k string) (interface{}, bool) {

	if i, ok := sctx.ContextParams[k]; ok {
		return i, true
	}

	return nil, false
}

func (sctx *serverContext) Add(k string, v interface{}) error {

	if _, ok := sctx.Context[k]; ok {
		log.Error().Str("key", k).Msg("server context already has value for key")
		return errors.New("server context already has value for key " + k)
	}

	sctx.Context[k] = v
	return nil
}

func (sctx *serverContext) Get(k string) (interface{}, bool) {

	if i, ok := sctx.Context[k]; ok {
		return i, true
	}

	return nil, false
}
