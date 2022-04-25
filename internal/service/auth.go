package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/vbetsun/todo-app/internal/domain"
	"github.com/vbetsun/todo-app/internal/repository"
)

const salt = "123salt123"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(u domain.User) (int, error) {
	u.Password = s.generateHash(u.Password)

	return s.repo.CreateUser(u)
}

func (s *AuthService) GenerateToken(u, pwd string) (string, error) {

	return s.repo.CreateUser(u)
}

func (s *AuthService) generateHash(pwd string) string {
	hash := sha1.New()
	hash.Write([]byte(pwd))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
