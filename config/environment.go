package config

import (
	"log"
	"os"
	"strconv"
)

const (
	defaultSessionTTLMinutes int    = 30
	defaultHashSecretSalt    string = ""
	defaultSQLTrace          bool   = false
)

var (
	sessionTTLMinutes = defaultSessionTTLMinutes
	hashSecretSalt    = defaultHashSecretSalt
	traceSQL          = defaultSQLTrace
)

func init() {
	var err error

	sessionTTLMinutes, err = strconv.Atoi(os.Getenv(EnvIDPSessionTTL))
	if err != nil {
		log.Fatalf("Failed to read %v: %v", EnvIDPSessionTTL, err.Error())
	}
	if sessionTTLMinutes == 0 {
		sessionTTLMinutes = defaultSessionTTLMinutes
	}

	hashSecretSalt = os.Getenv(EnvIDPSecretSalt)

	traceSQL, err = strconv.ParseBool(os.Getenv(EnvIDPSQLTrace))
	if err != nil {
		log.Printf("Failed to read %v: %v", EnvIDPSQLTrace, err.Error())
		traceSQL = defaultSQLTrace
	}
}

// SessionTTLMinutes returns a session TTL duration in minutes read from environment variables
func SessionTTLMinutes() int {
	return sessionTTLMinutes
}

// HashSecretSalt returns a hash's secret salt read from environment variables
func HashSecretSalt() string {
	return hashSecretSalt
}

// SQLTraceOn tells if SQLs should be logged
func SQLTraceOn() bool {
	return traceSQL
}
