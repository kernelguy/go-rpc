package gorpc

import (
    "testing"
)

func TestFactories(t *testing.T) {

	backup := GetFactory()
	SetFactory(&Factory{})
	f := GetFactory();
	
	a := f.MakeAddress("", "", nil)
	o := f.MakeTransportOptions()
	tr := f.MakeTransport(o)
	f.MakeConnection(tr, a)
	f.MakeController()
	f.MakeProtocol()
	f.MakeRequest(1, "Method", nil)
	f.MakeRequestWrapper()
	f.MakeResponse(1, "result", nil)  
	f.MakeRouter(nil)
	f.MakeRpcError(1, nil)
	
	SetFactory(backup)
}

