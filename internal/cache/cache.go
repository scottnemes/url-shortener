package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example.com/url-shortener/internal/config"
	"example.com/url-shortener/internal/model"
	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

func GetCacheClient() *redis.Client {
	cacheClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.CacheHost, config.CachePort),
		Password: config.CachePass,
		DB:       config.CacheDB,
	})
	return cacheClient
}

func GetCachedUrl(client *redis.Client, slug string) model.Url {
	url := model.Url{}
	// short circuit if cache is disabled
	if config.CacheEnabled != true {
		return url
	}
	result, err := client.Get(ctx, slug).Result()
	if err != nil {
		log.Printf("Error looking up cached target URL (slug: %v) (%v)", slug, err)
	}
	err = json.Unmarshal([]byte(result), &url)
	if err != nil {
		log.Printf("Error looking up cached target URL (slug: %v) (%v)", slug, err)
	}

	return url
}

func SetCachedUrl(client *redis.Client, url model.Url) {
	// short circuit if cache is disabled
	if config.CacheEnabled != true {
		return
	}
	json, err := json.Marshal(url)
	if err != nil {
		log.Printf("Error setting cached target URL (slug: %v) (%v)", url.Slug, err)
	}
	err = client.Set(ctx, url.Slug, json, config.CacheExpireHours*time.Hour).Err()
	if err != nil {
		log.Printf("Error setting cached target URL (slug: %v) (%v)", url.Slug, err)
	}
}

func DeleteCachedUrl(client *redis.Client, slug string) {
	// short circuit if cache is disabled
	if config.CacheEnabled != true {
		return
	}
	err := client.Del(ctx, slug).Err()
	if err != nil && err != redis.Nil {
		log.Printf("Error deleting cached target URL (slug: %v) (%v)", slug, err)
	}
}
