package jwt

import (
	"errors"
	"time"

	"webGL-720yun/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	RoleID   uint   `json:"role_id"`
	IsSuper  bool   `json:"is_super"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret         string
	accessTimeout  time.Duration
	refreshTimeout time.Duration
	issuer         string
}

func NewJWTService(cfg *config.JWTConfig) *JWTService {
	return &JWTService{
		secret:         cfg.Secret,
		accessTimeout:  time.Duration(cfg.AccessTimeout) * time.Minute,
		refreshTimeout: time.Duration(cfg.RefreshTimeout) * time.Hour,
		issuer:         cfg.Issuer,
	}
}

func (s *JWTService) GenerateAccessToken(userID uint, username string, role string, roleID uint, isSuper bool) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RoleID:   roleID,
		IsSuper:  isSuper,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTimeout)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) GenerateRefreshToken(userID uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTimeout)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) ValidateToken(tokenString string) error {
	_, err := s.ParseToken(tokenString)
	return err
}

func (s *JWTService) GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims.ExpiresAt == nil {
		return 0, errors.New("token has no expiration time")
	}

	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime < 0 {
		return 0, errors.New("token has expired")
	}

	return remainingTime, nil
}
