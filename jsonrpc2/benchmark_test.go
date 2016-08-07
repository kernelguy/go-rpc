package jsonrpc2

import (
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
//	beginTest("BenchmarkRpcEcho")

	gorpc.SetFactory(&testFactory{})
	
	f := gorpc.GetFactory()
	transport := f.MakeTransport(nil).(*gorpc.LoopbackTransport)
	transport.Start()

	conn := transport.CreateTestConnections()

	r, err := conn.Call("Echo", []interface{}{"Hello World"})

	b.ResetTimer()

	for i:=0; i < b.N; i++ {
		r, err = conn.Call("Echo", []interface{}{"Hello World"})
		if err != nil {
			break
		}
	}

	if err != nil {
		b.Errorf("We should not get an error here: %v", err)
	} else if r.(string) != "Hello World" {
		b.Errorf("Result Mismatch: %v", r)
	}

	transport.Stop()

//	endTest()
}
