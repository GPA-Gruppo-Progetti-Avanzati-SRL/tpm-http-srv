package health

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func init() {
	const semLogContext = "health-resource::init"
	log.Info().Msg(semLogContext)
	ra := httpsrv.GetApp()
	ra.RegisterGFactory(registerHealthEndpoints)
}

func registerHealthEndpoints(ctx httpsrv.ServerContext) []httpsrv.G {

	gs := make([]httpsrv.G, 0, 2)

	gs = append(gs, httpsrv.G{
		Name:    "Liveness endpoint",
		AbsPath: true,
		Path:    "health",
		Resources: []httpsrv.R{
			{
				Name:          "liveness",
				Path:          "liveness",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{func(c *gin.Context) { c.JSON(200, "OK") }},
			},
			{
				Name:          "readiness",
				Path:          "readiness",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{func(c *gin.Context) { c.JSON(200, "OK") }},
			},
		},
	})

	return gs
}
