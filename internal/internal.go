package internal

const (
	DATA_FOLDER = "data"
	TLS_FOLDER  = "certs"
	TYPE_CSV    = 1
	TYPE_JSON   = 2
	DATA_ARR    = 1
	DATA_OBJ    = 2
)

type Datasource struct {
	FileName     string
	FileType     int
	EndpointName string
	DataType     int
	Data         interface{}
}
