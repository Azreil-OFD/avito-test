package InfoService

import (
	"context"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/redis"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/errors"
	"github.com/ogen-go/ogen/json"
	"github.com/samber/lo"
	"strconv"
	"sync"
	"time"
)

type (
	InfoRepository interface {
		GetMerchByUserId(ctx context.Context, userId int) ([]domain.MerchStore, error)
		GetTransactionsInfoById(ctx context.Context, id int) ([]domain.TransactionFormat, error)
		GetUserById(ctx context.Context, id int) (domain.User, error)
	}

	InfoService struct {
		repository InfoRepository
		rdb        *redis.Redis
	}
)

func New(repository InfoRepository, rdb *redis.Redis) *InfoService {
	return &InfoService{
		repository: repository,
		rdb:        rdb,
	}
}

func (i *InfoService) GetInfo(ctx context.Context, userID int) (domain.UserProfile, error) {
	var (
		merch                              []domain.MerchStore
		transactions                       []domain.TransactionFormat
		user                               domain.User
		merchErr, transactionsErr, userErr error
		result                             domain.UserProfile
	)

	if cacheInfo, err := i.rdb.Client.Get(ctx, "user-info:"+strconv.Itoa(userID)).Result(); err == nil {
		if err := json.Unmarshal([]byte(cacheInfo), &result); err == nil {
			return result, nil
		}
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		merch, merchErr = i.repository.GetMerchByUserId(ctx, userID)
	}()
	go func() {
		defer wg.Done()
		transactions, transactionsErr = i.repository.GetTransactionsInfoById(ctx, userID)
	}()
	go func() {
		defer wg.Done()
		user, userErr = i.repository.GetUserById(ctx, userID)
	}()
	wg.Wait()

	if merchErr != nil || transactionsErr != nil || userErr != nil {
		return domain.UserProfile{}, errors.ErrInternal
	}

	inventory := lo.Map(merch, func(item domain.MerchStore, _ int) domain.InventoryItem {
		return domain.InventoryItem{Type: item.Type, Quantity: item.Quantity}
	})
	grouped := lo.GroupBy(inventory, func(item domain.InventoryItem) string {
		return item.Type
	})
	uniqueInventory := lo.MapToSlice(grouped, func(key string, group []domain.InventoryItem) domain.InventoryItem {
		totalQuantity := lo.SumBy(group, func(item domain.InventoryItem) int {
			return item.Quantity
		})
		return domain.InventoryItem{Type: key, Quantity: totalQuantity}
	})
	coinsHistory := domain.CoinHistory{}
	for _, item := range transactions {
		switch item.Type {
		case "sent":
			coinsHistory.Sent = append(coinsHistory.Sent, domain.CoinTransactionSent{FromUser: item.Username, Amount: item.Amount})
		case "received":
			coinsHistory.Received = append(coinsHistory.Received, domain.CoinTransactionReceived{ToUser: item.Username, Amount: item.Amount})
		}
	}
	result = domain.UserProfile{
		Coins:       user.Coins,
		Inventory:   uniqueInventory,
		CoinHistory: coinsHistory,
	}
	if resultJSON, err := json.Marshal(result); err == nil {
		i.rdb.Client.Set(ctx, "user-info:"+strconv.Itoa(userID), resultJSON, 1*time.Minute)
	}
	return result, nil
}
