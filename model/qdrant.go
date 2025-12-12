package model

import (
	"context"
	"fmt"
	"log"

	"github.com/kelindar/search"
	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	qc       *qdrant.Client
	embedder *search.Vectorizer
}

func InitializeQdrantClient() QdrantClient {
	m, err := search.NewVectorizer("./dist/all-minilm-l6-v2-q8_0.gguf", 1)
	if err != nil {
		fmt.Println("error setting up embedding client:", err)
		// handle error
		panic(err)
	}

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		fmt.Println("error setting up qdrant client:", err)
		panic(err)
	}

	qc := QdrantClient{
		qc:       client,
		embedder: m,
	}

	return qc
}

func (c *QdrantClient) GetOperation(text *string) string {
	if text != nil {
		embeddedMsg, err := c.embedder.EmbedText(*text)
		if err != nil {
			log.Printf("Error embedding text: %v", err)
			panic(err)
		}
		result, err := c.qc.Query(context.Background(), &qdrant.QueryPoints{
			CollectionName: "gc_operations",
			Query:          qdrant.NewQuery(embeddedMsg...),
			WithPayload:    qdrant.NewWithPayload(true),
		})
		if err != nil {
			fmt.Println("Error querying Qdrant:", err)
			panic(err)
		}
		payload := result[0].GetPayload()
		if operationValue, exists := payload["operation"]; exists {
			// The value is a *qdrant.Value - we need to get the string from it
			qdrantValue := operationValue

			// Check if it has a string value and extract it
			if qdrantValue.GetStringValue() != "" {
				operation := qdrantValue.GetStringValue()
				fmt.Println(operation) // Prints: Delete
				return operation
			}
		}
	}
	return ""
}
