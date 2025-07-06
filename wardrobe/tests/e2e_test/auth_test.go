package e2etest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"wardrobe/tests"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type ResponseAuthLogin struct {
	Data    dataAuthLogin `json:"data"`
	Message string        `json:"message"`
	Status  string        `json:"status"`
}

type dataAuthLogin struct {
	Role  string `json:"role"`
	Token string `json:"token"`
}

// API POST : Basic Login (Admin)
func TestSuccessPostBasicLoginWithValidAdminData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var res ResponseAuthLogin
	url := "http://127.0.0.1:9000/api/v1/auths/login"

	// Test Data
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	reqBody := map[string]interface{}{
		"email":    adminEmail,
		"password": adminPassword,
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
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
	assert.Equal(t, "Admin login", res.Message)

	// Check Auth Data
	assert.Equal(t, "admin", res.Data.Role)
	assert.IsType(t, "", res.Data.Token)
}

// API POST : Basic Login (User)
func TestSuccessPostBasicLoginWithValidUserData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var res ResponseAuthLogin
	url := "http://127.0.0.1:9000/api/v1/auths/login"

	// Test Data
	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")
	reqBody := map[string]interface{}{
		"email":    userEmail,
		"password": userPassword,
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
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
	assert.Equal(t, "User login", res.Message)

	// Check Auth Data
	assert.Equal(t, "user", res.Data.Role)
	assert.IsType(t, "", res.Data.Token)
}

// API POST : Basic Login (All Role)
func TestFailedPostBasicLoginWithUnregisteredAccount(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var res ResponseAuthLogin
	url := "http://127.0.0.1:9000/api/v1/auths/login"

	// Test Data
	reqBody := map[string]interface{}{
		"email":    "cobaadmin@gmail.com",
		"password": "nopass123",
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
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
	assert.Equal(t, "Account not found", res.Message)
}

func TestFailedPostBasicLoginWithInvalidPassword(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/auths/login"

	// Test Data
	userEmail := os.Getenv("USER_EMAIL")
	reqBody := map[string]interface{}{
		"email":    userEmail,
		"password": "",
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
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
	assert.Equal(t, "Password is required", res.Message[0].Error)
	assert.Equal(t, "Password", res.Message[0].Field)
}

func TestFailedPostBasicLoginWithInvalidCharLengthPassword(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/auths/login"

	// Test Data
	userEmail := os.Getenv("USER_EMAIL")
	reqBody := map[string]interface{}{
		"email":    userEmail,
		"password": "nopas",
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
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
	assert.Equal(t, "Password must be at least 6 characters long", res.Message[0].Error)
	assert.Equal(t, "Password", res.Message[0].Field)
}

func TestFailedPostBasicLoginWithInvalidEmail(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/auths/login"

	// Test Data
	reqBody := map[string]interface{}{
		"email":    "testeradmin.com",
		"password": "nopass123",
	}
	jsonValue, _ := json.Marshal(reqBody)

	// Exec
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
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
	assert.Equal(t, "Email is not valid", res.Message[0].Error)
	assert.Equal(t, "Email", res.Message[0].Field)
}
