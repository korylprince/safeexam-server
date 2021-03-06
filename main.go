package main

//go:generate go-bindata-assetfs static/...

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/korylprince/go-ad-auth"
	"github.com/korylprince/safeexam-server/api"
)

var static = []string{"/js", "/css", "/fonts", "/views", "/images"}

//middleware
func middleware(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout,
		handlers.CompressHandler(
			http.StripPrefix(config.Prefix,
				IndexHandler(
					ForwardedHandler(h)))))
}

func main() {
	ldapConfig := &auth.Config{
		Server:   config.LDAPServer,
		Port:     config.LDAPPort,
		BaseDN:   config.LDAPBaseDN,
		Security: config.ldapSecurity,
		Debug:    config.Debug,
	}
	c := &api.Context{
		Auth:          api.NewLDAPAuth(config.LDAPGroup, ldapConfig),
		CodeGenerator: api.NewRandomCodeGenerator(config.CodeLength, time.Duration(config.CodeInterval)*time.Minute),
		SessionStore:  api.NewMemorySessionStore(time.Duration(config.SessionDuration) * time.Minute),
	}

	r := mux.NewRouter()

	//static
	for _, route := range static {
		r.PathPrefix(route).Handler(http.FileServer(assetFS())).Methods("GET")
	}

	//index
	r.Handle("/", http.FileServer(assetFS())).Methods("GET")

	//api
	r.Handle("/api/2.0/auth", api.AuthHandler(c)).Methods("POST")
	r.Handle("/api/2.0/code", api.CodeHandler(c)).Methods("GET")
	r.Handle("/api/2.0/check", api.CheckHandler(c)).Methods("POST")

	//legacy api
	r.Handle("/check", api.CheckLegacyHandler(c)).Methods("GET")
	r.Handle("/default/check", api.CheckLegacyHandler(c)).Methods("GET")
	r.PathPrefix("/api").Handler(http.HandlerFunc(api.NotFoundHandler))

	log.Println("Listening on:", config.ListenAddr)
	log.Println(http.ListenAndServe(config.ListenAddr, middleware(r)))
}
