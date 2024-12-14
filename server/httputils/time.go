package httputils

import "time"

func ParseTime(date string) (time.Time, error) {
	format := "2006-01-02 15:04:05" // ? Golang use 2006 as format
	location, _ := time.LoadLocation("UTC")
	parse, err := time.ParseInLocation(format, date, location)

	return parse, err
}
