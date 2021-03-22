package shigoto

import (
	"errors"
	"log"
	"time"

	cronparser "github.com/KodepandaID/shigoto/pkg/cron-parser"
	"github.com/KodepandaID/shigoto/pkg/mongodb-connector"
)

// Config to set the configuration task scheduler
type Config struct {
	DB          string // The MongoDB uri
	DBName      string // Database name from MongoDB
	Timezone    string
	MaxRunJobs  int // Maximal to run job process at the same time, the default is 10
	MaxRunQueue int // Maximal to run task queue on the job process, the default is 100
	client      *mongodb.Connector
	parser      cronparser.Parser
}

var scheduleStorage = make(map[string]interface{})
var funcStorage = make(map[string]interface{})

// New to create task scheduler instance
func New(c *Config) (*Config, error) {
	client, e := mongodb.New(&mongodb.Connector{
		DB:     c.DB,
		DBName: c.DBName,
	})
	if e != nil {
		return &Config{}, e
	}

	if e := client.Ping(); e != nil {
		return &Config{}, errors.New("MongoDB not connected")
	}
	c.client = client

	// Cause I'm Indonesian I will be set the default timezone with Asia/Jakarta
	if c.Timezone == "" {
		c.Timezone = "Asia/Jakarta"
	}
	c.parser = cronparser.New(&cronparser.Parser{
		Timezone: c.Timezone,
	})

	return c, nil
}

// Command to create a new job process
func (c *Config) Command(jobName, funcName string, params ...interface{}) *Jobs {
	if jobName == "" {
		log.Fatal("The job's name cannot be empty")
	}

	return &Jobs{
		client:    c.client,
		parser:    &c.parser,
		JobName:   jobName,
		FuncName:  funcName,
		JobParams: params,
		Cron:      []string{"*", "*", "*", "*", "*"},
	}
}

// Register to register a function to call with the name
// The funcName should be the same with funcName at Command function
func (c *Config) Register(funcName string, jobFunc interface{}) {
	funcStorage[funcName] = jobFunc
}

// Run to run a background process to check the tasks
func (c *Config) Run() {
	go checkTask(c.Timezone)
	select {}
}

func checkTask(timezone string) {
	loc, _ := time.LoadLocation(timezone)
	for {
		tnow := time.Now().Local().In(loc).String()
		if scheduleStorage[tnow] != nil {
			ss := scheduleStorage[tnow].([]map[string]interface{})
			if len(ss) > 0 {
				for _, task := range ss {
					funcName := task["func_name"].(string)
					params := task["params"].([]interface{})
					if len(params) == 0 {
						callFunc(funcName)
					} else {
						callFuncWithParams(funcName, params)
					}
				}
				updateNextRun(tnow, timezone, ss[0]["cron"].([]string), ss)
			}
		}
	}
}

// To create the next schedule after the success running.
// The old schedule will be removed.
func updateNextRun(key, timezone string, cron []string, tasks []map[string]interface{}) {
	delete(scheduleStorage, key)

	parser := cronparser.New(&cronparser.Parser{
		Timezone: timezone,
	})
	loc, _ := time.LoadLocation(timezone)
	schedule, eFatal := parser.SetCurrentTime(time.Now().Local().In(loc)).Parse(cron)
	if eFatal != nil {
		panic(eFatal)
	}

	var newTasks []map[string]interface{}
	for _, task := range tasks {
		newTasks = append(newTasks, map[string]interface{}{
			"id":        task["id"].(string),
			"job_name":  task["job_name"].(string),
			"func_name": task["func_name"].(string),
			"params":    task["params"].([]interface{}),
			"cron":      task["cron"].([]string),
		})
	}

	scheduleStorage[schedule.Next.String()] = newTasks
}
