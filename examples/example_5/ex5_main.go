package main

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	mwerror2 "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwerror"
	mwtracing2 "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwtracing"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	jc, err := initGlobalTracer()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer jc.Close()

	s, err := httpsrv.NewServer(httpsrv.DefaultConfig,
		httpsrv.WithBindAddress("localhost"),
		httpsrv.WithListenPort(8080),
		httpsrv.WithShutdownTimeout(time.Duration(5)*time.Second),
		httpsrv.WithContextPath("/api"),
		httpsrv.WithMiddlewareHandlers(
			mwtracing2.MustNewTracingHandler(mwtracing2.DefaultTracingHandlerConfig).HandleFunc(),
			mwerror2.MustNewErrorHandler(mwerror2.DefaultErrorHandlerConfig).HandleFunc()))

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := s.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}
	defer s.Stop()

	for !s.IsReady() {
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

	sig := <-shutdownChannel
	log.Debug().Interface("signal", sig).Msg("got termination signal")
}

func initGlobalTracer() (io.Closer, error) {
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
