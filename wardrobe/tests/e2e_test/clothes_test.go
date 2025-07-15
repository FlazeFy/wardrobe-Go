package e2etest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/seeders"
	"wardrobe/tests"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type GetClothesLastHistory struct {
	LastAddedClothes   *string    `json:"last_added_clothes"`
	LastAddedDate      *time.Time `json:"last_added_date"`
	LastDeletedClothes *string    `json:"last_deleted_clothes"`
	LastDeletedDate    *time.Time `json:"last_deleted_date"`
}
type ResponseGetClothesLastHistory struct {
	Data    GetClothesLastHistory `json:"data"`
	Message string                `json:"message"`
	Status  string                `json:"status"`
}

type ResponseGetUsedHistory struct {
	Data    []models.ClothesUsedHistory `json:"data"`
	Message string                      `json:"message"`
	Status  string                      `json:"status"`
}

type ResponseGetClothesHeader struct {
	Data    []models.ClothesHeader `json:"data"`
	Message string                 `json:"message"`
	Status  string                 `json:"status"`
}

type ResponseGetClothesDetail struct {
	Data    []models.ClothesDetail `json:"data"`
	Message string                 `json:"message"`
	Status  string                 `json:"status"`
}

type ResponseGetClothesDeleted struct {
	Data    []models.ClothesDeleted `json:"data"`
	Message string                  `json:"message"`
	Status  string                  `json:"status"`
}

// API GET : Get Clothes Last History
// Test Case ID : TC-E2E-CL-001
func TestSuccessGetClothesLastHistoryWithValidData(t *testing.T) {
	var res ResponseGetClothesLastHistory
	url := "http://127.0.0.1:9000/api/v1/clothes/last_history"
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
	assert.Equal(t, "Clothes fetched", res.Message)
	assert.NotNil(t, res.Data)

	if res.Data.LastAddedClothes != nil {
		// Check Object
		assert.NotEqual(t, uuid.Nil, res.Data.LastAddedClothes)
		assert.NotEmpty(t, *res.Data.LastAddedDate)
		// Check Data Type
		assert.IsType(t, "", *res.Data.LastAddedClothes)
		assert.IsType(t, time.Time{}, *res.Data.LastAddedDate)
	}
	if res.Data.LastDeletedClothes != nil {
		// Check Object
		assert.NotEqual(t, uuid.Nil, res.Data.LastDeletedClothes)
		assert.NotEmpty(t, *res.Data.LastDeletedDate)
		// Check Data Type
		assert.IsType(t, "", *res.Data.LastAddedClothes)
		assert.IsType(t, time.Time{}, *res.Data.LastAddedDate)
	}
}

// Test Case ID : TC-E2E-CL-002
func TestFailedGetClothesLastHistoryWithEmptyData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	clothesRepo := repositories.NewClothesRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Precondition
	clothesRepo.DeleteAll()

	var res ResponseGetClothesLastHistory
	url := "http://127.0.0.1:9000/api/v1/clothes/last_history"
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
	assert.Equal(t, "Clothes not found", res.Message)

	// Seeder After Test
	seeders.SeedClothes(clothesRepo, userRepo, 200)
}

// Test Case ID : TC-E2E-CL-003
func TestFailedGetClothesLastHistoryWithForbiddenRole(t *testing.T) {
	var res ResponseGetAllError
	url := "http://127.0.0.1:9000/api/v1/clothes/last_history"
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

// API GET : Get Clothes Used History
// Test Case ID : TC-E2E-CL-004
func TestSuccessGetClothesUsedHistoryWithValidData(t *testing.T) {
	clothesId := "all"
	order := "desc"

	var res ResponseGetUsedHistory
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes_used/history/%s/%s", clothesId, order)
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
	assert.Equal(t, "Clothes used fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.ClothesName)
		assert.NotEmpty(t, dt.ClothesType)
		assert.NotEmpty(t, dt.UsedContext)
		assert.NotEmpty(t, dt.CreatedAt)
		if dt.ClothesNote != nil {
			assert.NotEmpty(t, dt.ClothesNote)
		}

		// Check Data Type
		assert.IsType(t, "", dt.ClothesName)
		assert.IsType(t, "", dt.ClothesType)
		assert.IsType(t, "", dt.UsedContext)
		assert.IsType(t, time.Time{}, dt.CreatedAt)
		if dt.ClothesNote != nil {
			assert.IsType(t, new(string), dt.ClothesNote)
		}
	}
}

// Test Case ID : TC-E2E-CL-005
func TestFailedGetClothesUsedHistoryWithEmptyData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	clothesRepo := repositories.NewClothesRepository(db)
	userRepo := repositories.NewUserRepository(db)
	clothesUsedRepo := repositories.NewClothesUsedRepository(db)

	// Precondition
	clothesUsedRepo.DeleteAll()

	var res ResponseGetUsedHistory
	clothesId := "all"
	order := "desc"

	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes_used/history/%s/%s", clothesId, order)
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
	assert.Equal(t, "Clothes used not found", res.Message)

	// Seeder After Test
	seeders.SeedClothesUseds(clothesUsedRepo, userRepo, clothesRepo, 25)
}

// Test Case ID : TC-E2E-CL-006
func TestFailedGetClothesUsedHistoryWithForbiddenRole(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "all"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes_used/history/%s/%s", clothesId, order)
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

// API GET : Get All Clothes Header
func TemplateSuccessGetAllClothesHeader(t *testing.T, resp *http.Response, res ResponseGetClothesHeader) {
	// Get Template Test
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "success", res.Status)
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Clothes fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.ClothesName)
		assert.NotEmpty(t, dt.ClothesColor)
		assert.NotEmpty(t, dt.ClothesQty)
		assert.NotEmpty(t, dt.ClothesType)
		assert.NotEmpty(t, dt.ClothesCategory)
		assert.NotEmpty(t, dt.ClothesSize)
		assert.NotEmpty(t, dt.ClothesGender)
		if dt.ClothesImage != nil {
			assert.NotEmpty(t, dt.ClothesImage)
		}

		// Check Data Type
		assert.IsType(t, "", dt.ClothesName)
		assert.IsType(t, "", dt.ClothesColor)
		assert.IsType(t, "", dt.ClothesType)
		assert.IsType(t, "", dt.ClothesCategory)
		assert.IsType(t, "", dt.ClothesSize)
		assert.IsType(t, "", dt.ClothesGender)
		assert.IsType(t, 1, dt.ClothesQty)
		assert.IsType(t, true, dt.IsFaded)
		assert.IsType(t, true, dt.HasWashed)
		assert.IsType(t, true, dt.HasIroned)
		assert.IsType(t, true, dt.IsFavorite)
		assert.IsType(t, true, dt.IsScheduled)
		if dt.ClothesImage != nil {
			assert.IsType(t, "", dt.ClothesImage)
		}
	}
}

// Test Case ID : TC-E2E-CL-007
func TestSuccessGetAllClothesHeaderWithValidData(t *testing.T) {
	var res ResponseGetClothesHeader
	category := "all"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/header/%s/%s", category, order)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Template Success Get All Clothes Header
	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)
	TemplateSuccessGetAllClothesHeader(t, resp, res)
}

// Test Case ID : TC-E2E-CL-010
func TestSuccessGetAllClothesHeaderWithValidClothesCategory(t *testing.T) {
	var res ResponseGetClothesHeader
	category := "bottom_body"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/header/%s/%s", category, order)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Template Success Get All Clothes Header
	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)
	TemplateSuccessGetAllClothesHeader(t, resp, res)
}

// Test Case ID : TC-E2E-CL-008
func TestFailedGetAllClothesHeaderWithEmptyData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	userRepo := repositories.NewUserRepository(db)
	clothesRepo := repositories.NewClothesRepository(db)

	// Precondition
	clothesRepo.DeleteAll()

	var res ResponseGetClothesHeader
	category := "all"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/header/%s/%s", category, order)
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
	assert.Equal(t, "Clothes not found", res.Message)

	// Seeder After Test
	seeders.SeedClothes(clothesRepo, userRepo, 200)
}

// Test Case ID : TC-E2E-CL-009
func TestFailedGetAllClothesHeaderWithInvalidClothesCategory(t *testing.T) {
	var res tests.ResponseSimple
	category := "clothes_source"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/header/%s/%s", category, order)
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
	assert.Equal(t, "Clothes category is not valid", res.Message)
}

// API GET : Get All Clothes Detail
func TemplateSuccessGetAllClothesDetail(t *testing.T, resp *http.Response, res ResponseGetClothesDetail) {
	// Get Template Test
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, res.Status)
	assert.Equal(t, "success", res.Status)
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Clothes fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.ClothesName)
		assert.NotEmpty(t, dt.ClothesColor)
		assert.NotEmpty(t, dt.ClothesQty)
		assert.NotEmpty(t, dt.ClothesType)
		assert.NotEmpty(t, dt.ClothesCategory)
		assert.NotEmpty(t, dt.ClothesSize)
		assert.NotEmpty(t, dt.ClothesGender)
		assert.NotEmpty(t, dt.ClothesPrice)
		assert.NotEmpty(t, dt.ClothesMadeFrom)
		assert.NotEmpty(t, dt.CreatedAt)
		if dt.ClothesBuyAt != nil {
			assert.NotEmpty(t, dt.ClothesBuyAt)
		}
		if dt.ClothesDesc != nil {
			assert.NotEmpty(t, dt.ClothesDesc)
		}
		if dt.ClothesMerk != nil {
			assert.NotEmpty(t, dt.ClothesMerk)
		}
		if dt.ClothesImage != nil {
			assert.NotEmpty(t, dt.ClothesImage)
		}
		if dt.UpdatedAt != nil {
			assert.NotEmpty(t, dt.UpdatedAt)
		}
		if dt.DeletedAt != nil {
			assert.NotEmpty(t, dt.DeletedAt)
		}

		// Check Data Type
		assert.IsType(t, "", dt.ClothesName)
		assert.IsType(t, "", dt.ClothesColor)
		assert.IsType(t, "", dt.ClothesType)
		assert.IsType(t, "", dt.ClothesCategory)
		assert.IsType(t, "", dt.ClothesSize)
		assert.IsType(t, "", dt.ClothesGender)
		assert.IsType(t, "", dt.ClothesMadeFrom)
		assert.IsType(t, 1, dt.ClothesQty)
		assert.IsType(t, true, dt.IsFaded)
		assert.IsType(t, true, dt.HasWashed)
		assert.IsType(t, true, dt.HasIroned)
		assert.IsType(t, true, dt.IsFavorite)
		assert.IsType(t, true, dt.IsScheduled)
		assert.IsType(t, time.Time{}, dt.CreatedAt)
		if dt.ClothesImage != nil {
			assert.IsType(t, "", *dt.ClothesImage)
		}
		if dt.ClothesDesc != nil {
			assert.IsType(t, "", *dt.ClothesDesc)
		}
		if dt.ClothesMerk != nil {
			assert.IsType(t, "", *dt.ClothesMerk)
		}
		if dt.ClothesPrice != nil {
			assert.IsType(t, 1, *dt.ClothesPrice)
		}
		if dt.UpdatedAt != nil {
			assert.IsType(t, time.Time{}, *dt.UpdatedAt)
		}
		if dt.DeletedAt != nil {
			assert.IsType(t, time.Time{}, *dt.DeletedAt)
		}
		if dt.ClothesBuyAt != nil {
			assert.IsType(t, time.Time{}, *dt.ClothesBuyAt)
		}
	}
}

// Test Case ID : TC-E2E-CL-011
func TestSuccessGetAllClothesDetailWithValidData(t *testing.T) {
	var res ResponseGetClothesDetail
	category := "all"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/detail/%s/%s", category, order)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Template Success Get All Clothes Detail
	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)
	TemplateSuccessGetAllClothesDetail(t, resp, res)
}

// Test Case ID : TC-E2E-CL-012
func TestSuccessGetAllClothesDetailWithValidClothesCategory(t *testing.T) {
	var res ResponseGetClothesDetail
	category := "upper_body"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/detail/%s/%s", category, order)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Template Success Get All Clothes Detail
	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)
	TemplateSuccessGetAllClothesDetail(t, resp, res)
}

// Test Case ID : TC-E2E-CL-013
func TestFailedGetAllClothesDetailWithEmptyData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	userRepo := repositories.NewUserRepository(db)
	clothesRepo := repositories.NewClothesRepository(db)

	// Precondition
	clothesRepo.DeleteAll()

	var res ResponseGetClothesHeader
	category := "all"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/detail/%s/%s", category, order)
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
	assert.Equal(t, "Clothes not found", res.Message)

	// Seeder After Test
	seeders.SeedClothes(clothesRepo, userRepo, 200)
}

// Test Case ID : TC-E2E-CL-014
func TestFailedGetAllClothesDetailWithInvalidClothesCategory(t *testing.T) {
	var res tests.ResponseSimple
	category := "clothes_source"
	order := "desc"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/detail/%s/%s", category, order)
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
	assert.Equal(t, "Clothes category is not valid", res.Message)
}

// API GET : Get Deleted Clothes
// Test Case ID : TC-E2E-CL-015
func TestSuccessGetDeletedClothesWithValidData(t *testing.T) {
	var res ResponseGetClothesDeleted
	url := "http://127.0.0.1:9000/api/v1/clothes/trash"
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
	assert.Equal(t, "Clothes fetched", res.Message)
	assert.NotNil(t, res.Data)

	for _, dt := range res.Data {
		// Check Object
		assert.NotEqual(t, uuid.Nil, dt.ID)
		assert.NotEmpty(t, dt.ClothesName)
		assert.NotEmpty(t, dt.ClothesColor)
		assert.NotEmpty(t, dt.ClothesQty)
		assert.NotEmpty(t, dt.ClothesType)
		assert.NotEmpty(t, dt.ClothesCategory)
		assert.NotEmpty(t, dt.ClothesSize)
		assert.NotEmpty(t, dt.ClothesGender)
		assert.NotEmpty(t, dt.DeletedAt)
		if dt.ClothesImage != nil {
			assert.NotEmpty(t, *dt.ClothesImage)
		}

		// Check Data Type
		assert.IsType(t, "", dt.ClothesName)
		assert.IsType(t, "", dt.ClothesColor)
		assert.IsType(t, "", dt.ClothesType)
		assert.IsType(t, "", dt.ClothesCategory)
		assert.IsType(t, 1, dt.ClothesQty)
		assert.IsType(t, "", dt.ClothesSize)
		assert.IsType(t, "", dt.ClothesGender)
		assert.IsType(t, time.Time{}, dt.DeletedAt)
		if dt.ClothesImage != nil {
			assert.IsType(t, "", *dt.ClothesImage)
		}
	}
}

// Test Case ID : TC-E2E-CL-016
func TestFailedGetDeletedClothesWithEmptyData(t *testing.T) {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	userRepo := repositories.NewUserRepository(db)
	clothesRepo := repositories.NewClothesRepository(db)

	// Precondition
	clothesRepo.DeleteAll()

	var res tests.ResponseSimple
	url := "http://127.0.0.1:9000/api/v1/clothes/trash"
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
	assert.Equal(t, "Clothes not found", res.Message)

	// Seeder After Test
	seeders.SeedClothes(clothesRepo, userRepo, 200)
}

// API PUT : Recover Deleted Clothes By Id
// Test Case ID : TC-E2E-CL-017
func TestSuccessRecoverDeletedClothesByIdWithValidId(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "d31ae086-18fa-42e5-b96e-2fbebe89ee06"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/recover/%s", clothesId)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("PUT", url, nil)
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
	assert.Equal(t, "Clothes recovered", res.Message)
}

// Test Case ID : TC-E2E-CL-018
func TestFailedRecoverDeletedClothesByIdWithInvalidId(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "93c7d7bf-3aa3-4859-916b-5415ac45bb4b"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/recover/%s", clothesId)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("PUT", url, nil)
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
	assert.Equal(t, "Clothes not found", res.Message)
}

// Test Case ID : TC-E2E-CL-019
func TestFailedRecoverDeletedClothesByIdWithInvalidUUID(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "93c7d7bf-3aa3-4859-916b"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/recover/%s", clothesId)
	token, _ := tests.TemplatePostBasicLogin(t, nil, nil, "user")

	// Exec
	req, err := http.NewRequest("PUT", url, nil)
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
	assert.Equal(t, "Invalid id", res.Message)
}

// API DELETE : Soft Delete Clothes By Id
// Test Case ID : TC-E2E-CL-020
func TestSuccessSoftDeleteClothesByIdWithValidId(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "06acbfe3-936f-4e61-a0d0-320c1685749d"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/%s", clothesId)
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Clothes deleted", res.Message)
}

// Test Case ID : TC-E2E-CL-021
func TestFailedSoftDeleteClothesByIdWithInvalidId(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "93c7d7bf-3aa3-4859-916b-5415ac45bb4b"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/%s", clothesId)
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Clothes not found", res.Message)
}

// Test Case ID : TC-E2E-CL-022
func TestFailedSoftDeleteClothesByIdWithInvalidUUID(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "93c7d7bf-3aa3-4859-916b"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/%s", clothesId)
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Invalid id", res.Message)
}

// API DELETE : Hard Delete Clothes By Id
// Test Case ID : TC-E2E-CL-023
func TestSuccessHardDeleteClothesByIdWithValidId(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "06acbfe3-936f-4e61-a0d0-320c1685749d"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/destroy/%s", clothesId)
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Clothes permanentally deleted", res.Message)
}

// Test Case ID : TC-E2E-CL-024
func TestFailedHardDeleteClothesByIdWithInvalidId(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "93c7d7bf-3aa3-4859-916b-5415ac45bb4b"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/destroy/%s", clothesId)
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Clothes not found", res.Message)
}

// Test Case ID : TC-E2E-CL-025
func TestFailedHardDeleteClothesByIdWithInvalidUUID(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "93c7d7bf-3aa3-4859-916b"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/destroy/%s", clothesId)
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "Invalid id", res.Message)
}

// Test Case ID : TC-E2E-CL-026
func TestFailedHardDeleteClothesByIdWithForbiddenRole(t *testing.T) {
	var res tests.ResponseSimple
	clothesId := "06acbfe3-936f-4e61-a0d0-320c1685749d"
	url := fmt.Sprintf("http://127.0.0.1:9000/api/v1/clothes/destroy/%s", clothesId)
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
	assert.NotEmpty(t, res.Message)
	assert.Equal(t, "access forbidden for this role", res.Message)
}
