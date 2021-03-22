package shigoto

import (
	"fmt"
	"math"
	"time"

	cronparser "github.com/KodepandaID/shigoto/pkg/cron-parser"
	"github.com/KodepandaID/shigoto/pkg/mongodb-connector"
)

func loadJobsFromPersistentStorage(c *Config) {
	jobs, e := c.client.GetJobCollection()
	if e != nil {
		panic(e)
	}

	for _, job := range jobs {
		loc, _ := time.LoadLocation(c.Timezone)
		tnow := time.Now().Local().In(loc)
		nextDate := tnow

		if tnow.Unix() > job.NextDate.Unix() {
			schedule, eFatal := c.parser.SetCurrentTime(tnow).Parse(job.CronFormat)
			if eFatal != nil {
				panic(eFatal)
			}
			nextDate = schedule.Next
		}

		tasks, e := c.client.GetTasks(job.ID)
		if e != nil {
			panic(e)
		}

		var newTasks []map[string]interface{}
		for _, task := range tasks {
			newTasks = append(newTasks, map[string]interface{}{
				"id":        task.JobId.Hex(),
				"job_name":  job.JobName,
				"func_name": job.FuncName,
				"params":    task.Params,
				"cron":      job.CronFormat,
			})
		}
		scheduleStorage[nextDate.String()] = newTasks
	}
}

func checkTask(c *Config) {
	loc, _ := time.LoadLocation(c.Timezone)
	for {
		tnow := time.Now().Local().In(loc)
		if scheduleStorage[tnow.String()] != nil {
			ss := scheduleStorage[tnow.String()].([]map[string]interface{})
			if len(ss) > 0 {
				var e error
				for _, task := range ss {
					funcName := task["func_name"].(string)
					params := task["params"].([]interface{})
					if len(params) == 0 {
						e = callFunc(funcName)
					} else {
						e = callFuncWithParams(funcName, params)
					}
					updateJob(c, tnow, task["job_name"].(string), e)
				}
				updateNextRun(tnow.String(), c.Timezone, ss[0]["cron"].([]string), ss)
			}
		}
	}
}

// To create the next schedule after the success running.
// The old schedule will be removed.
func updateNextRun(key, timezone string, cron []string, tasks []map[string]interface{}) time.Time {
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

	return schedule.Next
}

func updateJob(c *Config, tnow time.Time, jobName string, e error) {
	go func() {
		eInc := 0
		if e != nil {
			eInc = 1
		}

		job, e := c.client.GetOneJobCollection(jobName)
		if e == nil {
			successRate, errRate := countSuccessAndErrorRate(float64(job.TotalRun+1), float64(job.TotalError+eInc))

			schedule, eFatal := c.parser.SetCurrentTime(tnow).Parse(job.CronFormat)
			if eFatal != nil {
				panic(eFatal)
			}

			c.client.UpdateJobCollection(job.ID, &mongodb.JobCollection{
				NextDate:    schedule.Next,
				SuccessRate: successRate,
				ErrorRate:   errRate,
			}, eInc)
		}
		fmt.Println("TETEW")
	}()
}

func countSuccessAndErrorRate(totalRun, totalError float64) (success float64, err float64) {
	success = ((totalRun - totalError) / totalRun) * 100
	err = (totalError / totalRun) * 100

	success = math.Ceil(success*100) / 100
	err = math.Ceil(err*100) / 100

	return success, err
}
