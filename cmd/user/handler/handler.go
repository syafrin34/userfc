package handler

import (
	"net/http"
	"userfc/cmd/user/usecase"
	"userfc/infrastructure/logger"
	"userfc/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUseCase usecase.UserUsecase
}

func NewUserHandler(userUseCase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
	}
}
func (h *UserHandler) Register(c *gin.Context) {
	var param models.RegisterParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Logger.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid input parameter",
		})

		return
	}

	if len(param.Password) < 8 || len(param.ConfirmPassword) < 8 {
		logger.Logger.Info("Invalid Input")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Password must longer than 8 characters",
		})

		return
	}

	if param.Password != param.ConfirmPassword {
		logger.Logger.Info("Invalid Input")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Password and Confirm Password Not Match",
		})
		return
	}

	//cek status user apa sudah ada
	user, err := h.UserUseCase.GetUserByEmail(c.Request.Context(), param.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return

	}

	if user.ID != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "email already exists",
		})
		return
	}

	err = h.UserUseCase.RegisteUser(c.Request.Context(), &models.User{
		Name:     param.Name,
		Email:    param.Email,
		Password: param.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error meesgae": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User Successfully Registered",
	})

}

func (h *UserHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
func (h *UserHandler) Login(c *gin.Context) {
	var param models.LoginParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Logger.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_mmessage": "Invalid input parameter",
		})

	}

	if len(param.Password) < 8 {
		logger.Logger.Info("Invalid Input")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Password musst longer than 8 characters",
		})
		return
	}

	token, err := h.UserUseCase.Login(c.Request.Context(), &param)
	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email atau Password salah",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// extract user id
	userIDStr, isExists := c.Get("user_id")
	if !isExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unautorized 1",
		})
		return
	}

	userID, ok := userIDStr.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user id",
		})
	}

	user, err := h.UserUseCase.GetUserID(c.Request.Context(), int64(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error message ": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"email": user.Email,
	})

}
