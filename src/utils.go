package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func newBasicServer(addr string) *basicServer {
	serverUrl, err := url.Parse(addr)
	handleError(err)

	return &basicServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func handleError(err error) {
	if err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		servers:         servers,
		roundRobinCount: 0,
	}
}

func (l *LoadBalancer) getNextAvailableServer() Server {
	server := l.servers[l.roundRobinCount%len(l.servers)]

	for !server.IsAlive() {
		l.roundRobinCount++
		server = l.servers[l.roundRobinCount%len(l.servers)]
	}
	l.roundRobinCount++

	return server
}

func (l *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := l.getNextAvailableServer()
	log.Printf("forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(w, r)
}

func (s *basicServer) Address() string { return s.addr }

func (s *basicServer) IsAlive() bool { return true }

func (s *basicServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}
