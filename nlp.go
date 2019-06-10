package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContentRequest struct {
	Content string `json: content`
}

func main() {
	r := gin.Default()

	r.GET("/ping", ping)
	r.POST("nlp/finance/vanguard/categorize", categorize)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}

func categorize(c *gin.Context) {
	var req ContentRequest

	c.BindJSON(&req)
	c.JSON(http.StatusOK, req)
}
