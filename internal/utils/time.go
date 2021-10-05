package utils

import (
	"fmt"
	"strconv"
	"time"

	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/protobuf/types/known/durationpb"
)

// TimeToProtoDateTime
func TimeToProtoDateTime(t time.Time) (*datetime.DateTime, error) {
	dt := &datetime.DateTime{
		Year:    int32(t.Year()),
		Month:   int32(t.Month()),
		Day:     int32(t.Day()),
		Hours:   int32(t.Hour()),
		Minutes: int32(t.Minute()),
		Seconds: int32(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}

	// If the location is a UTC offset, encode it as such in the proto.
	loc := t.Location().String()
	if match := offsetRegexp.FindStringSubmatch(loc); len(match) > 0 {
		offsetInt, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
		dt.TimeOffset = &datetime.DateTime_UtcOffset{
			UtcOffset: &durationpb.Duration{Seconds: int64(offsetInt) * 3600},
		}
	} else if loc != "" {
		dt.TimeOffset = &datetime.DateTime_TimeZone{
			TimeZone: &datetime.TimeZone{Id: loc},
		}
	}

	return dt, nil
}

// ProtoDateTimeToTime
func ProtoDateTimeToTime(d *datetime.DateTime) (time.Time, error) {
	var err error

	// Determine the location.
	loc := time.UTC
	if tz := d.GetTimeZone(); tz != nil {
		loc, err = time.LoadLocation(tz.GetId())
		if err != nil {
			return time.Time{}, err
		}
	}
	if offset := d.GetUtcOffset(); offset != nil {
		hours := int(offset.GetSeconds()) / 3600
		loc = time.FixedZone(fmt.Sprintf("UTC%+d", hours), hours)
	}

	// Return the Time.
	return time.Date(
		int(d.GetYear()),
		time.Month(d.GetMonth()),
		int(d.GetDay()),
		int(d.GetHours()),
		int(d.GetMinutes()),
		int(d.GetSeconds()),
		int(d.GetNanos()),
		loc,
	), nil
}
