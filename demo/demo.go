package main

import (
	"fmt"

	xm "github.com/yilee/xiaomi-push"
)

var client = xm.NewClient("appSecret", "packageName")

func main() {
	msg := xm.NewMessage("hi baby", "hi ")
	result, err := client.Send(msg, "fake_reg_id_1")
	fmt.Println("result1, err1", result, err)

	result2, err := client.Stats("20160901", "20160902")
	fmt.Println("result2, err", result2, err)

	result3, err := client.Status("Xlm35b23474365994495Hu")
	fmt.Println("result3, err", result3, err)

}
