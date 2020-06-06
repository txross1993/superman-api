package errors

import "fmt"

// InvalidIP is the error type for bad IP geoencoding request
type InvalidIP struct {
	IP string
}

func (err *InvalidIP) Error() string {
	return fmt.Sprintf("invalid IP address format: %s", err.IP)
}
