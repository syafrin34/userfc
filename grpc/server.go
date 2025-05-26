package grpc

import (
	"context"
	"userfc/cmd/user/usecase"
	"userfc/proto/userpb"
)

type GRPCServer struct {
	userpb.UnimplementedUserServiceServer
	UserUsecase usecase.UserUsecase
}

func (s *GRPCServer) GetUserInfoByUserID(ctx context.Context, req *userpb.GetUserInfoRequest) (*userpb.GetUserInfoResult, error) {
	userInfo, err := s.UserUsecase.GetUserID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserInfoResult{
		Id:    userInfo.ID,
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Role:  userInfo.Role,
	}, nil
}
