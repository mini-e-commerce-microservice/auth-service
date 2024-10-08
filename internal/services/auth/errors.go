package auth

import "errors"

var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidToken = errors.New("invalid token")
var ErrRefreshTokenNotExistsInRedis = errors.New("refresh token not exists in redis")
