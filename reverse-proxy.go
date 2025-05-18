package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxyManager struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewReverseProxyManager(targetURL string) (*ReverseProxyManager, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = parsedURL.Host
		req.URL.Scheme = parsedURL.Scheme
		req.URL.Host = parsedURL.Host
		fmt.Printf("Proxying request to %s\n", req.URL.String())
	}

	return &ReverseProxyManager{
		target: parsedURL,
		proxy:  proxy,
	}, nil
}

func (rp *ReverseProxyManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Incoming request for host: %s\n", r.Host)
	rp.proxy.ServeHTTP(w, r)
}

func (rp *ReverseProxyManager) GetProxyInstance() *httputil.ReverseProxy {
	return rp.proxy
}
func main() {

	proxy, err := NewReverseProxyManager("http://example.com/")
	if err != nil {
		fmt.Print(err)
	}


	proxySite1, err := NewReverseProxyManager("http://localhost:4001")
	if err != nil {
		fmt.Print(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname := r.Host
		fmt.Printf("Incoming request for host: %s\n", hostname)
		if(r.Host == "indiastore1.duckdns.org"){
			proxySite1.ServeHTTP(w,r);
		}
		
		proxy.ServeHTTP(w,r);
	})

	fmt.Println("Reverse Proxy server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":80", nil))
}
