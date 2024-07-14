package logic

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/post-service/internal/domain"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/mq/producer"
)

type CommentService interface {
	Comment(context.Context, *postv1.CreateCommentRequest) (*postv1.CreateCommentResponse, error)
	GetComment(context.Context, *postv1.GetCommentRequest) (*postv1.GetCommentResponse, error)
}

type commentService struct {
	db *goqu.Database
}

func NewCommentService(
	db *goqu.Database,
	kafkaClient producer.Client,
) CommentService {
	return &commentService{
		db: db,
	}
}

func (c commentService) Comment(
	ctx context.Context,
	req *postv1.CreateCommentRequest,
) (*postv1.CreateCommentResponse, error) {

	comment := domain.Comment{
		PostId:        req.PostId,
		UserId:        req.Comment.UserId,
		CommentText:   req.Comment.CommentText,
		CommentLevel:  int(req.Comment.CommentLevel),
		CommentParent: uint64(req.Comment.CommentParent),
		Images:        strings.Join(req.Comment.Images, ", "),
		Tags:          strings.Join(req.Comment.Tags, ", "),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	query := c.db.
		Insert(domain.TabNameComment).
		Rows(comment)
	result, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &postv1.CreateCommentResponse{
		CommentId: uint64(id),
	}, nil
}

func (c commentService) GetComment(
	ctx context.Context,
	req *postv1.GetCommentRequest,
) (*postv1.GetCommentResponse, error) {

	// Create the raw SQL query
	sql := `
	WITH top_comments AS (
		SELECT id
		FROM comment
		WHERE comment_level = 0
			AND post_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	)
	SELECT
		sc.id AS sc_id,
		sc.user_id AS sc_user_id,
		sc.comment_text AS sc_comment_text,
		sc.comment_level AS sc_comment_level,
		sc.comment_parent AS sc_comment_parent,
		sc.images AS sc_images,
		sc.tags AS sc_tags,
		sc.created_at AS sc_created_at
	FROM top_comments tc
	LEFT JOIN comment sc ON sc.comment_parent = tc.id or sc.id = tc.id
	ORDER BY sc.comment_parent ASC, sc.created_at DESC
	`
	offset := (req.PageNumber - 1) * req.PageSize
	// Execute the raw SQL query with the parameters
	rows, err := c.db.Query(sql, req.PostId, req.PageSize, offset)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	var commentRes []*postv1.Comment
	// Process the results
	for rows.Next() {
		var scId uint64
		var scUserId string
		var scCommentText string
		var scCommentLevel int
		var scCommentParent uint64
		var scImages string
		var scTags string
		var scCreatedAt time.Time

		err := rows.Scan(&scId, &scUserId, &scCommentText, &scCommentLevel, &scCommentParent, &scImages, &scTags, &scCreatedAt)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		element := &postv1.Comment{
			CommentId:     scId,
			UserId:        scUserId,
			CommentText:   scCommentText,
			CommentLevel:  uint32(scCommentLevel),
			CommentParent: uint64(scCommentParent),
			Images:        []string{scImages},
			Tags:          []string{scTags},
			Date:          scCreatedAt.String(),
		}
		commentRes = append(commentRes, element)
	}

	return &postv1.GetCommentResponse{
		Comment: commentRes,
	}, nil
}
