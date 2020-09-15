package utils

import (
	"testing"
)

func TestUtilFunc(t *testing.T) {
	logger.Info(StrNumToIntArr("1,2", ",", 3))
}
