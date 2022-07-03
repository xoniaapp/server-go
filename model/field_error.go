package model

type ErrorsResponse struct {
	Errors []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error HttpError `json:"error"`
}

type HttpError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}