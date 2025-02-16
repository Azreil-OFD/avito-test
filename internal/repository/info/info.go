package InfoRepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/postgres"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
)

type InfoRepository struct {
	pg       postgres.PgxPool
	trGetter *trmpgx.CtxGetter
}

func New(pg postgres.PgxPool) *InfoRepository {
	return &InfoRepository{
		pg:       pg,
		trGetter: trmpgx.DefaultCtxGetter,
	}
}

func (s *InfoRepository) GetMerchByUserId(ctx context.Context, userId int) ([]domain.MerchStore, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "user_id", "item_type", "quantity").
		From("inventory").
		Where(squirrel.Eq{"user_id": userId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building SQL query: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var merchList []domain.MerchStore
	for rows.Next() {
		var merch domain.MerchStore
		if err := rows.Scan(&merch.ID, &merch.Type, &merch.Quantity); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		merchList = append(merchList, merch)
	}
	return merchList, nil
}

func (s *InfoRepository) GetTransactionsInfoById(ctx context.Context, id int) ([]domain.TransactionFormat, error) {
	// Строим SQL-запрос
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(
			"DISTINCT ON (t.id) t.id",
			"u.username",
			"t.amount",
			"CASE WHEN t.sender_id = $1 THEN 'sent' ELSE 'received' END AS type",
		).
		From("transactions t").
		Join("users u ON (t.sender_id = u.id OR t.receiver_id = u.id)").
		Where(squirrel.Or{
			squirrel.Eq{"t.sender_id": id},
			squirrel.Eq{"t.receiver_id": id},
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error building SQL query: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var transactions []domain.TransactionFormat
	for rows.Next() {
		var transaction domain.TransactionFormat
		if err := rows.Scan(&transaction.ID, &transaction.Username, &transaction.Amount, &transaction.Type); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	fmt.Println(transactions)
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return transactions, nil
}
func (s *InfoRepository) GetUserById(ctx context.Context, id int) (domain.User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "username", "password_hash", "coins", "is_deleted").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("error building SQL query: %w", err)
	}

	conn := s.trGetter.DefaultTrOrDB(ctx, s.pg)
	var user domain.User
	if err := conn.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins, &user.IsDeleted); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user not found")
		}
		return domain.User{}, fmt.Errorf("error executing query: %w", err)
	}
	return user, nil
}
