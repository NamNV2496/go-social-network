package entity

type Rule struct {
	Id          int64  `json:"id,omitempty"`
	Application string `json:"application,omitempty"`
	CommentText string `json:"comment_text,omitempty"`
	Visible     bool   `json:"visible,omitempty"`
}

type CreateCommentRuleRequest struct {
	Rule *Rule `json:"rule,omitempty"`
}

type CreateCommentRuleResponse struct {
	RuleId int64  `json:"rule_id,omitempty"`
	Status string `json:"status,omitempty"`
}

type GetCommentRuleByIdRequest struct {
	RuleId      int64    ` json:"rule_id,omitempty"`
	Application []string `json:"application,omitempty"`
	PageNumber  uint32   ` json:"page_number,omitempty"`
	PageSize    uint32   `json:"page_size,omitempty"`
}

type GetCommentRuleByIdResponse struct {
	Rule Rule `json:"rule,omitempty"`
}

type GetCommentRulesRequest struct {
	RuleId      int64    `json:"rule_id,omitempty"`
	Application []string `json:"application,omitempty"`
	PageNumber  uint32   `json:"page_number,omitempty"`
	PageSize    uint32   `json:"page_size,omitempty"`
}

type GetCommentRulesResponse struct {
	Rule Rule `json:"rule,omitempty"`
}

type UpdateCommentRuleRequest struct {
	RuleId int64 ` json:"rule_id,omitempty"`
	Rule   *Rule `json:"rule,omitempty"`
}

type UpdateCommentRuleResponse struct {
	Status string `json:"status,omitempty"`
}

type CommentRuleCheckRequest struct {
	Application []string `json:"application,omitempty"`
	CommentText string   `json:"comment_text,omitempty"`
}
