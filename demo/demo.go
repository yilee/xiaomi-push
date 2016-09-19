package main

import (
	"fmt"

	xm "github.com/yilee/xiaomi-push"
)

func main() {
	msg := xm.NewMessage("hi baby")
	client := xm.NewClient("xxxxxsad", "packagename")
	result, err := client.Push(msg, []string{"fake_reg_id_1", "fake_reg_id_2"})
	if err != nil {
		return
	}
	fmt.Printf("result:%#v\n", result)
}
