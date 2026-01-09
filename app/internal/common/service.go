package common

import "time"

type Service struct {
	Name     string
	URL      string
	IsUp     bool
	LastDown time.Time
}
