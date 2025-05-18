package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/entity"
	"github.com/namnv2496/post-service/internal/pkg"
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
}

func NewCommentService(
	commentRepository repository.ICommentRepository,
	commentRuleRepository repository.ICommentRuleRepository,
	trie pkg.ITrie,
	kafkaClient producer.Client,
) ICommentService {
	return &commentService{
		commentRepository:     commentRepository,
		commentRuleRepository: commentRuleRepository,
		trie:                  trie,
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
	// can cache this trie
	trieRoot, err := c.buildRuleTrie(ctx, req)
	if err != nil {
		return false, "", err
	}
	isViolent, violentWord := trieRoot.SearchSubstring(req.CommentText)
	if isViolent {
		return true, violentWord, nil
	}
	return false, "", nil
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
