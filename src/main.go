package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	serverList := []Server{
		newBasicServer("https://www.facebook.com"),
		newBasicServer("https://www.bing.com"),
		newBasicServer("https://www.duckduckgo.com"),
	}

	lb := NewLoadBalancer("8000", serverList)

	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Serving requests at 'localhost:%s'\n", lb.port)
	log.Fatal(http.ListenAndServe(":"+lb.port, nil))
}
