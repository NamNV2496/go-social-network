package grpc

import (
	"context"
	"fmt"

	userv1 "github.com/namnv2496/user-service/internal/handler/generated/user_core/v1"
	"github.com/namnv2496/user-service/internal/logic"
)

type GrpcHandler struct {
	userv1.UnimplementedAccountServiceServer
	userService logic.UserService
}

func NewGrpcHander(
	userService logic.UserService,
) userv1.AccountServiceServer {
	return &GrpcHandler{
		userService: userService,
	}
}

func (s GrpcHandler) CreateAccount(
	ctx context.Context,
	in *userv1.CreateAccountRequest,
) (*userv1.CreateAccountResponse, error) {

	fmt.Println("called createAccount")
	return &userv1.CreateAccountResponse{
		Account: &userv1.Account{},
	}, nil
}

func (s GrpcHandler) GetAccount(
	ctx context.Context,
	req *userv1.GetAccountRequest,
) (*userv1.GetAccountResponse, error) {

	user, err := s.userService.GetAccount(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &userv1.GetAccountResponse{
		Account: &userv1.Account{
			Id:        user.Id,
			Email:     user.Email,
			Name:      user.Name,
			Picture:   user.Picture,
			UserId:    user.UserId,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		},
	}, nil
}

func (s GrpcHandler) CreateSession(
	ctx context.Context,
	req *userv1.CreateSessionRequest,
) (*userv1.CreateSessionResponse, error) {

	token, err := s.userService.Login(ctx, req.UserId, req.Password)
	if err != nil {
		return nil, err
	}
	return &userv1.CreateSessionResponse{
		UserId: req.UserId,
		Token:  token,
	}, nil
}

func (s GrpcHandler) GetFollowing(
	ctx context.Context,
	req *userv1.GetFollowingRequest,
) (*userv1.GetFollowingResponse, error) {

	followings, err := s.userService.GetFollowing(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &userv1.GetFollowingResponse{
		UserId: followings,
	}, nil
}

func (s GrpcHandler) CreateFollowing(
	ctx context.Context,
	req *userv1.CheckFollowingRequest,
) (*userv1.CheckFollowingResponse, error) {

	result, err := s.userService.CreateFollowing(ctx, req.CurrentId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &userv1.CheckFollowingResponse{
		Following: result,
	}, nil
}

func (s GrpcHandler) CheckFollowing(
	ctx context.Context,
	req *userv1.CheckFollowingRequest,
) (*userv1.CheckFollowingResponse, error) {

	result, err := s.userService.CheckFollowing(ctx, req.CurrentId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &userv1.CheckFollowingResponse{
		Following: result,
	}, nil
}

func (s GrpcHandler) DeleteFollowing(
	ctx context.Context,
	req *userv1.CheckFollowingRequest,
) (*userv1.CheckFollowingResponse, error) {

	err := s.userService.Unfollowing(ctx, req.CurrentId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &userv1.CheckFollowingResponse{
		Following: false,
	}, nil
}
