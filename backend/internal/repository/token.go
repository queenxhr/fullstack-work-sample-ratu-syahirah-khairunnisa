package repository

import (
	"context"
)

type TokenRepository interface {
	RevokeToken(ctx context.Context, jti string, expSeconds int64) error
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
	ListRevokedTokens(ctx context.Context) ([]string, error)
}
