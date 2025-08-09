package models

import "time"

type Comment struct {
	Text      string    `bson:"text" json:"text"`
	CreatedBy string    `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type Task struct {
	ID            string     `bson:"_id,omitempty" json:"id"`
	Title         string     `bson:"title" json:"title"`
	Description   string     `bson:"description" json:"description"`
	Status        string     `bson:"status" json:"status"`
	Priority      string     `bson:"priority" json:"priority"`
	Assignee      *string    `bson:"assignee" json:"assignee"`
	CreatedBy     string     `bson:"createdBy" json:"createdBy"`
	CreatedAt     time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `bson:"updatedAt" json:"updatedAt"`
	Deadline      *time.Time `bson:"deadline" json:"deadline"`
	ScheduledDate *time.Time `bson:"scheduledDate" json:"scheduledDate"`
	Comments      []Comment  `bson:"comments" json:"comments"`
}
