package gorpc

import (

)


type RequestWrapper struct {
	requests []IRequest
	isBatch bool
}

func (p *RequestWrapper) AddRequest(request IRequest) {
	p.requests = append(p.requests, request)
	if len(p.requests) > 1 {
		p.isBatch = true
	}
}

func (p *RequestWrapper) SetBatchRequest(value bool) {
	p.isBatch = value
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
