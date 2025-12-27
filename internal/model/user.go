package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Preferences struct {
	Perspective     string `bson:"perspective" json:"perspective"`
	ShowCoordinates bool   `bson:"show_coordinates" json:"show_coordinates"`
	Theme           string `bson:"theme" json:"theme"`
}

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string        `bson:"email" json:"email"`
	Username     string        `bson:"username" json:"username"`
	PasswordHash string        `bson:"password_hash" json:"-"`
	Preferences  Preferences   `bson:"preferences" json:"preferences"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}

func NewUser(email, username, passwordHash string) *User {
	now := time.Now()
	return &User{
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		Preferences: Preferences{
			Perspective:     "white",
			ShowCoordinates: true,
			Theme:           "light",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
