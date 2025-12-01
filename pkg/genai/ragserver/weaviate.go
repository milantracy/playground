package ragserver

import (
	"context"
	"fmt"
	"log"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

const CollectionName = "Docs"

func newWeatherClient(ctx context.Context, config weaviate.Config) (*weaviate.Client, error) {
	cli, err := weaviate.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to init weaviate client: %v", err)
	}

	col := &models.Class{
		Class:      CollectionName,
		Vectorizer: "none",
	}
	exists, err := cli.Schema().ClassExistenceChecker().WithClassName(col.Class).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if class exists: %v", err)
	}
	if !exists {
		err = cli.Schema().ClassCreator().WithClass(col).Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create class: %v", err)
		}
	}
	return cli, nil
}

func decodeGraphQLResult(result *models.GraphQLResponse) ([]string, error) {
	log.Printf("graphql response: %v", result)
	data, ok := result.Data["Get"]
	if !ok {
		return nil, fmt.Errorf("invalid GraphQL response format")
	}
	doc, ok := data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid GraphQL data format")
	}
	l, ok := doc[CollectionName].([]any)
	if !ok {
		return nil, fmt.Errorf("invalid GraphQL collection format")
	}
	var out []string
	for _, item := range l {
		itemMap, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid GraphQL item format")
		}
		t, ok := itemMap["text"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid text field format")
		}
		out = append(out, t)
	}
	return out, nil
}

func combinedWeaviateError(result *models.GraphQLResponse, err error) error {
	if err != nil {
		return err
	}
	if len(result.Errors) != 0 {
		var ss []string
		for _, e := range result.Errors {
			ss = append(ss, e.Message)
		}
		return fmt.Errorf("weaviate error: %v", ss)
	}
	return nil
}
