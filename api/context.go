package api

import (
	"net/http"
)

//Context represents a group of services
type Context struct {
	Auth          Auth
	CodeGenerator CodeGenerator
	SessionStore  SessionStore
}

type contextHandler struct {
	HandleFunc func(*Context, http.ResponseWriter, *http.Request)
	Context    *Context
}

func (c contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.HandleFunc(c.Context, w, r)
}

//AuthHandler returns an Authentcation http.Handler with the given context
func AuthHandler(c *Context) http.Handler {
	return contextHandler{HandleFunc: authHandler, Context: c}
}

//CodeHandler returns an Code http.Handler with the given context
func CodeHandler(c *Context) http.Handler {
	return contextHandler{HandleFunc: codeHandler, Context: c}
}

//CheckHandler returns an Check http.Handler with the given context
func CheckHandler(c *Context) http.Handler {
	return contextHandler{HandleFunc: checkHandler, Context: c}
}

//CheckLegacyHandler returns an Check http.Handler with the given context
func CheckLegacyHandler(c *Context) http.Handler {
	return contextHandler{HandleFunc: checkLegacyHandler, Context: c}
}
