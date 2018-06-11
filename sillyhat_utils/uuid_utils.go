package sillyhat_utils

import (
	"strings"
	"github.com/google/uuid"
)

func GeneratorUUID() string {
	return strings.ToUpper(strings.Replace(uuid.New().String(),"-","",-1))
}