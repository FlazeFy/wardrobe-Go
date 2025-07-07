package e2etest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
	"wardrobe/models/others"
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

type ResponseAuthRegister struct {
	Data    dataAuthRegister `json:"data"`
	Message string           `json:"message"`
	Status  string           `json:"status"`
}

type dataAuthRegister struct {
	Token string `json:"token"`
}

type ResponseAuthMyProfile struct {
	Data    others.MyProfile `json:"data"`
	Message string           `json:"message"`
	Status  string           `json:"status"`
}

// API POST : Basic Login (Admin)
// Test Case ID : TC-E2E-AU-001
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
// Test Case ID : TC-E2E-AU-002
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
// Test Case ID : TC-E2E-AU-003
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

// Test Case ID : TC-E2E-AU-004
func TestFailedPostBasicLoginWithEmptyPassword(t *testing.T) {
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

// Test Case ID : TC-E2E-AU-005
func TestFailedPostBasicLoginWithWrongPassword(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var res ResponseAuthLogin
	url := "http://127.0.0.1:9000/api/v1/auths/login"

	// Test Data
	userEmail := os.Getenv("USER_EMAIL")
	reqBody := map[string]interface{}{
		"email":    userEmail,
		"password": "nopass1234",
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
	assert.Equal(t, "Invalid password", res.Message)
}

// Test Case ID : TC-E2E-AU-006
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

// Test Case ID : TC-E2E-AU-007
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

// API POST : Basic Register (User)
// Test Case ID : TC-E2E-AU-008
func TestSuccessPostBasicRegisterWithValidUserData(t *testing.T) {
	var res ResponseAuthRegister
	url := "http://127.0.0.1:9000/api/v1/auths/register"

	// Test Data
	reqBody := map[string]interface{}{
		"username":         "tester123",
		"password":         "nopass123",
		"email":            "tester123456@gmail.com",
		"telegram_user_id": "1317625123",
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
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "User register", res.Message)

	// Check Auth Data
	assert.IsType(t, "", res.Data.Token)
}

// Test Case ID : TC-E2E-AU-009
func TestFailedPostBasicRegisterWithShortCharUsername(t *testing.T) {
	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/auths/register"

	// Test Data
	reqBody := map[string]interface{}{
		"username":         "test",
		"password":         "nopass123",
		"email":            "tester123456@gmail.com",
		"telegram_user_id": "1317625123",
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
	assert.Equal(t, "Username must be at least 6 characters long", res.Message[0].Error)
	assert.Equal(t, "Username", res.Message[0].Field)
}

// Test Case ID : TC-E2E-AU-010
func TestFailedPostBasicRegisterWithInvalidEmail(t *testing.T) {
	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/auths/register"

	// Test Data
	reqBody := map[string]interface{}{
		"username":         "tester123",
		"password":         "nopass123",
		"email":            "tester123456.com",
		"telegram_user_id": "1317625123",
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

// Test Case ID : TC-E2E-AU-011
func TestFailedPostBasicRegisterWithEmptyPassword(t *testing.T) {
	var res tests.ResponseFailedValidation
	url := "http://127.0.0.1:9000/api/v1/auths/register"

	// Test Data
	reqBody := map[string]interface{}{
		"username":         "tester123",
		"password":         "",
		"email":            "tester123456.com",
		"telegram_user_id": "1317625123",
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

// Test Case ID : TC-E2E-AU-012
func TestFailedPostBasicRegisterWithUsedEmail(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	var res tests.ResponseSimple
	url := "http://127.0.0.1:9000/api/v1/auths/register"

	// Test Data
	userEmail := os.Getenv("USER_EMAIL")
	reqBody := map[string]interface{}{
		"username":         "testeruser",
		"password":         "nopass123",
		"email":            userEmail,
		"telegram_user_id": "1317625123",
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
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "failed", res.Status)
	assert.Equal(t, "Username or email has already been used", res.Message)
}

// API POST : Sign Out (Admin)
// Test Case ID : TC-E2E-AU-013
func TestSuccessPostSignOutWithValidAdminToken(t *testing.T) {
	var res tests.ResponseSimple
	url := "http://127.0.0.1:9000/api/v1/auths/signout"

	// Test Data
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "admin")

	// Exec
	req, err := http.NewRequest("POST", url, nil)
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
	assert.Equal(t, "Admin signed out", res.Message)
}

// API POST : Sign Out (User)
// Test Case ID : TC-E2E-AU-014
func TestSuccessPostSignOutWithValidUserToken(t *testing.T) {
	var res tests.ResponseSimple
	url := "http://127.0.0.1:9000/api/v1/auths/signout"

	// Test Data
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("POST", url, nil)
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
	assert.Equal(t, "User signed out", res.Message)
}

// API POST : Sign Out All Role
// Test Case ID : TC-E2E-AU-015
func TestFailedPostSignOutWithEmptyToken(t *testing.T) {
	var res tests.ResponseSimple
	url := "http://127.0.0.1:9000/api/v1/auths/signout"

	// Exec
	req, err := http.NewRequest("POST", url, nil)
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
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "invalid authorization header", res.Message)
}

// Test Case ID : TC-E2E-AU-016
func TestFailedPostSignOutWithExpiredToken(t *testing.T) {
	var res tests.ResponseSimple
	url := "http://127.0.0.1:9000/api/v1/auths/signout"

	// Exec
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTI0NDgxMjgsImlhdCI6MTc1MTg0MzMyOCwicm9sZSI6ImFkbWluIiwidXNlcl9pZCI6IjQ3MWMxZWViLWFmZTItNDIzOC1iMjgyLWIyYzEzNWIyNzg0OCJ9.sYP5eHVPo48OTBqZco_yYo7yXotFcXg_aszpSrrBxuo"
	req, err := http.NewRequest("POST", url, nil)
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
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "token already expired", res.Message)
}

// API GET : Get My Profile
// Test Case ID : TC-E2E-AU-017
func TestSuccessGetMyProfileWithValidDataAdminAndUser(t *testing.T) {
	roles := []string{"user", "admin"}

	for _, role := range roles {
		var res ResponseAuthMyProfile
		url := "http://127.0.0.1:9000/api/v1/auths/profile"
		token, _ := tests.TemplatePostBasicLogin(t, nil, nil, role)

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
		assert.Equal(t, "Profile fetched", res.Message)

		// Check Object
		assert.NotEmpty(t, res.Data.Username)
		assert.NotEmpty(t, res.Data.Email)
		assert.NotEmpty(t, res.Data.CreatedAt)

		// Check Data Type
		assert.IsType(t, "", res.Data.Email)
		assert.IsType(t, "", res.Data.Username)
		assert.IsType(t, true, res.Data.TelegramIsValid)
		if res.Data.TelegramUserId != nil {
			assert.IsType(t, "", *res.Data.TelegramUserId)
		}
		assert.IsType(t, time.Time{}, res.Data.CreatedAt)
	}
}

// Test Case ID : TC-E2E-AU-018
func TestFailedGetMyProfileWithInvalidTokenAdminAndUser(t *testing.T) {
	tokens := []string{"d91ue09sad09", "a89usd8u1a"}

	for _, token := range tokens {
		var res ResponseAuthMyProfile
		url := "http://127.0.0.1:9000/api/v1/auths/profile"

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
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.NotEmpty(t, res.Message)
		assert.Equal(t, "invalid or expired token", res.Message)
	}
}
