package dupan

import (
	"io/ioutil"
	"math/rand"
	"net/http"
)
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetHtml(urlStr string) (string, error) {
	request, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36")
	request.Header.Set("Connection", "keep-alive")

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(html), nil
}

