package main

import (
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/korylprince/go-ad-auth"
)

//Config represents options given in the environment
type Config struct {
	LDAPServer   string //required
	LDAPPort     int    //default: 389
	LDAPBaseDN   string //required
	LDAPGroup    string //optional
	LDAPSecurity string //default: none
	ldapSecurity auth.SecurityType

	CodeInterval int //in minutes; default: 12
	CodeLength   int //in characters; default: 3

	SessionDuration int //in minutes; default: 60

	ListenAddr string //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash

	Debug bool //default: false
}

var config = &Config{}

func checkEmpty(val, name string) {
	if val == "" {
		log.Fatalf("SAFEEXAM_%s must be configured\n", name)
	}
}

func init() {
	err := envconfig.Process("SAFEEXAM", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}
	checkEmpty(config.LDAPServer, "LDAPSERVER")

	if config.LDAPPort == 0 {
		config.LDAPPort = 389
	}

	checkEmpty(config.LDAPBaseDN, "LDAPBASEDN")

	switch strings.ToLower(config.LDAPSecurity) {
	case "", "none":
		config.ldapSecurity = auth.SecurityNone
	case "tls":
		config.ldapSecurity = auth.SecurityTLS
	case "starttls":
		config.ldapSecurity = auth.SecurityStartTLS
	default:
		log.Fatalln("Invalid SAFEEXAM_LDAPSECURITY:", config.LDAPSecurity)
	}

	if config.CodeInterval == 0 {
		config.CodeInterval = 12
	}

	if config.CodeLength == 0 {
		config.CodeLength = 3
	}

	if config.SessionDuration == 0 {
		config.SessionDuration = 60
	}

	checkEmpty(config.ListenAddr, "LISTENADDR")
}
