package shigoto

import (
	"errors"
	"log"

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

	loadJobsFromPersistentStorage(c)

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

// Delete to remove job from instance and persistent storage
func (c *Config) Delete(name string) {
	for key, jobs := range scheduleStorage {
		j := jobs.([]map[string]interface{})
		for index, task := range j {
			if name == task["job_name"].(string) {
				j = append(j[:index], j[index+1:]...)
				c.client.DeleteJobCollection(name)
			}
		}
		scheduleStorage[key] = j
	}
}

// Run n a background process to check the tasks
func (c *Config) Run() {
	go checkTask(c)
	select {}
}
