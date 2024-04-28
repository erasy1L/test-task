package models

import "time"

type Session struct {
	RefreshToken string    `bson:"refresh_token"`
	ExpiresAt    time.Time `bson:"expires_at"`
} // @name Session
