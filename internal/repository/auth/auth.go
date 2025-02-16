package AuthRepository

import (
	"context"
	"fmt"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/postgres"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
)

type AuthRepository struct {
	pg       postgres.PgxPool
	trGetter *trmpgx.CtxGetter
}

func New(pg postgres.PgxPool) *AuthRepository {
	return &AuthRepository{
		pg:       pg,
		trGetter: trmpgx.DefaultCtxGetter,
	}
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "username", "password_hash", "coins", "is_deleted").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("error building SQL query: %w", err)
	}

	conn := r.trGetter.DefaultTrOrDB(ctx, r.pg)

	var user domain.User
	err = conn.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins, &user.IsDeleted)
	if err != nil {
		return domain.User{}, fmt.Errorf("error executing QueryRow: %w", err)
	}

	return user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, username, passwordHash string, coins int) (domain.User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("users").
		Columns("username", "password_hash", "coins").
		Values(username, passwordHash, coins).
		Suffix("RETURNING id, username, password_hash, coins, is_deleted").
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("error building SQL query: %w", err)
	}

	conn := r.trGetter.DefaultTrOrDB(ctx, r.pg)

	var user domain.User
	err = conn.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins, &user.IsDeleted)
	if err != nil {
		return domain.User{}, fmt.Errorf("error executing QueryRow: %w", err)
	}

	return user, nil
}
