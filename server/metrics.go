package main

import "net/http"

func newMetricsServer(httpAddr string) *http.Server {
	httpSrv := &http.Server{Addr: httpAddr}
	m := http.NewServeMux()
	httpSrv.Handler = m
	return httpSrv
}
