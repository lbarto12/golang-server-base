package meilisearchapi

import (
	"errors"
	"fmt"
	"golang-server-base/api/meilisearchapi/models"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

var client meilisearch.ServiceManager

func EnvGetOptions() models.MeilisearchOptions {
	host, ok := os.LookupEnv("MEILI_HOST")
	if !ok {
		panic("MEILI_HOST environment variable not set.")
	}

	port, ok := os.LookupEnv("MEILI_PORT")
	if !ok {
		panic("MEILI_PORT environment variable not set.")
	}

	masterKey, ok := os.LookupEnv("MEILI_MASTER_KEY")
	if !ok {
		panic("MEILI_MASTER_KEY environment variable not set.")
	}

	return models.MeilisearchOptions{
		Host:   host,
		Port:   port,
		APIKey: masterKey,
	}
}

func Init(options models.MeilisearchOptions) {
	url := fmt.Sprintf("%s:%s", options.Host, options.Port)
	client = meilisearch.New(url, meilisearch.WithAPIKey(options.APIKey))
}

func Client() (meilisearch.ServiceManager, error) {
	if !client.IsHealthy() {
		return nil, errors.New("meilisearch service not available")
	}
	return client, nil
}
