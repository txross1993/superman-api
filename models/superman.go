package models

// Superman encapsulates the main Superman API response
type Superman struct {
	CurrentGeo           *Geography `json:"currentGeo"`
	TravelToSuspicious   bool       `json:"travelToCurrentGeoSuspicious"`
	TravelFromSuspicious bool       `json:"travelFromCurrentGeoSuspicious"`
	PrecedingIPAccess    *IPAccess  `json:"precedingIpAccess,omitempty"`
	SubsequentIPAccess   *IPAccess  `json:"subsequentIpAccess,omitempty"`
}

// SupermanOpt represents a functional option for building a Superman response
type SupermanOpt func(s *Superman)

// NewSuperman provides the builder pattern for building a Superman response
func NewSuperman(opts ...SupermanOpt) *Superman {
	s := &Superman{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithCurrentGeo provides the functional option for Superman.CurrentGeo
func WithCurrentGeo(geo *Geography) SupermanOpt {
	return func(s *Superman) {
		s.CurrentGeo = geo
	}
}

// WithPrecedingEvent provides the functional option for Superman.PrecedingIPAccess
// and Superman.TravelToSuspicious
func WithPrecedingEvent(event *IPAccess) SupermanOpt {
	return func(s *Superman) {
		s.PrecedingIPAccess = event

		if event != nil {
			if s.PrecedingIPAccess.Speed >= 500 {
				s.TravelToSuspicious = true
			}

		}
	}
}

// WithSubsequentEvent provides the functional option for Superman.SubsequentIPAccess
// and Superman.TravelFromSuspicious
func WithSubsequentEvent(event *IPAccess) SupermanOpt {
	return func(s *Superman) {
		s.SubsequentIPAccess = event
		if event != nil {
			if s.SubsequentIPAccess.Speed >= 500 {
				s.TravelFromSuspicious = true
			}

		}
	}
}
