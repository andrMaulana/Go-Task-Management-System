package service

import (
	"sync"
	"time"
)

type TokenService struct {
	blacklist map[string]time.Time
	mutex     sync.RWMutex
}

func NewTokenService() *TokenService {
	return &TokenService{
		blacklist: make(map[string]time.Time),
	}
}

func (s *TokenService) BlacklistToken(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.blacklist[token] = time.Now().Add(24 * time.Hour) // blacklist for 24 hours
}

func (s *TokenService) IsTokenBlacklisted(token string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	expiry, exists := s.blacklist[token]
	if !exists {
		return false
	}
	if time.Now().After(expiry) {
		delete(s.blacklist, token) // clean up expired tokens
		return false
	}
	return true
}
