package geolocate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txross1993/superman-api/errors"
)

const (
	validLatMin  = -90.0
	validLatMax  = 90.0
	validLongMin = -180.0
	validLongMax = 180.0
)

func TestGetCoordinatesByIP(t *testing.T) {

	geoSvc, err := NewGeoService()
	if err != nil {
		t.Fatal(err)
	}

	defer geoSvc.Close()

	tests := map[string]struct {
		IP          string
		ExpectedErr error
	}{
		"validIPv4": {
			IP:          "8.8.8.8",
			ExpectedErr: nil,
		},
		"invalidIPv4": {
			IP:          "300.300.300.300",
			ExpectedErr: &errors.InvalidIP{"300.300.300.300"},
		},
		"validIPv6": {
			IP:          "2001:4860:4860::8888",
			ExpectedErr: nil,
		},
		"invalidIPv6": {
			IP:          "2001:4860:4860:8888",
			ExpectedErr: &errors.InvalidIP{"2001:4860:4860:8888"},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case %s", name)
		geo, err := geoSvc.GetCoordinatesFromIP(test.IP)

		if err != nil {
			assert.Equal(t, test.ExpectedErr, err)
		} else {
			assert.True(t, isValidLat(geo.Latitude))
			assert.True(t, isValidLong(geo.Longitude))
		}
	}
}

func isValidLat(lat float64) bool {
	isValid := true
	if lat <= validLatMin || lat >= validLatMax {
		isValid = false
	}

	return isValid
}

func isValidLong(long float64) bool {
	isValid := true
	if long <= validLongMin || long >= validLongMax {
		isValid = false
	}

	return isValid
}
