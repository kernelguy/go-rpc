package gorpc

import (
	"strconv"
	"sync"
	"testing"
)

func BenchmarkRpcEcho(b *testing.B) {
	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	conn := transport.CreateTestConnections()

	r, err := conn.Call("Echo", []interface{}{"Hello World"})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
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
}

func BenchmarkRpcNoParams(b *testing.B) {
	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	conn := transport.CreateTestConnections()

	r, err := conn.Call("NoParams", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r, err = conn.Call("NoParams", nil)
		if err != nil {
			break
		}
	}

	if err != nil {
		b.Errorf("We should not get an error here: %v", err)
	} else if r != nil {
		b.Errorf("Result Mismatch: %v", r)
	}

	transport.quit <- true
}

func BenchmarkRpcTwoParams(b *testing.B) {
	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	conn := transport.CreateTestConnections()

	r, err := conn.Call("TwoParams", []interface{}{666, "Hello World"})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r, err = conn.Call("TwoParams", []interface{}{666, "Hello World"})
		if err != nil {
			break
		}
	}

	if err != nil {
		b.Errorf("We should not get an error here: %v", err)
	} else if r.(float64) != 42 {
		b.Errorf("Result Mismatch: %v", r)
	}

	transport.quit <- true
}

const cParallelTests = 16

func BenchmarkRpcParallel(b *testing.B) {
	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	var wg sync.WaitGroup
	
	var con [cParallelTests]IConnection
	
	for i:=0; i < cParallelTests ; i++ {
		con[i] = transport.addConnection(f.MakeAddress(strconv.Itoa(i+1), strconv.Itoa(i+1), nil))
	}
	
	m := func(con IConnection) {
		defer wg.Done();
		for i := 0; i < b.N; i++ {
			r, err := con.Call("TwoParams", []interface{}{666, "Hello World"})
			if err != nil {
				b.Errorf("Error should be nil: %v", err)
				break
			} else if r.(float64) != 42 {
				b.Errorf("Result Mismatch: %v", r)
				break
			}
		}
	}

	b.ResetTimer()
	for i:=0; i < cParallelTests ; i++ {
		wg.Add(1)
		go m(con[i])
	} 
	wg.Wait()

	transport.quit <- true
}
