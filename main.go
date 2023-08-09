package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
)

type server struct{
	address string
	proxy *httputil.ReverseProxy
}

func new_server(address string) *server{
	server_url, err := url.Parse(address)

	if err != nil {
		panic(err);
	}

	return &server{
		address: address,
		proxy: httputil.NewSingleHostReverseProxy(server_url),
	}
}

func main(){
	fmt.Println("Hello");
}