package geolocate

import (
	"net"

	geoDB "github.com/oschwald/geoip2-golang"

	"github.com/txross1993/superman-api/errors"
	"github.com/txross1993/superman-api/models"
)

// GeoService provides the service to geoencode IP addresses
type GeoService struct {
	db *geoDB.Reader
}

// NewGeoService initializes a new in-memory IP geoencoding service provided
// a path to the local database file
func NewGeoService(repositoryPath string) (*GeoService, error) {
	db, err := geoDB.Open(repositoryPath)
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
