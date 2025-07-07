package e2etest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"wardrobe/tests"

	"github.com/stretchr/testify/assert"
)

// API Delete : Delete History By Id With Valid Id
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
