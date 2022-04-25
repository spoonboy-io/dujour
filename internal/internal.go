package internal

import "time"

const (
	// server
	SRV_HOST = ""
	SRV_PORT = "18651"

	// data
	DATA_FOLDER = "data"

	// storage
	TYPE_CSV  = 1
	TYPE_JSON = 2
	DATA_ARR  = 1
	DATA_OBJ  = 2

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
	DataType     int
	Data         interface{}
}
