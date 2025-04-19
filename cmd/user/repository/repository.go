package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository struct {
	Database *gorm.DB
	Redis    *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{
		Redis:    redis,
		Database: db,
	}
}
