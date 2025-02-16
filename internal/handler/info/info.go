package InfoHandler

import (
	"context"
	api "github.com/Azreil-OFD/Avito-test/internal/generate"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/icefed/zlog"
	"github.com/samber/lo"
	"strconv"
)

type InfoUsecase interface {
	GetInfo(ctx context.Context, userID int) (domain.UserProfile, error)
}
type InfoHandler struct {
	log     zlog.Logger
	usecase InfoUsecase
}

func New(log zlog.Logger, usecase InfoUsecase) *InfoHandler {
	return &InfoHandler{
		log:     log,
		usecase: usecase,
	}
}

func (I *InfoHandler) APIInfoGet(ctx context.Context) (api.APIInfoGetRes, error) {
	userID, _ := strconv.Atoi(ctx.Value("userID").(string))
	info, err := I.usecase.GetInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &api.InfoResponse{
		Coins: api.NewOptInt(info.Coins),
		Inventory: lo.Map(info.Inventory, func(item domain.InventoryItem, index int) api.InfoResponseInventoryItem {
			return api.InfoResponseInventoryItem{
				Type:     api.NewOptString(item.Type),
				Quantity: api.NewOptInt(item.Quantity),
			}
		}),
		CoinHistory: api.OptInfoResponseCoinHistory{
			Value: api.InfoResponseCoinHistory{
				Received: lo.Map(info.CoinHistory.Received, func(item domain.CoinTransactionReceived, index int) api.InfoResponseCoinHistoryReceivedItem {
					return api.InfoResponseCoinHistoryReceivedItem{
						FromUser: api.OptString{
							Value: item.ToUser,
							Set:   item.ToUser != "",
						},
						Amount: api.OptInt{
							Value: item.Amount,
							Set:   item.Amount != 0,
						},
					}
				}),
				Sent: lo.Map(info.CoinHistory.Sent, func(item domain.CoinTransactionSent, index int) api.InfoResponseCoinHistorySentItem {
					return api.InfoResponseCoinHistorySentItem{
						ToUser: api.OptString{
							Value: item.FromUser,
							Set:   item.FromUser != "",
						},
						Amount: api.OptInt{
							Value: item.Amount,
							Set:   item.Amount != 0,
						},
					}
				}),
			},
			Set: true,
		},
	}, nil
}
