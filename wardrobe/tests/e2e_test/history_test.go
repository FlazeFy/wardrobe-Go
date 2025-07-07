package e2etest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"wardrobe/models"
	"wardrobe/tests"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type ResponseGetAllHistory struct {
	Data    []models.GetHistory `json:"data"`
	Message string              `json:"message"`
	Status  string              `json:"status"`
}

// API Delete : Delete History By Id With Valid Id
// Test Case ID : TC-E2E-HS-001
func TestSuccessDeleteHistoryByIdWithValidId(t *testing.T) {
	var res tests.ResponseSimple

	// Test Data
	id := "43d2a364-92c3-4832-b088-ae4c2c1a6e73"
	url := "http://127.0.0.1:9000/api/v1/histories/destroy/" + id
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("DELETE", url, nil)
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
	assert.Equal(t, "History permanentally deleted", res.Message)
}

// Test Case ID : TC-E2E-HS-002
func TestFailedDeleteHistoryByIdWithForbiddenRole(t *testing.T) {
	var res tests.ResponseSimple

	// Test Data
	id := "fdb86a68-3dc0-491e-baa6-c5743fd745ca"
	url := "http://127.0.0.1:9000/api/v1/histories/destroy/" + id
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "admin")

	// Exec
	req, err := http.NewRequest("DELETE", url, nil)
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
	assert.Equal(t, "access forbidden for this role", res.Message)
}

// Test Case ID : TC-E2E-HS-003
func TestFailedDeleteHistoryByIdWithInvalidUUID(t *testing.T) {
	var res tests.ResponseSimple

	// Test Data
	id := "fdb86a68"
	url := "http://127.0.0.1:9000/api/v1/histories/destroy/" + id
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("DELETE", url, nil)
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
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "failed", res.Status)
	assert.Equal(t, "Invalid id", res.Message)
}

// Test Case ID : TC-E2E-HS-004
func TestFailedDeleteHistoryByIdWithInvalidId(t *testing.T) {
	var res tests.ResponseSimple

	// Test Data
	id := "fdb86a68-3dc0-491e-baa6-c5743fd745ba"
	url := "http://127.0.0.1:9000/api/v1/histories/destroy/" + id
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("DELETE", url, nil)
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
	assert.Equal(t, "History not found", res.Message)
}

// API Get : Get History With Valid Data
// Test Case ID : TC-E2E-HS-005
func TestSuccessGetAllHistoryWithValidRoleUser(t *testing.T) {
	var res ResponseGetAllHistory

	// Test Data
	url := "http://127.0.0.1:9000/api/v1/histories"
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "History fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.HistoryContext)
		assert.NotEmpty(t, dt.HistoryType)
		assert.NotEmpty(t, dt.CreatedAt)

		// Check Data Type
		assert.IsType(t, "", dt.HistoryContext)
		assert.IsType(t, "", dt.HistoryType)
		assert.IsType(t, time.Time{}, dt.CreatedAt)
	}
}

// Test Case ID : TC-E2E-HS-006
func TestSuccessGetAllHistoryWithValidRoleAdmin(t *testing.T) {
	var res ResponseGetAllHistory

	// Test Data
	url := "http://127.0.0.1:9000/api/v1/histories"
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
	assert.Equal(t, "History fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.HistoryContext)
		assert.NotEmpty(t, dt.HistoryType)
		assert.NotEmpty(t, dt.CreatedAt)
		assert.NotEmpty(t, dt.Username)

		// Check Data Type
		assert.IsType(t, "", dt.HistoryContext)
		assert.IsType(t, "", dt.HistoryType)
		assert.IsType(t, "", dt.Username)
		assert.IsType(t, time.Time{}, dt.CreatedAt)
	}
}
