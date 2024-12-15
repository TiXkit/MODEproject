package errors

import "errors"

var (
	UserIsBot             = errors.New("пользователь является ботом")
	UserNotExist          = errors.New("данного пользователя не существует в базе")
	AlreadyStateExist     = errors.New("данное состояние изначально уже является активным для пользователя")
	CheckingStateIsNotSet = errors.New("состояние CheckingState не установлено")

	/*	UserIsNotBlocked = errors.New("данный пользователь не заблокирован")
		UserIsNotExistRedis   = errors.New("данный пользователь не имеет начального состояния, чтобы получить текущее")
	*/
)
