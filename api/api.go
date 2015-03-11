package api

import (
	"encoding/json"
	"time"
)

//AuthRequest is a client->server request for authentication
type AuthRequest struct {
	User   string
	Passwd string
}

//AuthResponse is a server->client response about authentication
type AuthResponse struct {
	SessionID string
}

//CodeResponse is a server->client response of the latest code
type CodeResponse struct {
	Code       string
	Expires    time.Time
	ServerTime time.Time
}

//MarshalJSON returns a custom json encoding of a CodeResponse
func (c CodeResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code       string
		Expires    int64
		ServerTime int64
	}{
		Code:       c.Code,
		Expires:    c.Expires.UTC().UnixNano(),
		ServerTime: c.ServerTime.UTC().UnixNano(),
	})
}

//CheckRequest is a client->server request for checking code validity
type CheckRequest struct {
	Code string
}

//CheckResponse is a server->client response about code validity
type CheckResponse struct {
	Status bool
}

//ErrorResponse is a server-client response indicating some kind of error
type ErrorResponse struct {
	Code  int
	Error string
}
