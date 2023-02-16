package swagger

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func init() {
	log.Info().Msg("swagger resources init function invoked")
	ra := httpsrv.GetApp()
	ra.RegisterGFactory(registerSwaggerhEndpoint)
}

func registerSwaggerhEndpoint(ctx httpsrv.ServerContext) []httpsrv.G {

	gs := make([]httpsrv.G, 0, 2)

	gs = append(gs, httpsrv.G{
		Name:    "Swagger endpoint",
		AbsPath: true,
		Path:    "swagger",
		Resources: []httpsrv.R{
			{
				Name:          "swagger",
				Path:          "*any",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{ginSwagger.WrapHandler(swaggerFiles.Handler)},
			},
		},
	})

	return gs
}
