package models

import "time"

type Statistics struct {
	TotalTasks      int `bson:"totalTasks" json:"totalTasks"`
	CompletedTasks  int `bson:"completedTasks" json:"completedTasks"`
	InProgressTasks int `bson:"inProgressTasks" json:"inProgressTasks"`
	OverdueTasks    int `bson:"overdueTasks" json:"overdueTasks"`
	CompletionRate  int `bson:"completionRate" json:"completionRate"`
}

type Report struct {
	ID              string     `bson:"_id,omitempty" json:"id"`
	UserID          string     `bson:"userId" json:"userId"`
	Type            string     `bson:"type" json:"type"`
	Period          string     `bson:"period" json:"period"`
	Title           string     `bson:"title" json:"title"`
	Content         string     `bson:"content" json:"content"`
	PolishedContent *string    `bson:"polishedContent" json:"polishedContent"`
	Tasks           []string   `bson:"tasks" json:"tasks"`
	Statistics      Statistics `bson:"statistics" json:"statistics"`
	CreatedAt       time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time  `bson:"updatedAt" json:"updatedAt"`
}
