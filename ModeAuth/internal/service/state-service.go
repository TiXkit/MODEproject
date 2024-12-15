package service

import (
	"ModeAuth/internal/repository"
	errors2 "ModeAuth/internal/service/errors"
	"ModeAuth/internal/shared/dto"
	"ModeAuth/pkg/logging"
	"context"
	"errors"
	"log"
)

var _ IStates = (*States)(nil)

type IStates interface {
	WorkState(ctx context.Context, userID string) (bool, error)
	SendingState(ctx context.Context, userID string) (bool, error)
	CheckingState(ctx context.Context, userID string) (bool, error)
	TakeForCheckingState(ctx context.Context, userID string) (*dto.Report, error)
}

type States struct {
	repo IStateRepository
}

func NewState(repo IStateRepository) *States {
	return &States{repo: repo}
}

func (s *States) WorkState(ctx context.Context, userID string) (bool, error) {
	log.Printf(logging.INFO+"Starting WorkState for userID %s\n", userID)

	log.Printf(logging.INFO+"Starting the process of setting the GoWork state for the user %s\n", userID)

	if err := s.repo.WorkState(ctx, userID); err != nil {
		if errors.Is(err, repository.AlreadyStateIsSet) {
			log.Printf(logging.WARN+"User %s already has an active GoWork state: %v\n", userID, err)

			return false, errors2.AlreadyStateExist
		}
		log.Printf(logging.ERROR+"Failed to set GoWork state for user %s: %v\n", userID, err)

		return false, err
	}
	log.Printf(logging.INFO+"GoWork state was successfully set for the user %s\n", userID)

	return true, nil
}

func (s *States) SendingState(ctx context.Context, userID string) (bool, error) {
	log.Printf(logging.INFO+"Starting SendingState for userID %s\n", userID)

	log.Printf(logging.INFO+"Starting the process of setting the Sending state for the user %s\n", userID)

	if err := s.repo.SendingState(ctx, userID); err != nil {
		if errors.Is(err, repository.AlreadyStateIsSet) {
			log.Printf(logging.WARN+"User %s already has an active Sending state: %v\n", userID, err)

			return false, errors2.AlreadyStateExist
		}
		log.Printf(logging.ERROR+"Failed to set Sending state for user %s: %v\n", userID, err)

		return false, err
	}
	log.Printf(logging.INFO+"Sending state was successfully set for the user %s\n", userID)

	return true, nil
}

func (s *States) CheckingState(ctx context.Context, userID string) (bool, error) {
	log.Printf(logging.INFO+"Starting CheckingState for userID %s\n", userID)

	log.Printf(logging.INFO+"Starting the process of setting the Checking state for the user %s\n", userID)

	if err := s.repo.CheckingState(ctx, userID); err != nil {
		if errors.Is(err, repository.AlreadyStateIsSet) {
			log.Printf(logging.WARN+"User %s already has an active Checking state: %v\n", userID, err)

			return false, errors2.AlreadyStateExist
		}
		log.Printf(logging.ERROR+"Failed to set Checking state for user %s: %v\n", userID, err)

		return false, err
	}
	log.Printf(logging.INFO+"Checking state was successfully set for the user %s\n", userID)

	return true, nil
}

func (s *States) TakeForCheckingState(ctx context.Context, userID string) (*dto.Report, error) {
	log.Printf(logging.INFO+"Starting TakeForCheckingState for userID %s\n", userID)

	log.Printf(logging.INFO+"Starting the process of receiving an unverified report for user %s\n", userID)

	report, err := s.repo.GetReport(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.StateIsNotSet) {
			log.Printf(logging.WARN+"Failed to receive report for user %s. The user does not have an active Checking state: %v\n", userID, err)

			return nil, errors2.CheckingStateIsNotSet
		}
		log.Printf(logging.ERROR+"Failed to get report for user %s: %v\n", userID, err)

		return nil, err
	}
	log.Printf(logging.INFO+"The report for user %s was successfully received\n", userID)

	return report, nil
}
