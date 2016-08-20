package msgpack

import (
	"github.com/kernelguy/gorpc"
	msgpack2 "gopkg.in/vmihailenco/msgpack.v2"
	log "github.com/Sirupsen/logrus"
)

type Request struct {
	gorpc.Request
}

var _ msgpack2.CustomEncoder = (*Request)(nil)
var _ msgpack2.CustomDecoder = (*Request)(nil)

func (this *Request) EncodeMsgpack(enc *msgpack2.Encoder) error {
	data := this.GetData()
	log.Debugf("Request.EncodeMsgpack: (%T)%v", data,data)
	return enc.Encode(data)
}

func (this *Request) DecodeMsgpack(dec *msgpack2.Decoder) error {
	var data map[string]interface{}
	err := dec.Decode(&data)
	if err != nil {
		return err
	}
	log.Debugf("Request.DecodeMsgpack: (%T)%v", data,data)
	this.Populate(data)
	return nil
}
