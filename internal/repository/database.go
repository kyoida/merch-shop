// internal/repository/postgres_repository.go
package repository

import (
	"database/sql"
	"errors"

	"merch-shop/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetUserByUsername(username string) (*model.User, error) {
	row := r.db.QueryRow("SELECT id, username, password, coins FROM users WHERE username = $1", username)
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Coins); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) CreateUser(user *model.User) error {
	query := "INSERT INTO users (username, password, coins, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id"
	return r.db.QueryRow(query, user.Username, user.Password, user.Coins).Scan(&user.ID)
}

func (r *PostgresRepository) UpdateUser(user *model.User) error {
	query := "UPDATE users SET password = $1, coins = $2 WHERE id = $3"
	_, err := r.db.Exec(query, user.Password, user.Coins, user.ID)
	return err
}

func (r *PostgresRepository) GetUserByID(userID int64) (*model.User, error) {
	row := r.db.QueryRow("SELECT id, username, password, coins FROM users WHERE id = $1", userID)
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Coins); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) CreateTransaction(t *model.Transaction) error {
	query := "INSERT INTO transactions (from_user_id, to_user_id, amount, type, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id"
	return r.db.QueryRow(query, t.FromUserID, t.ToUserID, t.Amount, t.Type).Scan(&t.ID)
}

func (r *PostgresRepository) CreatePurchase(p *model.Purchase) error {
	query := "INSERT INTO purchases (user_id, item, price, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id"
	return r.db.QueryRow(query, p.UserID, p.Item, p.Price).Scan(&p.ID)
}

func (r *PostgresRepository) GetPurchasesByUserID(userID int64) ([]*model.Purchase, error) {
	rows, err := r.db.Query("SELECT id, user_id, item, price, created_at FROM purchases WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*model.Purchase
	for rows.Next() {
		var p model.Purchase
		if err := rows.Scan(&p.ID, &p.UserID, &p.Item, &p.Price, &p.CreatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, &p)
	}
	return purchases, nil
}

func (r *PostgresRepository) GetTransactionsReceivedByUserID(userID int64) ([]*model.Transaction, error) {
	rows, err := r.db.Query("SELECT id, from_user_id, to_user_id, amount, type, created_at FROM transactions WHERE to_user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []*model.Transaction
	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Type, &tx.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, &tx)
	}
	return txs, nil
}

func (r *PostgresRepository) GetTransactionsSentByUserID(userID int64) ([]*model.Transaction, error) {
	rows, err := r.db.Query("SELECT id, from_user_id, to_user_id, amount, type, created_at FROM transactions WHERE from_user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []*model.Transaction
	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Type, &tx.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, &tx)
	}
	return txs, nil
}
