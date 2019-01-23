package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = vm
		} else {
			cp[k] = v
		}
	}
	return cp
}

/*
func After(value string, a string) string{
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
*/

func PadZeroRight(num int64) int64 {
	tmp := fmt.Sprintf("%-13d", num)
	tmp = strings.Replace(tmp, " ", "0", -1)
	paddedNum, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		panic(err)
	}
	return paddedNum
}

func CountDigits(i int64) (count int64) {
	for i != 0 {
		i /= 10
		count = count + 1
	}
	return count
}

func GetDateStr(timeFloat64 float64) string {
	timeMillis := int64(timeFloat64)
	timeMillisPadded := PadZeroRight(timeMillis)
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic(err)
	}
	return time.Unix(0, timeMillisPadded*int64(time.Millisecond)).In(location).String()
}
