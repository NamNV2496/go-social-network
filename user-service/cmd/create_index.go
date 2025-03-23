package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/domain"
	"github.com/spf13/cobra"

	"github.com/AlekSi/pointer"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/tokenchar"
	"github.com/moul/http2curl"
)

var createMappingIndex = &cobra.Command{
	Use:   "create_index",
	Short: "Create index with mapping",
	Long:  `Create index with mapping`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, _ := configs.NewConfig()
		log.Printf("[create_index] %v", conf)
		mapping := NewMapping(elasticsearch.Config{
			Addresses:     conf.ElasticSearch.Addr,
			Username:      conf.ElasticSearch.Username,
			Password:      conf.ElasticSearch.Password,
			RetryOnStatus: conf.ElasticSearch.RetryOnStatus,
			DisableRetry:  conf.ElasticSearch.DisableRetry,
			MaxRetries:    conf.ElasticSearch.MaxRetries,
			RetryBackoff: func(attempt int) time.Duration {
				return time.Duration(conf.ElasticSearch.RetryBackoff)
			},
		})

		if err := mapping.Start(); err != nil {
			panic(err)
		}

		log.Println("Create index with mapping successfully")
	},
}

var (
	TypeMapping = map[string]types.Property{
		"long":   types.NewLongNumberProperty(),
		"string": types.NewTextProperty(),
		"date":   types.NewDateProperty(),
	}
	settings = &types.IndexSettings{
		Analysis: &types.IndexSettingsAnalysis{
			Analyzer: map[string]types.Analyzer{
				"custome_vietnamese": types.CustomAnalyzer{
					Filter:    []string{"lowercase", "asciifolding", "vi_stopwords"},
					Type:      "custom",
					Tokenizer: "standard",
				},
				"autocomplete": types.CustomAnalyzer{
					Filter:    []string{"lowercase", "autocomplete_filter"},
					Type:      "custom",
					Tokenizer: "autocomplete",
				},
				"autocomplete_search": types.NewStandardAnalyzer(),
			},
			Tokenizer: map[string]types.Tokenizer{
				"autocomplete": types.EdgeNGramTokenizer{
					TokenChars: []tokenchar.TokenChar{{Name: "letter"}},
					Type:       "edge_ngram",
				},
			},
		},
		Mapping: &types.MappingLimitSettings{
			Coerce: pointer.ToBool(true),
		},
		Search: &types.SettingsSearch{
			Slowlog: &types.SlowlogSettings{
				Threshold: &types.SlowlogTresholds{
					Query: &types.SlowlogTresholdLevels{
						Debug: "200ms",
					},
				},
			},
		},
		RefreshInterval:  "30s",
		NumberOfShards:   "3",
		NumberOfReplicas: "2",
	}
)

type IMapping interface {
	Start() error
}

type Mapping struct {
	client *elasticsearch.TypedClient
}

func NewMapping(
	config elasticsearch.Config,
) IMapping {
	client, err := elasticsearch.NewTypedClient(config)
	if err != nil {
		panic(err)
	}
	return &Mapping{
		client: client,
	}
}

var _ IMapping = &Mapping{}

func (m *Mapping) Start() error {
	user, err := generateMappingProperties(&domain.User{})
	if err != nil {
		return err
	}
	userUser, err := generateMappingProperties(&domain.UserUser{})
	if err != nil {
		return err
	}
	var inp string
	log.Println("Do you want to create or alternate mapping of index? (0/1)")
	fmt.Scanf("%s", &inp)
	if inp == "0" {
		log.Println("==================================")
		if err := m.createIndexMapping("user_index", user); err != nil {
			return err
		}
	} else if inp == "1" {
		log.Println("==================================")
		if err := m.alternateIndexMapping("user_index", userUser); err != nil {
			return err
		}
	} else {
		panic(errors.New("unsupported method"))
	}

	log.Println("Copy and run curl in postman to update mapping")
	return nil
}

func (m *Mapping) createIndexMapping(index string, properties map[string]types.Property) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// create
	req := m.client.
		Indices.
		Create(index).
		Settings(settings).
		Mappings(
			&types.TypeMapping{
				Properties: properties,
			},
		)
	request, _ := req.HttpRequest(timeoutCtx)
	debug, err := http2curl.GetCurlCommand(request)
	log.Printf("run this curl in postman: %s", debug)
	// _, err := req.Do(timeoutCtx)

	return err
}

func (m *Mapping) alternateIndexMapping(index string, properties map[string]types.Property) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Edit mapping
	req := m.client.
		Indices.
		PutMapping(index).
		Properties(properties)
	request, _ := req.HttpRequest(timeoutCtx)
	debug, err := http2curl.GetCurlCommand(request)
	log.Printf("run this curl in postman: %s", debug)

	return err
}

func generateMappingProperties(record any) (map[string]types.Property, error) {
	if record == nil {
		return nil, nil
	}
	ift := reflect.TypeOf(record)
	if ift == nil {
		return nil, errors.New("record is nil")
	}
	if ift.Kind() == reflect.Ptr {
		ift = ift.Elem()
	}
	properties := make(map[string]types.Property)
	for i := 0; i < ift.NumField(); i++ {
		field := ift.Field(i)
		if field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) {
			subProperties, err := generateMappingProperties(reflect.New(field.Type).Elem().Interface())
			if err != nil {
				return nil, err
			}
			for k, v := range subProperties {
				if _, ok := properties[k]; ok {
					return nil, errors.New("field name is duplicated")
				}
				properties[k] = v
			}
		}
		esMappingField := field.Tag.Get("esMapping")
		if len(esMappingField) == 0 {
			return nil, errors.New("esMapping tag is empty")
		}
		if esMappingField == "-" {
			// skip fields which have type = "-"
			continue
		}
		key, esType := parseEsMappingField(esMappingField)
		if key != "" && esType != "" {
			properties[key] = TypeMapping[esType]
		}
	}
	return properties, nil
}

func parseEsMappingField(tag string) (string, string) {
	tags := strings.Split(tag, ",")
	if len(tags) < 1 {
		return "", ""
	}
	var key, esType string
	for _, tag := range tags {
		parts := strings.Split(tag, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == "key" {
			key = parts[1]
		} else if parts[0] == "type" {
			esType = parts[1]
		}
	}
	return key, esType
}
