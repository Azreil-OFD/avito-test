package BuyItemHandler

import (
	"context"
	api "github.com/Azreil-OFD/Avito-test/internal/generate"
	"github.com/icefed/zlog"
	"strconv"
)

type BuyItemUsecase interface {
	BuyItem(ctx context.Context, userID int, itemID string) error
}

type BuyItemHandler struct {
	log     zlog.Logger
	usecase BuyItemUsecase
}

func New(log zlog.Logger, usecase BuyItemUsecase) *BuyItemHandler {
	return &BuyItemHandler{
		log:     log,
		usecase: usecase,
	}
}

func (h *BuyItemHandler) APIBuyItemGet(ctx context.Context, params api.APIBuyItemGetParams) (api.APIBuyItemGetRes, error) {
	userID, _ := strconv.Atoi(ctx.Value("userID").(string))
	err := h.usecase.BuyItem(ctx, userID, params.Item)
	if err != nil {
		h.log.WarnContext(ctx, err.Error())
		return nil, err
	}
	return &api.APIBuyItemGetOK{}, nil
}
