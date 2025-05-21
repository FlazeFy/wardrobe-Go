package integrationtest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"wardrobe/models"

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
	url := "http://127.0.0.1:9000/api/v2/dictionary"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDgzODk1NjAsImlhdCI6MTc0Nzc4NDc2MCwidXNlcl9pZCI6IjU0YWYyOTEzLTVkMjEtNDgyOC05NDZhLTg2YTYyZTA1OGU1NCJ9.__Ag7Rcz-HInB4pjTzIN2KbBvGdulQIN8jmiU3Y9IZQ"

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
	assert.Equal(t, "dictionary fetched", res.Message)
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
