package BuyItemService

import (
	"context"
	"fmt"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/redis"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/errors"
	"strconv"
)

type (
	BuyItemRepository interface {
		GetMerchItem(ctx context.Context, itemName string) (domain.MerchStore, error)
		GetUserByID(ctx context.Context, userID int) (domain.User, error)
		DeductCoins(ctx context.Context, userID int, amount int) error
		AddToInventory(ctx context.Context, userID int, itemType string, quantity int) error
	}

	trManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}

	BuyItemService struct {
		repository BuyItemRepository
		trManager  trManager
		rdb        *redis.Redis
	}
)

func New(repository BuyItemRepository, trManager trManager, rdb *redis.Redis) *BuyItemService {
	return &BuyItemService{
		repository: repository,
		trManager:  trManager,
		rdb:        rdb,
	}
}

func (b *BuyItemService) BuyItem(ctx context.Context, userID int, itemName string) error {
	merchItem, err := b.repository.GetMerchItem(ctx, itemName)
	if err != nil {
		return errors.ErrBadRequest
	}

	user, err := b.repository.GetUserByID(ctx, userID)
	if err != nil {
		fmt.Println(err)
		return errors.ErrInternal
	}

	if user.Coins <= merchItem.Quantity {
		return errors.ErrBadRequest
	}

	err = b.trManager.Do(ctx, func(ctx context.Context) error {
		if err := b.repository.DeductCoins(ctx, userID, merchItem.Quantity); err != nil {
			return errors.ErrBadRequest
		}

		if err := b.repository.AddToInventory(ctx, userID, merchItem.Type, 1); err != nil {
			return errors.ErrInternal
		}
		return nil
	})

	if err == nil {
		if _, err_minor := b.rdb.Client.Get(ctx, "user-info:"+strconv.Itoa(userID)).Result(); err_minor != nil {
			b.rdb.Client.Del(ctx, "user-info:"+strconv.Itoa(userID))
		}
	}

	return err
}
