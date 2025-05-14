package utils

import (
	"fmt"
	"time"

	"github.com/LoaltyProgramm/to-do-app/internal/models"
)

func CheckingTheDateUsingATemplate(data string) bool {
	layoutInput := "02.01.2006"
	_, err := time.Parse(layoutInput, data)
	return err == nil
}

func CheckDate(task *models.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(Layout)
	}

	t, err := time.Parse(Layout, task.Date)
	if err != nil {
		return fmt.Errorf("Error parse string in time.Time: %w", err)
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	
	if ComparingDate(t, now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(Layout)
		} else {
			task.Date = next
		}
	}

	return nil
}

func ComparingDate(nowDate, nextDate time.Time) bool {
	nextDateStr := nextDate.Format(Layout)
	nowDateStr := nowDate.Format(Layout)

	if nextDateStr > nowDateStr {
		return true
	} else {
		return false
	}
}