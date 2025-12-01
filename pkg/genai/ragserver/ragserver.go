package ragserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"google.golang.org/api/option"
)

const (
	geminiFlashModel      = "gemini-2.5-flash"
	geminiFlashEmbedModel = "text-embedding-004"
	ragTemplateStr        = `
I will ask you a question and will provide some additional context information.
Assume this context information is factual and correct, as part of internal
documentation.
If the question relates to the context, answer it using the context.
If the question does not relate to the context, answer it as normal.

For example, let's say the context has nothing in it about tropical flowers;
then if I ask you about tropical flowers, just answer what you know about them
without referring to the context.

For example, if the context does mention minerology and I ask you about that,
provide information from the context along with general knowledge.

Question:
%s

Context:
%s
`
)

type RagServer struct {
	ctx            context.Context
	wvClient       *weaviate.Client
	genaiMdel      *genai.GenerativeModel
	embeddingModel *genai.EmbeddingModel
}

type Cleanup func()

type RagServerConfig struct {
	WeaviateConfig weaviate.Config
	GeminiOptions  option.ClientOption
}

func New(config RagServerConfig) (*RagServer, Cleanup, error) {
	ctx := context.Background()

	weaviateClient, err := newWeatherClient(ctx, config.WeaviateConfig)
	if err != nil {
		return nil, nil, err
	}

	genaiClient, err := genai.NewClient(ctx, config.GeminiOptions)
	if err != nil {
		return nil, nil, err
	}

	serverr := &RagServer{
		ctx:            ctx,
		wvClient:       weaviateClient,
		genaiMdel:      genaiClient.GenerativeModel(geminiFlashModel),
		embeddingModel: genaiClient.EmbeddingModel(geminiFlashEmbedModel),
	}
	cleanup := func() {
		genaiClient.Close()
	}
	return serverr, cleanup, nil
}

// Query implments net/http.Handler's ServeHTTP to
// handle query requests to the RAG server.
func (s *RagServer) Query(w http.ResponseWriter, r *http.Request) {
	type queryRequest struct {
		Query string
	}

	req := &queryRequest{}
	err := requestToJSON(r, req)
	log.Printf("Received query: %v", req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert query to an embedding.
	resp, err := s.embeddingModel.EmbedContent(s.ctx, genai.Text(req.Query))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// search in weaviate
	ql := s.wvClient.GraphQL()
	result, err := ql.Get().
		WithNearVector(ql.NearVectorArgBuilder().WithVector(resp.Embedding.Values)).
		WithClassName(CollectionName).
		WithFields(graphql.Field{Name: "text"}).
		WithLimit(3).
		Do(s.ctx)
	if werr := combinedWeaviateError(result, err); werr != nil {
		http.Error(w, werr.Error(), http.StatusInternalServerError)
		return
	}

	contents, err := decodeGraphQLResult(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ragQuery := fmt.Sprintf(ragTemplateStr, req.Query, strings.Join(contents, "\n"))
	// Let LLM answer the question with context.
	ragAnswer, err := s.genaiMdel.GenerateContent(s.ctx, genai.Text(ragQuery))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(ragAnswer.Candidates) != 1 {
		http.Error(w, "unexpected number of candidates", http.StatusInternalServerError)
		return
	}

	var answer []string
	for _, part := range ragAnswer.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			answer = append(answer, string(text))
		} else {
			http.Error(w, "unexpected part type", http.StatusInternalServerError)
			return
		}
	}
	writeJSONResponse(w, strings.Join(answer, "\n"))
}

// AddDocument implement net/http.Handler's ServeHTTP to
// add documents to the RAG server and store their embeddings in Weaviate.
func (s *RagServer) AddDocument(w http.ResponseWriter, r *http.Request) {
	type docuement struct {
		Text string
	}

	type addRequest struct {
		Documents []docuement
	}

	req := &addRequest{}

	err := requestToJSON(r, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%v documents added.", len(req.Documents))

	// converts documents to embeddings
	batch := s.embeddingModel.NewBatch()
	for _, doc := range req.Documents {
		batch.AddContent(genai.Text(doc.Text))
	}
	resp, err := s.embeddingModel.BatchEmbedContents(s.ctx, batch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(resp.Embeddings) != len(req.Documents) {
		http.Error(w, "embedding count mismatch", http.StatusInternalServerError)
		return
	}

	// store embeddings to weaviate
	objects := make([]*models.Object, len(req.Documents))
	for i, doc := range req.Documents {
		objects[i] = &models.Object{
			Class: CollectionName,
			Properties: map[string]interface{}{
				"text": doc.Text,
			},
			Vector: resp.Embeddings[i].Values,
		}
	}

	_, err = s.wvClient.Batch().ObjectsBatcher().WithObjects(objects...).Do(s.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
