package gorpc

import (
	"encoding/json"
	"fmt"
)


type RequestWrapper struct {
	requests []IRequest
	isBatch bool
}

func (this *RequestWrapper) AddRequest(request IRequest) IRequestWrapper {
	this.requests = append(this.requests, request)
	if len(this.requests) > 1 {
		this.isBatch = true
	}
	return this
}

func (this *RequestWrapper) Clear() IRequestWrapper {
	this.requests = make([]IRequest, 0)
	return this
}

func (this *RequestWrapper) SetBatchRequest(value bool) IRequestWrapper {
	this.isBatch = value
	return this
}

func (p *RequestWrapper) IsEmpty() bool {
	return (len(p.requests) == 0)
}

func (p *RequestWrapper) IsBatchRequest() bool {
	return (len(p.requests) > 1) || p.isBatch
}

func (p *RequestWrapper) GetRequest() IRequest {
	return p.requests[0]
}

func (p *RequestWrapper) GetBatchRequests() []IRequest {
	return p.requests
}

func (p *RequestWrapper) MarshalJSON() (result []byte, err error) {
	if p.IsBatchRequest() {
		result, err = json.Marshal(p.requests) 
	} else if !p.IsEmpty() {
		result, err = json.Marshal(p.requests[0]) 
	} else {
		err = fmt.Errorf("No requests in wrapper!!")
		result = nil
	}
	return
}

