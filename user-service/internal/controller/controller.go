package controller

import (
	"context"

	"github.com/namnv2496/user-service/internal/service"
	logic "github.com/namnv2496/user-service/internal/service"
	userv1 "github.com/namnv2496/user-service/pkg/user_core/v1"
)

type GrpcHandler struct {
	userv1.UnimplementedAccountServiceServer
	userService service.UserService
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

	id, err := s.userService.CreateAccount(ctx, in.Account)
	if err != nil {
		return &userv1.CreateAccountResponse{
			Id: 0,
		}, err
	}
	return &userv1.CreateAccountResponse{
		Id: id,
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
			Id:        uint64(user.Id),
			UserId:    user.UserId,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		},
	}, nil
}

func (s GrpcHandler) FindAccount(
	ctx context.Context,
	req *userv1.FindAccountRequest,
) (*userv1.FindAccountResponse, error) {
	users, err := s.userService.FindAccount(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	// users, err := s.userService.FindAccountByTemplate(ctx, req.UserId)
	// if err != nil {
	// 	return nil, err
	// }

	res := []*userv1.Account{}
	for _, user := range users {
		account := userv1.Account{
			Id:        uint64(user.Id),
			Email:     user.Email,
			Name:      user.Name,
			Picture:   user.Picture,
			UserId:    user.UserId,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		}
		res = append(res, &account)
	}

	return &userv1.FindAccountResponse{
		Account: res,
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
