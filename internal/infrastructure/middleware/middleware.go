package middleware

import (
	"context"
	api "github.com/Azreil-OFD/Avito-test/internal/generate"
)

type jwtToken interface {
	GetUserIDFromJWT(tokenString string) (string, error)
}

type Middleware struct {
	jwt jwtToken
}

func New(jwt jwtToken) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}

func (h *Middleware) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	id, err := h.jwt.GetUserIDFromJWT(t.GetToken())
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, "userID", id)
	return ctx, nil
}
