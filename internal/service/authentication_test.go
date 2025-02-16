package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"merch-shop/internal/model"
)

// fakeAuthRepository реализует Repository для тестирования аутентификации.
type fakeAuthRepository struct {
	users map[string]*model.User
}

func newFakeAuthRepository() *fakeAuthRepository {
	return &fakeAuthRepository{
		users: make(map[string]*model.User),
	}
}

func (r *fakeAuthRepository) GetUserByUsername(username string) (*model.User, error) {
	user, ok := r.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *fakeAuthRepository) CreateUser(user *model.User) error {
	if _, exists := r.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	user.ID = int64(len(r.users) + 1)
	r.users[user.Username] = user
	return nil
}

func (r *fakeAuthRepository) UpdateUser(user *model.User) error {
	r.users[user.Username] = user
	return nil
}

func (r *fakeAuthRepository) CreateTransaction(t *model.Transaction) error { return nil }
func (r *fakeAuthRepository) CreatePurchase(p *model.Purchase) error         { return nil }
func (r *fakeAuthRepository) GetUserByID(userID int64) (*model.User, error) {
	for _, user := range r.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
func (r *fakeAuthRepository) GetPurchasesByUserID(userID int64) ([]*model.Purchase, error) {
	return nil, nil
}
func (r *fakeAuthRepository) GetTransactionsReceivedByUserID(userID int64) ([]*model.Transaction, error) {
	return nil, nil
}
func (r *fakeAuthRepository) GetTransactionsSentByUserID(userID int64) ([]*model.Transaction, error) {
	return nil, nil
}

func TestAuthenticateUser_NewUser(t *testing.T) {
	repo := newFakeAuthRepository()
	req := model.AuthRequest{
		Username: "newuser",
		Password: "pass123",
	}

	user, err := AuthenticateUser(repo, req)
	assert.NoError(t, err)
	assert.Equal(t, "newuser", user.Username)
	assert.Equal(t, "pass123", user.Password)
	assert.Equal(t, 1000, user.Coins)
	assert.NotZero(t, user.ID)
}

func TestAuthenticateUser_ExistingUser_Success(t *testing.T) {
	repo := newFakeAuthRepository()
	// Создаем пользователя заранее.
	existing := &model.User{
		ID:       1,
		Username: "existing",
		Password: "secret",
		Coins:    1200,
	}
	repo.users[existing.Username] = existing

	req := model.AuthRequest{
		Username: "existing",
		Password: "secret",
	}
	user, err := AuthenticateUser(repo, req)
	assert.NoError(t, err)
	assert.Equal(t, existing, user)
}

func TestAuthenticateUser_InvalidCredentials(t *testing.T) {
	repo := newFakeAuthRepository()
	// Создаем пользователя с паролем "secret".
	existing := &model.User{
		ID:       1,
		Username: "existing",
		Password: "secret",
		Coins:    1200,
	}
	repo.users[existing.Username] = existing

	req := model.AuthRequest{
		Username: "existing",
		Password: "wrongpass",
	}
	_, err := AuthenticateUser(repo, req)
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}
