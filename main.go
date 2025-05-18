package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
)

type ReverseProxy struct {
    routes map[string]*httputil.ReverseProxy
}

func NewReverseProxy() *ReverseProxy {
    return &ReverseProxy{
        routes: make(map[string]*httputil.ReverseProxy),
    }
}

func (rp *ReverseProxy) AddRoute(hostname string, target string) {
    targetURL, err := url.Parse(target)
    if err != nil {
        log.Fatalf("Invalid target URL: %v", err)
    }
    rp.routes[hostname] = httputil.NewSingleHostReverseProxy(targetURL)
    log.Printf("Route added: %s -> %s", hostname, target)
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    hostname := strings.ToLower(r.Host)
    proxy, exists := rp.routes[hostname]
    
    if !exists {
        http.Error(w, "Host not found", http.StatusNotFound)
        return
    }

    log.Printf("Forwarding request for %s to %s", hostname, proxy)
    proxy.ServeHTTP(w, r)
}

func main() {
    proxy := NewReverseProxy()

    proxy.AddRoute("indiastore1.duckdns.com",  "http://localhost:4001")
    proxy.AddRoute("indiastore2.duckdns.com", "http://localhost:4002")
    proxy.AddRoute("indiastore3.duckdns.com", "http://localhost:4003")
    proxy.AddRoute("indiastore4.duckdns.com", "http://localhost:4004")
    proxy.AddRoute("indiastore5.duckdns.com", "http://localhost:4005")
    proxy.AddRoute("indiastore6.duckdns.com", "http://localhost:4006")
    proxy.AddRoute("indiastore7.duckdns.com", "http://localhost:4007")
    proxy.AddRoute("indiastore8.duckdns.com", "http://localhost:4008")
    proxy.AddRoute("indiastore9.duckdns.com", "http://localhost:4009")
    proxy.AddRoute("indiastore10.duckdns.com","http://localhost:4010")

    log.Println("Starting reverse proxy on :443")
    if err := http.ListenAndServe(":443", proxy); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
