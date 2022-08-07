package model

import (
	"context"
	"log"
	"time"

	"example.com/url-shortener/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Url struct {
	Slug    string `bson:"slug" json:"slug"`
	Target  string `bson:"target" json:"target"`
	Created uint64 `bson:"created" json:"created"`
	Hits    uint64 `bson:"hits" json:"hits"`
}

func GetDBClient() *mongo.Client {
	/*
		TODO:
		1. Setup database authentication
	*/

	// credentials := options.Credential{
	// 	Username: db_user,
	// 	Password: db_pass,
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//clientOptions := options.Client().ApplyURI(db_conn_string).SetAuth(credentials)
	clientOptions := options.Client().ApplyURI(config.DBConnString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	// verify that the database connection was established
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	return client
}

func InsertUrl(client *mongo.Client, url Url) error {
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, url)
	if err != nil {
		log.Println(err)
	}

	return err
}

func GetTargetUrl(client *mongo.Client, slug string) (Url, error) {
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"slug": bson.M{"$eq": slug}}
	opts := options.FindOne().SetSort(bson.M{"created": -1})
	url := Url{}

	err := collection.FindOne(ctx, filter, opts).Decode(&url)
	if err != nil {
		log.Printf("Error looking up target URL (slug: %v) (%v)", slug, err)
		return Url{}, err
	}

	// update the hit count for the given short URL
	err = UpdateUrlHits(client, slug)
	if err != nil {
		log.Println(err)
	}

	return url, err
}

func UpdateUrlHits(client *mongo.Client, slug string) error {
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"slug": slug},
		bson.M{"$inc": bson.M{"hits": 1}},
	)
	if err != nil {
		log.Println(err)
	}

	return err
}
