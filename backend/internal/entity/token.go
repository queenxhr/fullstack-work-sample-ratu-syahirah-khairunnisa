package entity

import "time"

type TokenClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	JTI      string `json:"jti"` // JWT Unique Identifier for revocation
	Exp      int64  `json:"exp"`
}

type RevokedToken struct {
	JTI       string    `json:"jti"`
	RevokedAt time.Time `json:"revoked_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
