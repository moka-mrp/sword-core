package db

import "time"

type SamTime time.Time
func (j SamTime) MarshalJSON() ([]byte, error) {
	if time.Time(j).IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"`+time.Time(j).Format("2006-01-02 15:04:05")+`"`), nil
}

