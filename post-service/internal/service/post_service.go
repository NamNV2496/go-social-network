package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/namnv2496/post-service/internal/domain"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg"
	"github.com/namnv2496/post-service/internal/repository"
	"github.com/namnv2496/post-service/internal/repository/mq"
	"github.com/namnv2496/post-service/internal/repository/mq/producer"
)

type IPostService interface {
	AddPost(context.Context, *postv1.CreatePostRequest) (*postv1.CreatePostResponse, error)
	GetPosts(context.Context, *postv1.GetPostRequest) (*postv1.GetPostResponse, error)
}

type postService struct {
	postRepository repository.IPostRepository
	kafkaClient    producer.Client
}

func NewPostService(
	postRepository repository.IPostRepository,
	kafkaClient producer.Client,
) IPostService {
	return &postService{
		postRepository: postRepository,
		kafkaClient:    kafkaClient,
	}
}

func (p postService) AddPost(
	ctx context.Context,
	req *postv1.CreatePostRequest,
) (*postv1.CreatePostResponse, error) {
	uuid, _ := uuid.GenerateUUID()
	post := domain.Post{
		Uuid:         uuid,
		User_id:      req.Post.UserId,
		Content_text: req.Post.ContentText,
		Images:       strings.Join(req.Post.Images, ","),
		Tags:         strings.Join(req.Post.Tags, ","),
		Visible:      req.Post.Visible,
	}
	// add post
	if err := p.postRepository.AddPost(ctx, post); err != nil {
		return nil, err
	}
	// get post
	posts, err := p.postRepository.GetPosts(ctx, domain.PostByUserId(req.Post.UserId), domain.PostByUuid(uuid), domain.PostOrderById())
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, errors.New("failed to create new post")
	}
	id := posts[0].Id
	// public to kafka
	if err := p.publishNewPost(ctx, *posts[0]); err != nil {
		log.Println("err: ", err)
	}
	return &postv1.CreatePostResponse{
		PostId: uint64(id),
	}, nil
}

func (p postService) GetPosts(
	ctx context.Context,
	req *postv1.GetPostRequest,
) (*postv1.GetPostResponse, error) {
	posts, err := p.postRepository.GetPosts(ctx, domain.PostByUserId(req.UserId), domain.PostOrderById())
	if err != nil {
		return nil, err
	}
	var postRes []*postv1.Post
	for _, post := range posts {
		var postElem *postv1.Post
		if err := pkg.Copy(&postElem, post); err != nil {
			log.Println("error: ", err)
			continue
		}
		postElem.Date = post.CreatedAt.String()
		postElem.Images = strings.Split(post.Images, ",")
		postElem.Tags = strings.Split(post.Tags, ",")
		postRes = append(postRes, postElem)
	}
	return &postv1.GetPostResponse{
		Post: postRes,
	}, nil
}

func (p postService) publishNewPost(ctx context.Context, post domain.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Error marshall data to send newsfeed")
	}
	fmt.Println("Call trigger to newsFeed a post: ", post)
	if err := p.kafkaClient.Produce(ctx, mq.TOPIC_POST_CONTENT, data); err != nil {
		log.Println("Error when send data to kafka: ", err)
		return err
	}
	return nil
}
