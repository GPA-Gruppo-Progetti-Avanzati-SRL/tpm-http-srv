package example_2

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func init() {
	log.Info().Msg("example_2 init function invoked")
	ra := httpsrv.GetApp()
	ra.RegisterG(registerGroup())
}

func registerGroup() httpsrv.G {
	g := httpsrv.G{
		Name:        "HelloWorld",
		Path:        "v1/test",
		Middlewares: []httpsrv.H{setLangHeader},
		Resources: []httpsrv.R{
			{
				Name:          "sayhello",
				Path:          "sayhello/:name",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{hello()},
			},
		},
	}
	return g
}

func hello() httpsrv.H {
	return func(c *gin.Context) {
		name := c.Param("name")
		c.String(200 /* httpsrv.StatusOK */, fmt.Sprintf("Hello %s", name))
	}
}

func setLangHeader(c *gin.Context) {
	c.Header("X-lang", "uk")
}
