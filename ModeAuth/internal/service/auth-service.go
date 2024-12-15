package service

import (
	errors2 "ModeAuth/internal/service/errors"
	"ModeAuth/internal/shared/dto"
	"ModeAuth/internal/shared/utils"
	"ModeAuth/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"time"
)

var _ IAuth = (*Authentication)(nil)

type IAuth interface {
	CheckUser(ctx context.Context, userID, userName string) (*dto.UserTransport, error)
	CheckUserIsBlocked(ctx context.Context, userID string) (*dto.DurationTransport, error)
}

type Authentication struct {
	repo IAuthRepository
}

func NewAuthentication(repo IAuthRepository) *Authentication {
	return &Authentication{repo: repo}
}

func (auth *Authentication) CheckUser(ctx context.Context, userID, userName string) (*dto.UserTransport, error) {
	log.Printf(logging.INFO+"Starting CheckUser for userID: %s\n", userID)

	log.Printf(logging.INFO+"Starting Checking if a user %s is a bot\n", userID)
	result, err := utils.UserIsBot(userID)
	if err != nil {
		log.Printf(logging.WARN+"Failed to check if user %s is a bot: %v\n", userID, err)
		return nil, err
	}
	if result {
		log.Printf(logging.INFO+"User %s is bot", userID)
		return nil, errors2.UserIsBot
	}
	log.Printf(logging.INFO+"User %s is not a bot", userID)

	var user dto.UserTransport

	log.Printf(logging.INFO+"[CheckUser] The user %s is searched in the database\n", userID)
	takenUser, err := auth.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			log.Printf(logging.WARN+"User %s not found in the database\n", userID)
			return nil, errors2.UserNotExist
		}

		log.Printf(logging.ERROR+"Failed to get user %s: %v\n", userID, err)
		return nil, err
	}
	log.Printf(logging.INFO+"[CheckUser] User %s was successfully found in the database\n", userID)

	// the method is logged below
	timeToUnlock, err := auth.CheckUserIsBlocked(ctx, userID)
	if err != nil {
		return nil, err
	}

	log.Printf(logging.INFO+"Checking if a user's %s username \"%s\" is new\n", userID, userName)
	if takenUser.UserName != userName {
		if err := auth.repo.ChangeUserName(ctx, takenUser); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf(logging.WARN+"User %s not found in the database\n", userID)

				return nil, err
			}
			log.Printf(logging.ERROR+"Failed to get user %s: %v\n", userID, err)

			return nil, err
		}
		log.Printf(logging.INFO+"The new UserName \"%s\" for the user %s was successfully saved to the database\n", userName, userID)
	}
	log.Printf(logging.INFO+"UserName \"%s\" of user %s is not new\n", userName, userID)

	user.ID = userID
	user.Role = takenUser.Role
	user.UserName = userName
	user.IsBlocked = timeToUnlock.IsBlocked
	user.BlockedEarlier = takenUser.BlockedEarlier
	user.TimeToUnlock = timeToUnlock.Duration
	user.CreatedAt = takenUser.CreatedAt
	user.UpdateAt = takenUser.UpdateAt

	log.Printf(logging.INFO+"CheckUser function. User %s verification completed successfully", userID)

	return &user, nil
}

func (auth *Authentication) CheckUserIsBlocked(ctx context.Context, userID string) (*dto.DurationTransport, error) {
	log.Printf(logging.INFO+"Starting CheckUserIsBlocked for userID: %s\n", userID)
	var durationBlock *dto.DurationTransport

	log.Printf(logging.INFO+"[CheckUserIsBlocked] The user %s is searched in the database\n", userID)
	user, err := auth.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf(logging.WARN+"User %s not found in the database\n", userID)

			return nil, errors2.UserNotExist
		}
		log.Printf(logging.ERROR+"Failed to get user %s: %v\n", userID, err)

		return nil, err
	}
	log.Printf(logging.INFO+"[CheckUserIsBlocked] User %s was successfully found in the database\n", userID)

	log.Printf(logging.INFO+"Getting a blocked user %s from the database\n", userID)
	blockUser, err := auth.repo.IsUserBlocked(ctx, userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf(logging.ERROR+"Error when trying to get a blocked user %s: %v\n", userID, err)

			return nil, fmt.Errorf("ошибка при попытке получить пользователя. error: %w", err)
		}
		log.Printf(logging.INFO+"Failed to retrieve blocked user from database. User %s is not blocked\n", userID)

		durationBlock.IsBlocked = false
		durationBlock.UserID = user.ID
		durationBlock.UserName = user.UserName
		durationBlock.UpdatedAt = user.UpdateAt

		log.Printf(logging.INFO+"User %s lock check successfully completed\n", userID)

		return durationBlock, nil
	}
	log.Printf(logging.INFO+"Blocked user %s received\n", userID)

	remainingTime := blockUser.ExpiresAt.Sub(time.Now())

	log.Printf(logging.INFO+"Checks the length of time before the user %s is unblocked\n", userID)

	if 0 > remainingTime {
		log.Printf(logging.INFO+"User %s block time has expired\n", userID)

		log.Printf(logging.INFO+"User %s unblock process in progress\n", userID)
		if err := auth.repo.UnBlockUser(ctx, userID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf(logging.WARN+"User %s not found in the database\n", userID)

				return nil, errors2.UserNotExist
			}
			if err == redis.Nil {
				log.Printf(logging.ERROR+"Failed to unblock user. User %s not found in the Redis\n", userID)
			}

			log.Printf(logging.ERROR+"Failed to unblock user %s: %v\n", userID, err)

			return nil, fmt.Errorf("ошибка при попытке разблокировать пользователя. error: %w", err)
		}
		log.Printf(logging.INFO+"User %s successfully unblocked\n", userID)

		durationBlock.IsBlocked = false
		durationBlock.UserID = user.ID
		durationBlock.UserName = user.UserName
		durationBlock.UpdatedAt = user.UpdateAt

	} else {

		durationBlock.IsBlocked = true
		durationBlock.UserID = user.ID
		durationBlock.UserName = user.UserName
		durationBlock.UpdatedAt = user.UpdateAt
		durationBlock.Duration = remainingTime

	}

	log.Printf(logging.INFO+"User %s lock duration check completed successfully\n", userID)

	return durationBlock, nil
}
