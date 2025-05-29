package service

import (
	"context"
	"userfc/cmd/user/repository"
	"userfc/models"
)

type UserService struct {
	UserRepo *repository.UserRepository
	//JWTSecret string
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (svc *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	user, err := svc.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (svc *UserService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	user, err := svc.UserRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (svc *UserService) CreateNewUser(ctx context.Context, user *models.User) (int64, error) {
	userID, err := svc.UserRepo.InsertNewUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
