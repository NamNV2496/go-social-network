package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/namnv2496/post-service/configs"
	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/entity"
	"github.com/namnv2496/post-service/internal/pkg"
	pkgLangfuse "github.com/namnv2496/post-service/internal/pkg/langfuse"
	"github.com/namnv2496/post-service/internal/repository"
	"github.com/namnv2496/post-service/internal/repository/mq/producer"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type ICommentService interface {
	AddComment(context.Context, *entity.CreateCommentRequest) (*entity.CreateCommentResponse, error)
	GetComment(context.Context, *entity.GetCommentRequest) (*entity.GetCommentResponse, error)
	CreateCommentRule(context.Context, *entity.CreateCommentRuleRequest) (int64, error)
	GetCommentRules(context.Context, *entity.GetCommentRulesRequest) ([]*entity.GetCommentRulesResponse, error)
	UpdateCommentRule(context.Context, *entity.UpdateCommentRuleRequest) (*entity.UpdateCommentRuleResponse, error)
}

type commentService struct {
	commentRepository     repository.ICommentRepository
	commentRuleRepository repository.ICommentRuleRepository
	trie                  pkg.ITrie
	trieRoot              *pkg.Trie
	cfg                   configs.Config
	langfuseClient        *pkgLangfuse.Client
}

func NewCommentService(
	commentRepository repository.ICommentRepository,
	commentRuleRepository repository.ICommentRuleRepository,
	trie pkg.ITrie,
	kafkaClient producer.Client,
	langfuseClient *pkgLangfuse.Client,
	cfg configs.Config,
) ICommentService {
	return &commentService{
		commentRepository:     commentRepository,
		commentRuleRepository: commentRuleRepository,
		trie:                  trie,
		cfg:                   cfg,
		langfuseClient:        langfuseClient,
	}
}

var _ ICommentService = &commentService{}

func (c *commentService) AddComment(ctx context.Context, req *entity.CreateCommentRequest) (*entity.CreateCommentResponse, error) {
	// comment rule check
	lowercaseCommentText := strings.ToLower(req.Comment.CommentText)
	lowercaseCommentText, _, err := transform.String(transform.Chain(norm.NFKC), lowercaseCommentText)
	if err != nil {
		return nil, err
	}
	violation, violenWords, err := c.commentRuleCheck(ctx, &entity.CommentRuleCheckRequest{
		CommentText: lowercaseCommentText,
		Application: req.Application,
	})
	if err != nil {
		return nil, err
	}
	if violation {
		return nil, fmt.Errorf("Comment rule violation with word: %s", violenWords)
	}
	comment := domain.Comment{
		PostId:        int64(req.PostId),
		UserId:        req.Comment.UserId,
		CommentText:   req.Comment.CommentText,
		CommentLevel:  int64(req.Comment.CommentLevel),
		CommentParent: int64(req.Comment.CommentParent),
		Images:        strings.Join(req.Comment.Images, ","),
		Tags:          strings.Join(req.Comment.Tags, ","),
	}
	if err := c.commentRepository.AddComment(ctx, comment); err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *commentService) GetComment(ctx context.Context, req *entity.GetCommentRequest) (*entity.GetCommentResponse, error) {
	comments, err := c.commentRepository.GetComment(ctx, domain.CommentByPostId(int64(req.PostId)))
	if err != nil {
		return nil, err
	}
	var resp []*entity.Comment
	for _, comment := range comments {
		resp = append(resp, &entity.Comment{
			CommentId:     uint64(comment.Id),
			UserId:        comment.UserId,
			CommentText:   comment.CommentText,
			CommentLevel:  uint32(comment.CommentLevel),
			CommentParent: uint64(comment.CommentParent),
			Images:        strings.Split(comment.Images, ","),
			Tags:          strings.Split(comment.Tags, ","),
			Date:          comment.CreatedAt.String(),
		})
	}
	return &entity.GetCommentResponse{
		Comment: resp,
	}, nil
}

func (c *commentService) CreateCommentRule(ctx context.Context, req *entity.CreateCommentRuleRequest) (int64, error) {
	commentRule := domain.CommentRule{
		CommentText: req.Rule.CommentText,
		Application: req.Rule.Application,
		Visible:     req.Rule.Visible,
	}
	return c.commentRuleRepository.AddCommentRule(ctx, commentRule)
}

func (c *commentService) GetCommentRules(ctx context.Context, req *entity.GetCommentRulesRequest) ([]*entity.GetCommentRulesResponse, error) {
	if req.RuleId != 0 {
		commentRules, err := c.commentRuleRepository.GetCommentRuleById(ctx, req.RuleId, req.Application)
		if err != nil {
			return nil, err
		}
		return []*entity.GetCommentRulesResponse{{
			Rule: entity.Rule{
				Id:          commentRules.Id,
				CommentText: commentRules.CommentText,
				Application: commentRules.Application,
				Visible:     commentRules.Visible,
			},
		},
		}, nil
	}
	commentRules, err := c.commentRuleRepository.GetCommentRules(ctx, req.Application, int32(req.PageNumber), int32(req.PageSize))
	if err != nil {
		return nil, err
	}
	var resp []*entity.GetCommentRulesResponse
	for _, commentRule := range commentRules {
		resp = append(resp, &entity.GetCommentRulesResponse{
			Rule: entity.Rule{
				Id:          commentRule.Id,
				CommentText: commentRule.CommentText,
				Application: commentRule.Application,
				Visible:     commentRule.Visible,
			},
		})
	}
	return resp, nil
}

func (c *commentService) UpdateCommentRule(ctx context.Context, req *entity.UpdateCommentRuleRequest) (*entity.UpdateCommentRuleResponse, error) {
	commentRule := domain.CommentRule{
		Id:          req.RuleId,
		CommentText: req.Rule.CommentText,
		Application: req.Rule.Application,
		Visible:     req.Rule.Visible,
	}
	if err := c.commentRuleRepository.UpdateCommentRule(ctx, commentRule); err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *commentService) commentRuleCheck(ctx context.Context, req *entity.CommentRuleCheckRequest) (bool, string, error) {
	if c.shouldUseAI() {
		violation, reason, err := c.checkWithAI(ctx, req)
		// Emit a trace for failed AI attempt so we can debug ingestion even when LLM call fails.
		if c.langfuseClient != nil {
			_ = pkgLangfuse.EmitTrace(c.langfuseClient, ctx, c.cfg.CommentRule.LangfuseTrace.TraceName, map[string]any{
				"comment_text": req.CommentText,
				"application":  req.Application,
			}, map[string]any{
				"error": "test",
			})
		}
		// Log the error for local visibility and continue with trie fallback.
		fmt.Printf("ai moderation call failed: %v; falling back to trie\n", err)
		if err == nil {
			if violation {
				return true, reason, nil
			}
			return false, "", nil
		}
	}

	if c.trieRoot == nil {
		trieRoot, err := c.buildRuleTrie(ctx, req)
		if err != nil {
			return false, "", err
		}
		c.trieRoot = trieRoot
	}
	if c.trieRoot == nil {
		return false, "", nil
	}
	isViolent, violentWord := c.trieRoot.SearchSubstring(req.CommentText)
	if isViolent {
		return true, violentWord, nil
	}
	return false, "", nil
}

func (c *commentService) shouldUseAI() bool {
	return strings.EqualFold(c.cfg.CommentRule.Mode, "ai") && c.cfg.CommentRule.AI.Enabled && c.cfg.CommentRule.AI.Endpoint != ""
}

func (c *commentService) checkWithAI(ctx context.Context, req *entity.CommentRuleCheckRequest) (bool, string, error) {

	var prompt string
	if c.langfuseClient != nil {
		// p, err := c.langfuseClient.GetPrompt(ctx, c.cfg.CommentRule.AI.PromptName)
		// if err != nil {
		// fallback to project+label lookup if name lookup fails
		if p2, err2 := c.langfuseClient.GetPromptByProjectAndName(ctx, c.cfg.CommentRule.AI.ProjectName, c.cfg.CommentRule.AI.PromptName); err2 == nil {
			prompt = p2
		} else {
			prompt = "You are a helpful assistant."
		}
		// } else {
		// 	prompt = p
		// }
	}

	// call LLM (OpenRouter-compatible) to evaluate the comment
	payload := map[string]any{
		"model": c.cfg.CommentRule.AI.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": prompt,
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("%s", req.CommentText),
			},
		},
		"temperature": 0.7,
		"max_tokens":  -1,
		"stream":      false,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return false, "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.CommentRule.AI.Endpoint, bytes.NewReader(body))
	if err != nil {
		return false, "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.cfg.CommentRule.AI.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.cfg.CommentRule.AI.APIKey)
	}

	// Use langfuse client's DoAndTrace to ensure every LLM query is logged and traced.
	status, respBody, err := c.langfuseClient.DoAndTrace(ctx, c.cfg.CommentRule.LangfuseTrace.TraceName, httpReq, map[string]any{
		"comment_text": req.CommentText,
		"model":        c.cfg.CommentRule.AI.Model,
		"prompt_name":  c.cfg.CommentRule.AI.PromptName,
	})
	if err != nil {
		return false, "", err
	}
	if status < 200 || status >= 300 {
		return false, "", fmt.Errorf("ai moderation request failed: %s", string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return false, "", err
	}
	var aiResult struct {
		Violation bool   `json:"violation"`
		Reason    string `json:"reason"`
	}

	// If the model returned no choices or empty content, treat as violation.
	if len(result.Choices) == 0 {
		aiResult.Violation = true
		aiResult.Reason = "ai moderation returned no choices"
	} else {
		content := result.Choices[0].Message.Content
		if strings.TrimSpace(content) == "" {
			aiResult.Violation = true
			aiResult.Reason = "ai moderation response empty"
		} else {
			if err := json.Unmarshal([]byte(content), &aiResult); err != nil {
				// Fallback: the model didn't return strict JSON. Apply a simple rule:
				// if the model output contains the word "test" (case-insensitive) it's a violation.
				lc := strings.ToLower(content)
				aiResult.Violation = strings.Contains(lc, "true")
				if aiResult.Violation {
					aiResult.Reason = "contains keyword 'test'"
				} else {
					aiResult.Reason = "no violation detected"
				}
			}
		}
	}

	if c.langfuseClient != nil {
		_ = pkgLangfuse.EmitTrace(c.langfuseClient, ctx, c.cfg.CommentRule.LangfuseTrace.TraceName, map[string]any{
			"comment_text": req.CommentText,
			"application":  req.Application,
		}, map[string]any{
			"violation": aiResult.Violation,
			"reason":    aiResult.Reason,
		})
	}

	return aiResult.Violation, aiResult.Reason, nil
}

func (c *commentService) buildRuleTrie(ctx context.Context, req *entity.CommentRuleCheckRequest) (*pkg.Trie, error) {
	totalRules, err := c.commentRuleRepository.CountCommentRules(ctx, req.Application)
	if err != nil {
		return nil, err
	}
	if totalRules == 0 {
		return nil, nil
	}
	trie := pkg.NewTrie()
	for i := 0; i < int(totalRules/20+1); i++ {
		commentRules, err := c.commentRuleRepository.GetCommentRules(ctx, req.Application, int32(1000*i), 1000)
		if err != nil {
			return nil, err
		}
		for _, commentRule := range commentRules {
			trie.Insert(commentRule.CommentText)
		}
	}
	return trie, nil
}
