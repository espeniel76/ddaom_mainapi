package tools

import (
	"strings"

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
