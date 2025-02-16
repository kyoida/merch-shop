package service

import (
	"testing"

	"merch-shop/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestPurchaseItem_Success(t *testing.T) {
	repo := newFakeRepository()
	repo.users["buyer"] = &model.User{ID: 1, Username: "buyer", Password: "pass", Coins: 1000}

	err := PurchaseItem(repo, "buyer", "t-shirt")
	assert.NoError(t, err, "purchase should succeed")

	buyer, _ := repo.GetUserByUsername("buyer")
	assert.Equal(t, 920, buyer.Coins, "Buyer balance should be reduced by 80")

	assert.Len(t, repo.purchases, 1)
	purchase := repo.purchases[0]
	assert.Equal(t, "t-shirt", purchase.Item)
	assert.Equal(t, 80, purchase.Price)

	assert.Len(t, repo.transactions, 1)
	tx := repo.transactions[0]
	assert.Equal(t, "purchase", tx.Type)
	assert.Equal(t, 80, tx.Amount)
}


func TestPurchaseItem_InsufficientCoins(t *testing.T) {
	repo := newFakeRepository()
	repo.users["buyer"] = &model.User{ID: 1, Username: "buyer", Password: "pass", Coins: 50}

	err := PurchaseItem(repo, "buyer", "t-shirt")
	assert.Error(t, err)
	assert.Equal(t, "insufficient coins", err.Error())

	buyer, _ := repo.GetUserByUsername("buyer")
	assert.Equal(t, 50, buyer.Coins)

	assert.Len(t, repo.purchases, 0)
	assert.Len(t, repo.transactions, 0)
}

func TestPurchaseItem_ItemNotFound(t *testing.T) {
	repo := newFakeRepository()
	repo.users["buyer"] = &model.User{ID: 1, Username: "buyer", Password: "pass", Coins: 1000}

	err := PurchaseItem(repo, "buyer", "nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
}
