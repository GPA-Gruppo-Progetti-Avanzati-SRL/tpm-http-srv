package example_3

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func init() {
	log.Info().Msg("example_3 init function invoked")
	ra := httpsrv.GetApp()
	ra.RegisterG(registerGroups()...)
}

func registerGroups() []httpsrv.G {

	gs := make([]httpsrv.G, 0, 2)

	gs = append(gs, httpsrv.G{
		Name:        "HelloWorldEn",
		Path:        "v1/en/test",
		Middlewares: []httpsrv.H{setLangHeader("uk")},
		Resources: []httpsrv.R{
			{
				Name:          "sayhello",
				Path:          "sayhello/:name",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{hello()},
			},
		},
	})

	gs = append(gs, httpsrv.G{
		Name:        "HelloWorldFr",
		Path:        "v1/fr/test",
		Middlewares: []httpsrv.H{setLangHeader("fr")},
		Resources: []httpsrv.R{
			{
				Name:          "sayhello",
				Path:          "sayhello/:name",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{helloFr()},
			},
		},
	})

	return gs
}

func hello() httpsrv.H {
	return func(c *gin.Context) {
		name := c.Param("name")
		c.String(200 /* httpsrv.StatusOK */, fmt.Sprintf("Hello %s", name))
	}
}

func helloFr() httpsrv.H {
	return func(c *gin.Context) {
		name := c.Param("name")
		c.String(200 /* httpsrv.StatusOK */, fmt.Sprintf("Bonjour %s", name))
	}
}

func setLangHeader(lang string) httpsrv.H {
	return func(c *gin.Context) {
		c.Header("X-lang", lang)
	}
}
