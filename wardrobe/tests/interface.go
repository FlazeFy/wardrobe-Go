package tests

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ResponseFailedValidation struct {
	Message []ValidationError `json:"message"`
	Status  string            `json:"status"`
}
