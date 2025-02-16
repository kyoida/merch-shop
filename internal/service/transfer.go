package service

import (
	"errors"
	"time"

	"merch-shop/internal/model"
	"merch-shop/internal/repository"
)

func TransferCoins(repo repository.Repository, senderUsername, recipientUsername string, amount int) error {
	sender, err := repo.GetUserByUsername(senderUsername)
	if err != nil {
		return err
	}
	if sender.Coins < amount {
		return errors.New("insufficient coins")
	}

	recipient, err := repo.GetUserByUsername(recipientUsername)
	if err != nil {
		return err
	}

	sender.Coins -= amount
	recipient.Coins += amount

	if err := repo.UpdateUser(sender); err != nil {
		return err
	}
	if err := repo.UpdateUser(recipient); err != nil {
		return err
	}

	tx := &model.Transaction{
		FromUserID: &sender.ID,
		ToUserID:   recipient.ID,
		Amount:     amount,
		Type:       "transfer",
		CreatedAt:  time.Now(),
	}
	if err := repo.CreateTransaction(tx); err != nil {
		return err
	}

	return nil
}
