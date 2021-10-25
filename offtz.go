package offtz

import (
	"errors"
	"fmt"
	"time"
)

// TimezonesFromOffset returns list of timezones according to the given offset from UTC in seconds
func TimezonesFromOffset(offset int) ([]string, error) {
	zones := timeZones[offset]
	if zones == nil {
		return nil, fmt.Errorf("unknown offset %d", offset)
	}
	return zones, nil
}

// OffsetFromTimezone returns timezone's short name and the offset from UTC
func OffsetFromTimezone(timezone string) (name string, offset string, offsetValue int, err error) {
	//for offset, tzs := range timeZones {
	//	for _, tz := range tzs {
	//		if tz == timezone {
	//			quotient, remainder := offset/3600, offset%3600
	//			return "", fmt.Sprintf("%+03d:%02d", quotient, remainder), offset, nil
	//		}
	//	}
	//}
	//return "", "", 0, errors.New("unknown time zone " + timezone)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", "", 0, errors.New("unknown time zone " + timezone)
	}

	tzName, tzOffset := time.Now().In(loc).Zone()
	quotient, remainder := tzOffset/3600, tzOffset%3600
	offset = fmt.Sprintf("%+03d:%02d", quotient, remainder)

	return tzName, offset, tzOffset, nil
}
