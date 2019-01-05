package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
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
	//req.URL.Host = "unix://"
	req.URL.Scheme = "http"
	req.URL.Host = "/var/run/docker.sock"

	if req.Header.Get("Upgrade") == "tcp" && req.Header.Get("Connection") == "Upgrade" {
		//hijack
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}

		b, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Print(string(b))

		cc, err := net.Dial("unix", "/var/run/docker.sock")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		conn, bufrw, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Don't forget to close the connection:
		defer conn.Close()

		cc.Write(b)
		go io.Copy(cc, conn)
		go io.Copy(bufrw, cc)
		select {
		case <-req.Context().Done():
		}
	} else {
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
}
