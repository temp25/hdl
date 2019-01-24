package tests

import (
	"github.com/temp25/hdl/helper"
	"reflect"
	"testing"
)

var testNumInt int64 = 1100015498
var testNumFloat float64 = float64(testNumInt)

func TestCopyMap(t *testing.T) {

	originalMap := map[string]interface{}{
		"Name":        "Scott",
		"Age":         25,
		"Is Employed": true,
	}

	clonedMap := helper.CopyMap(originalMap)

	mapEquality := reflect.DeepEqual(originalMap, clonedMap)

	if !mapEquality {
		t.Error("Expected", originalMap, " but got", clonedMap)
	}

}

func TestCountDigits(t *testing.T) {

	var expectedDigits int64 = 10
	actualDigits := helper.CountDigits(testNumInt)

	if expectedDigits != actualDigits {
		t.Errorf("\nExpected '%d' but got '%d'\n", expectedDigits, actualDigits)
	}

}

func TestPadZeroRight(t *testing.T) {

	var expectedPaddedNumber int64 = 1100015498000

	actualPaddedNumber := helper.PadZeroRight(testNumInt)

	if expectedPaddedNumber != actualPaddedNumber {
		t.Errorf("\nExpected '%d' but got '%d'\n", expectedPaddedNumber, actualPaddedNumber)
	}
}

func TestGetDateStr(t *testing.T) {
	expectedDateString := "1970-01-13 23:03:35.498 +0530 IST"

	actualDateString := helper.GetDateStr(testNumFloat)

	if expectedDateString != actualDateString {

		t.Errorf("\nExpected '%s' but got '%s'\n", expectedDateString, actualDateString)
	}
}
