package models

import "github.com/umahmood/haversine"

// Geography represents a lat,lon, and accuracy radius of the coordinates
type Geography struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Radius    uint16  `json:"radius"`
}

// MilesFrom calculates the miles between this coordinate and the provided point
func (g *Geography) MilesFrom(coordinate *Geography) float64 {
	if coordinate == nil {
		return 0.0
	}
	here := g.haversineCoord()
	there := coordinate.haversineCoord()

	distanceMiles, _ := haversine.Distance(here, there)
	return distanceMiles
}

// haversineCoord returns the point as a haversine coordinate for calculating
// distance
func (g *Geography) haversineCoord() haversine.Coord {
	return haversine.Coord{Lat: g.Latitude, Lon: g.Longitude}
}
