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

type IStateRepository interface {
	WorkState(ctx context.Context, userID string) error
	SendingState(ctx context.Context, userID string) error
	DeleteSendingState(ctx context.Context, userID string) error
	CheckingState(ctx context.Context, userID string) error
	GetReport(ctx context.Context, userID string) (*dto.Report, error)
	CheckState(ctx context.Context, userID string) (*dto.UserState, error)
}
