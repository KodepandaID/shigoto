package shigoto

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	cronparser "github.com/KodepandaID/shigoto/pkg/cron-parser"
	"github.com/KodepandaID/shigoto/pkg/mongodb-connector"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Jobs instance
type Jobs struct {
	client    *mongodb.Connector
	parser    *cronparser.Parser
	JobName   string
	FuncName  string
	JobParams []interface{}
	Cron      []string // Set run a jobs with periodic by second, minute and hour
}

// Do to run a schedule command
func (j *Jobs) Do() (id primitive.ObjectID, e error) {
	schedule, eFatal := j.parser.SetCurrentTime(time.Now()).Parse(j.Cron)
	if eFatal != nil {
		panic(eFatal)
	}

	totalTask := 0
	if len(j.JobParams) > 0 {
		totalTask = 1
	}

	id, e = j.client.InsertJobCollection(&mongodb.JobCollection{
		JobName:    j.JobName,
		FuncName:   j.FuncName,
		CronFormat: j.Cron,
		TotalTask:  totalTask,
	})

	if id != primitive.NilObjectID && e == nil || id != primitive.NilObjectID && e.Error() == "Jobs is already registered, use the different job name" {
		e = nil
		j.storedTask(id, schedule)
	}

	return id, e
}

func callFunc(funcName string) (e error) {
	f := reflect.ValueOf(funcStorage[funcName])
	if !f.IsValid() {
		return errors.New("Function invalid, check your function register")
	}

	values := f.Call([]reflect.Value{})
	e = handleErrFunc(values)

	return e
}

func callFuncWithParams(funcName string, params []interface{}) (e error) {
	f := reflect.ValueOf(funcStorage[funcName])

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	values := f.Call(in)
	e = handleErrFunc(values)

	return e
}

func handleErrFunc(values []reflect.Value) (e error) {
	for i, val := range values {
		if val.Type().String() == "error" && !val.IsNil() {
			e = fmt.Errorf("%s", values[i])
		}
	}

	return e
}

func (j *Jobs) storedTask(id primitive.ObjectID, schedule cronparser.Schedule) {
	if scheduleStorage[schedule.Next.String()] == nil {
		scheduleStorage[schedule.Next.String()] = []map[string]interface{}{
			{
				"id":        id.Hex(),
				"job_name":  j.JobName,
				"func_name": j.FuncName,
				"params":    j.JobParams,
				"cron":      j.Cron,
			},
		}
		j.client.InsertTask(id, j.JobParams...)
	} else {
		ss := scheduleStorage[schedule.Next.String()].([]map[string]interface{})

		var sameParams bool
		for _, task := range ss {
			// check if the task having the same params,
			// if they have the same params, it will be ignored.
			if reflect.DeepEqual(task["params"].([]interface{}), j.JobParams) && task["job_name"].(string) == j.JobName {
				sameParams = true
			}
		}

		if !sameParams {
			ss = append(ss, map[string]interface{}{
				"id":        id.Hex(),
				"job_name":  j.JobName,
				"func_name": j.FuncName,
				"params":    j.JobParams,
				"cron":      j.Cron,
			})
			j.client.InsertTask(id, j.JobParams...)
		}
	}
}
