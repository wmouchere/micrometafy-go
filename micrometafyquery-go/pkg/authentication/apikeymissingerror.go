package authentication

import "fmt"

// APIKeyMissingError handles queries for non existing keys
type APIKeyMissingError struct {
	Key string
}

func (e *APIKeyMissingError) Error() string {
	return fmt.Sprintf("Key for %s was missing in apikeys.yml", e.Key)
}
