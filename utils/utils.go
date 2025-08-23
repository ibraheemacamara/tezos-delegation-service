package utils

import "time"

func GetYearFromTimestamp(timestamp string) (string, error) {
	time, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "", err
	}
	return time.Format("2006"), nil
}
