package usecase

import (
	"context"
	"errors"
	"time"
	"userfc/cmd/user/service"
	"userfc/infrastructure/logger"
	"userfc/models"
	"userfc/utils"

	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
)

type UserUsecase struct {
	UserService *service.UserService
	JWTSecret   string
}

func NewUserUseCase(userService *service.UserService, jwtSecret string) *UserUsecase {
	return &UserUsecase{
		UserService: userService,
		JWTSecret:   jwtSecret,
	}
}

func (uc *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	tracer := otel.Tracer("userfc-Usecase")
	ctx, span := tracer.Start(ctx, "GetUserByEmail")
	defer span.End()
	user, err := uc.UserService.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}
	return user, nil
}
func (uc *UserUsecase) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	trace := otel.Tracer("userfc-UseCase")
	ctx, span := trace.Start(ctx, "GetUserByID")
	defer span.End()
	user, err := uc.UserService.GetUserByID(ctx, userID)

	if err != nil {
		logger.LogWithTrace(ctx)
		return nil, err
	}
	return user, nil
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, user *models.User) error {
	// hash password
	hashedPassword, err := utils.HashaPassword(user.Password)
	if err != nil {
		//log.Logger.WithFields(logrus.Fields{"email": user.Email}).Errorf("utils hashpassword() got error %v", err)
		return err
	}
	user.Password = hashedPassword
	_, err = uc.UserService.CreateNewUser(ctx, user)
	if err != nil {
		//log.Logger.WithFields(logrus.Fields{"email": user.Email, "name": user.name}).Errorf("uc.UserService.CreatedUser() got error %v", err)
		return err
	}
	return nil

}
func (uc *UserUsecase) Login(ctx context.Context, param models.LoginParameter, userID int64, storedPassword string) (string, error) {
	tracer := otel.Tracer("userfc-usecase")
	ctx, span := tracer.Start(ctx, "Login")
	defer span.End()

	isMatch, err := utils.CheckPasswordHash(storedPassword, param.Password)
	if err != nil {
		logger.Logger.WithFields(logger.Fields{
			"email": param.Email,
		}).Errorf("utils.CheckPasswordHash got error %v", err)
	}
	if !isMatch {
		return "", errors.New("email atau password salah")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"expire":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.JWTSecret))
	if err != nil {
		logger.Logger.WithFields(logger.Fields{
			"email": param.Email,
		}).Errorf("token.SignedString got error %v", err)
		return "", err
	}

	return tokenString, nil

	// user, err := uc.UserService.GetUserByEmail(ctx, param.Email)
	// if err != nil {
	// 	logger.Logger.WithFields(logger.Fields{
	// 		"email": user.Email,
	// 		"name":  user.Name,
	// 	}).Errorf("us.UserService.CreateNewUser() got error %v", err)

	// }
	// if user.ID == 0 {

	// 	return "", errors.New("email not found ")
	// }

	// isMatch, err := utils.CheckPasswordHash(user.Password, param.Password)
	// if err != nil {
	// 	fmt.Println("password", user.Password)
	// 	fmt.Println("password2", param.Password)
	// 	logger.Logger.WithFields(logger.Fields{
	// 		"email": param.Email,
	// 	}).Errorf("utils.CheckPasswordHash got error %v", err)
	// }

	// if !isMatch {
	// 	return "", errors.New("email atau password salah")
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"user_id": user.ID,
	// 	"expire":  time.Now().Add(time.Hour * 1).Unix(),
	// })

	// tokenString, err := token.SignedString([]byte(uc.JWTSecret))
	// if err != nil {
	// 	logger.Logger.WithFields(logger.Fields{
	// 		"email": param.Email,
	// 	}).Errorf("token.SignedString got error %v", err)
	// }
	// return tokenString, nil

}
