package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportType string

const (
	DailyReport   ReportType = "daily"
	WeeklyReport  ReportType = "weekly"
	MonthlyReport ReportType = "monthly"
)

type ReportStatistics struct {
	TotalTasks       int `json:"totalTasks" bson:"totalTasks"`
	CompletedTasks   int `json:"completedTasks" bson:"completedTasks"`
	InProgressTasks  int `json:"inProgressTasks" bson:"inProgressTasks"`
	OverdueTasks     int `json:"overdueTasks" bson:"overdueTasks"`
	CompletionRate   int `json:"completionRate" bson:"completionRate"`
}

type Report struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
	Type            ReportType         `json:"type" bson:"type" validate:"required"`
	Period          string             `json:"period" bson:"period" validate:"required"`
	Title           string             `json:"title" bson:"title" validate:"required"`
	Content         string             `json:"content" bson:"content" validate:"required"`
	PolishedContent *string            `json:"polishedContent" bson:"polishedContent"`
	Tasks           []primitive.ObjectID `json:"tasks" bson:"tasks"`
	Statistics      ReportStatistics   `json:"statistics" bson:"statistics"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
}