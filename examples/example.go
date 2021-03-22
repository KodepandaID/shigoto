package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KodepandaID/shigoto"
	"github.com/joho/godotenv"
)

func init() {
	wd, _ := os.Getwd()

	e := godotenv.Load(wd + "/.env")
	if e != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	client, _ := shigoto.New(&shigoto.Config{
		DB:     os.Getenv("MONGO_URI"),
		DBName: "jobs-scheduler",
	})
	client.Register("hello", hello)
	client.Register("helloWorld", helloWorld)
	client.Command("run-hello", "hello", "usman").EveryMinute().Do()
	client.Command("run-hello-world", "helloWorld").EveryMinute().Do()
	client.Run()
}

func hello(name string) error {
	fmt.Printf("Hello, %s\n", name)
	return nil
}

func helloWorld() error {
	fmt.Println("Hello World")
	return nil
}
