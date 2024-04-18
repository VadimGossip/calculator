package domain

import "errors"

var ErrUserNotFound = errors.New("user with such credentials not found")
var ErrRefreshTokenExpired = errors.New("refresh token expired")
var ErrEmptyAuthHeader = errors.New("empty auth header")
var ErrInvalidAuthHeader = errors.New("invalid auth header")
var ErrEmptyToken = errors.New("token is empty")
