package gorpc

import (
    "testing"
)

type mySimpleStruct struct {
	Name string
	Parents map[int]string
}

type myAdvancedStruct struct {
	I int
	M map[string]interface{}
	S []interface{}
	P *mySimpleStruct
	n int
	C chan int
}

func TestInitializeStruct(t *testing.T) {
	beginTest("TestInitializeStruct")

	mas := myAdvancedStruct{}
	
	InitializeStruct(&mas)

	mas.M["test"] = "Hello"
	if len(mas.S) != 0 {
		t.Errorf("Slice I is not correct: (%T)%v", mas,mas)
	}
	mas.P.Name = "Alfred"
	mas.P.Parents[0] = "Mother"
	mas.P.Parents[1] = "Father"

	go func() {
		mas.C <- 42
	}()
	if <-mas.C != 42 {
		t.Error("Channel is not working.")
	}
	
	endTest()
}

func TestFillStructFromMap(t *testing.T) {
	beginTest("TestFillStructFromMap")

	mas := myAdvancedStruct{}
	
	InitializeStruct(&mas)

	M := make(map[string]interface{})
	M["Hello"] = "World"
	M["Answer"] = 42

	parents := make(map[int]string)
	parents[0] = "Father"
	parents[1] = "Mother"
	
	P := make(map[string]interface{})
	P["Name"] = "Alfred"
	P["Parents"] = parents
	
    m := make(map[string]interface{})
    m["I"] = 42
    m["M"] = M
    m["S"] = []interface{}{1,42,"Hello","World"}
    m["P"] = P

	FillStructFromMap(&mas, m)

	if mas.I != 42 {
		t.Error("I is wrong: (%T)%v", mas.I, mas.I)
	}

	if mas.M["Answer"] != 42 {
		t.Error("M is wrong: (%T)%v", mas.M, mas.M)
	}

	if mas.S[1] != 42 {
		t.Error("S is wrong: (%T)%v", mas.S, mas.S)
	}

	if mas.P.Name != "Alfred" {
		t.Error("P is wrong: (%T)%v", mas.P, mas.P)
	}

	if mas.P.Parents[1] != "Mother" {
		t.Error("P.Parents is wrong: (%T)%v", mas.P.Parents, mas.P.Parents)
	}

	endTest()
}

func TestFillStructFromArray(t *testing.T) {
	beginTest("TestFillStructFromArray")

	mas := myAdvancedStruct{}
	
	InitializeStruct(&mas)

	M := make(map[string]interface{})
	M["Hello"] = "World"
	M["Answer"] = 42

	parents := make(map[int]string)
	parents[0] = "Father"
	parents[1] = "Mother"
	
	P := make(map[string]interface{})
	P["Name"] = "Alfred"
	P["Parents"] = parents
	
    m := make([]interface{}, 4)
    m[0] = 42
    m[1] = M
    m[2] = []interface{}{1,42,"Hello","World"}
    m[3] = P

	FillStructFromArray(&mas, m)

	if mas.I != 42 {
		t.Error("I is wrong: (%T)%v", mas.I, mas.I)
	}

	if mas.M["Answer"] != 42 {
		t.Error("M is wrong: (%T)%v", mas.M, mas.M)
	}

	if mas.S[1] != 42 {
		t.Error("S is wrong: (%T)%v", mas.S, mas.S)
	}

	if mas.P.Name != "Alfred" {
		t.Error("P is wrong: (%T)%v", mas.P, mas.P)
	}

	if mas.P.Parents[1] != "Mother" {
		t.Error("P.Parents is wrong: (%T)%v", mas.P.Parents, mas.P.Parents)
	}

	endTest()
}

func TestFillStructErrors(t *testing.T) {
	beginTest("TestFillStructErrors")

	mas := myAdvancedStruct{}
	
	InitializeStruct(&mas)

    m := make(map[string]interface{})
    m["Q"] = 42
    
	err := FillStructFromMap(&mas, m)
	if err == nil {
		t.Error("It should have failed on wrong field name.")
	}

    m = make(map[string]interface{})
    m["n"] = 42
    
	err = FillStructFromMap(&mas, m)
	if err == nil {
		t.Error("It should have failed on private field name.")
	}

    a := make([]interface{}, 5)
    a[0] = 0
    a[1] = m
    a[2] = []interface{}{0}
    a[3] = nil
    a[4] = 42

	err = FillStructFromArray(&mas, a)
	if err == nil {
		t.Error("It should have failed on private field name.")
	}
	if mas.P != nil {
		t.Error("P should have been nil.")
	}

	endTest()    
}	
