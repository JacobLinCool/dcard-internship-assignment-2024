package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createAdHandler(c *gin.Context) {
	var ad Ad
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ads.InsertOne(context.TODO(), ad)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func listAdsHandler(c *gin.Context) {
	var query AdQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := cache.filter(query)
	c.JSON(http.StatusOK, results)
}
