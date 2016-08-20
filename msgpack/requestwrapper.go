package msgpack

import (
	"github.com/kernelguy/gorpc"
	msgpack2 "gopkg.in/vmihailenco/msgpack.v2"
	"gopkg.in/vmihailenco/msgpack.v2/codes"
	"fmt"
)


type RequestWrapper struct {
	gorpc.RequestWrapper
}

func (p *RequestWrapper) MarshalMsgpack() (result []byte, err error) {
	if p.IsBatchRequest() {
		result, err = msgpack2.Marshal(p.GetBatchRequests()) 
	} else if !p.IsEmpty() {
		result, err = msgpack2.Marshal(p.GetRequest()) 
	} else {
		err = fmt.Errorf("No requests in wrapper!!")
		result = nil
	}
	return
}

func (this *RequestWrapper) DecodeMsgpack(dec *msgpack2.Decoder) error {
	c, err := dec.PeekCode()
	if err != nil {
		return err
	}
	if codes.IsFixedArray(c) || (c == codes.Array16) || (c == codes.Array32) {
		l := int(c & codes.FixedArrayMask)
		rw := make([]Request, l)
		err := dec.Decode(&rw)
		if err != nil {
			return err
		}
		for i:=0 ; i < l ; i++ {
			this.AddRequest(&rw[i])
		}
		this.SetBatchRequest(true)
	} else {
		var r Request
		err := dec.Decode(&r)
		if err != nil {
			return err
		}
		this.AddRequest(&r)
	}
	return nil
}
