package AuthService

import (
	"context"
	"fmt"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/redis"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/errors"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/utils"
	"strconv"
	"time"
)

type (
	AuthRepository interface {
		GetUserByUsername(ctx context.Context, username string) (domain.User, error)
		CreateUser(ctx context.Context, username, passwordHash string, coin int) (domain.User, error)
	}

	TokenManager interface {
		GenerateJWT(userID int) (string, error)
	}

	AuthService struct {
		repository AuthRepository
		tokenMgr   TokenManager
		rdb        *redis.Redis
	}
)

func New(storage AuthRepository, tokenMgr TokenManager, rdb *redis.Redis) *AuthService {
	return &AuthService{repository: storage, tokenMgr: tokenMgr, rdb: rdb}
}

func (u *AuthService) Auth(ctx context.Context, username, password string) (string, error) {
	user, err := u.repository.GetUserByUsername(ctx, username)
	if err != nil {
		passwordHash, err := utils.HashPassword(password)
		if err != nil {
			return "", errors.ErrBadRequest
		}

		user, err = u.repository.CreateUser(ctx, username, passwordHash, 1000)
		if err != nil {
			return "", errors.ErrInternal
		}
	}
	if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
		return "", errors.ErrInvalidPassword
	}
	cToken, err := u.rdb.Client.Get(ctx, "user-jwt:"+strconv.Itoa(user.ID)).Result()
	if err == nil && cToken != "" {
		return cToken, nil
	}

	token, err := u.tokenMgr.GenerateJWT(user.ID)
	if err != nil {
		fmt.Println(err)
		return "", errors.ErrInternal
	}
	u.rdb.Client.Set(ctx, "user-jwt:"+strconv.Itoa(user.ID), token, 1*time.Hour)
	return token, nil
}
