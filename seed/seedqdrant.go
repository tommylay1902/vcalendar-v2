package seed

import (
	"context"
	"fmt"

	"github.com/kelindar/search"
	"github.com/qdrant/go-client/qdrant"
)

func SeedGCOperations() {
	m, err := search.NewVectorizer("./dist/all-minilm-l6-v2-q8_0.gguf", 1)
	if err != nil {
		// handle error
		fmt.Println("hello")
		panic(err)
	}

	index := search.NewIndex[string]()

	embedUpdate, err := m.EmbedText("List events for today")
	if err != nil {
		fmt.Println("error with embedding tesxt")
		panic(err)
	}
	index.Add(embedUpdate, "List events for today")

	defer m.Close()

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		fmt.Println("error create collection")
		panic(err)
	}
	client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "gc_operations",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     384,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	operationInfo, err := client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "gc_operations",
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(4),
				Vectors: qdrant.NewVectors(embedUpdate...),
				Payload: qdrant.NewValueMap(map[string]any{"operation": "List"}),
			},
		},
	})
	if err != nil {
		fmt.Println("error insertin operation ")
		panic(err)

	}
	fmt.Println(operationInfo)
}
