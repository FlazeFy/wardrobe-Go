package e2etest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"wardrobe/models"
	"wardrobe/tests"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type ResponseGetAllDictionary struct {
	Data    []models.Dictionary `json:"data"`
	Message string              `json:"message"`
	Status  string              `json:"status"`
}

func TestGetAllDictionary(t *testing.T) {
	var res ResponseGetAllDictionary
	url := "http://127.0.0.1:9000/api/v1/dictionaries"
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Dictionary fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.DictionaryType)
		assert.NotEmpty(t, dt.DictionaryName)
		assert.NotEmpty(t, dt.CreatedAt)

		// Check Data Type
		assert.IsType(t, "", dt.DictionaryName)
		assert.IsType(t, "", dt.DictionaryType)
		assert.IsType(t, time.Time{}, dt.CreatedAt)
	}
}

func TestGetDictionaryByType(t *testing.T) {
	var res ResponseGetAllDictionary
	dictionary_type := "clothes_type"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/dictionaries/%s", dictionary_type)
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Dictionary fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.DictionaryType)
		assert.NotEmpty(t, dt.DictionaryName)
		assert.NotEmpty(t, dt.CreatedAt)

		// Check Data Type
		assert.IsType(t, "", dt.DictionaryName)
		assert.IsType(t, "", dt.DictionaryType)
		assert.IsType(t, time.Time{}, dt.CreatedAt)
	}
}
