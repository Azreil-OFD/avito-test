package SendCoinHandler

import (
	"context"
	stderror "errors"
	"github.com/AlekSi/pointer"
	api "github.com/Azreil-OFD/Avito-test/internal/generate"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/errors"
	"github.com/icefed/zlog"
	"strconv"
)

type SendCoinUsecase interface {
	SendCoin(ctx context.Context, toUser string, userID int, amount int) error
}

type SendCoinHandler struct {
	log     zlog.Logger
	usecase SendCoinUsecase
}

func New(log zlog.Logger, usecase SendCoinUsecase) *SendCoinHandler {
	return &SendCoinHandler{
		log:     log,
		usecase: usecase,
	}
}

func (h *SendCoinHandler) APISendCoinPost(ctx context.Context, req *api.SendCoinRequest) (api.APISendCoinPostRes, error) {
	userID, _ := strconv.Atoi(ctx.Value("userID").(string))
	if err := h.usecase.SendCoin(ctx, req.ToUser, userID, req.Amount); err != nil {
		if stderror.Is(err, errors.ErrUserNotFound) {
			return pointer.To(api.APISendCoinPostBadRequest(api.ErrorResponse{Errors: api.NewOptString(errors.ErrUserNotFound.Error())})), nil
		}
		if stderror.Is(err, errors.ErrInsufficientFound) {
			return pointer.To(api.APISendCoinPostBadRequest(api.ErrorResponse{Errors: api.NewOptString(errors.ErrInsufficientFound.Error())})), nil
		}
		if stderror.Is(err, errors.ErrSelfTransfer) {
			return pointer.To(api.APISendCoinPostBadRequest(api.ErrorResponse{Errors: api.NewOptString(errors.ErrSelfTransfer.Error())})), nil
		}
		h.log.WarnContext(ctx, err.Error())
		return pointer.To(api.APISendCoinPostInternalServerError(api.ErrorResponse{Errors: api.NewOptString(errors.ErrInternal.Error())})), nil
	}
	return &api.APISendCoinPostOK{}, nil
}
