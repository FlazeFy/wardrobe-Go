package e2etest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/seeders"
	"wardrobe/tests"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type ResponseGetAllUser struct {
	Data    []models.UserAnalytic `json:"data"`
	Message string                `json:"message"`
	Status  string                `json:"status"`
}

// API GET : Get All User
func TemplateSuccessGetAllUser(t *testing.T, resp *http.Response, res ResponseGetAllUser) {
	// Get Template Test
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "success", res.Status)
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "User fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.Username)
		assert.NotEmpty(t, dt.Email)
		assert.NotEmpty(t, dt.TotalClothes)
		assert.NotEmpty(t, dt.TotalOutfit)
		assert.NotEmpty(t, dt.CreatedAt)
		if dt.TelegramUserId != nil {
			assert.NotEmpty(t, dt.TelegramUserId)
		}

		// Check Data Type
		assert.IsType(t, "", dt.Username)
		assert.IsType(t, "", dt.Email)
		assert.IsType(t, time.Time{}, dt.CreatedAt)
		assert.IsType(t, 1, dt.TotalClothes)
		assert.IsType(t, 1, dt.TotalOutfit)
		assert.IsType(t, true, dt.TelegramIsValid)
		if dt.TelegramUserId != nil {
			assert.IsType(t, "", *dt.TelegramUserId)
		}
	}
}

// Test Case ID : TC-E2E-US-001
func TestSuccessGetAllUserWithValidData(t *testing.T) {
	var res ResponseGetAllUser
	username := "all"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/users/%s/%s", order, username)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "admin")

	// Exec
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Template Success Get All User
	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)
	TemplateSuccessGetAllUser(t, resp, res)
}

// Test Case ID : TC-E2E-US-002
func TestSuccessGetAllUserWithValidDataForSpecificUser(t *testing.T) {
	var res ResponseGetAllUser
	username := "tester"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/users/%s/%s", order, username)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "admin")

	// Exec
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Template Success Get All User
	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)
	TemplateSuccessGetAllUser(t, resp, res)
}

// Test Case ID : TC-E2E-US-003
func TestFailedGetAllUserWithForbiddenRole(t *testing.T) {
	var res tests.ResponseSimple
	username := "all"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/users/%s/%s", order, username)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Prepare Test
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &res)
	assert.NoError(t, err)

	// Get Template Test
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "access forbidden for this role", res.Message)
}

// Test Case ID : TC-E2E-US-004
func TestFailedGetAllUserWithEmptyData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	userRepo := repositories.NewUserRepository(db)

	// Precondition
	userRepo.DeleteAll()

	var res tests.ResponseSimple
	username := "all"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/users/%s/%s", order, username)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "admin")

	// Exec
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Prepare Test
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &res)
	assert.NoError(t, err)

	// Get Template Test
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "failed", res.Status)
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "User not found", res.Message)

	// Seeder After Test
	seeders.SeedUsers(userRepo, 200)
}
