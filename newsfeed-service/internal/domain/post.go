package domain

import (
	"encoding/json"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNamePost  = goqu.T("post")
	TabColUserId = "user_id"
)

type Post struct {
	Id           uint64    `db:"id" goqu:"omitnil"`
	User_id      string    `db:"user_id" goqu:"omitnil"`
	Content_text string    `db:"content_text" goqu:"omitnil"`
	Images       string    `db:"images" goqu:"omitnil"`
	Tags         string    `db:"tags" goqu:"omitnil"`
	Visible      bool      `db:"visible" goqu:"omitnil"`
	CreatedAt    time.Time `db:"created_at" goqu:"omitnil"`
	UpdatedAt    time.Time `db:"updated_at" goqu:"omitnil"`
}

func (p Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	alias := Alias(p)

	if alias.Content_text == "" {
		alias.Content_text = " "
	}
	if alias.Images == "" {
		alias.Images = " "
	}
	if alias.Tags == "" {
		alias.Tags = " "
	}

	return json.Marshal(alias)
}
