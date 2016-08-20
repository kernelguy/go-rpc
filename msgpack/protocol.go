package msgpack

import (
	"github.com/kernelguy/gorpc"
	msgpack2 "gopkg.in/vmihailenco/msgpack.v2"
)

type Protocol struct {
	gorpc.Protocol
}

func (this *Protocol) Encode(request gorpc.IRequestWrapper) ([]byte, error) {
	b, err := msgpack2.Marshal(request)
	return b, err
}

func (this *Protocol) Decode(data []byte) (gorpc.IRequestWrapper, error) {
	var result RequestWrapper
	err := msgpack2.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

