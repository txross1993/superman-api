package models

import (
	"encoding/json"
	"net"

	"github.com/txross1993/superman-api/errors"
)

// UserIPAccessEvent represents an instance of access from an IP address for a
// given username
type UserIPAccessEvent struct {
	EventUUID     string `json:"event_uuid" gorm:"primary_key"`
	Username      string `json:"username" gorm:"not null" sql:"index"`
	UnixTimestamp int64  `json:"unix_timestamp" gorm:"not null" sql:"index"`
	IPAddress     string `json:"ip_address" gorm:"not null"`
}

// UnmarshalJSON performs data validation on the ip address of the event
func (u *UserIPAccessEvent) UnmarshalJSON(data []byte) error {
	type Alias UserIPAccessEvent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	netIP := net.ParseIP(u.IPAddress)
	if netIP == nil {
		return &errors.InvalidIP{IP: u.IPAddress}
	}
	return nil
}

// AsIPAccess performs data translation from this model to the IPAccess model
func (u *UserIPAccessEvent) AsIPAccess() *IPAccess {
	return &IPAccess{
		IP:        u.IPAddress,
		Timestamp: u.UnixTimestamp,
	}
}
