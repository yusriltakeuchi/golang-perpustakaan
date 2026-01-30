package domain

import (
	"context"

	"github.com/yusriltakeuchi/golang-perpustakaan/dto"
)

type AuthService interface {
	Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error)
}
