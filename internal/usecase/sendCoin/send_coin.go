package SendCoinService

import (
	"context"
	stderrors "errors"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/redis"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/errors"
	"strconv"
)

type (
	SendCoinRepository interface {
		GetUserByUsername(ctx context.Context, username string) (domain.User, error)
		GetUserByID(ctx context.Context, userID int) (domain.User, error)
		DeductCoins(ctx context.Context, userID int, amount int) error
		AddCoins(ctx context.Context, userID int, amount int) error
		CreateTransaction(ctx context.Context, senderID, receiverID int, amount int) error
	}

	trManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}

	SendCoinService struct {
		repository SendCoinRepository
		trManager  trManager
		rdb        *redis.Redis
	}
)

func New(repository SendCoinRepository, trManager trManager, rdb *redis.Redis) *SendCoinService {
	return &SendCoinService{
		repository: repository,
		trManager:  trManager,
		rdb:        rdb,
	}
}

func (s *SendCoinService) SendCoin(ctx context.Context, toUser string, userID int, amount int) error {
	if toUser == "" {
		return errors.ErrBadRequest
	}

	receiver, err := s.repository.GetUserByUsername(ctx, toUser)
	if err != nil {
		if stderrors.Is(err, errors.ErrUserNotFound) {
			return errors.ErrUserNotFound
		}
		return errors.ErrBadRequest
	}

	if receiver.ID == userID {
		return errors.ErrSelfTransfer
	}

	sender, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		if stderrors.Is(err, errors.ErrUserNotFound) {
			return errors.ErrUserNotFound
		}
		return errors.ErrBadRequest
	}

	if sender.Coins < amount {
		return errors.ErrInsufficientFound
	}

	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		if err := s.repository.DeductCoins(ctx, userID, amount); err != nil {
			return errors.ErrBadRequest
		}

		if err := s.repository.AddCoins(ctx, receiver.ID, amount); err != nil {
			return errors.ErrInternal
		}

		if err := s.repository.CreateTransaction(ctx, userID, receiver.ID, amount); err != nil {
			return errors.ErrInternal
		}
		return nil
	})
	if err == nil {
		if _, err_minor := s.rdb.Client.Get(ctx, "user-info:"+strconv.Itoa(userID)).Result(); err_minor != nil {
			s.rdb.Client.Del(ctx, "user-info:"+strconv.Itoa(receiver.ID))
			s.rdb.Client.Del(ctx, "user-info:"+strconv.Itoa(sender.ID))
		}
	}
	return err
}
