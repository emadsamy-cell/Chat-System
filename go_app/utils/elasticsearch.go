package utils

import (
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

func GetElasticsearchClient() *elasticsearch.Client {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return es
}
