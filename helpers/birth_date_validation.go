package helpers

import (
	"fmt"
	"time"
)

func ValidateDate(dateStr string) (bool, time.Time, error) {
	birthDateStr := dateStr
	birthDate, err := time.Parse("02-01-2006", birthDateStr)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("failed to parse birth date: %s", err.Error())
	}

	return true, birthDate, nil
}
