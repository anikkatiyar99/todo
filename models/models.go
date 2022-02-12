package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          *string            `json:"name" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Token         *string            `json:"token"`
	UserType      *string            `json:"usertype" validate:"required,eq=ADMIN|eq=USER"`
	Refresh_token *string            `json:"refresh_token"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	User_id       string             `json:"user_id"`
}

// Task is the model that governs all notes objects retrived or inserted into the DB
type Task struct {
	ID          primitive.ObjectID `bson:"id"`
	UserID      string             `json:"userId" bson:"userId"`
	Title       *string            `json:"title" validate:"required,min=2,max=100" bson:"title"`
	Description *string            `json:"description" validate:"required,min=2,max=100" bson:"description"`
	DueDate     time.Time          `json:"dueDate" bson:"dueDate"`
	AlertBefore *int               `json:"alertBefore" bson:"alertBefore"`
	AlertEmail  string             `json:"alertEmail" bson:"alertEmail"`
	AlertAt     time.Time          `json:"alertAt" bson:"alertAt"`
	AlertStatus string             `json:"alertStatus,omitempty" bson:"alertStatus"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Status      string             `json:"status" bson:"status"`
	SubTasks    []SubTask          `json:"subTasks,omitempty" bson:"subTasks"`
}

// SubTask is the model that governs all notes objects retrived or inserted into the DB
type SubTask struct {
	Title       *string `json:"title" validate:"required,min=2,max=100" bson:"title"`
	Description *string `json:"description" validate:"required,min=2,max=1000" bson:"description"`
	Status      string  `json:"status" bson:"status"`
}
