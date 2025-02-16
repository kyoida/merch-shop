package service

import (
	"errors"
	"time"

	"merch-shop/internal/model"
	"merch-shop/internal/repository"
)

var merchCatalog = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

func PurchaseItem(repo repository.Repository, username, item string) error {
	price, exists := merchCatalog[item]
	if !exists {
		return errors.New("item not found")
	}

	user, err := repo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if user.Coins < price {
		return errors.New("insufficient coins")
	}

	user.Coins -= price
	if err := repo.UpdateUser(user); err != nil {
		return err
	}

	purchase := &model.Purchase{
		UserID:    user.ID,
		Item:      item,
		Price:     price,
		CreatedAt: time.Now(),
	}
	if err := repo.CreatePurchase(purchase); err != nil {
		return err
	}

	tx := &model.Transaction{
		FromUserID: nil, 
		ToUserID:   user.ID,
		Amount:     price,
		Type:       "purchase",
		CreatedAt:  time.Now(),
	}
	return repo.CreateTransaction(tx)
}
