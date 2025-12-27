package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// AuthSession represents an authenticated user session
type AuthSession struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
	Token     string        `bson:"token" json:"token"`
	ExpiresAt time.Time     `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

func NewAuthSession(userID bson.ObjectID, token string, maxAge int) *AuthSession {
	now := time.Now()
	return &AuthSession{
		UserID:    userID,
		Token:     token,
		ExpiresAt: now.Add(time.Duration(maxAge) * time.Second),
		CreatedAt: now,
	}
}

func (s *AuthSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
