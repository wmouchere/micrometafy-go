package main // import "github.com/wmouchere/micrometafy-go/micrometafyplaylist-go"

import (
	"flag"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wmouchere/micrometafy-go/micrometafyplaylist-go/pkg/resource"
)

func main() {
	flag.Parse()

	resourceSingleton := resource.GetResourceInstance()
	defer resourceSingleton.Close()
	router := gin.Default()
	router.Use(cors.New(cors.Config{AllowOrigins: []string{"http://localhost"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type"}}))
	api := router.Group("/micrometafy-playlist/api")
	api.POST("/playlist", resourceSingleton.NewPlaylistPOST)
	api.GET("/playlist/:id", resourceSingleton.PlaylistByIDGET)
	api.GET("/playlists", resourceSingleton.AllPlaylistsGET)
	api.DELETE("/playlist/:id", resourceSingleton.RemovePlaylistDELETE)
	api.PUT("/playlist/:id/add", resourceSingleton.AddTrackToPlaylistPUT)
	router.Run(":8080")
}
