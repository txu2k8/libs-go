package convert

import (
	"strconv"
	"strings"
	"testing"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

func ByteSize(bytes float64) string {
	unit := ""
	value := bytes

	switch {
	case bytes >= EXABYTE:
		unit = "Ei"
		value = value / EXABYTE
	case bytes >= PETABYTE:
		unit = "Pi"
		value = value / PETABYTE
	case bytes >= TERABYTE:
		unit = "Ti"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "Gi"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "Mi"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "Ki"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "Bi"
	case bytes == 0:
		return "0"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

func TestUtilFunc(t *testing.T) {
	// intIP := IP2Int("10.25.119.71")
	// logger.Infof("%v", intIP)

	// ip := Int2IP(169441095)
	// logger.Infof("%v", ip)

	// intArr := StrNumToIntArr("1,4", ",", 2)
	// logger.Infof("%v", intArr)

	byteStr := Byte2String(21902594568192)
	logger.Infof("%v", byteStr)

	byteStr = ByteSize(21902594568192)
	logger.Infof("%v", byteStr)

	// byteI := String2Byte("10.0 kB")
	// logger.Infof("%v", byteI)

	// h, m := To12Hour(16)
	// logger.Infof("%v, %v", h, m)
	// logger.Infof("%v", To24Hour(4, "pm"))

	// arr := []string{"aa", "bb", "cc"}
	// sArr := make([]interface{}, len(arr))
	// for i, v := range arr {
	// 	sArr[i] = v
	// }
	// rArr := ReverseArr(sArr)
	// logger.Info(rArr)
	// logger.Info(ReverseStringArr(arr))

	// logger.Info(TimeStrToInt64("2019-11-21 22:34:03"))
	// logger.Info(TimeInt64ToStr(1574346843))

}
