package repository

import (
	"context"
	"errors"

	"github.com/namnv2496/post-service/configs"
	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	AddComment(ctx context.Context, comment domain.Comment, opts ...database.QueryOption) error
	GetComment(ctx context.Context, opts ...database.QueryOption) ([]*domain.Comment, error)
}

type CommentRepository struct {
	database.ICRUDBase[*domain.Comment]
}

func NewCommentRepository(
	db *gorm.DB,
	config configs.Config,
) (*CommentRepository, error) {
	dbConfig := config.Database
	commentRepository := database.NewCRUDBase[*domain.Comment](db)
	if dbConfig.AutoMigrate {
		if err := db.Debug().AutoMigrate(&domain.Comment{}); err != nil {
			return nil, errors.New("cannot create like repository")
		}
	}
	return &CommentRepository{
		ICRUDBase: commentRepository,
	}, nil
}

func (_self *CommentRepository) AddComment(ctx context.Context, comment domain.Comment, opts ...database.QueryOption) error {
	if err := _self.Create(ctx, &comment); err != nil {
		return err
	}
	return nil
}

func (_self *CommentRepository) GetComment(ctx context.Context, opts ...database.QueryOption) ([]*domain.Comment, error) {
	comments, err := _self.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// 	// Create the raw SQL query
// 	sql := `
// 	WITH top_comments AS (
// 		SELECT id
// 		FROM comment
// 		WHERE comment_level = 0
// 			AND post_id = ?
// 		ORDER BY created_at DESC
// 		LIMIT ? OFFSET ?
// 	)
// 	SELECT
// 		sc.id AS sc_id,
// 		sc.user_id AS sc_user_id,
// 		sc.comment_text AS sc_comment_text,
// 		sc.comment_level AS sc_comment_level,
// 		sc.comment_parent AS sc_comment_parent,
// 		sc.images AS sc_images,
// 		sc.tags AS sc_tags,
// 		sc.created_at AS sc_created_at
// 	FROM top_comments tc
// 	LEFT JOIN comment sc ON sc.comment_parent = tc.id or sc.id = tc.id
// 	ORDER BY sc.comment_parent ASC, sc.created_at DESC
// 	`
// 	offset := (req.PageNumber - 1) * req.PageSize
// 	// Execute the raw SQL query with the parameters
// 	rows, err := c.commentRepository.Query(sql, req.PostId, req.PageSize, offset)
// 	if err != nil {
// 		log.Fatalf("Error executing query: %v", err)
// 	}
// 	defer rows.Close()

// 	var commentRes []*postv1.Comment
// 	// Process the results
// 	for rows.Next() {
// 		var scId uint64
// 		var scUserId string
// 		var scCommentText string
// 		var scCommentLevel int
// 		var scCommentParent uint64
// 		var scImages string
// 		var scTags string
// 		var scCreatedAt time.Time

// 		err := rows.Scan(&scId, &scUserId, &scCommentText, &scCommentLevel, &scCommentParent, &scImages, &scTags, &scCreatedAt)
// 		if err != nil {
// 			log.Fatalf("Error scanning row: %v", err)
// 		}

// 		element := &postv1.Comment{
// 			CommentId:     scId,
// 			UserId:        scUserId,
// 			CommentText:   scCommentText,
// 			CommentLevel:  uint32(scCommentLevel),
// 			CommentParent: uint64(scCommentParent),
// 			Images:        []string{scImages},
// 			Tags:          []string{scTags},
// 			Date:          scCreatedAt.String(),
// 		}
// 		commentRes = append(commentRes, element)
// 	}

//	return &postv1.GetCommentResponse{
//		Comment: commentRes,
//	}, nil
// return nil, nil
