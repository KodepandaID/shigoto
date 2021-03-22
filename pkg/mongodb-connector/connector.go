package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connector mongodb instance
type Connector struct {
	DB     string        // The MongoDB uri
	DBName string        // Database name from MongoDB
	client *mongo.Client // Mongodb client
}

type JobCollection struct {
	ID          primitive.ObjectID `bson:"_id"`
	JobName     string             `bson:"job_name"`
	FuncName    string             `bson:"func_name"`
	CronFormat  []string           `bson:"cron_format"`
	NextDate    time.Time          `bson:"next_date"`
	TotalTask   int                `bson:"total_task"`
	TotalRun    int                `bson:"total_run"`
	TotalError  int                `bson:"total_error"`
	SuccessRate float64            `bson:"success_rate"`
	ErrorRate   float64            `bson:"error_rate"`
}

type TaskCollection struct {
	JobId  primitive.ObjectID `bson:"job_id"`
	Params []interface{}      `bson:"params"`
}

var ctx context.Context

// New to create a new Mongodb connection
func New(c *Connector) (*Connector, error) {
	client, e := mongo.Connect(ctx, options.Client().ApplyURI(c.DB))
	if e != nil {
		return &Connector{}, e
	}

	return &Connector{
		DB:     c.DB,
		DBName: c.DBName,
		client: client,
	}, nil
}

// Ping to check connection status
func (c *Connector) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	e := c.client.Ping(ctx, readpref.Primary())

	return e
}

func (c *Connector) GetJobCollection() ([]JobCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var jobs []JobCollection
	cursor, e := c.client.Database(c.DBName).Collection("jobs").Find(ctx, bson.M{})
	if e != nil {
		return jobs, e
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &jobs)

	return jobs, nil
}

func (c *Connector) GetOneJobCollection(name string) (JobCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var jobs JobCollection
	e := c.client.Database(c.DBName).Collection("jobs").FindOne(ctx, bson.M{"job_name": name}).Decode(&jobs)

	return jobs, e
}

func (c *Connector) InsertJobCollection(payload *JobCollection) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var jobs bson.M
	c.client.Database(c.DBName).Collection("jobs").FindOne(ctx, bson.M{"job_name": payload.JobName}).Decode(&jobs)

	if jobs["job_name"] == payload.JobName {
		return jobs["_id"].(primitive.ObjectID), errors.New("Jobs is already registered, use the different job name")
	}

	res, e := c.client.Database(c.DBName).
		Collection("jobs").
		InsertOne(ctx, bson.D{{
			Key:   "job_name",
			Value: payload.JobName,
		}, {
			Key:   "func_name",
			Value: payload.FuncName,
		}, {
			Key:   "cron_format",
			Value: payload.CronFormat,
		}, {
			Key:   "total_task",
			Value: payload.TotalTask,
		}, {
			Key:   "total_run",
			Value: 0,
		}, {
			Key:   "total_error",
			Value: 0,
		}, {
			Key:   "next_date",
			Value: payload.NextDate,
		}, {
			Key:   "success_rate",
			Value: 0,
		}, {
			Key:   "error_rate",
			Value: 0,
		}})
	if e != nil {
		return primitive.NilObjectID, e
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (c *Connector) UpdateJobCollection(id primitive.ObjectID, payload *JobCollection, e int) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{
		"$inc": bson.M{"total_run": 1, "total_error": e},
		"$set": bson.M{
			"next_date":    payload.NextDate,
			"success_rate": payload.SuccessRate,
			"error_rate":   payload.ErrorRate,
		},
	}
	c.client.
		Database(c.DBName).
		Collection("jobs").UpdateOne(ctx, filter, update)
}

func (c *Connector) DeleteJobCollection(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var jobs bson.M
	c.client.Database(c.DBName).Collection("jobs").FindOne(ctx, bson.M{"job_name": name}).Decode(&jobs)

	c.client.Database(c.DBName).Collection("jobs").DeleteOne(ctx, bson.M{"job_name": name})
	c.client.Database(c.DBName).Collection("tasks").DeleteOne(ctx, bson.M{"job_id": bson.M{"$eq": jobs["_id"].(primitive.ObjectID)}})
}

func (c *Connector) GetTasks(id primitive.ObjectID) ([]TaskCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var tasks []TaskCollection
	filter := bson.M{"job_id": bson.M{"$eq": id}}
	cursor, e := c.client.Database(c.DBName).Collection("tasks").Find(ctx, filter)
	if e != nil {
		return tasks, e
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &tasks)

	return tasks, nil
}

func (c *Connector) InsertTask(id primitive.ObjectID, params ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var jobs bson.M
	filter := bson.M{"params": bson.M{"$eq": params}}
	c.client.Database(c.DBName).Collection("tasks").FindOne(ctx, filter).Decode(&jobs)

	if jobs == nil {
		if _, e := c.client.Database(c.DBName).
			Collection("tasks").
			InsertOne(ctx, bson.D{{
				Key:   "job_id",
				Value: id,
			}, {
				Key:   "params",
				Value: params,
			}}); e != nil {
			return e
		}

		// Update total_task at jobs
		filter = bson.M{"_id": bson.M{"$eq": id}}
		update := bson.M{"$inc": bson.M{"total_task": 1}}
		c.client.
			Database(c.DBName).
			Collection("jobs").UpdateOne(ctx, filter, update)
	}

	return nil
}
