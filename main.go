package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}
type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

func newSimpleServer(addr string) *simpleServer {
	serverurl, err := url.Parse(addr)
	handleError(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverurl),
	}
}

type LoadBalancer struct {
	port            string
	roundRobinIndex int
	serverList      []Server
}

func newLoadBalancer(port string, serverList []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinIndex: 0,
		serverList:      serverList,
	}
}
func handleError(err error) {
	if err != nil {
		fmt.Printf("error : %v\n", err)
		os.Exit(1)
	}
}
func(s *simpleServer) Address string {
	return s.addr
}
func(s *simpleServer) IsAlive bool {
	return true
}

func(s *simpleServer) Serve(rw http.ResponseWriter, rq *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

func (lb *LoadBalancer) getNextAvailableServer(Server) {
	server:=lb.serverList[lb.roundRobinIndex%len(lb.serverList)]
	for server.IsAlive() {
		lb.roundRobinIndex++
		server = lb.serverList[lb.roundRobinIndex%len(lb.serverList)]
	}
	lb.roundRobinIndex++
	return server

}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Redirecting to server : %s\n", targetServer.Address())
	targetServer.Serve(rw, req)
}

func main() {
	serverList := []Server{
		newSimpleServer("http://www.facebook.com"),
		newSimpleServer("http://www.google.com"),
		newSimpleServer("http://www.youtube.com"),
	}

	lb := newLoadBalancer("8080", serverList)
	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, req)
	}
	http.handleFunc("/", handleRedirect)
	fmt.Printf("Load Balancer started at :%s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
