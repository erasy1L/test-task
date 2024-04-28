package db

import (
	"context"
	"errors"
	"log"

	"github.com/erazr/test-task/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Collection) (*UserRepository, error) {
	ctx := context.Background()
	mod := mongo.IndexModel{
		Keys:    bson.M{"guid": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := db.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &UserRepository{db: db}, err
}

func (u *UserRepository) Create(ctx context.Context, user models.User) error {
	_, err := u.db.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("user already exists")
	}

	return err
}

func (u *UserRepository) GetByGUID(ctx context.Context, guid string) (models.User, error) {
	var user models.User
	err := u.db.FindOne(ctx, bson.D{{Key: "guid", Value: guid}}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepository) SetSession(ctx context.Context, guid string, session models.Session) error {
	_, err := u.db.UpdateOne(ctx, bson.D{{Key: "guid", Value: guid}}, bson.D{{Key: "$set", Value: bson.D{{Key: "session", Value: session}}}})

	return err
}
