package service

import (
	"errors"
	"testing"
	"time"

	"merch-shop/internal/model"

	"github.com/stretchr/testify/assert"
)

type fakeRepository struct {
	users        map[string]*model.User
	transactions []*model.Transaction
	purchases    []*model.Purchase
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		users:        make(map[string]*model.User),
		transactions: []*model.Transaction{},
		purchases:    []*model.Purchase{},
	}
}

func (r *fakeRepository) GetUserByUsername(username string) (*model.User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *fakeRepository) CreatePurchase(p *model.Purchase) error {
	p.ID = int64(len(r.purchases) + 1)
	p.CreatedAt = time.Now()
	r.purchases = append(r.purchases, p)
	return nil
}

func (r *fakeRepository) UpdateUser(user *model.User) error {
	r.users[user.Username] = user
	return nil
}

func (r *fakeRepository) CreateTransaction(t *model.Transaction) error {
	t.ID = int64(len(r.transactions) + 1)
	t.CreatedAt = time.Now()
	r.transactions = append(r.transactions, t)
	return nil
}

func (r *fakeRepository) CreateUser(user *model.User) error {
	return nil
}

func (r *fakeRepository) GetUserByID(userID int64) (*model.User, error) {
	return nil, nil
}

func (r *fakeRepository) GetPurchasesByUserID(userID int64) ([]*model.Purchase, error) {
	return nil, nil
}
func (r *fakeRepository) GetTransactionsReceivedByUserID(userID int64) ([]*model.Transaction, error) {
	return nil, nil
}
func (r *fakeRepository) GetTransactionsSentByUserID(userID int64) ([]*model.Transaction, error) {
	return nil, nil
}
func TestTransferCoins_Success(t *testing.T) {
	repo := newFakeRepository()
	repo.users["sender"] = &model.User{ID: 1, Username: "sender", Password: "pass", Coins: 1000}
	repo.users["recipient"] = &model.User{ID: 2, Username: "recipient", Password: "pass", Coins: 1000}

	err := TransferCoins(repo, "sender", "recipient", 100)
	assert.NoError(t, err, "перевод монет должен пройти успешно")

	sender, _ := repo.GetUserByUsername("sender")
	recipient, _ := repo.GetUserByUsername("recipient")
	assert.Equal(t, 900, sender.Coins)
	assert.Equal(t, 1100, recipient.Coins)

	assert.Len(t, repo.transactions, 1)
	tx := repo.transactions[0]
	assert.Equal(t, "transfer", tx.Type)
	assert.Equal(t, 100, tx.Amount)
	assert.NotNil(t, tx.FromUserID)
	assert.Equal(t, int64(1), *tx.FromUserID)
	assert.Equal(t, int64(2), tx.ToUserID)
}

func TestTransferCoins_InsufficientFunds(t *testing.T) {
	repo := newFakeRepository()
	repo.users["sender"] = &model.User{ID: 1, Username: "sender", Password: "pass", Coins: 50}
	repo.users["recipient"] = &model.User{ID: 2, Username: "recipient", Password: "pass", Coins: 1000}

	err := TransferCoins(repo, "sender", "recipient", 100)
	assert.Error(t, err)
	assert.Equal(t, "insufficient coins", err.Error())

	sender, _ := repo.GetUserByUsername("sender")
	recipient, _ := repo.GetUserByUsername("recipient")
	assert.Equal(t, 50, sender.Coins)
	assert.Equal(t, 1000, recipient.Coins)

	assert.Len(t, repo.transactions, 0)
}

func TestTransferCoins_RecipientNotFound(t *testing.T) {
	repo := newFakeRepository()
	repo.users["sender"] = &model.User{ID: 1, Username: "sender", Password: "pass", Coins: 1000}

	err := TransferCoins(repo, "sender", "nonexistent", 100)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}
