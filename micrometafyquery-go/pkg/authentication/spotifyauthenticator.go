package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// SpotifyAuthenticator is used to authenticate to the Spotify Web API
type SpotifyAuthenticator struct {
	apiKey              string
	token               string
	tokenExpirationDate int64
}

type spotifyAuthenticationReturn struct {
	Token               string `json:"access_token"`
	TokenExpirationDate int64  `json:"expires_in"`
}

var spotifyAuthenticatorInstance *SpotifyAuthenticator
var spotifyAuthenticatorError error
var spotifyAuthenticatorOnce sync.Once

// GetSpotifyAuthenticatorInstance returns the instance of the singleton SpotifyAuthenticator
func GetSpotifyAuthenticatorInstance() (*SpotifyAuthenticator, error) {
	spotifyAuthenticatorOnce.Do(func() {
		apiKeyLoader := GetAPIKeyLoaderInstance()
		apiKey, err := apiKeyLoader.GetAPIKey("Spotify")
		if err != nil {
			spotifyAuthenticatorError = err
			return
		}
		spotifyAuthenticatorInstance = &SpotifyAuthenticator{apiKey: apiKey}
		spotifyAuthenticatorInstance.RefreshToken()
	})
	if spotifyAuthenticatorError != nil {
		return nil, spotifyAuthenticatorError
	}
	return spotifyAuthenticatorInstance, nil
}

//TODO : Handle errors

// RefreshToken refreshes the authorization token for the Spotify API
func (a *SpotifyAuthenticator) RefreshToken() {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	form := url.Values{"grant_type": {"client_credentials"}}

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+a.apiKey)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var dto spotifyAuthenticationReturn
	err = json.Unmarshal(body, &dto)
	if err != nil {
		fmt.Print(err)
	}
	a.token = dto.Token
	a.tokenExpirationDate = dto.TokenExpirationDate
}

//TODO : Handle errors

// GetToken returns the authentication token
func (a *SpotifyAuthenticator) GetToken() string {
	if time.Unix(a.tokenExpirationDate, 0).After(time.Now()) {
		a.RefreshToken()
	}
	return a.token
}
