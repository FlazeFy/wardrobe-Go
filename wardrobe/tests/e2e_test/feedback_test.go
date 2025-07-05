package e2etest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"wardrobe/tests"

	"github.com/stretchr/testify/assert"
)

// API Post : Create Feedback
func TestFailedPostCreateFeedbackWithShortCharFeedbackBody(t *testing.T) {
	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/feedbacks"

	// Test Data
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")
	reqBody := map[string]interface{}{
		"feedback_rate": 4,
		"feedback_body": "t",
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
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

	// Check Validation Message
	assert.Equal(t, "FeedbackBody must be at least 3 characters long", res.Message[0].Error)
	assert.Equal(t, "FeedbackBody", res.Message[0].Field)
}

func TestFailedPostCreateFeedbackWithEmptyFeedbackBody(t *testing.T) {
	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/feedbacks"

	// Test Data
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")
	reqBody := map[string]interface{}{
		"feedback_rate": 4,
		"feedback_body": "",
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
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

	// Check Validation Message
	assert.Equal(t, "FeedbackBody is required", res.Message[0].Error)
	assert.Equal(t, "FeedbackBody", res.Message[0].Field)
}

func TestFailedPostCreateFeedbackWithInvalidFeedbackRate(t *testing.T) {
	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/feedbacks"

	// Test Data
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")
	reqBody := map[string]interface{}{
		"feedback_rate": 6,
		"feedback_body": "lorem ipsun",
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
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

	// Check Validation Message
	assert.Equal(t, "FeedbackRate must be at most 5 characters long", res.Message[0].Error)
	assert.Equal(t, "FeedbackRate", res.Message[0].Field)
}
