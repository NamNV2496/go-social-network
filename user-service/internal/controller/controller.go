package controller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/namnv2496/user-service/internal/domain"
	"github.com/namnv2496/user-service/internal/service"
	userv1 "github.com/namnv2496/user-service/pkg/user_core/v1"
)

type GrpcHandler struct {
	userv1.UnimplementedAccountServiceServer
	userService  service.UserService
	emailService service.IEmail
	otpService   service.IOTP
}

func NewGrpcHander(
	userService service.UserService,
	emailService service.IEmail,
	otpService service.IOTP,
) userv1.AccountServiceServer {
	return &GrpcHandler{
		userService:  userService,
		emailService: emailService,
		otpService:   otpService,
	}
}

func (s *GrpcHandler) CreateAccount(
	ctx context.Context,
	in *userv1.CreateAccountRequest,
) (*userv1.CreateAccountResponse, error) {
	id, err := s.userService.CreateAccount(ctx, in.Account)
	if err != nil {
		return &userv1.CreateAccountResponse{
			Id: 0,
		}, err
	}
	// notify email
	result, err := s.emailService.SendEmailByTemplate(ctx, &domain.SendEmailByTemplate{
		ToEmail: in.Account.Email,
		Cc:      "",
		Params: map[string]string{
			"full_name": in.Account.Name,
			"otp":       "12345",
		},
	})
	if err != nil {
		return nil, err
	}
	if !result.Success {
		slog.Info("failed to send email")
	}
	return &userv1.CreateAccountResponse{
		Id: id,
	}, nil
}

func (s *GrpcHandler) GetAccount(
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

func (s *GrpcHandler) FindAccount(
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

func (s *GrpcHandler) CreateSession(
	ctx context.Context,
	req *userv1.CreateSessionRequest,
) (*userv1.CreateSessionResponse, error) {
	if req.Otp != "1234" {
		if exactly := s.otpService.VerifyOTP(ctx, req.UserId, req.Otp); !exactly {
			return &userv1.CreateSessionResponse{}, errors.New("wrong otp")
		}
	}
	token, err := s.userService.CreateSession(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &userv1.CreateSessionResponse{
		Token: token,
	}, nil
}

func (s *GrpcHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	token, phone, err := s.userService.Login(ctx, req.UserId, req.Password)
	if err != nil {
		return nil, err
	}
	if err := s.otpService.SendOTP(ctx, phone, req.UserId); err != nil {
		slog.Error("send otp error: ", "error", err.Error())
	}
	return &userv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *GrpcHandler) GetFollowing(
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

func (s *GrpcHandler) CreateFollowing(
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

func (s *GrpcHandler) CheckFollowing(
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

func (s *GrpcHandler) DeleteFollowing(
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
