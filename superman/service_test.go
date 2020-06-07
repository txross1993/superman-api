package superman

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txross1993/superman-api/models"
	"github.com/txross1993/superman-api/testdata"
)

func TestSuperman(t *testing.T) {
	testEvent := testdata.GenerateCurrentEvent()

	tests := map[string]struct {
		testParams             testParams
		expectedFromSupsicious bool
		expectedToSupsicious   bool
	}{
		"Valid": {
			testParams: testParams{
				validPreceding:  true,
				validSubsequent: true,
			},
			expectedFromSupsicious: false,
			expectedToSupsicious:   false,
		},
		"No Preceding Event": {
			testParams: testParams{
				noPreceding:     true,
				validSubsequent: true,
			},
			expectedFromSupsicious: false,
			expectedToSupsicious:   false,
		},
		"No Subsequent Event": {
			testParams: testParams{
				validPreceding: true,
				noSubsequent:   true,
			},
			expectedFromSupsicious: false,
			expectedToSupsicious:   false,
		},
		"Same Timestamp": {
			testParams: testParams{
				sameTimestamp: true,
			},
			expectedFromSupsicious: true,
			expectedToSupsicious:   true,
		},
		"Suspicious Preceding": {
			testParams: testParams{
				suspiciousPreceding: true,
				validSubsequent:     true,
			},
			expectedFromSupsicious: false,
			expectedToSupsicious:   true,
		},
		"Suspicious Subsequent": {
			testParams: testParams{
				validPreceding:       true,
				suspiciousSubsequent: true,
			},
			expectedFromSupsicious: true,
			expectedToSupsicious:   false,
		},
	}

	geo := &mockGeo{}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		db := &mockDB{test.testParams}
		superman := NewService(geo, db)

		resp, err := superman.AnalyzeEvent(testEvent)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatal("Nil response")
		}

		t.Logf("Response: %v", resp)
		assert.Equal(t, test.expectedFromSupsicious, resp.TravelFromSuspicious)
		assert.Equal(t, test.expectedToSupsicious, resp.TravelToSuspicious)
	}
}

// TestSupermanUtils tests the speed and distance functions
func TestSupermanUtils(t *testing.T) {
	t.Run("Calc speed tests", func(t *testing.T) {
		t.Parallel()
		tests := map[string]struct {
			distance  float64
			timedelta int64
			want      int64
		}{
			"0 time": {
				distance:  4300,
				timedelta: 0,
				want:      math.MaxInt64,
			},
			"default": {
				distance:  3600,
				timedelta: 3600,
				want:      3600,
			},
		}
		for name, test := range tests {
			t.Logf("Running test case: %s", name)
			got := calculateSpeedMPH(test.distance, test.timedelta)
			if got != test.want {
				t.Fatalf("GOT %d, WANT %d", got, test.want)
			}
		}
	})

	t.Run("Cal timedelta tests", func(t *testing.T) {
		t.Parallel()
		tests := map[string]struct {
			t1   int64
			t2   int64
			want int64
		}{
			"negative delta": {
				t1:   200,
				t2:   2000,
				want: 1800,
			},
			"default": {
				t1:   3500,
				t2:   2000,
				want: 1500,
			},
		}
		for name, test := range tests {
			t.Logf("Running test case: %s", name)
			got := calculateTimedelta(test.t1, test.t2)
			if got != test.want {
				t.Fatalf("GOT %d, WANT %d", got, test.want)
			}
		}
	})
}

type mockGeo struct{}

func (m *mockGeo) GetCoordinatesFromIP(ip string) (*models.Geography, error) {
	switch ip {
	case testdata.TestCurrentIP:
		return &models.Geography{
			Latitude:  34.7725,
			Longitude: 113.7266,
			Radius:    50,
		}, nil
	case testdata.TestPecedingIP:
		return &models.Geography{
			Latitude:  45.4998,
			Longitude: -122.9586,
			Radius:    5,
		}, nil
	case testdata.TestSubsequentIP:
		return &models.Geography{
			Latitude:  37.4627,
			Longitude: 118.4917,
			Radius:    1,
		}, nil
	}
	return nil, nil
}

type testParams struct {
	validPreceding       bool
	validSubsequent      bool
	noPreceding          bool
	noSubsequent         bool
	sameTimestamp        bool
	suspiciousPreceding  bool
	suspiciousSubsequent bool
}

type mockDB struct {
	testParams
}

func (m *mockDB) FindOrCreateUserIPAccessEvent(e *models.UserIPAccessEvent) error {
	return nil
}

func (m *mockDB) FindPrecedingIPAccessEvent(e *models.UserIPAccessEvent) (*models.UserIPAccessEvent, error) {
	switch {
	case m.validPreceding:
		return testdata.GeneratePreviousEvent(false, false), nil
	case m.noPreceding:
		return nil, nil
	case m.sameTimestamp:
		return testdata.GeneratePreviousEvent(false, true), nil
	case m.suspiciousPreceding:
		return testdata.GeneratePreviousEvent(true, false), nil
	}
	return nil, nil
}

func (m *mockDB) FindSubsequentIPAccessEvent(e *models.UserIPAccessEvent) (*models.UserIPAccessEvent, error) {
	switch {
	case m.validSubsequent:
		return testdata.GenerateSubsequentEvent(false, false), nil
	case m.noSubsequent:
		return nil, nil
	case m.sameTimestamp:
		return testdata.GenerateSubsequentEvent(false, true), nil
	case m.suspiciousSubsequent:
		return testdata.GenerateSubsequentEvent(true, false), nil
	}
	return nil, nil
}
