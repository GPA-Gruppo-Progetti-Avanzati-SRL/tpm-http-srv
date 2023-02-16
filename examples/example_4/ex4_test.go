package example_4_test

import (
	_ "GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/examples/example_4"
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {

	s, err := httpsrv.NewServer(httpsrv.DefaultConfig,
		httpsrv.WithBindAddress("localhost"),
		httpsrv.WithListenPort(8080),
		httpsrv.WithShutdownTimeout(time.Duration(5)*time.Second),
		httpsrv.WithContextPath("/api"))

	if err != nil {
		panic(err.Error())
	}

	if err := s.Start(); err != nil {
		panic(err.Error())
	}
	defer s.Stop()

	for !s.IsReady() {
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

	resp, err := http.Get("http://:10003/api/v1/en/test/sayhello/gotest")
	if nil != err {
		t.Fatal(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		t.Fatal(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Wrong Status Code")
	assert.Equal(t, "Hello gotest", string(body), "Wrong Response")
	assert.Equal(t, "uk", resp.Header.Get("X-lang"), "Wrong Header")

	respfr, err := http.Get("http://:10003/api/v1/fr/test/sayhello/gotest")
	if nil != err {
		t.Fatal(err.Error())
	}
	defer respfr.Body.Close()

	bodyfr, err := ioutil.ReadAll(respfr.Body)
	if nil != err {
		t.Fatal(err.Error())
	}
	assert.Equal(t, http.StatusOK, respfr.StatusCode, "Wrong Status Code")
	assert.Equal(t, "Bonjour gotest", string(bodyfr), "Wrong Response")
	assert.Equal(t, "fr", respfr.Header.Get("X-lang"), "Wrong Header")

}
