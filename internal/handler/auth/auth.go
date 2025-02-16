package AuthHandler

import (
	"context"
	api "github.com/Azreil-OFD/Avito-test/internal/generate"
	"github.com/icefed/zlog"
)

type AuthUsecase interface {
	Auth(ctx context.Context, username, password string) (string, error)
}
type AuthHandler struct {
	log     zlog.Logger
	usecase AuthUsecase
}

func New(log zlog.Logger, usecase AuthUsecase) *AuthHandler {
	return &AuthHandler{
		log:     log,
		usecase: usecase,
	}
}

func (h *AuthHandler) APIAuthPost(ctx context.Context, req *api.AuthRequest) (api.APIAuthPostRes, error) {
	token, err := h.usecase.Auth(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		h.log.WarnContext(ctx, err.Error())
		return nil, err
	}
	return &api.AuthResponse{Token: api.NewOptString(token)}, nil
}
