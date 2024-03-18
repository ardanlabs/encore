package usergrp

import (
	"context"
	"errors"

	"encore.dev/middleware"
	"encore.dev/types/uuid"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

type ctxUserKey int

const (
	userIDKey ctxUserKey = iota + 1
	userKey
)

func SetUserID(req middleware.Request, userID uuid.UUID) middleware.Request {
	ctx := context.WithValue(req.Context(), userIDKey, userID)
	return req.WithContext(ctx)
}

func getUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found")
	}

	return v, nil
}

func SetUser(req middleware.Request, usr user.User) middleware.Request {
	ctx := context.WithValue(req.Context(), userKey, usr)
	return req.WithContext(ctx)
}

func getUser(ctx context.Context) (user.User, error) {
	v, ok := ctx.Value(userKey).(user.User)
	if !ok {
		return user.User{}, errors.New("user not found")
	}

	return v, nil
}
