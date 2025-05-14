package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	Layout = "20060102"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	parseRepeat := strings.Split(repeat, " ")

	if len(parseRepeat) < 2 && parseRepeat[0] == "d" {
		return "", errors.New("invalid repeat format: expected 'd <number>")
	}

	parseTimeDstart, err := time.Parse(Layout, dstart)
	if err != nil {
		return "", fmt.Errorf("error parse dstart: %w", err)
	}
	
	var next time.Time

	switch parseRepeat[0] {
	case "d":
		days, err := strconv.Atoi(parseRepeat[1])
		if err != nil {
			return "", fmt.Errorf("error convertation: %w", err)
		}

		if days > 400 {
			return "", errors.New("the maximum value of the day can be 400")
		}
		
		for {
			next = parseTimeDstart.AddDate(0, 0, days)

			if ComparingDate(now, next) {
				break
			} else {
				parseTimeDstart = next
			}
		}

	case "y":
		for {
			next = parseTimeDstart.AddDate(1, 0, 0)
			if ComparingDate(now, next) {
				break
			} else {
				parseTimeDstart = next
			}
		}
	default:
		return "", fmt.Errorf("unsupported format %w", parseRepeat[0])
	}

	return fmt.Sprint(next.Format(Layout)), nil
}