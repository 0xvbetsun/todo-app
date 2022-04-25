package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vbetsun/todo-app/internal/domain"
	"github.com/vbetsun/todo-app/internal/repository"
)

const (
	salt       = "123salt123"
	signingKey = "123signing123"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}
type TokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(u domain.User) (int, error) {
	u.Password = s.generateHash(u.Password)

	return s.repo.CreateUser(u)
}

func (s *AuthService) GenerateToken(uname, pwd string) (string, error) {
	user, err := s.repo.GetUser(uname, s.generateHash(pwd))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.ID,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, errors.New("claims are not of type *TokenClaims")
	}
	return claims.UserID, nil
}

func (s *AuthService) generateHash(pwd string) string {
	hash := sha1.New()
	hash.Write([]byte(pwd))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
