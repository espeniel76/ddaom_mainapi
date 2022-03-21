package tools

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CutString(str string, size int) string {
	if len(str) > size {
		str = str[:size-2] + "..."
	}
	return str
}

func MakeUUID() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}

func TodayFormattedDateFull() string {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}

func TodayFormattedDate() string {
	t := time.Now()
	formatted := fmt.Sprintf("%d%02d%02d",
		t.Year(), t.Month(), t.Day(),
	)
	return formatted
}
