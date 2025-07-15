package e2etest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"wardrobe/config"
	"wardrobe/models/others"
	"wardrobe/repositories"
	"wardrobe/seeders"
	"wardrobe/tests"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type ResponseGetMostContext struct {
	Data    []others.StatsContextTotal `json:"data"`
	Message string                     `json:"message"`
	Status  string                     `json:"status"`
}

type TestDataGetMostContext struct {
	TargetCol string
	Module    string
	Message   string
}

type TestDataGetMonthlyContext struct {
	ID      string
	Module  string
	Message string
	Year    int
}

// API GET : Get Most Context Clothes
func TestSuccessGetMostContextWithValidData(t *testing.T) {
	var testData = []TestDataGetMostContext{
		// Test Case ID : TC-E2E-ST-001
		{TargetCol: "clothes_category", Module: "clothes", Message: "Clothes fetched"},
		// Test Case ID : TC-E2E-ST-002
		{TargetCol: "used_context", Module: "clothes_used", Message: "Clothes used fetched"},
		// Test Case ID : TC-E2E-ST-003
		{TargetCol: "day", Module: "schedule", Message: "Schedule fetched"},
		// Test Case ID : TC-E2E-ST-004
		{TargetCol: "wash_type", Module: "wash", Message: "Wash fetched"},
		// Test Case ID : TC-E2E-ST-005
		{TargetCol: "weather_condition", Module: "user_weather", Message: "User weather fetched"},
	}

	for _, td := range testData {
		var res ResponseGetMostContext
		url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/stats/most_context/%s/%s", td.Module, td.TargetCol)
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
		assert.Equal(t, td.Message, res.Message)
		assert.NotNil(t, res.Data)

		for _, dt := range res.Data {
			// Check Object
			assert.NotEmpty(t, dt.Context)
			assert.NotEmpty(t, dt.Total)

			// Check Data Type
			assert.IsType(t, "", dt.Context)
			assert.IsType(t, 0, dt.Total)
		}
	}
}

func TestFailedGetMostContextWithInvalidTargetCol(t *testing.T) {
	var testData = []TestDataGetMostContext{
		// Test Case ID : TC-E2E-ST-006
		{TargetCol: "clothes_category_invalid", Module: "clothes"},
		// Test Case ID : TC-E2E-ST-007
		{TargetCol: "used_context_invalid", Module: "clothes_used"},
		// Test Case ID : TC-E2E-ST-008
		{TargetCol: "day_invalid", Module: "schedule"},
		// Test Case ID : TC-E2E-ST-009
		{TargetCol: "wash_type_invalid", Module: "wash"},
		// Test Case ID : TC-E2E-ST-010
		{TargetCol: "weather_condition_invalid", Module: "user_weather"},
	}

	for _, td := range testData {
		var res tests.ResponseSimple
		url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/stats/most_context/%s/%s", td.Module, td.TargetCol)
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
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.NotEmpty(t, res.Status)
		assert.Equal(t, "failed", res.Status)
		assert.NotEmpty(t, res.Message)
		assert.Equal(t, "Target col is not valid", res.Message)
	}
}

func TestFailedGetMostContextWithForbiddenRole(t *testing.T) {
	var testData = []TestDataGetMostContext{
		// Test Case ID : TC-E2E-ST-011
		{TargetCol: "clothes_category_invalid", Module: "clothes"},
		// Test Case ID : TC-E2E-ST-012
		{TargetCol: "used_context_invalid", Module: "clothes_used"},
		// Test Case ID : TC-E2E-ST-013
		{TargetCol: "day_invalid", Module: "schedule"},
		// Test Case ID : TC-E2E-ST-014
		{TargetCol: "wash_type_invalid", Module: "wash"},
		// Test Case ID : TC-E2E-ST-015
		{TargetCol: "weather_condition_invalid", Module: "user_weather"},
	}

	for _, td := range testData {
		var res tests.ResponseSimple
		url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/stats/most_context/%s/%s", td.Module, td.TargetCol)
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
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		assert.NotEmpty(t, res.Message)
		assert.Equal(t, "access forbidden for this role", res.Message)
	}
}

func TestFailedGetMostContextWithEmptyData(t *testing.T) {
	var testData = []TestDataGetMostContext{
		// Test Case ID : TC-E2E-ST-016
		{TargetCol: "clothes_category", Module: "clothes", Message: "Clothes not found"},
		// Test Case ID : TC-E2E-ST-017
		{TargetCol: "used_context", Module: "clothes_used", Message: "Clothes used not found"},
		// Test Case ID : TC-E2E-ST-018
		{TargetCol: "day", Module: "schedule", Message: "Schedule not found"},
		// Test Case ID : TC-E2E-ST-019
		{TargetCol: "wash_type", Module: "wash", Message: "Wash not found"},
		// Test Case ID : TC-E2E-ST-020
		{TargetCol: "weather_condition", Module: "user_weather", Message: "User weather not found"},
	}

	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	userRepo := repositories.NewUserRepository(db)
	clothesRepo := repositories.NewClothesRepository(db)
	clothesUsedRepo := repositories.NewClothesUsedRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	washRepo := repositories.NewWashRepository(db)
	userWeatherRepo := repositories.NewUserWeatherRepository(db)

	// Precondition
	clothesRepo.DeleteAll()
	clothesUsedRepo.DeleteAll()
	scheduleRepo.DeleteAll()
	washRepo.DeleteAll()
	userWeatherRepo.DeleteAll()

	for _, td := range testData {
		var res ResponseGetMostContext
		url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/stats/most_context/%s/%s", td.Module, td.TargetCol)
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
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.NotEmpty(t, res.Status)
		assert.Equal(t, "failed", res.Status)
		assert.NotEmpty(t, res.Message)
		assert.Equal(t, td.Message, res.Message)
	}

	// Seeder After Test
	seeders.SeedClothes(clothesRepo, userRepo, 200)
	seeders.SeedClothesUseds(clothesUsedRepo, userRepo, clothesRepo, 20)
	seeders.SeedSchedules(scheduleRepo, userRepo, clothesRepo, 7)
	seeders.SeedWashs(washRepo, userRepo, clothesRepo, 10)
	seeders.SeedUserWeathers(userWeatherRepo, userRepo, 100)
}

// API GET : Get Monthly Context
func TestSuccessGetMonthlyContextWithValidData(t *testing.T) {
	var testData = []TestDataGetMonthlyContext{
		// Test Case ID : TC-E2E-ST-021
		{ID: "all", Module: "clothes_used", Message: "Clothes used fetched", Year: 2025},
		// Test Case ID : TC-E2E-ST-022
		{ID: "all", Module: "outfit_used", Message: "Outfit used fetched", Year: 2025},
		// Test Case ID : TC-E2E-ST-023
		{ID: "all", Module: "wash", Message: "Wash fetched", Year: 2025},
	}

	for _, td := range testData {
		var res ResponseGetMostContext
		url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/stats/monthly/%s/%s/%d", td.Module, td.ID, td.Year)
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
		assert.Equal(t, td.Message, res.Message)
		assert.NotNil(t, res.Data)

		for _, dt := range res.Data {
			// Check Object
			assert.NotEmpty(t, dt.Context)
			assert.NotEmpty(t, dt.Total)

			// Check Data Type
			assert.IsType(t, "", dt.Context)
			assert.IsType(t, 0, dt.Total)
		}
	}
}

func TestFailedGetMonthlyContextSpecificWithInvalidUUID(t *testing.T) {
	var testData = []TestDataGetMonthlyContext{
		// Test Case ID : TC-E2E-ST-024
		{ID: "12345-ABCDE", Module: "clothes_used", Message: "Invalid clothes id", Year: 2025},
		// Test Case ID : TC-E2E-ST-025
		{ID: "12345-ABCDE", Module: "outfit_used", Message: "Invalid outfit id", Year: 2025},
		// Test Case ID : TC-E2E-ST-026
		{ID: "12345-ABCDE", Module: "wash", Message: "Invalid clothes id", Year: 2025},
	}

	for _, td := range testData {
		var res tests.ResponseSimple
		url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/stats/monthly/%s/%s/%d", td.Module, td.ID, td.Year)
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
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.NotEmpty(t, res.Status)
		assert.Equal(t, "failed", res.Status)
		assert.NotEmpty(t, res.Message)
		assert.Equal(t, td.Message, res.Message)
	}
}
