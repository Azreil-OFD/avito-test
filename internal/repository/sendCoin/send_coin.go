package SendCoinRepository

import (
	"context"
	"fmt"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/postgres"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/errors"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
)

type SendCoinRepository struct {
	pg       postgres.PgxPool
	trGetter *trmpgx.CtxGetter
}

func New(pg postgres.PgxPool) *SendCoinRepository {
	return &SendCoinRepository{
		pg:       pg,
		trGetter: trmpgx.DefaultCtxGetter,
	}
}

func (s *SendCoinRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "username", "password_hash", "coins", "is_deleted").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("error sql build: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)

	var user domain.User
	err = conn.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins, &user.IsDeleted)
	if err != nil {
		return domain.User{}, fmt.Errorf("error QueryRow: %w", err)
	}
	return user, nil
}

func (s *SendCoinRepository) GetUserByID(ctx context.Context, userID int) (domain.User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "username", "password_hash", "coins", "is_deleted").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("error sql build: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)

	var user domain.User
	err = conn.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins, &user.IsDeleted)
	if err != nil {
		return domain.User{}, errors.ErrUserNotFound
	}

	return user, nil
}

func (s *SendCoinRepository) DeductCoins(ctx context.Context, userID int, amount int) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("users").
		Set("coins", squirrel.Expr("coins - ?", amount)).
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error sql build: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error Exec: %w", err)
	}

	return nil
}

func (s *SendCoinRepository) AddCoins(ctx context.Context, userID int, amount int) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("users").
		Set("coins", squirrel.Expr("coins + ?", amount)).
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error sql build: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error Exec: %w", err)
	}

	return nil
}

func (s *SendCoinRepository) CreateTransaction(ctx context.Context, senderID, receiverID int, amount int) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("transactions").
		Columns("sender_id", "receiver_id", "amount").
		Values(senderID, receiverID, amount).
		ToSql()
	if err != nil {
		return fmt.Errorf("error sql build: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error Exec: %w", err)
	}

	return nil
}
