package test

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/KodepandaID/shigoto"
	"github.com/joho/godotenv"
)

func init() {
	wd, _ := os.Getwd()
	p := filepath.Dir(wd)

	e := godotenv.Load(p + "/.env")
	if e != nil {
		log.Fatal("Error loading .env file")
	}
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

func TestDeleteJob(t *testing.T) {
	client, e := shigoto.New(&shigoto.Config{
		DB:     os.Getenv("MONGO_URI"),
		DBName: "jobs-scheduler",
	})
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}
	client.Delete("run-hello")
	if len(shigoto.ScheduleStorage) > 0 {
		t.Fatal("Schedule storage should be nil")
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
	client.Delete("run-hello-without-params")
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
	client.Register("hello-err", helloErr)
	_, e = client.Command("run-hello-err", "hello-err", "usman").Daily().Do()
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}

	client.Delete("run-hello-err")
}

func TestMultipleDo(t *testing.T) {
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

	client.Register("hello-err", helloErr)
	_, e = client.Command("run-hello-err", "hello-err", "usman").Daily().Do()
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}

	client.Delete("run-hello-without-params")
	client.Delete("run-hello-err")
}

func TestRun(t *testing.T) {
	client, e := shigoto.New(&shigoto.Config{
		DB:      os.Getenv("MONGO_URI"),
		DBName:  "jobs-scheduler",
		Timeout: time.Minute,
	})
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}

	client.Register("hello", helloRun)
	_, e = client.Command("hello-run", "hello").EveryMinute().Do()
	if e != nil {
		t.Fatal(e)
		t.Fail()
	}

	client.Run()
	client.Delete("hello-run")
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

func helloRun() error {
	return nil
}
