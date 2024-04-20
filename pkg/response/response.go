package response

import "fmt"

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	statusOK                   = "OK"
	statusError                = "Error"
	internalServerErrorMessage = "Internal server error"
	badRequestErrorMessage     = "Bad request"
)

func OK() Response {
	return Response{
		Status: statusOK,
	}
}

func InternalError() Response {
	return Error(internalServerErrorMessage)
}

func BadRequest(msg string) Response {
	return Error(fmt.Sprintf("%s: %s", badRequestErrorMessage, msg))
}

func Error(msg string) Response {
	return Response{
		Status: statusError,
		Error:  msg,
	}
}
