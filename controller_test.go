package gorpc

import (
    "testing"
    "time"
)

type embedObj struct {
	S1 string
	S2 int
}

type funcParams struct {
	P1 string
	P2 int
	P3 embedObj
}

type RpcController struct {
	Result funcParams 
}

type MyRootController struct {
	Controller
	Cont1 RpcController
}

func (this *RpcController) RPC_NotifyTest(params funcParams) {
	this.Result = params
}


func TestController(t *testing.T) {
	beginTest("TestController")

	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	transport.addConnection(f.MakeAddress("1", "2", nil))
	transport.addConnection(f.MakeAddress("2", "1", nil))

	conn1, _ := transport.getConnection("1")
	conn1.SetRootController(&MyRootController{})

	param := funcParams{P1:"Test", P2:42, P3:embedObj{"Wow", 13}}
	
	conn2, _ := transport.getConnection("2")
	conn2.Notify("Cont1.NotifyTest", param)

	time.Sleep(time.Millisecond * 1)

	if conn1.RootController().(*MyRootController).Cont1.Result != param {
		t.Errorf("Params was not passed correctly: %v != %v", param, conn1.RootController().(*MyRootController).Cont1.Result) 
	}

	transport.quit <- true

	endTest()
}


func TestAdvancedParams(t *testing.T) {
	beginTest("TestAdvancedParams")

	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	transport.addConnection(f.MakeAddress("1", "2", nil))
	transport.addConnection(f.MakeAddress("2", "1", nil))

	conn1, _ := transport.getConnection("1")
	conn1.SetRootController(&MyRootController{})

	param := []interface{}{"Test", 42, embedObj{"Wow", 13}}
	
	conn2, _ := transport.getConnection("2")
	conn2.Notify("Cont1.NotifyTest", param)

	time.Sleep(time.Millisecond * 1)

	p := funcParams{P1:"Test", P2:42, P3:embedObj{"Wow", 13}}
	
	if conn1.RootController().(*MyRootController).Cont1.Result != p {
		t.Errorf("Params was not passed correctly: %v != %v", p, conn1.RootController().(*MyRootController).Cont1.Result) 
	}

	transport.quit <- true

	endTest()
}


func TestClose(t *testing.T) {
	beginTest("TestRpcEcho")

	f := GetFactory()
	options := f.MakeTransportOptions()
	if _, ok := options.(ITransportOptions); !ok {
		t.Errorf("options should be of type ITransportOption, not (%T)%v", options, options)
	}
	
	transport := f.MakeTransport(options).(*LoopbackTransport)

	transport.addConnection(f.MakeAddress("1", "2", nil))

	conn, _ := transport.getConnection("1")
	if len(transport.connections) != 1 {
		t.Error("There is not exactly one connection.")
	}

	conn.Close()
	if len(transport.connections) != 0 {
		t.Error("There is not exactly zero connections.")
	}

	endTest()
}
