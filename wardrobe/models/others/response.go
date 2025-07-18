package others

import "wardrobe/models"

type (
	// Auth : BasicLogin
	ResponsePostBasicLogin struct {
		Message string    `json:"message" example:"User login"`
		Status  string    `json:"status" example:"success"`
		Data    LoginData `json:"data"`
	}
	LoginData struct {
		Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
		Role  string `json:"role" example:"user"`
	}
	// Auth : BasicSignOut
	ResponsePostBasicSignOut struct {
		Message string `json:"message" example:"User signed out"`
		Status  string `json:"status" example:"success"`
	}
	ResponseBadRequestBasicSignOut struct {
		Message string `json:"message" example:"missing authorization header"`
		Status  string `json:"status" example:"failed"`
	}
	// History : GetHistory
	ResponseGetHistory struct {
		Message  string              `json:"message" example:"History fetched"`
		Status   string              `json:"status" example:"success"`
		Data     []models.GetHistory `json:"data"`
		Metadata Metadata            `json:"metadata"`
	}
	// History : HardDeleteHistoryById
	ResponseHardDeleteHistoryById struct {
		Message string `json:"message" example:"history permanentally deleted"`
		Status  string `json:"status" example:"success"`
	}
	// For Response
	ResponseBadRequest struct {
		Message string `json:"message" example:"email is not valid"`
		Status  string `json:"status" example:"failed"`
	}
	ResponseBadRequestInvalidUserId struct {
		Message string `json:"message" example:"invalid user id"`
		Status  string `json:"status" example:"failed"`
	}
	ResponseNotFound struct {
		Message string `json:"message" example:"account not found"`
		Status  string `json:"status" example:"failed"`
	}
	ResponseInternalServerError struct {
		Message string `json:"message" example:"something went wrong"`
		Status  string `json:"status" example:"error"`
	}
)
