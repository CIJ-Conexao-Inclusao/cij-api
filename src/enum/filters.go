package enum

type PeriodFilterEnum string

const (
	LastThreeMonths PeriodFilterEnum = "last_three_months"
	LastSixMonths   PeriodFilterEnum = "last_six_months"
	LastYear        PeriodFilterEnum = "last_year"
)

func (e *PeriodFilterEnum) IsValid() bool {
	switch *e {
	case LastThreeMonths, LastSixMonths, LastYear:
		return true
	}

	return false
}

func (e *PeriodFilterEnum) String() string {
	return string(*e)
}

func GetPeriodFilterEnum(value string) PeriodFilterEnum {
	switch value {
	case "last_three_months":
		return LastThreeMonths
	case "last_six_months":
		return LastSixMonths
	case "last_year":
		return LastYear
	}

	return ""
}
