package service

import (
	"ModeAuth/internal/shared/dto"
	"context"
)

type IAuthRepository interface {
	GetUserByID(ctx context.Context, userID string) (*dto.User, error)
	IsUserBlocked(ctx context.Context, userID string) (*dto.BlockedUser, error)
	UnBlockUser(ctx context.Context, userID string) error
	ChangeUserName(ctx context.Context, user *dto.User) error
}
