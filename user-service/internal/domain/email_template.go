package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameEmailTemplate = goqu.T("email_template")
)

type EmailTemplate struct {
	Id         int64     `db:"id" goqu:"omitnil"`
	TemplateId string    `db:"template_id" goqu:"omitnil"`
	Template   string    `db:"template" goqu:"omitnil"`
	CreatedAt  time.Time `db:"created_at" goqu:"omitnil"`
	UpdatedAt  time.Time `db:"updated_at" goqu:"omitnil"`
}
