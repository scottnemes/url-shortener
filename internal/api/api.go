package api

import (
	"context"
	"net/http"
	"time"

	"example.com/url-shortener/internal/model"
	"example.com/url-shortener/internal/util"
	"github.com/gin-gonic/gin"
)

func Server() {
	counterFile := "counter_range.dat"

	// if counterFile exists, used the saved range
	// if not, get a new range
	cnt := util.Counter{}
	if util.FileExists(counterFile) {
		cnt.Counter, cnt.CounterEnd = util.LoadCounterRange(counterFile)
	} else {
		cnt.Counter, cnt.CounterEnd = cnt.GetNewRange()
	}
	// save counter range to file on non-fatal exit
	defer util.SaveCounterRange(counterFile, &cnt)

	c := model.GetDBClient()
	// close the database connection before exit
	defer func() {
		if err := c.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()

	router.GET("/api/url", func(gc *gin.Context) {
		json := model.Url{}
		if err := gc.ShouldBindJSON(&json); err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		target, err := model.GetTargetUrl(c, json.Slug)
		if err != nil {
			gc.JSON(http.StatusNotFound, gin.H{
				"message": "Short URL not found.",
			})
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"target": target.Target,
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

		err := model.InsertUrl(c, newUrl)
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"message": "Error creating new short URL entry.",
			})
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"message": "New short URL created.",
				"slug":    newUrl.Slug,
			})
		}
	})

	router.Run(":8080")
}
