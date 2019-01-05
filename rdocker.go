package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (conn net.Conn, e error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
	err = http.Serve(l, proxy)
	if err != nil {
		log.Fatal(err)
	}
}
