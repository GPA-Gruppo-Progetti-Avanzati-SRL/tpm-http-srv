package httpsrv

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv/embedstatic"
	"context"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/middleware"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Server interface {
	Start() error
	Stop()
	IsReady() bool
	Add2ServerContext(k string, v interface{}) error
}

type serverImpl struct {
	cfg *Config

	engine  *gin.Engine
	isReady bool
	srv     *http.Server

	serverCtx ServerContext
}

// NewServer returns a serverImpl with the passed configuration
func NewServer(cfg Config, opts ...CfgOption) (Server, error) {

	for _, opt := range opts {
		opt(&cfg)
	}

	sctx := NewServerContext(cfg)

	return &serverImpl{
		isReady:   false,
		cfg:       &cfg,
		serverCtx: sctx,
	}, nil
}

// Start starts the httpsrv
func (s *serverImpl) Start() error {

	gin.SetMode(s.cfg.ServerMode)
	log.Info().Str("server-mode", s.cfg.ServerMode).Msg("http server instantiation...")

	mhs := make([]H, 0, len(s.cfg.MwUse)+len(s.cfg.mwHandlers))
	for _, s := range s.cfg.MwUse {
		f := middleware.GetHandlerFunc(s)
		if f != nil {
			mhs = append(mhs, f)
		} else {
			log.Warn().Str("mw-name", s).Msg("the requested handler cannot be found in handler registry")
		}
	}

	mhs = append(mhs, s.cfg.mwHandlers...)
	r := newRouter(s.serverCtx, mhs, s.cfg.PathsNotToLog)
	r.RedirectTrailingSlash = false

	for _, docRoot := range s.cfg.Statics {
		if "" != docRoot.DocumentRoot {
			r.Use(static.Serve(docRoot.UrlPrefix, static.LocalFile(docRoot.DocumentRoot, docRoot.Indexes)))
		} else {
			r.Use(embedstatic.ServeEmbedded(docRoot.UrlPrefix, embedstatic.EmbedFile(docRoot.EmbedFileSystem, docRoot.Indexes)))
		}
	}

	if "" != s.cfg.HtmlContent {
		r.LoadHTMLGlob(s.cfg.HtmlContent)
	}

	s.engine = r

	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.BindAddress, s.cfg.ListenPort),
		Handler: s.engine,
	}

	log.Info().Str("server-mode", s.cfg.ServerMode).Msg("http server starting...")

	go func() {
		s.isReady = true
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Int("port", s.cfg.ListenPort).Msg("http server error on start up")
		}
	}()
	log.Info().Msgf("http server started on port %d", s.cfg.ListenPort)

	return nil
}

// IsReady returns the httpsrv readiness
func (s *serverImpl) Add2ServerContext(k string, v interface{}) error {
	return s.serverCtx.Add(k, v)
}

// IsReady returns the httpsrv readiness
func (s *serverImpl) IsReady() bool {
	return s.isReady
}

func (s *serverImpl) Stop() {
	log.Info().Msg("http server stopping...")
	s.isReady = false

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("http server stopping ... error")
	}

	log.Info().Msg("http server stopped")
}
