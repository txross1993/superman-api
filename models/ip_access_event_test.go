package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txross1993/superman-api/errors"
)

func TestUnmarshalValidation(t *testing.T) {

	invalidIP := "300.300.300.300"
	b, _ := json.Marshal(&UserIPAccessEvent{
		IPAddress:     invalidIP,
		Username:      "bob",
		EventUUID:     "abc",
		UnixTimestamp: 1514764800,
	})

	u := UserIPAccessEvent{}
	err := json.Unmarshal(b, &u)
	expected := &errors.InvalidIP{invalidIP}
	assert.EqualError(t, err, expected.Error())
}
