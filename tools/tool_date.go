package tools

import (
	"time"
)

func IsNight() bool {
	isNight := false
	now := time.Now()
	if now.Hour() >= 9 && now.Hour() <= 20 {
		isNight = false
	} else {
		isNight = true
	}
	return isNight
}
