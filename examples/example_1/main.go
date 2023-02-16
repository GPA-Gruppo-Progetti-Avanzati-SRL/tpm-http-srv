package main

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/middleware"
	"github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

func main() {

	c, err := initTracer()
	if err != nil {
		panic(err)
	}

	defer c.Close()

	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("hello logger")

	r := gin.New()
	r.Use(ginzerolog.Logger("gin"),
		middleware.MustNewTracingHandler(middleware.DefaultTracingHandlerConfig).HandleFunc(),
		middleware.MustNewErrorHandler(middleware.DefaultErrorHandlerConfig).HandleFunc())

	// if the prefixes match it takes the first... apparently....
	r.Use(static.Serve("/static2", static.LocalFile("/Users/marioa.imperato/projects/tpm/http/tpm-http-srv/examples/example_1", true)))
	r.Use(static.Serve("/static2", static.LocalFile("/Users/marioa.imperato/projects/tpm/http/tpm-http-srv/examples/example_2", true)))

	r.GET("/ping", func(c *gin.Context) {
		// c.AbortWithError(500, errors.New("Ciao"))
		c.Error(middleware.NewAppError())
		/*
			c.JSON(200, gin.H{
				"message": "pong",
			})
		*/
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func initTracer() (io.Closer, error) {
	var tracer opentracing.Tracer
	var closer io.Closer

	jcfg := jaegercfg.Configuration{
		ServiceName: "gintest",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: 1.0,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := jcfg.NewTracer(
		jaegercfg.Logger(&jlogger{}),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if nil != err {
		log.Error().Err(err).Msg("Error in NewTracer")
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}

type jlogger struct{}

func (l *jlogger) Error(msg string) {
	log.Error().Msg("(jaeger) " + msg)
}

func (l *jlogger) Infof(msg string, args ...interface{}) {
	log.Info().Msgf("(jaeger) "+msg, args...)
}
