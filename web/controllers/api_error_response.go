package controllers

import "encoding/json"

// CreateApiError creates ApiErrorResponse struct instance with error message.
func createApiError(err error) ApiErrorResponse {
	return ApiErrorResponse{ Error: true, Message: err.Error() }
}

// AsJson generates JSON serialized string of ApiErrorResponse object
func (aer ApiErrorResponse) AsJson() string {
	data, e := json.Marshal(aer)
	if e != nil {
		return string(data)
	} else {
		return ""
	}
}
