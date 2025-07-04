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

type ResponseGetAllQuestion struct {
	Data    []models.Question `json:"data"`
	Message string            `json:"message"`
	Status  string            `json:"status"`
}

func TestGetAllQuestion(t *testing.T) {
	var res ResponseGetAllQuestion
	url := "http://127.0.0.1:9000/api/v1/questions"

	// Exec
	resp, err := http.Get(url)
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
	assert.Equal(t, "Question fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, item := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, item.ID)
		assert.NotEmpty(t, item.Question)
		assert.NotEmpty(t, item.CreatedAt)
		if item.Answer != nil {
			assert.NotEmpty(t, *item.Answer)
		}

		// Check Data Type
		assert.IsType(t, true, item.IsShow)
		assert.IsType(t, "", item.Question)
		assert.IsType(t, time.Time{}, item.CreatedAt)
		if item.Answer != nil {
			assert.IsType(t, "", *item.Answer)
		}
	}
}
