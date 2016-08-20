package msgpack

import (
	"bytes"
    "testing"
    "github.com/kernelguy/gorpc"
)

type testFactory struct {
	Factory
}

func (this *testFactory) MakeTransport(options gorpc.ITransportOptions) gorpc.ITransport {
	t := &gorpc.LoopbackTransport{Transport:gorpc.Transport{Options: options, Protocol: this.MakeProtocol()}}
	t.SetFactory(this)
	t.Init(nil, nil)
	return t
}

func BenchmarkRpcEcho(b *testing.B) {
	gorpc.SetFactory(&testFactory{})
	
	f := gorpc.GetFactory()
	transport := f.MakeTransport(nil).(*gorpc.LoopbackTransport)
	transport.Start()

	conn := transport.CreateTestConnections()

	r, err := conn.Call("Echo", []interface{}{[]byte{1,2,3,4,5}})

	b.ResetTimer()

	for i:=0; i < b.N; i++ {
		r, err = conn.Call("Echo", []interface{}{[]byte{1,2,3,4,5}})
		if err != nil {
			break
		}
	}

	if err != nil {
		b.Errorf("We should not get an error here: %v", err)
	} else if !bytes.Equal([]byte(r.(string)), []byte{1,2,3,4,5}) {
		b.Errorf("Result Mismatch: %v", r)
	}

	transport.Stop()
}
