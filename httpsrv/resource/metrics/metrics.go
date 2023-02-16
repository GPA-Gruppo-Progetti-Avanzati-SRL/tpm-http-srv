package metrics

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

const (
	MiddlewarePromHttpDefaultEndpoint = "/metrics"
)

func init() {
	log.Info().Msg("metrics resources init function invoked")
	ra := httpsrv.GetApp()
	ra.RegisterGFactory(registerMetricsEndpoints)
}

func registerMetricsEndpoints(ctx httpsrv.ServerContext) []httpsrv.G {

	gs := make([]httpsrv.G, 0, 2)

	p := MiddlewarePromHttpDefaultEndpoint
	if prop, ok := ctx.GetConfig(httpsrv.ServerContextMetricsEndpointProperty); ok {
		p = prop.(string)
	}

	port := 8080
	if prop, ok := ctx.GetConfig(httpsrv.SRVCTX_SERVER_PORT); ok {
		port = prop.(int)
	}

	connStr := fmt.Sprintf("http://127.0.0.1:%d/%s", port, strings.TrimPrefix(p, "/"))

	gs = append(gs, httpsrv.G{
		Name:    "Metrics endpoint",
		AbsPath: true,
		Path:    p,
		Resources: []httpsrv.R{
			{
				Name:          "export-metrics",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{exportMetricsHandleFunc()},
			},

			{
				Name:          "export-metrics-json",
				Method:        http.MethodGet,
				Path:          "json",
				RouteHandlers: []httpsrv.H{JsonMetrics(connStr)},
			},

			{
				Name:          "export-metrics-json",
				Method:        http.MethodGet,
				Path:          "json/metric/:metricName",
				RouteHandlers: []httpsrv.H{JsonMetricsByName(connStr)},
			},
		},
	})

	return gs
}

func exportMetricsHandleFunc() gin.HandlerFunc {

	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

/*
returns observability in a json format
*/
func JsonMetrics(connStr string) httpsrv.H {

	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		mfChan := make(chan *dto.MetricFamily, 1024)

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		go func() {
			err := prom2json.FetchMetricFamilies(connStr, mfChan, transport)
			if err != nil {
				// TODO: should be fatal?
				log.Fatal().Err(err).Send()
			}
		}()

		result := []*prom2json.Family{}
		for mf := range mfChan {
			result = append(result, prom2json.NewFamily(mf))
		}

		c.JSON(http.StatusOK, result)
	}
}

/*
returns a metric in a json format
with the specified name
*/
func JsonMetricsByName(connStr string) httpsrv.H {

	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		if metricName, metricNameNotNil := c.Params.Get("metricName"); metricNameNotNil {

			log.Info().Msgf("MetricHandler::JsonMetricByName retrieving metric having name: %s", metricName)

			mfChan := make(chan *dto.MetricFamily, 1024)

			transport := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			go func() {
				// connStr := fmt.Sprintf("http://127.0.0.1:%d/observability", config.ServerPort())
				err := prom2json.FetchMetricFamilies(connStr, mfChan, transport)
				if err != nil {
					// TODO: should be fatal?
					log.Fatal().Err(err).Send()
				}
			}()

			result := []*prom2json.Family{}
			for mf := range mfChan {
				if metric := prom2json.NewFamily(mf); metric.Name == metricName {
					result = append(result, prom2json.NewFamily(mf))
				}
			}

			c.JSON(http.StatusOK, result)
		} else {

			c.JSON(http.StatusNotFound, gin.H{})
		}
	}
}
