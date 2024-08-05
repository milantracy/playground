package main

import (
	"log"
	"net/rpc"

	"github.com/milantracy/playground/internal/ts"
)

func main() {
	var reply int64
	args := ts.TimeServerArgs{}

	// Try localhost:8080
	client, err := rpc.DialHTTP("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	err = client.Call("TimeServer.ServerTime", args, &reply)
	if err != nil {
		log.Fatal("RPC error: ", err)
	}

	// Just print.
	log.Printf("%d", reply)
}
