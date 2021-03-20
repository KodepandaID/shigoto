package test

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/KodepandaID/shigoto"
	"github.com/joho/godotenv"
)

var mongoURI string

func init() {
	wd, _ := os.Getwd()
	p := filepath.Dir(wd)

	e := godotenv.Load(p + "/.env")
	if e != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI = os.Getenv("MONGO_URI")
}

func TestCreateInstance(t *testing.T) {
	if _, e := shigoto.New(&shigoto.Config{
		DB:     os.Getenv("MONGO_URI"),
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

func TestDoSchedule(t *testing.T) {
	client, e := shigoto.New(&shigoto.Config{
		DB:     os.Getenv("MONGO_URI"),
		DBName: "jobs-scheduler",
	})
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}
	client.Register("hello", hello)
	_, e = client.Command("run-hello", "hello", "usman").Daily().Do()
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}
}

func TestDoScheduleWithoutParams(t *testing.T) {
	client, e := shigoto.New(&shigoto.Config{
		DB:     os.Getenv("MONGO_URI"),
		DBName: "jobs-scheduler",
	})
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}
	client.Register("hello", helloWithoutParams)
	_, e = client.Command("run-hello-without-params", "hello").Daily().Do()
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}
}

func TestDoScheduleError(t *testing.T) {
	client, e := shigoto.New(&shigoto.Config{
		DB:     os.Getenv("MONGO_URI"),
		DBName: "jobs-scheduler",
	})
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}
	client.Register("hello", helloErr)
	_, e = client.Command("run-hello", "hello", "usman").Daily().Do()
	if e == nil {
		t.Log("This test should be fail")
		t.Fail()
	}
}

func hello(name string) error {
	return nil
}

func helloErr(name string) error {
	return errors.New("Test with error")
}

func helloWithoutParams() error {
	return nil
}
