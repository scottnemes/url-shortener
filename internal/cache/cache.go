package cache

import (
	"context"
	"fmt"
	"log"

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
	targetUrl, err := client.Get(ctx, slug).Result()
	url.Target = targetUrl
	if err != nil {
		log.Printf("Error looking up cached target URL (slug: %v) (%v)", slug, err)
	}

	return url
}

func SetCachedUrl(client *redis.Client, url model.Url) {
	err := client.Set(ctx, url.Slug, url.Target, 0).Err()
	if err != nil {
		log.Printf("Error setting cached target URL (slug: %v) (%v)", url.Slug, err)
	}
}

func DeleteCachedUrl(client *redis.Client, slug string) {
	err := client.Del(ctx, slug).Err()
	if err != nil && err != redis.Nil {
		log.Printf("Error deleting cached target URL (slug: %v) (%v)", slug, err)
	}
}
