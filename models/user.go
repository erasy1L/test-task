package models

type User struct {
	Session Session `json:"session,omitempty" bson:"session,omitempty" swaggerignore:"true"`
	GUID    string  `json:"guid" bson:"guid"`
} // @name User
