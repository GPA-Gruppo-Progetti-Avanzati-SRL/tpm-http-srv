package httputil

import (
	"bytes"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

func newHarRequest(c *gin.Context) (*har.Request, error) {
	const semLogContext = "http-server::har-request-from-http-request"

	var err error
	var bodyContent []byte
	if c.Request.ContentLength > 0 {
		bodyContent, err = io.ReadAll(c.Request.Body)
		log.Trace().Msgf(semLogContext+" - Lunghezza body :  %d, - Body : %s ,", len(bodyContent), string(bodyContent))
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		// Put back for subsequent validation phase.
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyContent))
	}

	req, err := har.NewRequest(
		c.Request.Method,
		c.Request.URL.String(),
		bodyContent,
		c.Request.Header,
		HarQueryStringParamsFromGinParams(c),
		HarParamsFromGinParams(c))
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return req, nil
}

// HarParamsFromGinParams introduced because of sonar cognitive complexity.
func HarParamsFromGinParams(c *gin.Context) []har.Param {
	var pars []har.Param
	if len(c.Params) > 0 {
		pars = make([]har.Param, 0)
		for _, h := range c.Params {
			pars = append(pars, har.Param{Name: h.Key, Value: h.Value})
		}
	}

	return pars
}

func HarQueryStringParamsFromGinParams(c *gin.Context) har.NameValuePairs {
	var pars har.NameValuePairs
	if len(c.Request.URL.Query()) > 0 {
		pars = make(har.NameValuePairs, 0)
		for n, q := range c.Request.URL.Query() {
			pars = append(pars, har.NameValuePair{Name: n, Value: strings.Join(q, ",")})
		}
	}

	return pars
}
