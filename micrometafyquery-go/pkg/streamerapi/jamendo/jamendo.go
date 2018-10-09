package jamendo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/authentication"
	"github.com/wmouchere/micrometafy-go/micrometafyquery-go/pkg/model"
)

// JamendoClient is a client to do requests on the jamendo web API
type JamendoClient struct {
	Client *http.Client
}

type jamendoSearchRetrieveDTO struct {
	Tracks []jamendoTrackRetrieveDTO `json:"results"`
}

type jamendoTrackRetrieveDTO struct {
	Name     string `json:"name"`
	Artist   string `json:"artist_name"`
	Duration int64  `json:"duration"`
	URL      string `json:"audio"`
}

// NewJamendoClient builds a new JamendoClient object
func NewJamendoClient() *JamendoClient {
	return &JamendoClient{&http.Client{Timeout: time.Second * 3}}
}

// TODO handle errors

// SearchTrack queries to the Jamendo Web API
func (c *JamendoClient) SearchTrack(queryString string) ([]model.Track, error) {
	token, err := authentication.GetAPIKeyLoaderInstance().GetAPIKey("Jamendo")
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req, _ := http.NewRequest("GET", "https://api.jamendo.com/v3.0/tracks", nil)
	req.Close = true

	q := req.URL.Query()
	q.Add("client_id", token)
	q.Add("format", "json")
	q.Add("search", queryString)
	q.Add("limit", "10")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var dto jamendoSearchRetrieveDTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var result []model.Track
	for _, item := range dto.Tracks {
		result = append(result, model.Track{Name: item.Name, Author: item.Artist, URL: item.URL, Duration: item.Duration * 1000, Origin: "Jamendo"})
	}

	return result, nil
}
