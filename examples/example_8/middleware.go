package main

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"reflect"
)

func init() {
	const semLogContext = "example-8-middleware::init"
	log.Info().Msg(semLogContext)
	mwregistry.RegisterHandlerFactory(HarTracingHandlerId, NewHarTracingHandler)
}

const (
	HarTracingHandlerId   = "gin-mw-har-tracing"
	HarTracingHandlerKind = "mw-kind-har-tracing"
)

type HarTracingHandler struct {
	config *HarTracingHandlerConfig
}

type HarTracingHandlerConfig struct {
}

var DefaultHarTracingHandlerConfig = HarTracingHandlerConfig{}

// NewErrorHandler builds an Error Handler with the following options:

func NewHarTracingHandler(cfg interface{}) (mws.MiddlewareHandler, error) {

	const semLogContext = "new-har-tracing-handler"
	tcfg := DefaultHarTracingHandlerConfig

	if cfg != nil && !reflect.ValueOf(cfg).IsNil() {
		if mapCfg, ok := cfg.(mws.MiddlewareHandlerConfig); ok {
			err := mapstructure.Decode(mapCfg, &tcfg)
			if err != nil {
				return nil, err
			}
		} else {
			log.Warn().Msg(semLogContext + " unmarshal issue for tracing handler config")
		}
	} else {
		log.Info().Str("mw-id", HarTracingHandlerId).Msg(semLogContext + " config null...reverting to default values")
	}

	log.Info().Str("mw-id", HarTracingHandlerId).Interface("cfg", tcfg).Msg(semLogContext + " handler loaded config")

	return &HarTracingHandler{config: &tcfg}, nil
}

func (t *HarTracingHandler) GetKind() string {
	return HarTracingHandlerKind
}

func (t *HarTracingHandler) HandleFunc() gin.HandlerFunc {

	const semLogContext = "har-tracing-handler"
	return func(c *gin.Context) {

		log.Trace().Str("requestPath", c.Request.RequestURI).Msg(semLogContext)

		if nil != c {
			c.Next()
		}

	}
}
