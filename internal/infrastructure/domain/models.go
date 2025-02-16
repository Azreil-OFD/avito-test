package domain

import "time"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Coins        int    `json:"coins"`
	IsDeleted    bool   `json:"is_deleted"`
}

type Inventory struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ItemType  string `json:"item_type"`
	Quantity  int    `json:"quantity"`
	IsDeleted bool   `json:"is_deleted"`
}

type Transaction struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Amount     int       `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
	IsDeleted  bool      `json:"is_deleted"`
}
type TransactionFormat struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Amount   int    `json:"amount"`
	Type     string `json:"type"`
}
type MerchStore struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfile struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []CoinTransactionReceived `json:"received"`
	Sent     []CoinTransactionSent     `json:"sent"`
}

type CoinTransactionReceived struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
type CoinTransactionSent struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}
