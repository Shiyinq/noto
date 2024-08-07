package model

import (
	"time"
)

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Name      string    `json:"name" bson:"name"`
	PhotoURL  string    `json:"photo_url" bson:"photoUrl"`
	CreatedAt time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" bson:"updatedAt"`
}

type AuthToken struct {
	Token string `json:"token"`
}
