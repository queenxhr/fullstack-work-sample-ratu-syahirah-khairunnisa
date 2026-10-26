package auth

import (
	"net/http"

	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/usecase"
	"fullstack-backend/utils"
	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	ListBlacklist(ctx *gin.Context)
}

type AuthHandler struct {
	authUc usecase.IAuthUseCase
}

func NewHandler(authUc usecase.IAuthUseCase) IAuthHandler {
	return &AuthHandler{
		authUc: authUc,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var reqBody LoginRequestDto
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, err.Error()))
		return
	}

	token, claims, err := h.authUc.Login(ctx.Request.Context(), reqBody.Email, reqBody.Password)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, gin.H{
		"token":  token,
		"claims": claims,
	}, "Login berhasil! JWT token diterbitkan dengan JTI unik.")
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	tokenString, exists := ctx.Get("jwt_token")
	if !exists {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrUnauthorized, "Token not found in context"))
		return
	}

	err := h.authUc.Logout(ctx.Request.Context(), tokenString.(string))
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, nil, "Logout berhasil! Token JTI telah dimasukkan ke Redis Revocation Blacklist dengan TTL otomatis.")
}

func (h *AuthHandler) GetProfile(ctx *gin.Context) {
	claims, exists := ctx.Get("user_claims")
	if !exists {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrUnauthorized, "Unauthorized"))
		return
	}

	utils.JSON(ctx, http.StatusOK, claims.(*entity.TokenClaims), "Protected profile accessed successfully")
}

func (h *AuthHandler) ListBlacklist(ctx *gin.Context) {
	blacklisted, err := h.authUc.ListBlacklistedTokens(ctx.Request.Context())
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, gin.H{
		"revoked_count": len(blacklisted),
		"revoked_jtis":  blacklisted,
	}, "Daftar JWT Token yang sedang aktif di Blacklist (menunggu natural expiry)")
}
