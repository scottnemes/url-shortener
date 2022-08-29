package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/url-shortener/internal/model"
	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

/*
	Returns a valid cache client for use by other functions.
*/
func GetCacheClient(host string, port string, db int, pass string) *redis.Client {
	cacheClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: pass,
		DB:       db,
	})
	return cacheClient
}

/*
	Adds or updates the target URL in the cache based on the provided short URL slug.
*/
func SetCachedUrl(f *os.File, debug bool, expireHours time.Duration, client *redis.Client, url model.Url) error {
	log.SetOutput(f)
	json, err := json.Marshal(url)
	if err != nil {
		log.Printf("Error marshalling cached URL (slug: %v) (%v)", url.Slug, err)
	}
	err = client.Set(ctx, url.Slug, json, expireHours*time.Hour).Err()
	if err != nil {
		log.Printf("Error setting cached URL (slug: %v) (%v)", url.Slug, err)
	}

	if debug {
		log.Printf("[DEBUG] Inserted URL in cache (slug: %v) (target: %v)", url.Slug, url.Target)
	}
	return err
}

/*
	Checks the cache for the provided short URL slug and returns the target URL if available.
*/
func GetCachedUrl(f *os.File, debug bool, client *redis.Client, slug string) (model.Url, error) {
	log.SetOutput(f)
	url := model.Url{}
	result, err := client.Get(ctx, slug).Result()
	if err == redis.Nil {
		if debug {
			log.Printf("[DEBUG] Attempted to get missing URL from cache (slug: %v)", slug)
		}
		return url, err
	} else if err != nil && err != redis.Nil {
		log.Printf("Error getting cached URL (slug: %v) (%v)", slug, err)
		return url, err
	}
	err = json.Unmarshal([]byte(result), &url)
	if err != nil {
		log.Printf("Error unmarshalling cached URL (slug: %v) (%v)", slug, err)
	}

	if debug {
		log.Printf("[DEBUG] Got URL from cache (slug: %v) (target: %v)", url.Slug, url.Target)
	}

	return url, err
}

/*
	Removes the cached URL record for the provided short URL slug if present.
*/
func DeleteCachedUrl(f *os.File, debug bool, client *redis.Client, slug string) error {
	log.SetOutput(f)
	err := client.Del(ctx, slug).Err()
	if err != nil && err != redis.Nil {
		log.Printf("Error deleting cached URL (slug: %v) (%v)", slug, err)
	}

	if debug {
		log.Printf("[DEBUG] Deleted URL from cache (slug: %v)", slug)
	}
	return err
}
