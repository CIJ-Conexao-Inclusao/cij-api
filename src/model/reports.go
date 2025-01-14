package model

type DisabilityTotals struct {
	Visual       int `json:"visual"`
	Hearing      int `json:"hearing"`
	Physical     int `json:"physical"`
	Intellectual int `json:"intellectual"`
	Psychosocial int `json:"psychosocial"`
}

type DisabilityTotalsByNeighborhood = DisabilityTotals

type CountActivitiesByPeriod struct {
	ActivityType string         `json:"activityType"`
	MonthsCount  map[string]int `json:"monthsCount"`
}
