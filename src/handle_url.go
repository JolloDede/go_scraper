package src

import (
	"net/http"
	"slices"
)

var dic map[string]int

func HandleUrl(url string, recursiv bool) map[string]int {
	dic = make(map[string]int)

	checkPath(url)

	return dic
}

func checkPath(url string) {
	res, err := http.Get(url)

	if err != nil {
		return
	}

	if slices.Contains([]int{300, 303}, res.StatusCode) {
		locations := res.Header["Location"]
		for i := 0; i < len(locations); i++ {
			checkPath(locations[i])
		}
	}

	if slices.Contains([]int{301, 302, 305, 307, 308}, res.StatusCode) {
		location := res.Request.Response.Header["Location"]
		checkPath(location[0])
	}

	if dic[url] == 0 {
		dic[url] = res.StatusCode
	}
}
