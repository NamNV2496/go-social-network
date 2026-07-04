package langfuse

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/namnv2496/post-service/configs"
)

func TestGetPromptByProjectAndNameReturnsDetail(t *testing.T) {
	// Mock Langfuse server
	h := http.NewServeMux()
	h.HandleFunc("/api/public/v2/prompts", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		project := q.Get("project")
		name := q.Get("name")
		if project == "test_project" && name == "test_prompt" {
			resp := map[string]any{
				"data": []map[string]any{{
					"name":          "test_prompt",
					"tags":          []string{},
					"lastUpdatedAt": "2026-07-04T07:29:11.731Z",
					"versions":      []int{1},
					"labels":        []string{"latest", "production"},
					"type":          "text",
					"lastConfig":    map[string]any{},
				}},
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	h.HandleFunc("/api/public/v2/prompts/test_prompt", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"name":    "test_prompt",
			"version": 1,
			"prompt":  "test_prompt_content",
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	// Accept OTLP POSTs to /api/public/otel to avoid exporter errors
	h.HandleFunc("/api/public/otel", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(h)
	defer srv.Close()

	cfg := configs.CommentRule{
		Langfuse: configs.Langfuse{
			Enabled:   true,
			Endpoint:  srv.URL,
			PublicKey: "pk-lf-76a7e93b-ef39-47b6-9c8c-2224c6a60c77",
			SecretKey: "sk-lf-3190eec3-32f7-47f4-91ca-41fd7a338a6e",
		},
		LangfuseTrace: configs.LangfuseTrace{
			Endpoint:  srv.URL,
			PublicKey: "pk",
			SecretKey: "sk",
			TraceName: "comment_rule_ai_check",
		},
	}

	client := NewClient(cfg)
	if client == nil {
		t.Fatalf("NewClient returned nil")
	}

	p, err := client.GetPromptByProjectAndName(context.Background(), "test_project", "test_prompt")
	if err != nil {
		t.Fatalf("GetPromptByProjectAndName error: %v", err)
	}
	if p != "test_prompt_content" {
		t.Fatalf("unexpected prompt: %s", p)
	}
}
