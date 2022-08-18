package model

import (
	"context"
	"log"
	"time"

	"example.com/url-shortener/internal/config"
	"example.com/url-shortener/internal/logging"
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
	log.SetOutput(logging.F)
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, url)
	if err != nil {
		log.Printf("Error creating new short URL (slug: %v) (%v)", url.Slug, err)
	}

	if config.DebugMode {
		log.Printf("[DEBUG] Inserted URL in database (slug: %v) (target: %v)", url.Slug, url.Target)
	}

	return err
}

func GetUrl(client *mongo.Client, slug string) (Url, error) {
	log.SetOutput(logging.F)
	url := Url{}
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"slug": bson.M{"$eq": slug}}
	opts := options.FindOne().SetSort(bson.M{"created": -1})

	err := collection.FindOne(ctx, filter, opts).Decode(&url)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Error looking up URL (slug: %v) (%v)", slug, err)
	}

	if config.DebugMode {
		if err == mongo.ErrNoDocuments {
			log.Printf("[DEBUG] Attempted to get missing URL from database (slug: %v)", slug)
		} else {
			log.Printf("[DEBUG] Got URL from database (slug: %v) (target: %v)", url.Slug, url.Target)
		}
	}

	return url, err
}

func GetUrls(client *mongo.Client) ([]Url, error) {
	log.SetOutput(logging.F)
	urls := []Url{}
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Error retrieving all URLs (%v)", err)
		return []Url{}, err
	}

	count := 0
	for cur.Next(ctx) {
		var url Url
		err := cur.Decode(&url)
		if err != nil {
			log.Printf("Error retrieving all URLs (%v)", err)
			return []Url{}, err
		}
		urls = append(urls, url)
		count += 1
	}

	if config.DebugMode {
		log.Printf("[DEBUG] Got URLs from database (count: %v)", count)
	}

	return urls, err
}

func UpdateUrl(client *mongo.Client, url Url) error {
	log.SetOutput(logging.F)
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"slug": url.Slug},
		bson.M{"$set": bson.M{"target": url.Target}},
	)
	if err != nil {
		log.Printf("Error updating target URL (slug: %v) (%v)", url.Slug, err)
	}

	if config.DebugMode {
		log.Printf("[DEBUG] Updated URL in database (slug: %v) (target: %v)", url.Slug, url.Target)
	}

	return err
}

func DeleteUrl(client *mongo.Client, slug string) error {
	log.SetOutput(logging.F)
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(
		ctx,
		bson.M{"slug": slug},
	)
	if err != nil {
		log.Printf("Error deleting URL (slug: %v) (%v)", slug, err)
	}

	if config.DebugMode {
		log.Printf("[DEBUG] Deleted URL from database (slug: %v)", slug)
	}

	return err
}

func UpdateUrlHits(client *mongo.Client, slug string) error {
	log.SetOutput(logging.F)
	collection := client.Database(config.DBDatabase).Collection(config.DBCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"slug": slug},
		bson.M{"$inc": bson.M{"hits": 1}},
	)
	if err != nil {
		log.Printf("Error updating URL hits (slug: %v) (%v)", slug, err)
	}

	if config.DebugMode {
		log.Printf("[DEBUG] Updated URL hits in database (slug: %v)", slug)
	}

	return err
}
