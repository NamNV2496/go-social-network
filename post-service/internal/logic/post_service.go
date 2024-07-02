package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/post-service/internal/domain"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/mq"
	"github.com/namnv2496/post-service/internal/mq/producer"
)

type UserService interface {
	Post(context.Context, *postv1.CreatePostRequest) (*postv1.CreatePostResponse, error)
}

type userService struct {
	db          *goqu.Database
	kafkaClient producer.Client
}

func NewUserService(
	db *goqu.Database,
	kafkaClient producer.Client,
) UserService {
	return &userService{
		db:          db,
		kafkaClient: kafkaClient,
	}
}

func (u userService) Post(
	ctx context.Context,
	req *postv1.CreatePostRequest,
) (*postv1.CreatePostResponse, error) {

	post := domain.Post{
		User_id:      req.Post.UserId,
		Content_text: req.Post.ContentText,
		Images:       strings.Join(req.Post.Images, ", "),
		Tags:         strings.Join(req.Post.Tags, ", "),
		Visible:      req.Post.Visible,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	query := u.db.
		Insert(domain.TabNamePost).
		Rows(post)

	result, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// public to kafka
	data, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Error marshall data to send newsfeed")
	}
	go func() {
		fmt.Println("Call trigger to newsFeed a post: ", post)
		if err := u.kafkaClient.Produce(context.Background(), mq.TOPIC_POST_CONTENT, data); err != nil {
			log.Println("Error when send data to kafka: ", err)
			return
		}
	}()

	return &postv1.CreatePostResponse{
		PostId: uint64(id),
	}, nil
}
