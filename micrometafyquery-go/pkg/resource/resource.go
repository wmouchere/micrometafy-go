package resource

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/streamerapi"
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/streamerapi/deezer"
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/streamerapi/jamendo"
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/streamerapi/spotify"
)

var apis = []streamerapi.StreamerAPI{
	spotify.NewSpotifyClient(),
	deezer.NewDeezerClient(),
	jamendo.NewJamendoClient(),
}

type trackDTO struct {
	Name     string `json:"name"`
	Author   string `json:"author"`
	URL      string `json:"url"`
	Duration int64  `json:"duration"`
	Origin   string `json:"origin"`
}

// SearchGET is a HTTP GET request handler that returns a list of Tracks retrieved from all the services in the "apis" list
func SearchGET(c *gin.Context) {
	query := c.Param("query")
	var result []trackDTO

	var wg sync.WaitGroup
	wg.Add(3)
	responseChan := make(chan []trackDTO, 3)

	for _, api := range apis {
		go func(api streamerapi.StreamerAPI, wg *sync.WaitGroup, c chan []trackDTO) {
			defer wg.Done()
			var innerResult []trackDTO
			tracks, err := api.SearchTrack(query)
			if err != nil {
				return
			}
			for _, track := range tracks {
				innerResult = append(innerResult, trackDTO{track.Name, track.Author, track.URL, track.Duration, track.Origin})
			}
			c <- innerResult
		}(api, &wg, responseChan)
	}

	wg.Wait()
	close(responseChan)
	for resp := range responseChan {
		result = append(result, resp...)
	}

	c.JSON(http.StatusOK, result)
}
