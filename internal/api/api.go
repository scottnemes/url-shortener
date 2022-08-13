package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/url-shortener/internal/cache"
	"example.com/url-shortener/internal/config"
	"example.com/url-shortener/internal/model"
	"example.com/url-shortener/internal/util"
	"github.com/gin-gonic/gin"
)

func Start() {
	// if counterFile exists, used the saved range
	// if not, get a new range
	cnt := util.Counter{}
	if util.FileExists(config.CounterFile) {
		cnt.Counter, cnt.CounterEnd = util.LoadCounterRange(config.CounterFile)
	} else {
		cnt.Counter, cnt.CounterEnd = cnt.GetNewRange()
	}
	// save counter range to file on non-fatal exit
	defer util.SaveCounterRange(config.CounterFile, &cnt)

	dbClient := model.GetDBClient()
	// close the database connection before exit
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	cacheClient := cache.GetCacheClient()

	if config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// create new short URL
	router.POST("/urls", func(gc *gin.Context) {
		url := model.Url{}
		if err := gc.ShouldBindJSON(&url); err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if url.Target == "" {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": "Error: missing target URL.",
			})
			return
		}

		url.Slug = util.GenerateUrlSlug(&cnt)
		url.Created = uint64(time.Now().Unix())
		url.Hits = 1

		err := model.InsertUrl(dbClient, url)
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": "Error creating new short URL entry.",
			})
			return
		}

		cache.SetCachedUrl(cacheClient, url)
		gc.JSON(http.StatusOK, gin.H{
			"message": "New short URL created.",
			"urls":    url,
		})
	})

	// get target URL from slug
	router.GET("/urls/:slug", func(gc *gin.Context) {
		slug := gc.Param("slug")

		// check cache (if enabled) for the provided slug
		// if found, return
		// if not found, continue on to check the database
		url := cache.GetCachedUrl(cacheClient, slug)
		if url.Target != "" {
			gc.JSON(http.StatusOK, gin.H{
				"slug":   url.Slug,
				"target": url.Target,
			})
			return
		}

		// check database for the provided slug
		url, err := model.GetTargetUrl(dbClient, slug)
		if err != nil {
			gc.JSON(http.StatusNotFound, gin.H{
				"message": "Short URL not found.",
			})
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"slug":   url.Slug,
				"target": url.Target,
			})
		}
	})

	// get all URLs
	router.GET("/urls", func(gc *gin.Context) {
		urls, err := model.GetUrls(dbClient)
		if err != nil {
			gc.JSON(http.StatusNotFound, gin.H{
				"message": "Error retrieving all URLs.",
			})
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"urls": urls,
			})
		}
	})

	// update target URL from slug
	router.PUT("/urls/:slug", func(gc *gin.Context) {
		slug := gc.Param("slug")
		url := model.Url{}
		if err := gc.ShouldBindJSON(&url); err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if url.Target == "" {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": "Error: missing target URL.",
			})
			return
		}

		url.Slug = slug
		url.Created = uint64(time.Now().Unix())
		url.Hits = 1

		err := model.UpdateUrl(dbClient, url)
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": "Error updating URL entry.",
			})
			return
		}

		cache.SetCachedUrl(cacheClient, url)
		gc.JSON(http.StatusOK, gin.H{
			"message": "Short URL updated.",
			"slug":    url.Slug,
			"target":  url.Target,
		})
	})

	// delete short URL by slug
	router.DELETE("/urls/:slug", func(gc *gin.Context) {
		slug := gc.Param("slug")

		// delete URL from cache (if enabled)
		cache.DeleteCachedUrl(cacheClient, slug)

		// delete URL from database
		err := model.DeleteUrl(dbClient, slug)
		if err != nil {
			gc.JSON(http.StatusNotFound, gin.H{
				"message": "Short URL not found.",
			})
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"message": "Short URL deleted.",
			})
		}
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.GinPort),
		Handler: router,
	}

	// start server in goroutine to allow graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Listen: %s\n", err)
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// wait 5 seconds before shutting down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server (%v).", err)
	}

	log.Println("Server exiting.")
}
