package dupan

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
	"time"
)



// 获取验证码地址
func (this *Pan) GetQrImg(saveTo string) (qrUrl string, err error) {
	type T struct {
		Errno  int    `json:"errno"`
		ImgUrl string `json:"imgurl"`
		Sign   string `json:"sign"`
	}

	this.guid = fmt.Sprintf("%s-%s-%s-%s-%s", RandStr(7), RandStr(4), RandStr(4), RandStr(4), RandStr(12))
	this.guid = strings.ToUpper(this.guid)
	html, err := this.getHtml(fmt.Sprintf("https://passport.baidu.com/v2/api/getqrcode?lp=pc&gid=%s&apiver=v3&tt=%d&tpl=netdisk&_=%d",
		this.guid,
		time.Now().UnixNano()/1000000,
		time.Now().UnixNano()/1000000,
	))

	if err != nil {
		return
	}

	var t T
	err = json.Unmarshal([]byte(html), &t)
	if err != nil {
		return
	}

	this.sign = t.Sign
	qrUrl = t.ImgUrl

	if saveTo != "" {
		err = this.DownloadFile("https://"+t.ImgUrl, saveTo)
	}
	return
}

func _getVmValue(vm *otto.Otto, name string) (value otto.Value, err error) {
	value, err = vm.Get(name)
	if err != nil {
		err = errors.New("获取" + name + "失败")
		return
	}

	return
}

func getVmInt(vm *otto.Otto, name string) (n int64, err error) {
	value, err := _getVmValue(vm, name)
	if err != nil {
		return
	}
	n, err = value.ToInteger()
	if err != nil {
		err = errors.New(name + "转数字失败")
		return
	}
	return
}

func getVmStr(vm *otto.Otto, name string) (str string, err error) {
	value, err := _getVmValue(vm, name)
	if err != nil {
		return
	}
	str, err = value.ToString()
	if err != nil {
		err = errors.New(name + "转字符串失败")
		return
	}
	return
}

// ok true表示扫码登陆成功，false表示还未确定登陆
func (this *Pan) UniCast() (ok bool, err error) {
	uri := fmt.Sprintf("https://passport.baidu.com/channel/unicast?channel_id=%s&tpl=netdisk&gid=%s&callback=&apiver=v3&tt=%d&_=%d",
		this.sign,
		this.guid,
		time.Now().UnixNano()/1000000,
		time.Now().UnixNano()/1000000,
	)
	fmt.Println(uri)
	html, err := this.getHtml(uri)
	if err != nil {
		return
	}

	vm := otto.New()

	str := fmt.Sprintf(`
		var text = %s;
		var errno = text.errno
		var status = -1
		var v = ""
		if (errno == 0){
			var t1 = JSON.parse(text.channel_v)
			if (t1) {
				status = t1.status
				v = t1.v
			}
		}
	`, html)
	vm.Run(str)
	if err != nil {
		err = errors.New("执行js失败")
		return
	}

	errno, err := getVmInt(vm, "errno")
	if err != nil {
		return
	}

	status, err := getVmInt(vm, "status")
	if err != nil {
		return
	}

	if errno == 0 && status == 0 { // 扫描点确定了
		this.bduss, err = getVmStr(vm, "v")
		if err != nil {
			return
		}
		ok = true
	}
	return
}

func (this *Pan) Login() (err error) {
	uri := fmt.Sprintf("https://passport.baidu.com/v3/login/main/qrbdusslogin?v=%d&bduss=%s&u=%s&loginVersion=v4&qrcode=1&tpl=netdisk&apiver=v3&tt=%d",
		time.Now().UnixNano()/1000000,
		this.bduss,
		"https%253A%252F%252Fpan.baidu.com%252Fdisk%252Fhome",
		time.Now().UnixNano()/1000000,
	)
	_,err= this.getHtml(uri)
	if err != nil {
		return
	}

	html, err := this.getHtml("https://pan.baidu.com/disk/home")
	if err != nil {
		return
	}

	reg, err := regexp.Compile(`(var context={"loginstate"[^\n]+)`)
	if err != nil {
		return
	}
	arr := reg.FindStringSubmatch(html)
	if len(arr) != 2 {
		err = errors.New("获取登陆有的页面用户信息失败，可能登陆失败")
		return
	}
	err = this.getUserInfo(arr[1])
	return
}

//获取用户信息
/**
XDUSS: "pansec_DCb740ccc5511e5e8fedcff06b081203-YKddAQmEETw1p%2BwxpxO56%2F69NwpBG03TU4zQ0KWyYA93x1kPHouql5%2Blprki69fOVp5cv93EIusvQ4tQ9rr1wRzWNFQlb5J%2FTkBJWpojwyWA8%2BhZ%2F4ynth2z8wftAg2bBySLUZZmX%2Bz7ZEY7aai6FAUxQ%2BMkhQTqhdOe3dP%2Be7GPGNs5r%2BA%2B0dvNwce1BNPWHXRm0prug7htBV96EoumTGOrESlO%2F8jN%2FvBgr3ZW6x3agmLIW3X6VnNDR972TBAcRMt23Bx4sYo6LXMUliXdzg%3D%3D"
activity_end_time: 0
activity_status: 0
applystatus: 1
bdstoken: "a65502953c80857f15dd751816754e08"
bt_paths: null
curr_activity_code: 0
face_status: 0
file_list: null
flag: 1
is_auto_svip: 0
is_evip: 0
is_svip: 0
is_vip: 0
is_year_vip: 0
loginstate: 1
need_tips: null
pansuk: "Sa0GAZQi_IHLK5Dd7y5ltw"
photo: "https://ss0.bdstatic.com/7Ls0a8Sm1A5BphGlnYG/sys/portrait/item/3ea7be0a.jpg"
sampling: {expvar: Array(6)}
sharedir: 0
show_vip_ad: 0
sign1: "a0a5eb8ece1d116ca7e62cb146853dbc01c5285c"
sign2: "function s(j,r){var a=[];var p=[];var o="";var v=j.length;for(var q=0;q<256;q++){a[q]=j.substr((q%v),1).charCodeAt(0);p[q]=q}for(var u=q=0;q<256;q++){u=(u+p[q]+a[q])%256;var t=p[q];p[q]=p[u];p[u]=t}for(var i=u=q=0;q<r.length;q++){i=(i+1)%256;u=(u+p[i])%256;var t=p[i];p[i]=p[u];p[u]=t;k=p[((p[i]+p[u])%256)];o+=String.fromCharCode(r.charCodeAt(q)^k)}return o};"
sign3: "e8c7d729eea7b54551aa594f942decbe"
srv_ts: 1540357696
task_key: "72541c9e3cf0fef89c6e54fbc183daa8a02b0d1d"
task_time: 1540357696
third: 0
timeline_status: 1
timestamp: 1540357696
token: "132edmjOwDGW2EWSRNJcSQsw/OIRRo4kngcJOXTFVatTjVMz2tbB0jnFA4miTGOoL9QlKJeVISof3nrcUJdb6i0ZlWBXQ4bNnImm71gUj2rfcaVLNe8pookKaw1eGy1PMvwfvCSajX5EUcVkTLazkn6kDbg8YKluarEjXBdTqFT5VydD58sgbY+o9W6DnnruPXMQrJPRMNkB00Vd9ZTijiUe8SyqPDSxLzPvZ0A7f0ncfY7AVS20AG1yd+qvZ96oOMxpvQXp/zjlu3HzSYcjZNQUMuzXIcjW"
uk: 2802741345
urlparam: []
username: "173126019"
vip_end_time: null
vol_autoup: null
 */
func (this *Pan) getUserInfo(str string) (err error) {
	vm := otto.New()
	names:=[]string{"token","bdstoken"}
	for _,v:=range names{
		str += "var "+v+" = context."+v+";"
	}

	vm.Run(str)
	if err != nil {
		err = errors.New("执行js失败")
		return
	}

	this.token, err = getVmStr(vm, "token")
	if err != nil {
		return
	}
	this.bdstoken, err = getVmStr(vm, "bdstoken")
	if err != nil {
		return
	}

	return nil
}
