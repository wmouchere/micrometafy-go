package authentication

// APIAuthenticator gives primitives to unify authentication to any music service requiring "complex" authentication
type APIAuthenticator interface {
	GetToken()
	RefreshToken()
}
