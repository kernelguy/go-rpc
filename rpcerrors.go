package gorpc

import (
	"fmt"
)

const cParseError = -32700  	// An error occurred on the server while parsing the JSON text.
const cInvalidRequest = -32600 	// The JSON sent is not a valid Request object.
const cMethodNotFound = -32601 	// The method does not exist / is not available.
const cInvalidParams = -32602 	// Invalid method parameter(s).
const cInternalError = -32603	// Internal error
const cNotAllowed = -32000		// Not allowed custom exception


type RpcError struct {
	code int
	message string
	data error
}


func NewRpcError(code int, previous error) IRpcError {
	var message string
	switch code {
		case cParseError:
			message = "Parse Error"

		case cInvalidRequest:
			message = "Invalid Request"

		case cMethodNotFound:
			message = "Method Not Found"

		case cInvalidParams:
			message = "Invalid Params"

		case cNotAllowed:
			message = "Not Allowed"
			
		default:
			message = "Internal Error"
			code = cInternalError
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

