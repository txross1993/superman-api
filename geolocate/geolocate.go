package geolocate

import (
	"net"
	"path"
	"runtime"

	geoDB "github.com/oschwald/geoip2-golang"

	"github.com/txross1993/superman-api/errors"
	"github.com/txross1993/superman-api/models"
)

const repository = "./GeoLite2-City_20200602/GeoLite2-City.mmdb"

// GeoService provides the service to geoencode IP addresses
type GeoService struct {
	db *geoDB.Reader
}

// NewGeoService initializes a new in-memory IP geoencoding service
func NewGeoService() (*GeoService, error) {
	_, filename, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(filename), repository)
	db, err := geoDB.Open(path)

	if err != nil {
		return nil, err
	}

	return &GeoService{db}, nil
}

// GetCoordinatesFromIP parses the input IP and queries the database for
// latitude and longitude
func (g GeoService) GetCoordinatesFromIP(ip string) (*models.Geography, error) {
	var geo models.Geography
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return nil, &errors.InvalidIP{IP: ip}
	}

	record, err := g.db.City(netIP)
	if err != nil {
		return nil, err
	}

	geo.Latitude = record.Location.Latitude
	geo.Longitude = record.Location.Longitude
	geo.Radius = record.Location.AccuracyRadius
	return &geo, nil

}

// Close closes the GeoService repository
func (g GeoService) Close() error {
	return g.db.Close()

}
