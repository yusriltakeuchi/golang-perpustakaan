package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yusriltakeuchi/golang-perpustakaan/domain"
	"github.com/yusriltakeuchi/golang-perpustakaan/dto"
	"github.com/yusriltakeuchi/golang-perpustakaan/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	conf           *config.Config
	userRepository domain.UserRepository
}

func NewAuth(cnf *config.Config, userRepository domain.UserRepository) domain.AuthService {
	return &authService{
		conf:           cnf,
		userRepository: userRepository,
	}
}

// Login implements [domain.AuthService].
func (a *authService) Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if user.Id == "" {
		return dto.AuthResponse{}, errors.New("authentication failed")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("authentication failed")
	}

	// Generate JWT Token
	claim := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	return dto.AuthResponse{
		Token: tokenStr,
	}, nil
}
