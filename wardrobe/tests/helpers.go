package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// For E2E Test
func TemplatePostBasicLogin(t *testing.T, email, password *string, roleAccount string) (string, string) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var emailTest string
	var passwordTest string
	if email != nil && *email != "" {
		emailTest = *email
	} else {
		if roleAccount == "admin" {
			emailTest = os.Getenv("ADMIN_EMAIL")

		} else if roleAccount == "user" {
			emailTest = os.Getenv("USER_EMAIL")
		}
	}
	if password != nil && *password != "" {
		passwordTest = *password
	} else {
		passwordTest = "nopass123"
	}

	payload := map[string]string{
		"email":    emailTest,
		"password": passwordTest,
	}
	fmt.Println(payload)
	jsonPayload, _ := json.Marshal(payload)

	url := "http://127.0.0.1:9000/api/v1/auths/login"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "User login", result["message"])

	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	token, ok := data["token"].(string)
	assert.True(t, ok, "token should be a string")
	assert.NotEmpty(t, token)

	role, ok := data["role"].(string)
	assert.True(t, ok, "role should be a string")
	assert.NotEmpty(t, role)

	return token, role
}
