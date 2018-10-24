package dupan

type Pan struct {
	client   *Client
	guid     string
	sign     string
	bduss    string
	token    string
	bdstoken string
}

func NewPan() *Pan {
	return &Pan{
		client: newClient(),
	}
}

func (this *Pan) getHtml(uri string) (html string, err error) {
	return this.client.getHtml(uri)
}

func (this *Pan) DownloadFile(uri,saveTo string)(err error) {
	return this.client.DownloadFile(uri,saveTo)
}
