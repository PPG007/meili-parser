package main

import (
	"github.com/meilisearch/meilisearch-go"
)

func main() {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: "test_key",
	})
	_, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "ppg007-github-io",
		PrimaryKey: "id",
	})
	if err != nil {
		panic(err)
	}
	baseUrl := "https://ppg007.github.io"
	index, err := client.GetIndex("ppg007-github-io")
	if err != nil {
		panic(err)
	}
	index.AddDocuments(map[string]interface{}{
		"a":  1,
		"id": "123",
	})
	parseFromStartPage(baseUrl, index)
}
