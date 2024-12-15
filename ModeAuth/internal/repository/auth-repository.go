package repository

import (
	"ModeAuth/internal/shared/dto"
	"ModeAuth/pkg/logging"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"time"
)

type Repository struct {
	DB  *gorm.DB
	rDB *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{DB: db, rDB: rdb}
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (*dto.User, error) {
	var user dto.User

	result := r.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			return nil, fmt.Errorf("Пользователь не найден в базе данных: %w ", result.Error)
		}

		return nil, fmt.Errorf("Database query failed: %w ", result.Error)
	}

	return &user, nil
}

func (r *Repository) IsUserBlocked(ctx context.Context, userID string) (*dto.BlockedUser, error) {
	blockKey := fmt.Sprintf("blocked_user:%s", userID)
	log.Printf(logging.INFO + "BlockKey has been Created")

	var userBlocked dto.BlockedUser

	data, err := r.rDB.Get(ctx, blockKey).Result()
	if err == nil {

		if err := json.Unmarshal([]byte(data), &userBlocked); err != nil {
			log.Printf(logging.ERROR+"[Repository] json.Unmarshal failed: %v", err)

			return nil, fmt.Errorf("не удалось десериализовать объект в JSON. err: %w", err)
		}

		return &userBlocked, nil

	} else if err != redis.Nil {
		return nil, fmt.Errorf("ошибка чтения из Redis: %w", err)
	}

	if err := r.DB.Where("id = ?", userID).First(&userBlocked).Error; err != nil {
		return nil, err
	}

	if err := r.UpdateRedisCache(ctx, blockKey, userBlocked, 5*time.Second); err != nil {
		return nil, err
	}

	return &userBlocked, nil
}

func (r *Repository) UpdateRedisCache(ctx context.Context, key string, object any, ttl time.Duration) error {
	data, err := json.Marshal(object)
	if err != nil {
		log.Printf(logging.ERROR+"[Repository] json.Unmarshal failed: %v", err)

		return fmt.Errorf("не удалось сериализовать объект из JSON. err: %w", err)
	}

	if err := r.rDB.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("не удалось сохранить данные в кэш Redis. err: %w", err)
	}

	return nil
}

func (r *Repository) UnBlockUser(ctx context.Context, userID string) error {
	_, err := r.IsUserBlocked(ctx, userID)
	if err != nil {
		return err
	}

	if err := r.DB.Delete(&dto.BlockedUser{}, userID).Error; err != nil {
		return err
	}

	log.Printf(logging.INFO + "BlockKey has been Created")
	blockKey := fmt.Sprintf("blocked_user:%s", userID)

	if err := r.rDB.Del(ctx, blockKey).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) ChangeUserName(ctx context.Context, user *dto.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}
