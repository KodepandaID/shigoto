# Shigoto - Task Scheduling
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/KodepandaID/shigoto)
![GitHub](https://img.shields.io/github/license/KodepandaID/shigoto)
![](https://github.com/KodepandaID/shigoto/workflows/Go/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/KodepandaID/shigoto/badge.svg?branch=main)](https://coveralls.io/github/KodepandaID/shigoto?branch=main)

Shigoto is a task scheduling for Golang with friendly API. Shigoto used MongoDB for persistent storage and used Cron format to set the schedule.

## Installation

```bash
go get github.com/KodepandaID/shigoto
```

## Example Usage
### Init
```go
package main

import "github.com/KodepandaID/shigoto"

func main() {
    client, e := shigoto.New(&shigoto.Config{
		DB:      "mongodb://localhost:27017",
		DBName:  "jobs-scheduler",
        Timezone: "Asia/Jakarta",
	})
    if e != nil {
        log.Fatal(e)
    }

    client.Run()
}
```

### Register a function to call
```go
func main() {
    client, e := shigoto.New(&shigoto.Config{
		DB:      "mongodb://localhost:27017",
		DBName:  "jobs-scheduler",
        Timezone: "Asia/Jakarta",
	})
    if e != nil {
        log.Fatal(e)
    }
    
    client.Register("hello", hello)
    client.Register("hello-params", helloParams, "message here")
    client.Run()
}

func hello() error {
    return nil
}

func helloParams(message string) error {
    return nil
}
```

### Set Job
For more scheduling information, you can read at [Go References](https://pkg.go.dev/github.com/KodepandaID/shigoto).
```go
func main() {
    client, e := shigoto.New(&shigoto.Config{
		DB:      "mongodb://localhost:27017",
		DBName:  "jobs-scheduler",
        Timezone: "Asia/Jakarta",
	})
    if e != nil {
        log.Fatal(e)
    }
    
    client.Register("hello", hello)
    if _, e := client.
        Command("job-name-here", "hello").
        EveryMinute().Do(); e != nil {
        log.Fatal(e)
    }
    client.Run()
}

func hello() error {
    return nil
}
```

### Remove a Job
```go
func main() {
    client, e := shigoto.New(&shigoto.Config{
		DB:      "mongodb://localhost:27017",
		DBName:  "jobs-scheduler",
        Timezone: "Asia/Jakarta",
	})
    if e != nil {
        log.Fatal(e)
    }
    
    client.Register("hello", hello)
    client.Delete("hello")
    client.Run()
}
```


## License

Copyright [Yudha Pratama Wicaksana](https://github.com/LordAur), Licensed under [MIT](./LICENSE).