package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/milantracy/playground/internal/ts"
)

func main() {
	timeServer := ts.TimeServer{}
	rpc.Register(&timeServer)
	rpc.HandleHTTP()
	// Listen for the incoming requests on port 8080.
	listerner, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Listener err: ", err)
	}
	http.Serve(listerner, nil)
}
