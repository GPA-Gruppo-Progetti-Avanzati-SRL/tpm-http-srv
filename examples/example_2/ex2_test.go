package example_2_test

import (
	_ "GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/examples/example_2"
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"fmt"
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

	resp, err := http.Get("http://:8080/api/v1/test/sayhello/gotest")
	if nil != err {
		panic(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		panic(err.Error())
	}

	fmt.Printf("%s\n", body)

}
