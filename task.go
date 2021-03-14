package shigoto

import (
	"errors"
	"log"

	"github.com/KodepandaID/shigoto/pkg/mongodb-connector"
)

// Config to set the configuration task scheduler
type Config struct {
	DB          string // The MongoDB uri
	DBName      string // Database name from MongoDB
	Timezone    string
	MaxRunJobs  int // Maximal to run job process at the same time, the default is 10
	MaxRunQueue int // Maximal to run task queue on the job process, the default is 100
}

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

	return c, nil
}

// Command to create a new job process
func (c *Config) Command(jobName string, jobFunc interface{}, params ...interface{}) *Jobs {
	if jobName == "" {
		log.Fatal("The job's name cannot be empty")
	}

	return &Jobs{
		JobName:   jobName,
		JobFunc:   jobFunc,
		JobParams: params,
		Cron:      []string{"*", "*", "*", "*", "*"},
	}
}
