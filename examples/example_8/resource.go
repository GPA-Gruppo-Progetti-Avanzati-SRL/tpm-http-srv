package main

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

/*
 * This example is identical to the example_3. The difference is the way of registration that is done 'postponed' via a factory method invoked
 * when the server gets started....
 */
func init() {
	log.Info().Msg("example_7 init function invoked")
	ra := httpsrv.GetApp()
	ra.RegisterGFactory(registerGroups)
}

func registerGroups(_ httpsrv.ServerContext) []httpsrv.G {

	gs := make([]httpsrv.G, 0, 2)

	gs = append(gs, httpsrv.G{
		Name:        "HelloWorldEn",
		Path:        ":site/:lang",
		UseSysMw:    true,
		Middlewares: []httpsrv.H{setLangHeader("uk")},
		Resources: []httpsrv.R{
			{
				Name:          "home",
				Path:          "",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{example()},
			},
			{
				Name:          "proxy-to-app-home",
				Path:          ":appName",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{example()},
			},
			{
				Name:          "proxy-to-app",
				Path:          ":appName/*proxyPath",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{example()},
			},
		},
	})

	return gs
}

func example() httpsrv.H {
	return func(c *gin.Context) {
		site := c.Param("site")
		lang := c.Param("lang")
		appName := c.Param("appName")
		log.Info().Str("site", site).Str("lang", lang).Str("appName", appName).Str("target-path", c.Param("proxyPath")).Msg("route found")

		remote, err := url.Parse("http://localhost:3000")
		if err != nil {
			panic(err)
		}

		if lang == "it" {
			// c.Error(middleware.NewAppError(middleware.AppErrorWithStatusCode(350), middleware.AppErrorWithText("error text")))
			c.AbortWithStatusJSON(350, middleware.NewAppError(middleware.AppErrorWithStatusCode(350), middleware.AppErrorWithText("error text")))
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = "/r3ds9-auth/user"
		}

		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

func setLangHeader(lang string) httpsrv.H {
	return func(c *gin.Context) {
		site := c.Param("site")
		lang := c.Param("lang")
		appName := c.Param("appName")
		log.Info().Str("site", site).Str("lang", lang).Str("appName", appName).Str("target-path", c.Param("proxyPath")).Msg("middleware")

		c.Header("X-lang", lang)
		c.Next()
	}
}
