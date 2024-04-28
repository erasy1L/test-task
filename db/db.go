package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

func (m *MongoDB) Client() *mongo.Client {
	if m.client == nil {
		return nil
	}
	return m.client
}

func (m *MongoDB) Database() *mongo.Database {
	if m.database == nil {
		return nil
	}
	return m.database
}

func (m *MongoDB) ConnectDB(ctx context.Context, dbName, connectionString string) (db *MongoDB, err error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return db, err
	}

	log.Println("Connected to MongoDB")

	go func() {
		<-ctx.Done()
		client.Disconnect(ctx)
	}()

	return &MongoDB{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.client != nil {
		return m.client.Disconnect(ctx)
	}
	return nil
}
