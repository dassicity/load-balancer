package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Serve interface{
	Address() string
	Is_alive()	bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type Server struct{
	address 	string
	proxy 		*httputil.ReverseProxy
}

type Load_Balancer struct{
	port 				string
	round_robin_count 	int
	servers 			[]Server
}

func handle_error(err error){
	if err != nil {
		panic(err);
	}
}

func new_server(address string) *Server{
	server_url, err := url.Parse(address)

	handle_error(err);

	return &Server{
		address: address,
		proxy: httputil.NewSingleHostReverseProxy(server_url),
	}
}




func main(){
	fmt.Println("Hello");
}