package dupan

import (
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {
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

	pan.Login()

}
