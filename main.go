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

func (s *Server) Address() string{return s.address}

func (s *Server) Is_alive() bool{return true}

func (s *Server) Serve(rw http.ResponseWriter, req *http.Request){
	s.proxy.ServeHTTP(rw, req);
}

func new_server(address string) *Server{
	server_url, err := url.Parse(address)

	handle_error(err);

	return &Server{
		address: address,
		proxy: httputil.NewSingleHostReverseProxy(server_url),
	}
}

func New_Load_balancer(port string, servers []Server) *Load_Balancer{
	return &Load_Balancer{
		port: port,
		round_robin_count: 0,
		servers: servers,
	}
}

func (lb *Load_Balancer) get_next_server() Server{

	for i := 0; i < len(lb.servers); i++ {
		server := &lb.servers[lb.round_robin_count%len(lb.servers)]

		if server.Is_alive() {
			lb.round_robin_count++;
			return *server;
		}

		lb.round_robin_count++;
	}

	// Fallback to the last server if none are alive
	return lb.servers[len(lb.servers)-1]
}

func (lb *Load_Balancer) serve_proxy(rw http.ResponseWriter, req *http.Request){
	// fmt.Println("Here");
	target_server := lb.get_next_server();
	fmt.Printf("Forwarding request to %q\n", target_server.Address());
	target_server.Serve(rw, req);
} 

func main(){
	servers := []Server{
		*new_server("https://www.google.com"),
		*new_server("https://www.bing.com"),
		*new_server("https://www.duckduckgo.com"),
	}

	lb := New_Load_balancer("8000", servers);
	fmt.Println(lb.round_robin_count);

	handle_redirect := func(rw http.ResponseWriter, req *http.Request){
		// fmt.Println("Here");
		lb.serve_proxy(rw, req);
	}

	http.HandleFunc("/", handle_redirect);
	fmt.Printf("Serving at 'Localhost: %s'\n", lb.port);
	http.ListenAndServe(":"+lb.port, nil);
}