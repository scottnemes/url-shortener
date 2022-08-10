package api

import (
	"context"
	"errors"
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

	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.GET(":slug", func(gc *gin.Context) {
		slug := gc.Param("slug")
		/*
			TODO:
			1. Add error checking for slug
		*/

		// check if the URL is in the cache if enabled
		url := cache.GetCachedUrl(cacheClient, slug)
		if url.Target != "" {
			gc.Redirect(http.StatusTemporaryRedirect, url.Target)
			return
		}

		// if URL is not in cache or cache is disabled, check the database
		url, err := model.GetTargetUrl(dbClient, slug)
		if err != nil {
			gc.Redirect(http.StatusTemporaryRedirect, "/")
		}
		gc.Redirect(http.StatusTemporaryRedirect, url.Target)
	})

	router.GET("/api/url", func(gc *gin.Context) {
		json := model.Url{}
		if err := gc.ShouldBindJSON(&json); err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		// check if the URL is in the cache if enabled
		url := cache.GetCachedUrl(cacheClient, json.Slug)
		if url.Target != "" {
			gc.JSON(http.StatusOK, gin.H{
				"target": url.Target,
			})
			return
		}

		// if URL is not in cache or cache is disabled, check the database
		url, err := model.GetTargetUrl(dbClient, json.Slug)
		if err != nil {
			gc.JSON(http.StatusNotFound, gin.H{
				"message": "Short URL not found.",
			})
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"target": url.Target,
			})
		}
	})

	router.POST("/api/url", func(gc *gin.Context) {
		json := model.Url{}
		if err := gc.ShouldBindJSON(&json); err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if json.Target == "" {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": "Error: missing target URL.",
			})
			return
		}

		newUrl := model.Url{
			Slug:    util.GenerateUrlSlug(&cnt),
			Target:  json.Target,
			Created: uint64(time.Now().Unix()),
			Hits:    1,
		}

		err := model.InsertUrl(dbClient, newUrl)
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": "Error creating new short URL entry.",
			})
			return
		}

		cache.SetCachedUrl(cacheClient, newUrl)
		gc.JSON(http.StatusOK, gin.H{
			"message": "New short URL created.",
			"slug":    newUrl.Slug,
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
