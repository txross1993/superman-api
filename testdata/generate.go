package testdata

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"time"

	"github.com/txross1993/superman-api/models"
)

const (
	TestUser             = "bob"
	TestCurrentIP        = "42.222.21.19"
	TestPecedingIP       = "73.11.21.110" // distance from first: 5858 miles
	TestSubsequentIP     = "27.202.31.1"  // distance from first: 324 miles
	TestCurrentTimestmap = int64(1514764800)
)

var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateCurrentEvent creates a current test user ip access event
func GenerateCurrentEvent() *models.UserIPAccessEvent {
	return &models.UserIPAccessEvent{
		EventUUID:     getUUID(),
		Username:      TestUser,
		UnixTimestamp: TestCurrentTimestmap,
		IPAddress:     TestCurrentIP,
	}
}

// GeneratePreviousEvent creates a preceding event that may flag
// the superman api if either suspicious or instantaneous flags are true
func GeneratePreviousEvent(suspcious, instantaneous bool) *models.UserIPAccessEvent {
	return &models.UserIPAccessEvent{
		EventUUID:     getUUID(),
		Username:      TestUser,
		UnixTimestamp: precedingTS(TestCurrentTimestmap, suspcious, instantaneous),
		IPAddress:     TestPecedingIP,
	}
}

// GenerateSubsequentEvent creates a subsequent event that may flag
// the superman api if either suspicious or instantaneous flags are true
func GenerateSubsequentEvent(suspcious, instantaneous bool) *models.UserIPAccessEvent {
	return &models.UserIPAccessEvent{
		EventUUID:     getUUID(),
		Username:      TestUser,
		UnixTimestamp: subsequentTS(TestCurrentTimestmap, suspcious, instantaneous),
		IPAddress:     TestSubsequentIP,
	}
}

func precedingTS(currentTS int64, suspicious bool, instantaneous bool) int64 {
	if instantaneous {
		return currentTS
	}

	supsiciousThreshold := int64((5858 * 3600) / 500)
	timeDelta := getRandTimeDelta(supsiciousThreshold - 1)
	if suspicious {
		// the time delta must be < suspiciousThreshold
		return currentTS - timeDelta
	}

	// the time delta must be >= suspiciousThreshold
	return currentTS - (supsiciousThreshold + timeDelta)
}

func subsequentTS(currentTS int64, suspicious bool, instantaneous bool) int64 {
	if instantaneous {
		return currentTS
	}

	supsiciousThreshold := int64((324 * 3600) / 500)
	timeDelta := getRandTimeDelta(supsiciousThreshold - 1)
	if suspicious {
		// the time delta must be < suspiciousThreshold
		return currentTS + timeDelta
	}

	// the time delta must be >= suspiciousThreshold
	return currentTS + (supsiciousThreshold + timeDelta)
}

func getRandTimeDelta(max int64) int64 {
	return r.Int63n(max)
}

func getUUID() string {
	b := make([]byte, 16)
	crand.Read(b)

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
