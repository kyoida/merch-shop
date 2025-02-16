package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Coins     int       `json:"coins"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID         int64     `json:"id"`
	FromUserID *int64    `json:"from_user_id,omitempty"` 
	ToUserID   int64     `json:"to_user_id"`
	Amount     int       `json:"amount"`
	Type       string    `json:"type"` // "transfer" или "purchase"
	CreatedAt  time.Time `json:"created_at"`
}

type Purchase struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Item      string    `json:"item"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistoryReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type CoinHistorySent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []CoinHistoryReceived `json:"received"`
	Sent     []CoinHistorySent     `json:"sent"`
}

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required,gt=0"`
}
