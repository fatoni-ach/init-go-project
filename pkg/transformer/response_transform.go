package transformer

import "net/http"

type ResponseSuccess struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type ResponseFailed struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
}

func TransformResponse(response *ResponseSuccess, data interface{}) {
	response.Status = http.StatusOK
	response.Data = data
}

func TransformResponseFailed(response *ResponseFailed, statusCode int, message string, err string) {
	response.Status = statusCode
	response.Data = nil
	response.Message = message
	response.Error = err
}
