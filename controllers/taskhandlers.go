package controllers

import (
	"context"

	"github.com/anikkatiyar99/todo/models"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/anikkatiyar99/todo/database"

	helper "github.com/anikkatiyar99/todo/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskCollection *mongo.Collection = database.OpenCollection(database.Client, "task")

// CreateTask is the api used to create a task for a user
func CreateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var task models.Task

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(task)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		task.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		task.ID = primitive.NewObjectID()
		task.Status = "pending"
		task.AlertAt = task.DueDate.Add(time.Duration(*task.AlertBefore) * time.Hour * -1)
		task.AlertStatus = "created"
		task.UserID = c.GetString("uid")
		if task.AlertEmail == "" {
			task.AlertEmail = c.GetString("email")
		}
		for i := range task.SubTasks {
			task.SubTasks[i].Status = "pending"
		}

		_, insertErr := taskCollection.InsertOne(ctx, task)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "task item was not created"})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

// UpdateTask is the api used to get update a task for a task
func UpdateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var task models.Task

		taskId := c.Param("task_id")
		userId := c.GetString("uid")

		primTaskId, _ := primitive.ObjectIDFromHex(taskId)
		err := taskCollection.FindOne(ctx, bson.M{"id": primTaskId}).Decode(&task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if task.UserID != userId {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(task)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		primId, _ := primitive.ObjectIDFromHex(taskId)

		task.UpdatedAt = time.Now()
		task.ID = primId
		task.UserID = c.GetString("uid")

		task.AlertAt = task.DueDate.Add(time.Duration(*task.AlertBefore) * time.Hour * -1)
		if task.Status == "completed" {
			for i := range task.SubTasks {
				task.SubTasks[i].Status = "completed"
			}
		}

		updatedResult, err := taskCollection.ReplaceOne(ctx, bson.M{"id": primId}, task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if updatedResult.ModifiedCount == 0 {
			c.JSON(http.StatusExpectationFailed, gin.H{"error": "task could not be updated"})
			return
		}

		c.JSON(http.StatusOK, task)

	}
}

// GetTasks is the api used to get all tasks
func GetTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var task models.Task
		userId := c.GetString("uid")
		filter := bson.M{}

		filter["userId"] = userId

		// filter for task status
		if c.Request.Header.Get("type") == "all" {
		} else if c.Request.Header.Get("type") == "completed" {
			filter["status"] = "completed"
		} else if c.Request.Header.Get("type") == "" {
			filter["status"] = "pending"
		}

		// filter for searching a task
		if c.Request.Header.Get("search") != "" {
			regexstr := ".*" + c.Request.Header.Get("search") + ".*"
			filter["$or"] = []bson.M{}
			filter["$or"] = append(filter["$or"].([]bson.M),
				bson.M{"title": bson.M{"$regex": regexstr, "$options": "i"}},
				bson.M{"description": bson.M{"$regex": regexstr, "$options": "i"}},
			)
		}

		// filter for find pending tasks based on due date
		if c.Request.Header.Get("dueDate") == "Today" {
			filter["$and"] = []bson.M{}
			filter["$and"] = append(filter["$and"].([]bson.M),
				bson.M{"dueDate": bson.M{"$gte": time.Now().UTC(), "$lt": time.Now().UTC().AddDate(0, 0, 1)}},
				bson.M{"status": bson.M{"$ne": "completed"}},
			)
		} else if c.Request.Header.Get("dueDate") == "This week" {
			filter["$and"] = []bson.M{}
			filter["$and"] = append(filter["$and"].([]bson.M),
				bson.M{"dueDate": bson.M{"$gte": time.Now().UTC(), "$lt": time.Now().UTC().AddDate(0, 0, 7)}},
				bson.M{"status": bson.M{"$ne": "completed"}},
			)
		} else if c.Request.Header.Get("dueDate") == "Next week" {
			filter["$and"] = []bson.M{}
			filter["$and"] = append(filter["$and"].([]bson.M),
				bson.M{"dueDate": bson.M{"$gte": time.Now().UTC().AddDate(0, 0, 7), "$lt": time.Now().UTC().AddDate(0, 0, 14)}},
				bson.M{"status": bson.M{"$ne": "completed"}},
			)
		} else if c.Request.Header.Get("dueDate") == "Overdue" {
			filter["$and"] = []bson.M{}
			filter["$and"] = append(filter["$and"].([]bson.M),
				bson.M{"dueDate": bson.M{"$lt": time.Now().UTC()}},
				bson.M{"status": bson.M{"$ne": "completed"}},
			)
		}

		// findOptions for sorting the tasks
		findOptions := options.Find()
		findOptions.SetSort(map[string]int{"dueDate": 1})

		err := helper.MatchUserTypeToUid(c, userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cursor, err := taskCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		taskList := []models.Task{}
		for cursor.Next(ctx) {
			err = cursor.Decode(&task)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			taskList = append(taskList, task)
		}

		c.JSON(http.StatusOK, taskList)

	}
}
