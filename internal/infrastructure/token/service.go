package token

import "time"

// TokenService centralizes access/refresh generation and validation
type TokenService struct {
	AccessMaker  *JWTMaker
	RefreshMaker *JWTMaker
	// optional: persistence for refresh tokens can be added later
}

func NewTokenService(accessSecret string, accessTTL time.Duration, refreshSecret string, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		AccessMaker:  &JWTMaker{Secret: accessSecret, TTL: accessTTL},
		RefreshMaker: &JWTMaker{Secret: refreshSecret, TTL: refreshTTL},
	}
}

func (s *TokenService) GenerateAccess(userID int64) (string, error) {
	return s.AccessMaker.Generate(userID)
}

func (s *TokenService) GenerateRefresh(userID int64) (string, error) {
	return s.RefreshMaker.Generate(userID)
}

func (s *TokenService) ValidateAccess(token string) (int64, error) {
	return s.AccessMaker.Validate(token)
}

func (s *TokenService) ValidateRefresh(token string) (int64, error) {
	return s.RefreshMaker.Validate(token)
}
