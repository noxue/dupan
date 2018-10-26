package dupan

import (
	"fmt"
	"strings"
)

type Pan struct {
	Username string
	client   *Client
	Guid     string
	sign     string
	bduss    string
	token    string
	bdstoken string
}

func NewPan() *Pan {
	guid := fmt.Sprintf("%s-%s-%s-%s-%s", RandStr(7), RandStr(4), RandStr(4), RandStr(4), RandStr(12))
	return &Pan{
		Guid:   strings.ToUpper(guid),
		client: newClient(),
	}
}

func (this *Pan) GetHtml(uri string) (html []byte, err error) {
	return this.client.getHtml(uri)
}

func (this *Pan) DownloadFile(uri, saveTo string) (err error) {
	return this.client.DownloadFile(uri, saveTo)
}
