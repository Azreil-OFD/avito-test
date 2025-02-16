package BuyItemRepository

import (
	"context"
	"fmt"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/postgres"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
)

type BuyItemRepository struct {
	pg       postgres.PgxPool
	trGetter *trmpgx.CtxGetter
}

func New(pg postgres.PgxPool) *BuyItemRepository {
	return &BuyItemRepository{
		pg:       pg,
		trGetter: trmpgx.DefaultCtxGetter,
	}
}

func (s *BuyItemRepository) GetMerchItem(ctx context.Context, itemName string) (domain.MerchStore, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "item", "price", "created_at", "updated_at").
		From("merch_store").
		Where(squirrel.Eq{"item": itemName}).
		ToSql()
	if err != nil {
		return domain.MerchStore{}, fmt.Errorf("error building SQL: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)

	var merch domain.MerchStore
	err = conn.QueryRow(ctx, query, args...).Scan(&merch.ID, &merch.Type, &merch.Quantity, &merch.CreatedAt, &merch.UpdatedAt)
	if err != nil {
		return domain.MerchStore{}, fmt.Errorf("error executing query: %w", err)
	}

	return merch, nil
}

func (s *BuyItemRepository) GetUserByID(ctx context.Context, userID int) (domain.User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "username", "password_hash", "coins", "is_deleted").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("error building SQL: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)

	var user domain.User
	err = conn.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins, &user.IsDeleted)
	if err != nil {
		return domain.User{}, fmt.Errorf("error executing query: %w", err)
	}

	return user, nil
}

func (s *BuyItemRepository) DeductCoins(ctx context.Context, userID int, amount int) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("users").
		Set("coins", squirrel.Expr("coins - ?", amount)).
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func (s *BuyItemRepository) AddToInventory(ctx context.Context, userID int, itemType string, quantity int) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("inventory").
		Columns("user_id", "item_type", "quantity").
		Values(userID, itemType, quantity).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}
