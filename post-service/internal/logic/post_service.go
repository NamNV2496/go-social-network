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

type PostService interface {
	Post(context.Context, *postv1.CreatePostRequest) (*postv1.CreatePostResponse, error)
	GetPosts(context.Context, *postv1.GetPostRequest) (*postv1.GetPostResponse, error)
}

type postService struct {
	db          *goqu.Database
	kafkaClient producer.Client
}

func NewPostService(
	db *goqu.Database,
	kafkaClient producer.Client,
) PostService {
	return &postService{
		db:          db,
		kafkaClient: kafkaClient,
	}
}

func (p postService) Post(
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
	query := p.db.
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
	post.Id = uint64(id)
	data, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Error marshall data to send newsfeed")
	}
	go func() {
		fmt.Println("Call trigger to newsFeed a post: ", post)
		if err := p.kafkaClient.Produce(context.Background(), mq.TOPIC_POST_CONTENT, data); err != nil {
			log.Println("Error when send data to kafka: ", err)
			return
		}
	}()

	return &postv1.CreatePostResponse{
		PostId: uint64(id),
	}, nil
}

func (p postService) GetPosts(
	ctx context.Context,
	req *postv1.GetPostRequest,
) (*postv1.GetPostResponse, error) {

	query := p.db.
		From(domain.TabNamePost).
		Where(
			goqu.C(domain.TabColUserId).Eq(req.UserId),
		).
		Order(goqu.I(domain.TabColCreatedAt).Desc())
	var posts []domain.Post
	err := query.Executor().ScanStructsContext(ctx, &posts)
	if err != nil {
		return nil, err
	}
	postRes := make([]*postv1.Post, 0)
	for _, post := range posts {
		element := &postv1.Post{
			PostId:      post.Id,
			UserId:      post.User_id,
			ContentText: post.Content_text,
			Tags:        strings.Split(post.Tags, ","),
			Images:      strings.Split(post.Images, ","),
			Visible:     post.Visible,
			Date:        post.CreatedAt.String(),
		}
		postRes = append(postRes, element)
	}
	return &postv1.GetPostResponse{
		Post: postRes,
	}, nil
}
