package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"text/template"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/namnv2496/user-service/internal/domain"
)

var (
	MapQueryTemplates = map[string]*template.Template{
		// "tmplGetItemByShopId":        template.Must(template.ParseFiles(filepath.Join("templates", "tmplGetItemByShopId.json"))),
		// "tmplFindItemByNameMatch":    template.Must(template.ParseFiles(filepath.Join("templates", "tmplFindItemByNameMatch.json"))),
		"tmplFindItemByNameRegex": template.Must(template.ParseFiles(filepath.Join("templates", "tmplFindItemByNameRegex.json"))),
		// "tmplFindItemByNameWildcard": template.Must(template.ParseFiles(filepath.Join("templates", "tmplFindItemByNameWildcard.json"))),
	}
)

func (u *elasticSearch) TemplateQuery(ctx context.Context, user_id, template string, page, size int64) ([]*domain.User, error) {
	if u.esClient == nil {
		return nil, errors.New("unsupported server: Elasticsearch client is not initialized")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	params := struct {
		UserId string `json:"user_id,omitempty"`
		From   int64  `json:"from,omitempty"`
		Size   int64  `json:"size,omitempty"`
	}{
		UserId: user_id,
		From:   page * size,
		Size:   size,
	}

	queryBuffer, err := buildQuery(MapQueryTemplates["tmplFindItemByNameRegex"], params)
	if err != nil {
		fmt.Println("Error building query:", err)
		return nil, err
	}

	var query types.Query
	if err := json.Unmarshal(queryBuffer.Bytes(), &query); err != nil {
		return nil, err
	}

	result, err := u.esClient.Search().
		Index("user").
		Query(&query).
		Header("Content-Type", "application/json").
		Do(timeoutCtx)

	if err != nil {
		return nil, err
	}

	var resp []*domain.User
	for _, h := range result.Hits.Hits {
		var user *domain.User
		if err := json.Unmarshal(h.Source_, &user); err != nil {
			return nil, err
		}
		resp = append(resp, user)
	}
	return resp, nil
}

func buildQuery(tmpl *template.Template, params any) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return nil, err
	}
	return &buf, nil
}
