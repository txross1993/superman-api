package superman

import (
	"math"

	"github.com/txross1993/superman-api/models"
)

type database interface {
	FindOrCreateUserIPAccessEvent(*models.UserIPAccessEvent) error
	FindPrecedingIPAccessEvent(*models.UserIPAccessEvent) (*models.UserIPAccessEvent, error)
	FindSubsequentIPAccessEvent(*models.UserIPAccessEvent) (*models.UserIPAccessEvent, error)
}

type geoservice interface {
	GetCoordinatesFromIP(string) (*models.Geography, error)
}

// Service uses an ip geoencoder service and a persistence mechanism
// to store, query, and analyze user ip access events
type Service struct {
	geoSvc geoservice
	db     database
}

// NewService creates a new service instance to process user ip access
// login events
func NewService(geo geoservice, db database) *Service {
	return &Service{
		geoSvc: geo,
		db:     db,
	}
}

// AnalyzeEvent inspects the current user ip access login event and compares
// the login event to prior and subsequent login events for the same user
// to evaluate suspicious login activity
func (s *Service) AnalyzeEvent(event *models.UserIPAccessEvent) (*models.Superman, error) {
	var superman *models.Superman
	var supermanOpts []models.SupermanOpt

	applyOpts := func() {
		superman = models.NewSuperman(supermanOpts...)
	}

	// Inspect current event
	if err := s.db.FindOrCreateUserIPAccessEvent(event); err != nil {
		applyOpts()
		return superman, err
	}
	currentAccess := event.AsIPAccess()
	currentGeoOpt, err := s.inspectCurrent(currentAccess)
	if err != nil {
		applyOpts()
		return superman, err
	}
	supermanOpts = append(supermanOpts, currentGeoOpt)

	// Inspect preceding event
	preceding, err := s.db.FindPrecedingIPAccessEvent(event)
	if err != nil {
		applyOpts()
		return superman, err
	}

	if preceding != nil {
		precedingIPAccess := preceding.AsIPAccess()
		precedingOpt, err := s.inspectPreceding(currentAccess, precedingIPAccess)
		if err != nil {
			applyOpts()
			return superman, err
		}
		supermanOpts = append(supermanOpts, precedingOpt)
	}

	// Inspect subsequent event
	subsequent, err := s.db.FindSubsequentIPAccessEvent(event)
	if err != nil {
		applyOpts()
		return superman, err
	}

	if subsequent != nil {
		subsequentAccess := subsequent.AsIPAccess()
		subsequentOpt, err := s.inspectSubsequentEvent(currentAccess, subsequentAccess)
		if err != nil {
			applyOpts()
			return superman, err
		}
		supermanOpts = append(supermanOpts, subsequentOpt)

	}

	applyOpts()
	return superman, nil
}

// inspectCurrent geoencodes the current ip access event and provides the
// builder option for the Superman response
func (s *Service) inspectCurrent(current *models.IPAccess) (models.SupermanOpt, error) {
	current, err := s.geoencode(current)
	return models.WithCurrentGeo(current.Geography), err
}

// inspectPreceding geoencodes the preceding ip access event and determines
// the distance and speed from the current ip access event
func (s *Service) inspectPreceding(current, preceding *models.IPAccess) (models.SupermanOpt, error) {
	preceding, err := s.geoencode(preceding)
	preceding = s.analyzeEventSequence(current, preceding)
	return models.WithPrecedingEvent(preceding), err
}

// inspectSubsequentEvent geoencodes the subsequent ip access event and determines
// the distance and speed from the current ip access event
func (s *Service) inspectSubsequentEvent(current, subsequent *models.IPAccess) (models.SupermanOpt, error) {
	subsequent, err := s.geoencode(subsequent)
	subsequent = s.analyzeEventSequence(current, subsequent)
	return models.WithSubsequentEvent(subsequent), err
}

// analyzeEventSequence compares the current event to an alternate event
// to determine the distance and speed of access between events
func (s *Service) analyzeEventSequence(current, alt *models.IPAccess) *models.IPAccess {
	if alt == nil {
		return nil
	}

	var distanceMiles float64
	if alt.Geography != nil {
		distanceMiles = current.Geography.MilesFrom(alt.Geography)
	}

	timedelta := calculateTimedelta(current.Timestamp, alt.Timestamp)
	alt.Speed = calculateSpeedMPH(distanceMiles, timedelta)
	return alt
}

// geoencode applies the geoencoding service to the ip access event ip address
func (s *Service) geoencode(event *models.IPAccess) (*models.IPAccess, error) {
	geo, err := s.geoSvc.GetCoordinatesFromIP(event.IP)
	event.Geography = geo
	return event, err
}

// calculateTimedelta expects two unix epoch timestamps and returns the absolute
// value of the difference
func calculateTimedelta(t1, t2 int64) int64 {
	delta := t1 - t2
	if delta < 0 {
		return -delta
	}
	return delta
}

// calculateSpeedMPH expects distance in miles and timedelta in seconds to
// calculate miles per hour
func calculateSpeedMPH(distance float64, timedelta int64) int64 {
	if timedelta == 0 {
		return math.MaxInt64
	}
	milesPerSecond := math.Round(distance) / float64(timedelta)
	return int64(milesPerSecond * 3600)
}
