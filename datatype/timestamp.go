package datatype

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimeStamp time.Time

func (date *TimeStamp) Scan(value any) error {
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)
	*date = TimeStamp(nullTime.Time)
	return err
}

func (date TimeStamp) Value() (driver.Value, error) {
	t := time.Time(date)
	if t.IsZero() || t.UnixMicro() == 0 {
		return nil, nil
	}
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Time(date).Location()), nil
}

// GormDataType gorm common data type
func (date TimeStamp) GormDataType() string {
	return "timestamp"
}

func (date TimeStamp) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *TimeStamp) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

func (date TimeStamp) MarshalJSON() ([]byte, error) {
	t := time.Time(date)
	if t.IsZero() {
		return []byte("0"), nil
	}
	return []byte(fmt.Sprintf("%d", t.UnixMicro()/1e3)), nil
}

func (date *TimeStamp) UnmarshalJSON(b []byte) error {
	if len(b) < 1 {
		return nil
	}
	str := string(b)
	if len(str) < 1 {
		return nil
	}
	t, err := strconv.ParseInt(str, 10, 64)
	if t > 0 && err == nil {
		d := time.UnixMilli(t)
		if d.IsZero() {
			return nil
		}
		*date = TimeStamp(d)
		return nil
	}
	if strings.Contains(str, "/") {
		str = strings.ReplaceAll(str, "/", "-")
	}
	d, err := time.Parse("2006-01-02 15:04:05", str)
	if err == nil && !d.IsZero() {
		*date = TimeStamp(d)
		return nil
	}

	return nil
}

func (date TimeStamp) GetTime() time.Time {
	t := time.Time(date)
	return t
}
