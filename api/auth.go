package api

import "github.com/korylprince/go-ad-auth"

// Auth is an interface for an arbitrary authentication backend.
type Auth interface {
	// Login returns whether or not the given username or password is valid.
	// If the backend malfunctions, status will be false and error will be non-nil.
	Login(username, password string) (status bool, err error)
}

//LDAPAuth represents an Auth that uses an Active Directory backend
type LDAPAuth struct {
	group  string
	config *auth.Config
}

//Login returns whether or not the given username or password is valid.
//If the backend malfunctions, status will be false and error will be non-nil.
func (a *LDAPAuth) Login(username, password string) (status bool, err error) {
	return auth.Login(username, password, a.group, a.config)
}

//NewLDAPAuth returns a new LDAPAuth with the given config, restricting logins to group, if given.
func NewLDAPAuth(group string, config *auth.Config) *LDAPAuth {
	return &LDAPAuth{
		group:  group,
		config: config,
	}
}
