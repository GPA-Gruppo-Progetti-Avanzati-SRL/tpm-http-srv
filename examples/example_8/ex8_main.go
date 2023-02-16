package main

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	_ "GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv/resource/health"
	_ "GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv/resource/metrics"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing/logzerotracer"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwhartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwerror"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwmetrics"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwtracing"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type LogConfig struct {
	Level      int  `yaml:"level"`
	EnableJSON bool `yaml:"enablejson"`
}

type ExampleConfig struct {
	Log    LogConfig `yaml:"log"`
	Config AppConfig `yaml:"config"`
}

type AppConfig struct {
	Http       httpsrv.Config
	MwRegistry mwregistry.HandlerCatalogConfig `yaml:"mw-handler-registry" mapstructure:"mw-handler-registry"`
}

func (m *AppConfig) PostProcess() error {
	return nil
}

/*
func (m *AppConfig) GetDefaults() []configuration.VarDefinition {

	vd := make([]configuration.VarDefinition, 0, 20)
	vd = append(vd, httpsrv.GetConfigDefaults()...)
	vd = append(vd, middleware.GetConfigDefaults("config.mw-handler-registry")...)
	return vd
}
*/

//go:embed config.yml
var configFile []byte

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	t, _ := logzerotracer.NewTracer()
	hartracing.SetGlobalTracer(t)

	exampleCfg := ExampleConfig{}

	cfgFile, err := os.ReadFile("examples/example_8/config.yml")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = yaml.Unmarshal(cfgFile, &exampleCfg)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	appCfg := exampleCfg.Config
	b, _ := json.Marshal(appCfg)
	fmt.Println(string(b))

	if appCfg.MwRegistry != nil {
		if err := mwregistry.InitializeHandlerRegistry(appCfg.MwRegistry, appCfg.Http.MwUse); err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	jc, err := initGlobalTracer()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer jc.Close()

	s, err := httpsrv.NewServer(appCfg.Http, httpsrv.WithListenPort(9090), httpsrv.WithDocumentRoot("/www", "/tmp", false))
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
