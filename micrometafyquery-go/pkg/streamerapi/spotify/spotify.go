package spotify

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/authentication"
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/model"
)

// SpotifyClient is a client to do requests on the spotify web API
type SpotifyClient struct {
	Client *http.Client
}

type spotifySearchRetrieveDTO struct {
	Tracks spotifyTracksRetrieveDTO `json:"tracks"`
}

type spotifyTracksRetrieveDTO struct {
	Items []spotifyTrackRetrieveDTO `json:"items"`
}

type spotifyTrackRetrieveDTO struct {
	Name     string                     `json:"name"`
	Artists  []spotifyArtistRetrieveDTO `json:"artists"`
	Duration int64                      `json:"duration_ms"`
	URL      string                     `json:"preview_url"`
}

type spotifyArtistRetrieveDTO struct {
	Name string `json:"name"`
}

// NewSpotifyClient builds a new SpotifyClient object
func NewSpotifyClient() *SpotifyClient {
	return &SpotifyClient{&http.Client{Timeout: time.Second * 3}}
}

// TODO handle (all) possible errors

// SearchTrack queries to the Spotify Web API
func (c *SpotifyClient) SearchTrack(queryString string) ([]model.Track, error) {
	authenticator, err := authentication.GetSpotifyAuthenticatorInstance()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	token := authenticator.GetToken()

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/search", nil)
	req.Close = true

	q := req.URL.Query()
	q.Add("q", queryString)
	q.Add("type", "track")
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
	var dto spotifySearchRetrieveDTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var result []model.Track
	for _, item := range dto.Tracks.Items {
		result = append(result, model.Track{Name: item.Name, Author: item.Artists[0].Name, URL: item.URL, Duration: item.Duration, Origin: "Spotify"})
	}

	return result, nil
}
