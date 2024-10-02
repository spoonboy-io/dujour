package internal

import (
	"errors"
	"time"
)

const (
	// server
	SRV_HOST = ""
	SRV_PORT = "18651"

	// data
	DATA_FOLDER = "data"

	// storage
	TYPE_CSV      = 1
	TYPE_JSON     = 2
	TYPE_DB_QUERY = 3

	// tls configuration
	TLS_FOLDER    = "certs"
	TLS_ORG       = "Spoon Boy"
	TLS_VALID_FOR = 365 * 24 * time.Hour
)

// Datasource contains both the data and metadata of a discovered and validated datasource
type Datasource struct {
	FileName     string
	FileType     int
	EndpointName string
	Data         interface{}
}

var (
	// TODO should do this for all error messages in the app
	ERR_NO_USER     = errors.New("username is required")
	ERR_NO_PORT     = errors.New("port number is required")
	ERR_NO_HOST     = errors.New("hostname or ip address is required")
	ERR_NO_PASSWORD = errors.New("password is required for authentication")
	ERR_NO_QUERY    = errors.New("SQL query is required")
	ERR_NOT_ACTIVE  = errors.New("DB config file is not active, skipping")
)
