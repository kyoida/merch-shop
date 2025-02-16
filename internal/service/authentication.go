package service

import (
	"errors"

	"merch-shop/internal/model"
	"merch-shop/internal/repository"
)

func AuthenticateUser(repo repository.Repository, req model.AuthRequest) (*model.User, error) {
	user, err := repo.GetUserByUsername(req.Username)
	if err != nil {
		if err.Error() == "user not found" {
			newUser := &model.User{
				Username: req.Username,
				Password: req.Password,
				Coins:    1000,
			}
			if err := repo.CreateUser(newUser); err != nil {
				return nil, err
			}
			return newUser, nil
		}
		return nil, err
	}

	if user.Password != req.Password {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
