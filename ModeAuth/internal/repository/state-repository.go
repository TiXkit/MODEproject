package repository

import (
	"ModeAuth/internal/shared/dto"
	"ModeAuth/pkg/logging"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var (
	AlreadyStateIsSet = errors.New("данное состояние изначально уже является активным для пользователя")
	StateIsNotSet     = errors.New("состояние CheckingState не установлено")
)

func (r *Repository) WorkState(ctx context.Context, userID string) error {
	key := fmt.Sprintf("keyState:%s", userID)
	log.Printf(logging.INFO + "KeyState has been Created")

	userState, err := r.CheckState(ctx, userID)
	if err != nil {
		return err
	}

	if userState.WorkState {
		return AlreadyStateIsSet
	}

	if userState.SendingState {
		if err := r.DeleteSendingState(ctx, userID); err != nil {
			if err == redis.Nil {

				log.Fatal(logging.FATAL+"SendingState, for user %w, is missing in Redis: ", userID, err)
			}

			return err
		}
	}

	userState.WorkState = true
	userState.SendingState = false
	userState.CheckingState = false

	data, err := json.Marshal(userState)
	if err != nil {
		log.Printf(logging.ERROR+"[Repository] json.Marshal failed: %v", err)

		return fmt.Errorf("не удалось сериализовать объект в JSON. err: %w", err)
	}

	if err := r.rDB.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		return err
	}

	log.Printf("для userID:%s, по ключу:%s. Установлено новое состояние - WorkState(true).", userID, key)

	return nil
}

func (r *Repository) SendingState(ctx context.Context, userID string) error {
	key := fmt.Sprintf("keyState:%s", userID)
	log.Printf(logging.INFO + "KeyState has been Created")

	userState, err := r.CheckState(ctx, userID)
	if err != nil {
		return err
	}

	if userState.SendingState {
		return AlreadyStateIsSet
	}

	userState.WorkState = false
	userState.SendingState = true
	userState.CheckingState = false

	data, err := json.Marshal(userState)
	if err != nil {
		log.Printf(logging.ERROR+"[Repository] json.Unmarshal failed: %v", err)

		return fmt.Errorf("не удалось сериализовать объект в JSON. err: %w", err)
	}

	if err := r.rDB.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteSendingState(ctx context.Context, userID string) error {
	key := fmt.Sprintf("keySendingState:%s", userID)

	if err := r.rDB.Del(ctx, key).Err(); err != nil {
		if err == redis.Nil {
			return err
		}
		return fmt.Errorf("ошибка при попытке удаления состояния по ключу, из Redis. err:%w", err)
	}

	return nil
}

func (r *Repository) CheckingState(ctx context.Context, userID string) error {
	key := fmt.Sprintf("keyState:%s", userID)
	log.Printf(logging.INFO + "KeyState has been Created")

	userState, err := r.CheckState(ctx, userID)
	if err != nil {
		return err
	}

	if userState.CheckingState {
		return AlreadyStateIsSet
	}

	if userState.SendingState {
		if err := r.DeleteSendingState(ctx, userID); err != nil {
			if err == redis.Nil {

				log.Fatal(logging.FATAL+"SendingState, for user %w, is missing in Redis: ", userID, err)
			}

			return err
		}
	}

	userState.WorkState = false
	userState.SendingState = false
	userState.CheckingState = true

	data, err := json.Marshal(userState)
	if err != nil {
		log.Printf(logging.ERROR+"[Repository] json.Marshal failed: %v", err)

		return fmt.Errorf("не удалось сериализовать объект в JSON. err: %w", err)
	}

	if err := r.rDB.Set(ctx, key, string(data), 24*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetReport(ctx context.Context, userID string) (*dto.Report, error) {
	userTaken, err := r.CheckState(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !userTaken.CheckingState {
		return nil, StateIsNotSet
	}

	var report dto.Report

	result := r.DB.Where("status = ?", false).First(&report)
	if result.Error != nil {
		return nil, err
	}

	return &report, nil
}

func (r *Repository) CheckState(ctx context.Context, userID string) (*dto.UserState, error) {
	key := fmt.Sprintf("keyState:%s", userID)

	var userState dto.UserState

	dataR, err := r.rDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return &userState, nil
		}
		return nil, fmt.Errorf("ошибка при попытке получения состояния пользователя по ключу, из Redis. err:%w", err)
	}

	if err := json.Unmarshal([]byte(dataR), &userState); err != nil {
		log.Printf(logging.ERROR+"[Repository] json.Unmarshal failed: %v", err)

		return nil, fmt.Errorf("не удалось сериализовать объект из JSON. err: %w", err)
	}

	return &userState, nil
}
