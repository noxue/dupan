package dupan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type BaiduUrlInfo struct {
	Errno      int    `json:"errno"`
	Request_id int64  `json:"request_id"`
	Shareid    int64  `json:"shareid"`
	Link       string `json:"link"`
	Shorturl   string `json:"shorturl"`
	Ctime      int32  `json:"ctime"`
	Premis     bool   `json:"premis"`
	Password   string
}

type BaiduPanInfo struct {
	Errno          int   `json:"errno"`
	Is_show_window int   `json:"is_show_window"`
	Total          int64 `json:"total"`
	Free           int64 `json:"free"`
	Request_id     int64 `json:"request_id"`
	Expire         bool  `json:"expire"`
	Used           int64 `json:"used"`
}

type BdItem struct {
	Local_mtime     int64  `json:"local_mtime"`
	Path            string `json:"path"`
	Server_mtime    int64  `json:"server_mtime"`
	Server_ctime    int64  `json:"server_ctime"`
	Isdir           int    `json:"isdir"`
	Server_filename string `json:"server_filename"`
	Fs_id           int64  `json:"fs_id"`
	Unlist          int    `json:"unlist"`
	Dir_empty       int    `json:"dir_empty"`
	Oper_id         int    `json:"oper_id"`
	Category        int    `json:"category"`
	Size            int64  `json:"size"`
}

type BdList struct {
	Errno    int      `json:"errno"`
	Has_more int      `json:"has_more"`
	List     []BdItem `json:"list"`
}

//获取文件列表
func (this *Pan) List(path string, page, num int, order string, isDesc bool) (*BdList, error) {

	//默认第一页
	if page <= 0 {
		page = 1
	}

	//不让一次载入过多的文件
	if num <= 0 || num > 100 {
		num = 100
	}
	if order == "" {
		order = "time"
	}
	desc := 1
	if !isDesc {
		desc = 0
	}
	listUrl := fmt.Sprintf("https://pan.baidu.com/api/list?dir=%s&bdstoken=%s&logid==&num=%d&order=%s&desc=%d&clienttype=0&showempty=0&web=1&page=%d&channel=chunlei&web=1&app_id=250528",
		url.QueryEscape(path),
		this.bdstoken,
		num,
		order,
		desc,
		page,
	)
	html, err := this.getHtml(listUrl)
	if err != nil {
		return nil, err
	}
	dat := &BdList{}
	if err := json.Unmarshal([]byte(html), dat); err == nil {
		return dat, nil
	} else {
		return nil, err
	}
}

//获取网盘空间大小信息
func (this *Pan) Quota() (*BaiduPanInfo, error) {
	urlStr := "https://pan.baidu.com/api/quota?checkexpire=1&checkfree=1&channel=chunlei&web=1&app_id=250528&bdstoken=" + this.bdstoken + "&logid==&clienttype=0"
	html, err := this.getHtml(urlStr)
	if err != nil {
		return nil, err
	}

	bdInfo := &BaiduPanInfo{}

	err = json.Unmarshal([]byte(html), bdInfo)
	if err != nil {
		return nil, err
	}

	return bdInfo, nil
}

//分享链接,文件编号，分享文件密码
func (this *Pan) Share(fs_ids []int64, pwd string) (*BaiduUrlInfo, error) {

	//如果不等于空，表示需要密码，就限制为4个字符
	if pwd != "" {
		if len(pwd) != 4 {
			return nil, errors.New("提取密码字符长度必须为4个字符")
		}
	}

	urlStr := "https://pan.baidu.com/share/set?channel=chunlei&clienttype=0&web=1&channel=chunlei&web=1&app_id=250528&bdstoken=" + this.bdstoken + "&logid=MTU0MDM2NDc5OTA2NTAuMDkxNDA4MzY3ODM3MDYzODU=&clienttype=0"

	//把文件数字编号列表转成字符列表
	fs_idstr := make([]string, 0, 10)

	for _, d := range fs_ids {
		fs_idstr = append(fs_idstr, strconv.FormatInt(d, 10))
	}

	postData := &url.Values{
		"fid_list":     {"[" + strings.Join(fs_idstr, ",") + "]"},
		"schannel":     {"0"},
		"period":       {"0"},
		"channel_list": {"[]"},
	}

	//如果是生成带密码的链接
	if pwd != "" {
		postData.Set("schannel", "4")
		postData.Add("pwd", pwd)
	}

	request, err := http.NewRequest("POST", urlStr, strings.NewReader(postData.Encode()))
	if err != nil {
		return nil, err
	}

	this.client.commonHeader(request)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := this.client.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bdUrlInfo := &BaiduUrlInfo{}
	bdUrlInfo.Password = pwd
	err = json.Unmarshal(html, bdUrlInfo)
	if err != nil {
		return nil, err
	}

	return bdUrlInfo, nil
}

//取消分享，shareid_list 分享后获取到的id号列表
func (this *Pan) CancelShare(shareid_list []int64) {
	urlStr := "https://pan.baidu.com/share/cancel?bdstoken=" + this.bdstoken + "&channel=chunlei&clienttype=0&web=1"

	//把文件数字编号列表转成字符列表
	shareid_list_str := make([]string, 0, 10)

	for _, d := range shareid_list {
		shareid_list_str = append(shareid_list_str, strconv.FormatInt(d, 10))
	}
	postData := &url.Values{
		"shareid_list": {"[" + strings.Join(shareid_list_str, ",") + "]"},
		"type":         {"1"},
	}

	request, err := http.NewRequest("POST", urlStr, strings.NewReader(postData.Encode()))

	this.client.commonHeader(request)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := this.client.client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
}

//根据关键词查询, 返回文件列表信息和错误信息
func (this *Pan) Search(key string, page int, num int) (*BdList, error) {

	//默认第一页
	if page <= 0 {
		page = 1
	}

	//不让一次载入过多的文件
	if num <= 0 || num > 100 {
		num = 100
	}

	searchUrl := "https://pan.baidu.com/api/search?recursion=1&order=time&desc=1&showempty=0&web=1&page=" + strconv.Itoa(page) + "&num=" + strconv.Itoa(num) + "&key=" + key + "&t=0.9965368690407208&channel=chunlei&web=1&app_id=250528&bdstoken=" + this.bdstoken + "&logid==&clienttype=0"
	html, err := this.getHtml(searchUrl)
	if err != nil {
		return nil, err
	}

	bdList := &BdList{}
	err = json.Unmarshal([]byte(html), bdList)
	if err != nil {
		return nil, err
	}

	return bdList, nil
}

//判断分享链接是否有效
func (this *Pan) CheckShareUrl(url string) (err error) {
	html, err := GetHtml(url)

	if err != nil {
		return
	}

	pos := strings.Index(html, "share_nofound_des")
	if pos != -1 {
		err = errors.New("链接已失效")
	}
	return
}
