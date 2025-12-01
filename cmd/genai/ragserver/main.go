package main

import (
	"log"
	"net/http"
	"os"

	"github.com/milantracy/playground/pkg/genai/ragserver"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"google.golang.org/api/option"
)

func main() {
	conf := ragserver.RagServerConfig{
		WeaviateConfig: weaviate.Config{
			Host:   "localhost:8090",
			Scheme: "http",
		},
		GeminiOptions: option.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
	}
	server, cleanup, err := ragserver.New(conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	http.HandleFunc("POST /add/", server.AddDocument)
	http.HandleFunc("POST /query/", server.Query)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
