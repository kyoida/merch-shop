package service

import (
	"merch-shop/internal/model"
	"merch-shop/internal/repository"
)

func GetInfo(repo repository.Repository, userID int64) (*model.InfoResponse, error) {
	user, err := repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	purchases, err := repo.GetPurchasesByUserID(userID)
	if err != nil {
		return nil, err
	}
	inventoryMap := make(map[string]int)
	for _, p := range purchases {
		inventoryMap[p.Item]++
	}
	var inventory []model.InventoryItem
	for item, count := range inventoryMap {
		inventory = append(inventory, model.InventoryItem{Type: item, Quantity: count})
	}

	receivedTxs, err := repo.GetTransactionsReceivedByUserID(userID)
	if err != nil {
		return nil, err
	}
	var received []model.CoinHistoryReceived
	for _, tx := range receivedTxs {
		if tx.FromUserID != nil {
			sender, err := repo.GetUserByID(*tx.FromUserID)
			if err != nil {
				continue
			}
			received = append(received, model.CoinHistoryReceived{
				FromUser: sender.Username,
				Amount:   tx.Amount,
			})
		}
	}

	sentTxs, err := repo.GetTransactionsSentByUserID(userID)
	if err != nil {
		return nil, err
	}
	var sent []model.CoinHistorySent
	for _, tx := range sentTxs {
		recipient, err := repo.GetUserByID(tx.ToUserID)
		if err != nil {
			continue
		}
		sent = append(sent, model.CoinHistorySent{
			ToUser: recipient.Username,
			Amount: tx.Amount,
		})
	}

	info := &model.InfoResponse{
		Coins:       user.Coins,
		Inventory:   inventory,
		CoinHistory: model.CoinHistory{Received: received, Sent: sent},
	}
	return info, nil
}
