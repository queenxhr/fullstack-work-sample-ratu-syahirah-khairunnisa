package usecase

import (
	"context"
	"fmt"
	"time"

	"fullstack-backend/config"
	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

type IAuthUseCase interface {
	Login(ctx context.Context, email, password string) (string, *entity.TokenClaims, error)
	Logout(ctx context.Context, tokenString string) error
	ValidateToken(ctx context.Context, tokenString string) (*entity.TokenClaims, error)
	ListBlacklistedTokens(ctx context.Context) ([]string, error)
}

type AuthUseCase struct {
	tokenRepo repository.TokenRepository
	cfg       *config.Config
}

func NewAuthUseCase(tokenRepo repository.TokenRepository, cfg *config.Config) IAuthUseCase {
	return &AuthUseCase{
		tokenRepo: tokenRepo,
		cfg:       cfg,
	}
}

// Login memvalidasi kredensial lalu menerbitkan JWT stateless
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, *entity.TokenClaims, error) {
	if email == "" || password == "" {
		return "", nil, appErrors.New(appErrors.ErrBadRequest, "Email dan password wajib diisi")
	}
	// Kredensial demo untuk assessment (production: lookup ke DB + bcrypt.CompareHashAndPassword)
	if password != "password123" && password != "secret" {
		return "", nil, appErrors.New(appErrors.ErrUnauthorized, "Email atau password salah (demo: gunakan password123)")
	}

	// Terbitkan JWT dengan JTI unik (UUID v4) untuk mendukung revokasi per-token
	jti := uuid.New().String()
	expTime := time.Now().Add(15 * time.Minute) // Short-lived: 15 menit (production best practice)

	claims := CustomClaims{
		UserID: "usr-" + uuid.New().String()[:8],
		Email:  email,
		Role:   "senior_engineer",
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "fullstack-work-sample",
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(uc.cfg.JWTSecret))
	if err != nil {
		return "", nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("gagal menandatangani JWT: %v", err))
	}

	tokenClaims := &entity.TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
		JTI:    claims.JTI,
		Exp:    expTime.Unix(),
	}

	return signedToken, tokenClaims, nil
}

// Logout merevoke token dengan memasukkan JTI ke blacklist (Redis-style in-memory TTL cache)
func (uc *AuthUseCase) Logout(ctx context.Context, tokenString string) error {
	claims, err := uc.parseToken(tokenString)
	if err != nil {
		return err
	}

	remainingSeconds := claims.ExpiresAt.Unix() - time.Now().Unix()
	if remainingSeconds <= 0 {
		return nil
	}

	return uc.tokenRepo.RevokeToken(ctx, claims.JTI, remainingSeconds)
}

// ValidateToken memverifikasi signature JWT dan cek JTI di blacklist
func (uc *AuthUseCase) ValidateToken(ctx context.Context, tokenString string) (*entity.TokenClaims, error) {
	claims, err := uc.parseToken(tokenString)
	if err != nil {
		return nil, appErrors.New(appErrors.ErrUnauthorized, fmt.Sprintf("token tidak valid: %v", err))
	}

	isRevoked, err := uc.tokenRepo.IsTokenRevoked(ctx, claims.JTI)
	if err != nil {
		return nil, appErrors.New(appErrors.ErrInternal, "gagal memeriksa revocation store")
	}
	if isRevoked {
		return nil, appErrors.New(appErrors.ErrUnauthorized, "Token telah dicabut (Revoked / Logged out). Autentikasi ditolak!")
	}

	return &entity.TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
		JTI:    claims.JTI,
		Exp:    claims.ExpiresAt.Unix(),
	}, nil
}

func (uc *AuthUseCase) ListBlacklistedTokens(ctx context.Context) ([]string, error) {
	return uc.tokenRepo.ListRevokedTokens(ctx)
}

func (uc *AuthUseCase) parseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(uc.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
