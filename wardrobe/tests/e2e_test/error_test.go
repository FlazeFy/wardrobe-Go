package e2etest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/models/others"
	"wardrobe/repositories"
	"wardrobe/seeders"
	"wardrobe/tests"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type ResponseGetAllError struct {
	Data     []models.ErrorAudit `json:"data"`
	Message  string              `json:"message"`
	Status   string              `json:"status"`
	Metadata others.Metadata     `json:"metadata"`
}

// API GET : Get All Error
// Test Case ID : TC-E2E-ER-001
func TestSuccessGetAllErrorWithValidData(t *testing.T) {
	var res ResponseGetAllError
	url := "http://127.0.0.1:9000/api/v1/errors"
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "success", res.Status)
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Error fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEmpty(t, dt.Message)
		assert.NotEmpty(t, dt.Total)
		assert.NotEmpty(t, dt.CreatedAt)

		// Check Data Type
		assert.IsType(t, "", dt.Message)
		assert.IsType(t, 0, dt.Total)
		assert.IsType(t, "", dt.CreatedAt)
	}

	// Pagination
	tests.TemplatePagination(t, res.Metadata)
}

// Test Case ID : TC-E2E-ER-002
func TestFailedGetAllErrorWithEmptyData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	errorRepo := repositories.NewErrorRepository(db)

	// Precondition
	errorRepo.DeleteAll()

	var res ResponseGetAllError
	url := "http://127.0.0.1:9000/api/v1/errors"
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
	assert.Equal(t, "Error not found", res.Message)

	// Seeder After Test
	seeders.SeedErrors(errorRepo, 25)
}

// Test Case ID : TC-E2E-ER-003
func TestFailedGetAllErrorWithForbiddenRole(t *testing.T) {
	var res ResponseGetAllError
	url := "http://127.0.0.1:9000/api/v1/errors"
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
