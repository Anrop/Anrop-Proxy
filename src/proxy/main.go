package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
)

const host = "anrop.se.preview.citynetwork.se"

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Host = host
			return r, nil
		})
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host == "" {
			fmt.Fprintln(w, "Cannot handle requests without Host header, e.g., HTTP 1.0")
			return
		}
		req.URL.Scheme = "http"
		req.URL.Host = host
		proxy.ServeHTTP(w, req)
	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":"+port, proxy))
}
