package main

import (
	"fmt"

	"github.com/yilee/xiaomi-push/xiaomipush"
)

func main() {
	msg := xiaomipush.NewMessage("hi baby")
	client := xiaomipush.NewClient("xxxxxsad", "packagename")
	result, err := client.Push(msg, []string{"fake_reg_id_1", "fake_reg_id_2"})
	if err != nil {
		return
	}
	fmt.Printf("result:%#v\n", result)
}
