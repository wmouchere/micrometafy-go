package main // import "github.com/wmouchere/micrometafy-go/micrometafyquery-go"

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/resource"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{AllowOrigins: []string{"http://localhost"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Origin"}}))
	api := router.Group("/micrometafy-query/api")
	api.GET("/search/:query", resource.SearchGET)
	router.Run(":8080")
}
