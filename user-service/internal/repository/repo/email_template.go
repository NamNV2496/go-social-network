package repo

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/user-service/internal/domain"
)

type IEmailTemplateRepo interface {
	GetEmailTemplateById(context.Context, *domain.GetEmailTemplateRequest) ([]domain.EmailTemplate, error)
	GetEmailTemplateByTemplateId(context.Context, *domain.GetEmailTemplateByTemplateIdRequest) ([]domain.EmailTemplate, error)
	AddEmailTemplate(ctx context.Context, req *domain.AddEmailTemplateRequest) (string, error)
	UpdateEmailTemplate(ctx context.Context, req *domain.UpdateEmailTemplateRequest) (string, error)
}
type EmailTemplateRepo struct {
	db *goqu.Database
}

func NewEmailTemplateRepository(
	db *goqu.Database,
) IEmailTemplateRepo {
	return &EmailTemplateRepo{
		db: db,
	}
}

func (u *EmailTemplateRepo) GetEmailTemplateById(ctx context.Context, req *domain.GetEmailTemplateRequest) ([]domain.EmailTemplate, error) {
	query := u.db.
		From(domain.TabNameEmailTemplate).
		Where(
			goqu.C(domain.TabColId).Eq(req.Id),
		)
	fmt.Println(query.ToSQL())
	var emailTemplate []domain.EmailTemplate
	err := query.Executor().ScanStructsContext(ctx, &emailTemplate)
	if err != nil {
		return nil, err
	}
	return emailTemplate, nil
}

func (u *EmailTemplateRepo) GetEmailTemplateByTemplateId(ctx context.Context, req *domain.GetEmailTemplateByTemplateIdRequest) ([]domain.EmailTemplate, error) {
	query := u.db.
		From(domain.TabNameEmailTemplate).
		Where(
			goqu.C(domain.TabTemplateId).Eq(req.TemplateId),
		)
	fmt.Println(query.ToSQL())
	var emailTemplate []domain.EmailTemplate
	err := query.Executor().ScanStructsContext(ctx, &emailTemplate)
	if err != nil {
		return nil, err
	}
	return emailTemplate, nil
}

func (u *EmailTemplateRepo) AddEmailTemplate(ctx context.Context, req *domain.AddEmailTemplateRequest) (string, error) {
	insert := u.db.
		Insert(domain.TabNameEmailTemplate).
		Rows(goqu.Record{
			domain.TabColId: req.Template,
		}).
		Returning(goqu.C(domain.TabColId))
	fmt.Println(insert.ToSQL())
	var insertedId uint64
	_, err := insert.Executor().ScanValContext(ctx, &insertedId)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(http.StatusOK), nil
}

func (u *EmailTemplateRepo) UpdateEmailTemplate(ctx context.Context, req *domain.UpdateEmailTemplateRequest) (string, error) {
	update := u.db.
		Update(domain.TabNameEmailTemplate).
		Set(goqu.Record{
			domain.TabColId: req.Template,
		}).
		Where(goqu.C(domain.TabColId).Eq(req.Template.Id))
	fmt.Println(update.ToSQL())
	_, err := update.Executor().ExecContext(ctx)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(http.StatusOK), nil
}
