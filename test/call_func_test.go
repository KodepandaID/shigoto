package test

import (
	"errors"
	"testing"

	"github.com/KodepandaID/shigoto"
)

var funcStorage = make(map[string]interface{})

func TestCallFunc(t *testing.T) {
	shigoto.FuncStorage["hello-call-func"] = helloCallFunc
	if e := shigoto.CallFunc("hello-call-func"); e != nil {
		t.Fatal(e)
		t.Fail()
	}
}

func TestCallFuncError(t *testing.T) {
	shigoto.FuncStorage["hello-call-func-err"] = helloCallFuncError
	if e := shigoto.CallFunc("hello-call-func-err"); e == nil {
		t.Fatal("This test should be error")
		t.Fail()
	}
}

func TestCallFuncParams(t *testing.T) {
	shigoto.FuncStorage["hello-call-func-params"] = helloCallFuncParams
	if e := shigoto.CallFuncWithParams("hello-call-func-params", []interface{}{"usman"}); e != nil {
		t.Fatal(e)
		t.Fail()
	}
}

func TestCallFuncParamsError(t *testing.T) {
	shigoto.FuncStorage["hello-call-func-params-err"] = helloCallFuncParamsError
	if e := shigoto.CallFuncWithParams("hello-call-func-params-err", []interface{}{"usman"}); e == nil {
		t.Fatal("This test should be error")
		t.Fail()
	}
}

func helloCallFunc() error {
	return nil
}

func helloCallFuncError() error {
	return errors.New("Error")
}

func helloCallFuncParams(name string) error {
	return nil
}

func helloCallFuncParamsError(name string) error {
	return errors.New("Error")
}
