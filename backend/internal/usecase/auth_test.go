package usecase_test

import (
	"context"
	"testing"

	"fullstack-backend/config"
	tokenRepo "fullstack-backend/infrastructure/repository/token"
	"fullstack-backend/internal/usecase"
)

func TestAuthUseCase_Login_And_TokenRevocation(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key-12345",
	}
	repo := tokenRepo.NewInMemoryRevocationStore()
	authUc := usecase.NewAuthUseCase(repo, cfg)

	ctx := context.Background()

	// 1. Test Login dengan kredensial valid
	token, claims, err := authUc.Login(ctx, "ratu.syahirah@example.com", "password123")
	if err != nil {
		t.Fatalf("expected login berhasil, got: %v", err)
	}
	if token == "" || claims.JTI == "" {
		t.Fatalf("expected token dan JTI valid, got empty")
	}
	t.Logf("✅ Login berhasil. Email: %s | Role: %s | JTI: %s", claims.Email, claims.Role, claims.JTI)

	// 2. Test Login dengan password salah (harus gagal 401)
	_, _, err = authUc.Login(ctx, "ratu.syahirah@example.com", "wrongpassword")
	if err == nil {
		t.Fatalf("expected login GAGAL dengan password salah, tapi berhasil!")
	}
	t.Logf("✅ Login dengan password salah berhasil ditolak: %v", err)

	// 3. Test Validasi Token sebelum revokasi (harus berhasil)
	validatedClaims, err := authUc.ValidateToken(ctx, token)
	if err != nil {
		t.Fatalf("expected validasi token berhasil, got: %v", err)
	}
	if validatedClaims.Email != "ratu.syahirah@example.com" {
		t.Fatalf("expected email ratu.syahirah@example.com, got: %s", validatedClaims.Email)
	}
	t.Logf("✅ Token valid sebelum logout: user=%s role=%s", validatedClaims.Email, validatedClaims.Role)

	// 4. Test Logout (Token Revocation via Redis-style JTI Blacklist)
	err = authUc.Logout(ctx, token)
	if err != nil {
		t.Fatalf("expected logout berhasil, got: %v", err)
	}
	t.Logf("✅ Logout berhasil. JTI %s dimasukkan ke revocation blacklist", claims.JTI)

	// 5. Test Validasi Token SETELAH Logout (HARUS DITOLAK dengan 401!)
	_, err = authUc.ValidateToken(ctx, token)
	if err == nil {
		t.Fatalf("expected validasi GAGAL setelah logout, tapi berhasil! → Bug di revocation logic!")
	}
	t.Logf("✅ Token dengan JTI %s BERHASIL DITOLAK setelah logout: %v", claims.JTI, err)

	// 6. Verifikasi JTI ada di daftar blacklist
	blacklist, _ := authUc.ListBlacklistedTokens(ctx)
	found := false
	for _, jti := range blacklist {
		if jti == claims.JTI {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("JTI %s tidak ditemukan di blacklist setelah logout!", claims.JTI)
	}
	t.Logf("✅ JTI ditemukan di blacklist store. Total revoked tokens: %d", len(blacklist))
}
