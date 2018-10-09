package authentication

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// APIKeyLoader lazily loads the API keys when needed
type APIKeyLoader struct {
	keys map[string]string
}

var apiKeyLoaderInstance *APIKeyLoader
var apiKeyLoaderOnce sync.Once

type ymlConf struct {
	Spotify ymlSpotify `yaml:"Spotify"`
	Deezer  ymlDeezer  `yaml:"Deezer"`
	Jamendo ymlJamendo `yaml:"Jamendo"`
}

type ymlSpotify struct {
	User string `yaml:"user"`
	Key  string `yaml:"key"`
}

type ymlDeezer struct {
	Token string `yaml:"token"`
}

type ymlJamendo struct {
	Key string `yaml:"key"`
}

// GetAPIKeyLoaderInstance returns the instance of the singleton APIKeyLoader
func GetAPIKeyLoaderInstance() *APIKeyLoader {
	apiKeyLoaderOnce.Do(func() {
		ymlFile, err := ioutil.ReadFile("apikeys.yml")
		if err != nil {
			log.Printf("Error reading apikeys.yml : %v", err)
			os.Exit(1)
		}
		var ymlConf ymlConf
		err = yaml.Unmarshal(ymlFile, &ymlConf)
		if err != nil {
			log.Printf("Error unmarshaling apikeys.yml : %v", err)
			os.Exit(1)
		}
		apiKeyLoaderInstance = &APIKeyLoader{map[string]string{
			"Spotify": base64.StdEncoding.EncodeToString([]byte(ymlConf.Spotify.User + ":" + ymlConf.Spotify.Key)),
			"Deezer":  ymlConf.Deezer.Token,
			"Jamendo": ymlConf.Jamendo.Key,
		}}
	})
	return apiKeyLoaderInstance
}

// GetAPIKey returns the key for each music service
func (a *APIKeyLoader) GetAPIKey(apiName string) (string, error) {
	v, ok := a.keys[apiName]
	if ok {
		return v, nil
	}
	return "", &APIKeyMissingError{apiName}
}
