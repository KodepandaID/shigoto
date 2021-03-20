package shigoto

import (
	"errors"
	"fmt"
	"reflect"

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
	if len(j.JobParams) == 0 {
		if e = callFunc(j.FuncName); e != nil {
			return id, e
		}
	} else {
		if e = callFuncWithParams(j.FuncName, j.JobParams); e != nil {
			return id, e
		}
	}

	schedule, eFatal := j.parser.Parse(j.Cron)
	if eFatal != nil {
		panic(eFatal)
	}

	id, e = j.client.InsertJobCollection(&mongodb.JobCollection{
		JobName:     j.JobName,
		FuncName:    j.FuncName,
		CronFormat:  j.Cron,
		NextDate:    schedule.Next,
		TotalTask:   1,
		SuccessRate: 0,
		ErrorRate:   0,
	})
	if e == nil && len(j.JobParams) > 0 || id != primitive.NilObjectID && len(j.JobParams) > 0 {
		e = nil
		e = j.client.InsertTask(id, j.JobParams...)
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
