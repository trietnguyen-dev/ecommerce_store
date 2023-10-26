package database

import (
	"context"
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UserColl  string = "users"
	AdminColl string = "admin"
)

// InitDB :
func InitDB(conf *config.Config) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(conf.DBUri)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(conf.DBName)

	if err = CreateIndexes(db); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateIndexes :
func CreateIndexes(db *mongo.Database) error {
	userCollection := db.Collection(UserColl)
	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: -1},
				{Key: "phone_number", Value: 1},
			},
		},
	}

	ctx, cancel := utils.NewCtx()
	defer cancel()

	_, err := userCollection.Indexes().CreateMany(ctx, userIndexes)

	if err != nil {
		return err
	}
	return nil
}
