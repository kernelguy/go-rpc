package gorpc

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

const ErrParseError = -32700     // An error occurred on the server while parsing the JSON text.
const ErrInvalidRequest = -32600 // The JSON sent is not a valid Request object.
const ErrMethodNotFound = -32601 // The method does not exist / is not available.
const ErrInvalidParams = -32602  // Invalid method parameter(s).
const ErrInternalError = -32603  // Internal error
const ErrNotAllowed = -32000     // Not allowed custom exception

type RpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
	return p.Data.(error)
}

func (p *RpcError) getDataString() interface{} {
	if p.Data != nil {
		return p.Data.(error).Error()
	}
	return nil
}

func (p *RpcError) Error() string {
	var s string
	if p.Data != nil {
		s = p.Data.(error).Error()
	} else {
		s = "nil"
	}
	return fmt.Sprintf("Code: %d, Message: %s, Data: %s", p.Code, p.Message, s)
}

func (this *RpcError) MarshalJSON() (result []byte, err error) {
	e := *this
	
	if e.Data != nil {
		e.Data = e.getDataString()
	}
	result, err = json.Marshal(e)
	log.Debugf("RpcError Encoded: %s, %v", string(result), err)
	return result, err
}
