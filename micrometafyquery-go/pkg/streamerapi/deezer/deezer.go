package deezer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/authentication"
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/model"
)

// DeezerClient is a client to do requests on the deezer web API
type DeezerClient struct {
	Client *http.Client
}

type deezerSearchRetrieveDTO struct {
	Tracks []deezerTrackRetrieveDTO `json:"data"`
}

type deezerTrackRetrieveDTO struct {
	Name     string                  `json:"title"`
	Artist   deezerArtistRetrieveDTO `json:"artist"`
	Duration int64                   `json:"duration"`
	URL      string                  `json:"preview"`
}

type deezerArtistRetrieveDTO struct {
	Name string `json:"name"`
}

// NewDeezerClient builds a new DeezerClient object
func NewDeezerClient() *DeezerClient {
	return &DeezerClient{&http.Client{Timeout: time.Second * 3}}
}

// TODO handle errors

// SearchTrack queries to the Deezer Web API
func (c *DeezerClient) SearchTrack(queryString string) ([]model.Track, error) {
	token, err := authentication.GetAPIKeyLoaderInstance().GetAPIKey("Deezer")
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req, _ := http.NewRequest("GET", "https://api.deezer.com/search/track", nil)
	req.Close = true

	q := req.URL.Query()
	q.Add("q", queryString)
	q.Add("limit", "10")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var dto deezerSearchRetrieveDTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var result []model.Track
	for _, item := range dto.Tracks {
		result = append(result, model.Track{Name: item.Name, Author: item.Artist.Name, URL: item.URL, Duration: item.Duration * 1000, Origin: "Deezer"})
	}

	return result, nil
}
