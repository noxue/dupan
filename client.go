package dupan

import (
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
)

type Client struct {
	client *http.Client
}

func newClient() *Client {
	jar, err := cookiejar.New(nil)

	if err != nil {
		glog.Error(err.Error())
	}
	//proxy, err := url.Parse("http://127.0.0.1:8888")
	//if err != nil {
	//	glog.Error(err.Error())
	//}
	return &Client{
		client: &http.Client{
			Jar: jar,
			//Transport: &http.Transport{
			//	Proxy: http.ProxyURL(proxy),
			//},
		},
	}
}

func (this *Client) commonHeader(request *http.Request) {
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36")
	request.Header.Set("Referer", "https://pan.baidu.com")
	request.Header.Set("Connection", "keep-alive")
}

func (this *Client) getHtml(uri string) (string, error) {

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
	return "", err
	}

	this.commonHeader(request)

	res, err := this.client.Do(request)
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

func (this *Client) DownloadFile(uri,saveTo string)(err error) {
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	this.commonHeader(request)
	res, err := this.client.Do(request)
	if err != nil {
		return
	}
	defer res.Body.Close()
	f, err := os.Create(saveTo)
	if err != nil {
		return
	}
	defer f.Close()
	_,err = io.Copy(f, res.Body)
	return
}
