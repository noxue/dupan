package dupan

import (
	"fmt"
	"testing"
)

func TestPan_List(t *testing.T) {
	pan := NewPan()
	uri, _ := pan.GetQrImg("d:/1.png")

	fmt.Println(uri)

	for i:=0;i<10 ;i++  {
		ok,err:=pan.UniCast()
		if err != nil {
			panic(err)
		}
		if ok {
			break
		}
	}

	err:=pan.Login()
	if err!=nil {
		panic(err)
	}

	info,err:=pan.Quota()
	if err!=nil {
		panic(err)
	}
	fmt.Println(*info)
	list,err:=pan.List("/",1,100,"",true)
	if err!= nil {
		panic(err)
	}
	for _,v:=range list.List{
		fmt.Println(v)

		fds:=[]int64{v.Fs_id}
		info,err:=pan.Share(fds,"1234")
		fmt.Println(err, *info)
		break
	}
}
