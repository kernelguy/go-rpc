package gorpc

import (
    "testing"
)

func BenchmarkRpcEcho(b *testing.B) {
//	beginTest("BenchmarkRpcEcho")

	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
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

	transport.quit <- true

//	endTest()
}
