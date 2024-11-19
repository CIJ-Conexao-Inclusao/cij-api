package utils

import "cij_api/src/enum"

func PeriodToDays(period enum.PeriodFilterEnum) int {
	switch period {
	case enum.LastThreeMonths:
		return 90
	case enum.LastSixMonths:
		return 180
	case enum.LastYear:
		return 365
	}

	return 0
}
