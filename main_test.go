package gorpc

import (
	"fmt"
	"github.com/kernelguy/gorpc/test"
	"testing"
)

type testFactory struct {
	Factory
}

type testController struct {
	Controller
}

func (this *testFactory) MakeTransport(options ITransportOptions) ITransport {
	t := &LoopbackTransport{Transport:Transport{Options: options, Protocol: this.MakeProtocol()}}
	t.SetFactory(this)
	t.Init(nil, nil)
	return t
}

func (this *testFactory) MakeController() IController {
	return &testController{}
}


func (this *testController) RPC_NoParams() {
}

func (this *testController) RPC_FailWithError() {
	panic(fmt.Errorf("Chained Error"))
}

func (this *testController) RPC_FailWithRPCError() {
	e := GetFactory().MakeRpcError(ErrInternalError, nil)
	panic(*e.(*RpcError))
}

func (this *testController) RPC_FailWithString() {
	panic("Chained String")
}

func (this *testController) RPC_FailWithInt() {
	panic(0)
}

func (this *testController) RPC_IllegalParams(p1 int, p2 string, p3 float64) {
}

type TwoParams struct {
	V1 int
	V2 string
}
func (this *testController) RPC_TwoParams(params TwoParams) int {
	return 42
}


func beginTest(name string) {
	test.Begin(name, 3)
}

func endTest() {
	test.End()
}

func TestMain(m *testing.M) {
	test.Init(m)
	SetFactory(&testFactory{})
	test.Run(m)
}
