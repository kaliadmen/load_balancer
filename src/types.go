package main

import (
	"net/http"
	"net/http/httputil"
)

type basicServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}

type LoadBalancer struct {
	port            string
	servers         []Server
	roundRobinCount int
}
