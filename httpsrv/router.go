package httpsrv

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwzerologger"
	ginzgzip "github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type H = gin.HandlerFunc

type GRegistry []G
type GFactory func(ctx ServerContext) []G

type App interface {
	RegisterG(rr ...G)
	RegisterGFactory(rr ...GFactory)
}

func GetApp() App {
	return application
}

var application = &app{
	gRegistry:              make(GRegistry, 0),
	gRegistrationFactories: make([]GFactory, 0),
}

type app struct {
	gRegistry              GRegistry
	gRegistrationFactories []GFactory
}

func (ra *app) RegisterG(rr ...G) {
	ra.gRegistry = append(ra.gRegistry, rr...)
}

func (ra *app) RegisterGFactory(rr ...GFactory) {
	ra.gRegistrationFactories = append(ra.gRegistrationFactories, rr...)
}

const (
	MethodAny = "ANY" // Used when regiistering resources for all the methods.
)

type G struct {
	Name        string
	UseSysMw    bool
	AbsPath     bool
	Path        string
	Middlewares []H
	Resources   []R
}

type R struct {
	Name          string
	Method        string
	Path          string
	Middlewares   []H
	RouteHandlers []H
}

func newRouter(serverContext ServerContext, mws []H, pathsNotToLog []string, gzipCfg *GzipConfig) *gin.Engine {

	const semLogContext = "http-srv::new-router"
	r := gin.New()
	r.Use(mwzerologger.ZeroLogger("gin", pathsNotToLog...))
	r.Use(gin.Recovery())

	// Gzip globale: si applica a tutte le route SPA e ai file statici.
	// Va registrato subito dopo Recovery, prima di qualsiasi route o r.Use(static.Serve(...)),
	// in modo che il ResponseWriter compresso sia attivo per tutta la catena successiva.
	// Per escludere path specifici usare excluded-paths (match esatto) o
	// excluded-paths-regexs (regexp) in config.yml.
	if gzipCfg != nil && gzipCfg.Enabled {
		level := gzipCfg.Level
		if level == 0 {
			level = ginzgzip.DefaultCompression
		}
		var gzipOpts []ginzgzip.Option
		if len(gzipCfg.ExcludedPaths) > 0 {
			gzipOpts = append(gzipOpts, ginzgzip.WithExcludedPaths(gzipCfg.ExcludedPaths))
		}
		if len(gzipCfg.ExcludedPathsRegexs) > 0 {
			gzipOpts = append(gzipOpts, ginzgzip.WithExcludedPathsRegexs(gzipCfg.ExcludedPathsRegexs))
		}
		r.Use(ginzgzip.Gzip(level, gzipOpts...))
		log.Info().Int("level", level).
			Strs("excluded-paths", gzipCfg.ExcludedPaths).
			Strs("excluded-paths-regexs", gzipCfg.ExcludedPathsRegexs).
			Msg(semLogContext + " - gzip compression enabled")
	}

	withPprof := false
	if prf, ok := serverContext.GetConfig("with-pprof"); ok {
		if b, ok := prf.(bool); ok && b {
			withPprof = true
			log.Info().Msg(semLogContext + " pprof handler enabled")
			pprof.Register(r, pprof.DefaultPrefix)
		}
	}

	if !withPprof {
		log.Info().Msg(semLogContext + " pprof handler NOT enabled (if needed set the with-pprof param to true in server-context.context")
	}

	for _, gdef := range application.gRegistry {
		if gdef.UseSysMw {
			gdef.Middlewares = append(gdef.Middlewares, mws...)
		}
		addGroup2Engine(r, serverContext.GetContextPath(), gdef)
	}

	for _, gfact := range application.gRegistrationFactories {
		gs := gfact(serverContext)
		for _, gdef := range gs {
			if gdef.UseSysMw {
				gdef.Middlewares = append(gdef.Middlewares, mws...)
			}
			addGroup2Engine(r, serverContext.GetContextPath(), gdef)
		}
	}

	return r
}

func addGroup2Engine(eng *gin.Engine, contextPath string, gdef G) {

	// The path provided is prefixed with context path unless declared abssolute in configuration.
	gpath := gdef.Path
	if !gdef.AbsPath {
		gpath = fmt.Sprintf("%s/%s", contextPath, gdef.Path)
	}

	g := eng.Group(gpath, gdef.Middlewares...)
	log.Trace().Str("path", gpath).Msg("registering group")
	for _, r := range gdef.Resources {

		log.Trace().Str("method", r.Method).Str("path", r.Path).Msg("registering resource")

		switch r.Method {
		/* Confusing.... is a group middleware configured at the level of single route.

		case "USE":
			g.Use(r.RouteHandlers...)

		*/
		case http.MethodGet:
			g.GET(r.Path, r.RouteHandlers...)
		case http.MethodPost:
			g.POST(r.Path, r.RouteHandlers...)
		case http.MethodPut:
			g.PUT(r.Path, r.RouteHandlers...)
		case http.MethodDelete:
			g.DELETE(r.Path, r.RouteHandlers...)
		case http.MethodPatch:
			g.PATCH(r.Path, r.RouteHandlers...)
		case "ANY":
			g.Any(r.Path, r.RouteHandlers...)
		}
	}
}
