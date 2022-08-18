package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example.com/url-shortener/internal/config"
	"example.com/url-shortener/internal/logging"
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
	log.SetOutput(logging.F)
	url := model.Url{}
	result, err := client.Get(ctx, slug).Result()
	if err == redis.Nil {
		if config.DebugMode {
			log.Printf("[DEBUG] Attempted to get missing URL from cache (slug: %v)", slug)
		}
		return url
	} else if err != nil && err != redis.Nil {
		log.Printf("Error getting cached URL (slug: %v) (%v)", slug, err)
		return url
	}
	err = json.Unmarshal([]byte(result), &url)
	if err != nil {
		log.Printf("Error unmarshalling cached URL (slug: %v) (%v)", slug, err)
	}

	if config.DebugMode {
		log.Printf("[DEBUG] Got URL from cache (slug: %v) (target: %v)", url.Slug, url.Target)
	}

	return url
}

func SetCachedUrl(client *redis.Client, url model.Url) {
	log.SetOutput(logging.F)
	json, err := json.Marshal(url)
	if err != nil {
		log.Printf("Error marshalling cached URL (slug: %v) (%v)", url.Slug, err)
	}
	err = client.Set(ctx, url.Slug, json, config.CacheExpireHours*time.Hour).Err()
	if err != nil {
		log.Printf("Error setting cached URL (slug: %v) (%v)", url.Slug, err)
	}

	if config.DebugMode {
		log.Printf("[DEBUG] Inserted URL in cache (slug: %v) (target: %v)", url.Slug, url.Target)
	}
}

func DeleteCachedUrl(client *redis.Client, slug string) {
	log.SetOutput(logging.F)
	err := client.Del(ctx, slug).Err()
	if err != nil && err != redis.Nil {
		log.Printf("Error deleting cached URL (slug: %v) (%v)", slug, err)
	}

	if config.DebugMode {
		log.Printf("[DEBUG] Deleted URL from cache (slug: %v)", slug)
	}
}
