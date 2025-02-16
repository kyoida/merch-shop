package repository

import "merch-shop/internal/model"

type Repository interface {
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	GetUserByID(userID int64) (*model.User, error)

	CreateTransaction(t *model.Transaction) error
	CreatePurchase(p *model.Purchase) error

	GetPurchasesByUserID(userID int64) ([]*model.Purchase, error)
	GetTransactionsReceivedByUserID(userID int64) ([]*model.Transaction, error)
	GetTransactionsSentByUserID(userID int64) ([]*model.Transaction, error)
}
