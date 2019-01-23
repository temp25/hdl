package helper

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

var requestHeaders = map[string]string{
	"Hotstarauth": GenerateHotstarAuth(),
	"X-Country-Code": "IN",
	"X-Platform-Code": "JIO",
}

func checkIfError(err error) {
	if err!=nil {
		panic(err)
	}
}

func GetPageContents(url string, hasHeaders bool) string{
	
	request, err := http.NewRequest("GET", url, nil)
	
	checkIfError(err)

	if hasHeaders {
		for headerName, headerValue := range requestHeaders {
			request.Header.Set(headerName, headerValue)
		}
	}

	response, err := http.DefaultClient.Do(request)

	checkIfError(err)

	defer response.Body.Close()

	htmlBytes, err := ioutil.ReadAll(response.Body)

	checkIfError(err)

	htmlContent := fmt.Sprintf("%s", htmlBytes)

	return htmlContent

}
