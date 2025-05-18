package controller

import (
	"context"

	"github.com/namnv2496/post-service/internal/entity"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg"
)

func (c *Controller) CreateComment(
	ctx context.Context,
	req *postv1.CreateCommentRequest,
) (*postv1.CreateCommentResponse, error) {
	var request entity.CreateCommentRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	result, err := c.commentService.AddComment(ctx, &request)
	if err != nil {
		return nil, err
	}
	var resp postv1.CreateCommentResponse
	if err := pkg.Copy(&resp, result); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Controller) GetComment(
	ctx context.Context,
	req *postv1.GetCommentRequest,
) (*postv1.GetCommentResponse, error) {
	var request entity.GetCommentRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	result, err := c.commentService.GetComment(ctx, &request)
	if err != nil {
		return nil, err
	}
	var resp postv1.GetCommentResponse
	if err := pkg.Copy(&resp, result); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Controller) CreateCommentRule(ctx context.Context, req *postv1.CreateCommentRuleRequest) (*postv1.CreateCommentRuleResponse, error) {
	if req.Rule == nil {
		return nil, nil
	}
	var request entity.CreateCommentRuleRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	id, err := c.commentService.CreateCommentRule(ctx, &request)
	if err != nil {
		return nil, err
	}
	return &postv1.CreateCommentRuleResponse{
		RuleId: id,
		Status: "success",
	}, nil
}

func (c *Controller) GetCommentRuleById(ctx context.Context, req *postv1.GetCommentRuleByIdRequest) (*postv1.GetCommentRuleByIdResponse, error) {
	var request entity.GetCommentRulesRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	result, err := c.commentService.GetCommentRules(ctx, &request)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &postv1.GetCommentRuleByIdResponse{
		Rule: &postv1.Rule{
			Id:          result[0].Rule.Id,
			Application: result[0].Rule.Application,
			CommentText: result[0].Rule.CommentText,
		},
	}, nil
}

func (c *Controller) GetCommentRules(ctx context.Context, req *postv1.GetCommentRulesRequest) (*postv1.GetCommentRulesResponse, error) {
	var request entity.GetCommentRulesRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	result, err := c.commentService.GetCommentRules(ctx, &request)
	if err != nil {
		return nil, err
	}
	var resp []*postv1.Rule
	for _, rule := range result {
		resp = append(resp, &postv1.Rule{
			Id:          rule.Rule.Id,
			Application: rule.Rule.Application,
			CommentText: rule.Rule.CommentText,
		})
	}
	return &postv1.GetCommentRulesResponse{
		Rules: resp,
	}, nil
}

func (c *Controller) UpdateCommentRule(ctx context.Context, req *postv1.UpdateCommentRuleRequest) (*postv1.UpdateCommentRuleResponse, error) {
	var request entity.UpdateCommentRuleRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	_, err := c.commentService.UpdateCommentRule(ctx, &request)
	if err != nil {
		return nil, err
	}
	return &postv1.UpdateCommentRuleResponse{
		Status: "success",
	}, nil
}
