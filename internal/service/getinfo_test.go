package service

import (
	"errors"
	"testing"
	"time"

	"merch-shop/internal/model"

	"github.com/stretchr/testify/assert"
)

type fakeInfoRepository struct {
	users       map[int64]*model.User
	purchases   map[int64][]*model.Purchase
	receivedTxs map[int64][]*model.Transaction
	sentTxs     map[int64][]*model.Transaction
}

func newFakeInfoRepository() *fakeInfoRepository {
	return &fakeInfoRepository{
		users:       make(map[int64]*model.User),
		purchases:   make(map[int64][]*model.Purchase),
		receivedTxs: make(map[int64][]*model.Transaction),
		sentTxs:     make(map[int64][]*model.Transaction),
	}
}

func (r *fakeInfoRepository) GetUserByUsername(username string) (*model.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *fakeInfoRepository) CreateUser(user *model.User) error {
	user.ID = int64(len(r.users) + 1)
	r.users[user.ID] = user
	return nil
}

func (r *fakeInfoRepository) UpdateUser(user *model.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *fakeInfoRepository) CreateTransaction(t *model.Transaction) error {
	// Для теста не требуется реализация.
	return nil
}

func (r *fakeInfoRepository) CreatePurchase(p *model.Purchase) error {
	r.purchases[p.UserID] = append(r.purchases[p.UserID], p)
	return nil
}

func (r *fakeInfoRepository) GetUserByID(userID int64) (*model.User, error) {
	user, exists := r.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *fakeInfoRepository) GetPurchasesByUserID(userID int64) ([]*model.Purchase, error) {
	return r.purchases[userID], nil
}

func (r *fakeInfoRepository) GetTransactionsReceivedByUserID(userID int64) ([]*model.Transaction, error) {
	return r.receivedTxs[userID], nil
}

func (r *fakeInfoRepository) GetTransactionsSentByUserID(userID int64) ([]*model.Transaction, error) {
	return r.sentTxs[userID], nil
}

func TestGetInfo_Success(t *testing.T) {
	repo := newFakeInfoRepository()

		user := &model.User{
		ID:       1,
		Username: "user1",
		Password: "pass",
		Coins:    900,
	}
	repo.users[user.ID] = user

	repo.purchases[user.ID] = []*model.Purchase{
		{ID: 1, UserID: user.ID, Item: "t-shirt", Price: 80, CreatedAt: time.Now()},
		{ID: 2, UserID: user.ID, Item: "t-shirt", Price: 80, CreatedAt: time.Now()},
		{ID: 3, UserID: user.ID, Item: "cup", Price: 20, CreatedAt: time.Now()},
	}

	sender := &model.User{
		ID:       2,
		Username: "sender",
		Password: "pass",
		Coins:    1000,
	}
	repo.users[sender.ID] = sender
	repo.receivedTxs[user.ID] = []*model.Transaction{
		{ID: 1, FromUserID: &sender.ID, ToUserID: user.ID, Amount: 50, Type: "transfer", CreatedAt: time.Now()},
	}

	recipient := &model.User{
		ID:       3,
		Username: "recipient",
		Password: "pass",
		Coins:    1100,
	}
	repo.users[recipient.ID] = recipient
	repo.sentTxs[user.ID] = []*model.Transaction{
		{ID: 2, FromUserID: &user.ID, ToUserID: recipient.ID, Amount: 30, Type: "transfer", CreatedAt: time.Now()},
	}

	info, err := GetInfo(repo, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 900, info.Coins)

	invMap := make(map[string]int)
	for _, item := range info.Inventory {
		invMap[item.Type] = item.Quantity
	}
	assert.Equal(t, 2, invMap["t-shirt"])
	assert.Equal(t, 1, invMap["cup"])

	assert.Len(t, info.CoinHistory.Received, 1)
	assert.Equal(t, "sender", info.CoinHistory.Received[0].FromUser)
	assert.Equal(t, 50, info.CoinHistory.Received[0].Amount)

	assert.Len(t, info.CoinHistory.Sent, 1)
	assert.Equal(t, "recipient", info.CoinHistory.Sent[0].ToUser)
	assert.Equal(t, 30, info.CoinHistory.Sent[0].Amount)
}

func TestGetInfo_UserNotFound(t *testing.T) {
	repo := newFakeInfoRepository()
	_, err := GetInfo(repo, 999) 
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}
