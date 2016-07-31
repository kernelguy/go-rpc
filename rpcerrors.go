package gorpc

import (
	"fmt"
)

const ErrParseError = -32700     // An error occurred on the server while parsing the JSON text.
const ErrInvalidRequest = -32600 // The JSON sent is not a valid Request object.
const ErrMethodNotFound = -32601 // The method does not exist / is not available.
const ErrInvalidParams = -32602  // Invalid method parameter(s).
const ErrInternalError = -32603  // Internal error
const ErrNotAllowed = -32000     // Not allowed custom exception

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    error  `json:"data"`
}

func NewRpcError(code int, previous error) IRpcError {
	var message string
	switch code {
	case ErrParseError:
		message = "Parse Error"

	case ErrInvalidRequest:
		message = "Invalid Request"

	case ErrMethodNotFound:
		message = "Method Not Found"

	case ErrInvalidParams:
		message = "Invalid Params"

	case ErrNotAllowed:
		message = "Not Allowed"

	default:
		message = "Internal Error"
		code = ErrInternalError
	}
	return &RpcError{Code: code, Message: message, Data: previous}
}

func (p *RpcError) GetCode() int {
	return p.Code
}

func (p *RpcError) GetMessage() string {
	return p.Message
}

func (p *RpcError) GetData() error {
	return p.Data
}

func (p *RpcError) Error() string {
	var data interface{}
	if p.Data != nil {
		data = p.Data.Error()
	} else {
		data = ""
	}
	return fmt.Sprintf("Code: %d, Message: %s, Data: %s", p.Code, p.Message, data)
}

//func (this *RpcError) MarshalJSON() (result []byte, err error) {
//	result, err = json.Marshal(this.data)
//	log.Debugf("Request Encoded: %s, %v", string(result), err)
//	return result, err
//}
