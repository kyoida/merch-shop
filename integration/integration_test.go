package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"merch-shop/internal/database"
	"merch-shop/internal/handlers"
	"merch-shop/internal/middleware"
	"merch-shop/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)
	if err := godotenv.Load(); err != nil {
		t.Log("No .env file found, using system environment variables")
	}

	os.Setenv("DB_HOST", getEnv("DB_HOST", "localhost"))
	os.Setenv("DB_PORT", getEnv("DB_PORT", "5432"))
	os.Setenv("DB_USER", getEnv("DB_USER", "postgres"))
	os.Setenv("DB_PASSWORD", getEnv("DB_PASSWORD", "1234"))
	os.Setenv("DB_NAME", getEnv("DB_NAME", "merch_shop"))
	os.Setenv("PORT", getEnv("PORT", "8080"))

	db, err := database.Connect()
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	repo := repository.NewPostgresRepository(db)

	router := gin.Default()
	router.POST("/api/auth", handlers.AuthHandler(repo))

	authGroup := router.Group("/api")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		authGroup.GET("/info", handlers.InfoHandler(repo))
		authGroup.POST("/sendCoin", handlers.SendCoinHandler(repo))
		authGroup.GET("/buy/:item", handlers.BuyHandler(repo))
	}
	return router
}

func TestIntegration_AuthAndInfo(t *testing.T) {
	router := setupTestRouter(t)
	ts := httptest.NewServer(router)
	defer ts.Close()

	authPayload := map[string]string{
		"username": "integrationUser",
		"password": "testpassword",
	}
	authPayloadJSON, _ := json.Marshal(authPayload)
	resp, err := http.Post(ts.URL+"/api/auth", "application/json", bytes.NewBuffer(authPayloadJSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var authResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	assert.NoError(t, err)
	token, ok := authResp["token"]
	assert.True(t, ok)
	assert.NotEmpty(t, token)
	resp.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", ts.URL+"/api/info", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var infoResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&infoResp)
	assert.NoError(t, err)
	resp.Body.Close()

	coins, ok := infoResp["coins"].(float64)
	assert.True(t, ok)
	assert.Equal(t, 1000.0, coins)
}

func TestIntegration_TransferCoins(t *testing.T) {
	router := setupTestRouter(t)
	ts := httptest.NewServer(router)
	defer ts.Close()
	client := &http.Client{Timeout: 5 * time.Second}

	senderPayload := map[string]string{
		"username": "senderUser",
		"password": "senderPass",
	}
	senderPayloadJSON, _ := json.Marshal(senderPayload)
	resp, err := http.Post(ts.URL+"/api/auth", "application/json", bytes.NewBuffer(senderPayloadJSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var senderAuthResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&senderAuthResp)
	assert.NoError(t, err)
	senderToken := senderAuthResp["token"]
	resp.Body.Close()

	recipientPayload := map[string]string{
		"username": "recipientUser",
		"password": "recipientPass",
	}
	recipientPayloadJSON, _ := json.Marshal(recipientPayload)
	resp, err = http.Post(ts.URL+"/api/auth", "application/json", bytes.NewBuffer(recipientPayloadJSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var recipientAuthResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&recipientAuthResp)
	assert.NoError(t, err)
	recipientToken := recipientAuthResp["token"]
	resp.Body.Close()

	transferPayload := map[string]interface{}{
		"toUser": "recipientUser",
		"amount": 100,
	}
	transferPayloadJSON, _ := json.Marshal(transferPayload)
	req, err := http.NewRequest("POST", ts.URL+"/api/sendCoin", bytes.NewBuffer(transferPayloadJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+senderToken)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	req, err = http.NewRequest("GET", ts.URL+"/api/info", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+senderToken)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var senderInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&senderInfo)
	assert.NoError(t, err)
	resp.Body.Close()

	req, err = http.NewRequest("GET", ts.URL+"/api/info", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+recipientToken)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var recipientInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&recipientInfo)
	assert.NoError(t, err)
	resp.Body.Close()

	senderCoins, ok := senderInfo["coins"].(float64)
	assert.True(t, ok)
	recipientCoins, ok := recipientInfo["coins"].(float64)
	assert.True(t, ok)
	assert.Equal(t, 900.0, senderCoins)
	assert.Equal(t, 1100.0, recipientCoins)
}
