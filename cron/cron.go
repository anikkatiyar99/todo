package cron

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/anikkatiyar99/todo/database"
	helper "github.com/anikkatiyar99/todo/helpers"
	"github.com/anikkatiyar99/todo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskCollection *mongo.Collection = database.OpenCollection(database.Client, "task")

func CronEmailRunner() {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var task models.Task
	fmt.Println("Cron Email Runner is Listening...", time.Now())

	filter := bson.M{}

	filter["alertAt"] = bson.M{"$lte": time.Now()}
	filter["alertStatus"] = "created"

	// findOptions for sorting the tasks
	findOptions := options.Find()
	findOptions.SetSort(map[string]int{"dueDate": 1})

	cursor, err := taskCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Println(err)
	}

	for cursor.Next(ctx) {
		err := cursor.Decode(&task)
		if err != nil {
			log.Println(err)
		}
		err = helper.SendEmail(&task)
		if err != nil {
			log.Println(err)
		}
		task.AlertStatus = "sent"
		taskCollection.ReplaceOne(ctx, bson.M{"id": task.ID}, task)
		fmt.Println("Email sent to: ", task.AlertEmail)
	}
}
