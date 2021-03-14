package test

import (
	"testing"

	"github.com/KodepandaID/shigoto"
)

func TestCreateInstance(t *testing.T) {
	if _, e := shigoto.New(&shigoto.Config{
		DB:     "mongodb://localhost:27017",
		DBName: "jobs-scheduler",
	}); e != nil {
		t.Fatal(e)
		t.Fail()
	}
}

func TestErrorCreateInstance(t *testing.T) {
	if _, e := shigoto.New(&shigoto.Config{
		DB:     "mongodb://123",
		DBName: "jobs-scheduler",
	}); e == nil {
		t.Error("Test should be fail")
		t.Fail()
	}
}
