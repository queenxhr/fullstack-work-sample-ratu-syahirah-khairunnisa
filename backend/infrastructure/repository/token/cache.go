package token

import (
	"context"
	"sync"
	"time"

	"fullstack-backend/internal/repository"
)

// InMemoryRedisRevocationStore simulates a Redis Blacklist with TTL
// Key: jti (JWT ID), Value: Expiration Time
type inMemoryRedisRevocationStore struct {
	mu        sync.RWMutex
	blacklist map[string]time.Time
}

func NewInMemoryRevocationStore() repository.TokenRepository {
	store := &inMemoryRedisRevocationStore{
		blacklist: make(map[string]time.Time),
	}
	// Background garbage collector for expired tokens (mimicking Redis automatic TTL eviction)
	go store.startEvictionLoop()
	return store
}

func (s *inMemoryRedisRevocationStore) RevokeToken(ctx context.Context, jti string, expSeconds int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiresAt := time.Now().Add(time.Duration(expSeconds) * time.Second)
	s.blacklist[jti] = expiresAt
	return nil
}

func (s *inMemoryRedisRevocationStore) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exp, exists := s.blacklist[jti]
	if !exists {
		return false, nil
	}
	// If already past expiry, treat as naturally expired
	if time.Now().After(exp) {
		return false, nil
	}
	return true, nil
}

func (s *inMemoryRedisRevocationStore) ListRevokedTokens(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	var tokens []string
	for jti, exp := range s.blacklist {
		if now.Before(exp) {
			tokens = append(tokens, jti)
		}
	}
	return tokens, nil
}

func (s *inMemoryRedisRevocationStore) startEvictionLoop() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for jti, exp := range s.blacklist {
			if now.After(exp) {
				delete(s.blacklist, jti)
			}
		}
		s.mu.Unlock()
	}
}
