package models

import "time"

// User mirrors Mongo document structure
// Indexes (unique) for username, email should be ensured via MongoDB

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"-"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
