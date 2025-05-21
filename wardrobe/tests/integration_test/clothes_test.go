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

func TestGetClothesLastHistory(t *testing.T) {
	var res ResponseGetClothesLastHistory
	url := "http://127.0.0.1:9000/api/v2/clothes/last_history"
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
	assert.Equal(t, "clothes last history fetched", res.Message)
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

func TestGetClothesUsedHistory(t *testing.T) {
	var res ResponseGetUsedHistory
	url := "http://127.0.0.1:9000/api/v2/clothes/history/all/desc"
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
	assert.Equal(t, "clothes fetched", res.Message)
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
			assert.IsType(t, "", dt.ClothesNote)
		}
	}
}

func TestGetAllClothesHeader(t *testing.T) {
	var res ResponseGetClothesHeader
	url := "http://127.0.0.1:9000/api/v2/clothes/header/all/desc"
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
	assert.Equal(t, "clothes fetched", res.Message)
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
