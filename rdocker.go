package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type rdocker struct {
	client http.Client
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	err = http.Serve(l, &rdocker{
		client: http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (conn net.Conn, e error) {
					return net.Dial("unix", "/var/run/docker.sock")
				},
			},
			Timeout: 60 * time.Second,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (r *rdocker) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//req.URL.Host = "unix:///var/run/docker.sock"
	req.RequestURI = ""
	resp, err := r.client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	for k := range resp.Header {
		w.Header().Add(k, resp.Header.Get(k))
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
	_ = resp.Body.Close()
}
