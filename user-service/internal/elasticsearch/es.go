package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ElasticSearchClient interface {
	CreateIndex(ctx context.Context, indexName string) error
	AddDataToIndex(ctx context.Context, indexName string, data map[string]interface{}) error
	SearchDataFromIndex(ctx context.Context, indexName string, query string) (map[string]interface{}, error)
}

type elasticSearch struct {
	client *elasticsearch7.Client
}

func NewElasticSearch() (ElasticSearchClient, error) {

	esUrl := os.Getenv("ELASTIC_SEARCH_URL")
	esUser := os.Getenv("ELASTICSEARCH_USERNAME")
	esPassword := os.Getenv("ELASTICSEARCH_PASSWORD")
	log.Println("config ES: ", esUrl, esUser, esPassword)
	var cfg elasticsearch7.Config
	if esUrl != "" && esUser != "" && esPassword != "" {
		cfg = elasticsearch7.Config{
			Addresses: []string{
				esUrl,
			},
			Username: esUser,
			Password: esPassword,
		}
	} else {
		cfg = elasticsearch7.Config{
			Addresses: []string{
				"http://localhost:9200",
			},
			Username: "elastic",
			Password: "admin",
		}
	}

	es, err := elasticsearch7.NewClient(cfg)
	if err != nil {
		log.Fatalln("Error creating the client: %w", err)
		return nil, err
	}

	return &elasticSearch{
		client: es,
	}, nil
}

func (es *elasticSearch) CreateIndex(
	ctx context.Context,
	indexName string,
) error {

	// Check if the index already exists
	res, err := es.client.Indices.Exists([]string{indexName})
	if err != nil {
		log.Println("Error checking if index exists: ", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		// Index does not exist, create it
		// ================= Way 1: user elastic v7 =================
		// res, err = es.client.Indices.Create(
		// 	indexName,
		// 	es.client.Indices.Create.WithBody(strings.NewReader(initData)),
		// )

		// ================= Way 2: use esapi =================
		req := esapi.IndexRequest{
			Index:      indexName,
			DocumentID: "", // null to auto generate
			// Body:       strings.NewReader(initData),
			Refresh: "true",
		}
		res, err := req.Do(ctx, es.client)
		if err != nil {
			log.Println("Error indexing document: ", err)
			return err
		}
		defer res.Body.Close()

		fmt.Println("Index created: ", indexName)
	} else {
		fmt.Println("Index already exists: ", indexName)
	}
	return nil
}

func (es *elasticSearch) AddDataToIndex(
	ctx context.Context,
	indexName string,
	data map[string]interface{},
) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data: ", err)
		return err
	}

	// ================= Way 1: user elastic v7 =================
	// res, err := es.client.Index(
	// 	indexName,
	// 	bytes.NewReader(dataJSON),
	// )
	// ================= Way 2: use esapi =================
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: "", // null to auto generate
		Body:       bytes.NewReader(dataJSON),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, es.client)
	if err != nil {
		log.Println("Error indexing data: ", err)
		return err
	}
	defer res.Body.Close()
	fmt.Println("New data indexed in ", indexName)
	return nil
}

func (es *elasticSearch) SearchDataFromIndex(
	ctx context.Context,
	indexName string,
	query string,
) (map[string]interface{}, error) {

	// ================= Way 1: user elastic v7 =================
	// res, err = es.client.Indices.Create(
	// 	indexName,
	// 	es.client.Indices.Create.WithBody(strings.NewReader(initData)),
	// )

	// ================= Way 2: use esapi =================

	// query := `{
	//     "query": {
	//         "match_all": {}
	//     }
	// }`

	// Create a new SearchRequest
	req := esapi.SearchRequest{
		Index:          []string{indexName},
		Body:           strings.NewReader(query),
		TrackTotalHits: true,
		Pretty:         true,
	}
	res, err := req.Do(ctx, es.client)
	if err != nil {
		log.Println("Error searching data: ", err)
		return nil, err
	}
	defer res.Body.Close()
	// Check for errors in the response
	if res.IsError() {
		log.Println("Error response: ", res.String())
		return nil, fmt.Errorf("Error response: %w", res.String())
	}

	// Decode the JSON response body
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Println("Error parsing the response body: ", err)
		return nil, fmt.Errorf("Error parsing the response body: %w", err)
	}

	// // Print query results
	// fmt.Println("Response: ", res.Body)
	// for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
	// 	fmt.Println(hit.(map[string]interface{})["_source"])
	// }

	return r, nil
}
