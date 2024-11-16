package httputil

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/gin-gonic/gin"
	"strings"
)

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
