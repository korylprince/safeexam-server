package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

//handleError returns a json response for the given code and logs the error
func handleError(w http.ResponseWriter, code int, err error) {
	log.Println(err)
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	encErr := e.Encode(ErrorResponse{Code: code, Error: http.StatusText(code)})
	if encErr != nil {
		panic(encErr)
	}
}

//NotFoundHandler returns a json 401 response
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	handleError(w, http.StatusNotFound, errors.New("handler not found"))
}

//authHandler will return a sessionID if the credentials are valid
//or an HTTP 401 Error if not.
func authHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var aReq AuthRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&aReq)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("Error decoding json: %v", err))
		return
	}

	status, err := c.Auth.Login(aReq.User, aReq.Passwd)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Error authenticating: %v", err))
		return
	}
	if status {
		sessionID, err := c.SessionStore.Create()
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error creating session key: %v", err))
			return
		}

		e := json.NewEncoder(w)
		err = e.Encode(AuthResponse{SessionID: sessionID})
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error encoding json: %v", err))
		}
		return
	}
	handleError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
}

//codeHandler will return the current code if the sessionID is valid
//or an HTTP 401 Error if not.
func codeHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	key := r.Header.Get("X-Session-Key")
	if key == "" {
		handleError(w, http.StatusBadRequest, errors.New("X-Session-Key header empty"))
		return
	}

	status, err := c.SessionStore.Check(key)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Error checking session key: %v", err))
		return
	}
	if status {
		code, expires, err := c.CodeGenerator.Generate()
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error getting code: %v", err))
			return
		}

		e := json.NewEncoder(w)
		err = e.Encode(CodeResponse{Code: code, Expires: expires, ServerTime: time.Now()})
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error encoding json: %v", err))
			return
		}
		return
	}
	handleError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
}

//checkHandler validates a given code
func checkHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cReq CheckRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&cReq)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("Error decoding json: %v", err))
		return
	}

	code, _, err := c.CodeGenerator.Generate()
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Error getting code: %v", err))
		return
	}

	//401 with {status:false} response
	if cReq.Code != code {
		w.WriteHeader(http.StatusUnauthorized)
		e := json.NewEncoder(w)
		encErr := e.Encode(CheckResponse{Status: false})
		if encErr != nil {
			panic(encErr)
		}
		return
	}

	e := json.NewEncoder(w)
	err = e.Encode(CheckResponse{Status: true})
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Error encoding json: %v", err))
		return
	}
}

//checkLegacyHandler returns "true"/"false" if the given "pass" is correct
//of an HTTP 400 if not
func checkLegacyHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	pass := r.FormValue("pass")
	if pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	code, _, err := c.CodeGenerator.Generate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	if pass != code {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	}
	w.Write([]byte("true"))
}
