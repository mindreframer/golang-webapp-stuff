package util

import (
	"strconv"
)

type ReturnData struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error_message"`
	JsonData []byte `json:"json_data"`
	Status   string `json:"status"`
}

func (return_data *ReturnData) GetSuccess() string {
	return strconv.FormatBool(return_data.Success)
}

func (return_data *ReturnData) GetJsonData() string {
	if len(return_data.JsonData) == 0 {
		return "{}"
	}
	return string(return_data.JsonData)
}

func (return_data *ReturnData) GetErrorMessage() string {
	return return_data.ErrorMsg
}

func (return_data *ReturnData) GetStatus() string {
	return return_data.Status
}

func (return_data *ReturnData) ToString() string {
	str := "{\"data\":{" +
		"\"success\":\"" + return_data.GetSuccess() + "\"" +
		",\"error_msg\":\"" + return_data.GetErrorMessage() + "\"" +
		",\"json_data\":" + return_data.GetJsonData() +
		",\"status\":\"" + return_data.GetStatus() + "\"" +
		"}}"
	return str
}
