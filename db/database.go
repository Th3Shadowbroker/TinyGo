package db

import (
	"context"
	"tiny-go/util"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client  *mongo.Client
	Context context.Context
	Cancel  context.CancelFunc

	Collections = map[string]*mongo.Collection{}
)

func ConnectToDatabase(config *util.Configuration) {
	var err error

	Client, Context, Cancel, err = inititalizeConnection(&config.Database)
	if err != nil {
		panic(err)
	}

	// Link collection
	collection := Client.Database(config.Database.Database).Collection("links")
	collection.Indexes().CreateOne(Context, mongo.IndexModel{
		Keys:    map[string]int{"shortId": 1},
		Options: options.Index().SetUnique(true),
	})
	Collections["links"] = collection
}

func GetLinkCollection() *mongo.Collection {
	return Collections["links"]
}

func inititalizeConnection(config *util.DatabaseConfiguration) (*mongo.Client, context.Context, context.CancelFunc, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.ConnectionString()))
	if err != nil {
		return nil, nil, nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, cancel, err
	}

	return client, ctx, cancel, nil
}
