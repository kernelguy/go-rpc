package gorpc

import (
	"fmt"
)

const ErrParseError = -32700  		// An error occurred on the server while parsing the JSON text.
const ErrInvalidRequest = -32600 	// The JSON sent is not a valid Request object.
const ErrMethodNotFound = -32601 	// The method does not exist / is not available.
const ErrInvalidParams = -32602 	// Invalid method parameter(s).
const ErrInternalError = -32603		// Internal error
const ErrNotAllowed = -32000		// Not allowed custom exception


type RpcError struct {
	code int
	message string
	data error
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
	return &RpcError{code: code, message: message, data: previous}
}


func (p *RpcError) GetCode() int {
	return p.code
}

func (p *RpcError) GetMessage() string {
	return p.message
}

func (p *RpcError) GetData() error {
	return p.data
}

func (p *RpcError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Data: %s", p.code, p.message, p.data.Error());
}

