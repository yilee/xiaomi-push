package xiaomipush

import (
	"fmt"
	"testing"
)

var packageName string = "com.xiaomi.mipushdemo"

var client = NewClient("yourappSecret", packageName)

var msg1 *Message = NewAndroidMessage("hi baby1", "hi1").SetRestrictedPackageNames([]string{packageName}).SetPayload("this is payload1").SetPassThrough(0)
var msg2 *Message = NewAndroidMessage("hi baby2", "hi2 ").SetRestrictedPackageNames([]string{packageName}).SetPayload("this is payload2").SetPassThrough(0)

var regID1 string = "WFioJi0fiIco7vOrI4dnxxjeKAUqR7fjugoGkHUgxeo="
var regID2 string = "52Pe7fPIRXWsXhzn4eYJ1njYhBhN8Lcp8IJPOMjThdk="

var alias1 string = "alias1"
var alias2 string = "alias2"

var account1 string = "account1"
var account2 string = "account2"

var topic1 string = "topic1"
var topic2 string = "topic2"

func TestMiPush_Send(t *testing.T) {
	result, err := client.Send(msg1, regID1)
	fmt.Println("result1, err1", result, err)
}

func TestMiPush_SendToList(t *testing.T) {
	result, err := client.SendToList(msg1, []string{regID1, regID2})
	fmt.Println("result1, err1", result, err)
}

// Not Finished
func TestMiPush_SendTargetMessageList(t *testing.T) {
	msgList := []*TargetedMessage{NewTargetedMessage(msg1, regID1, TargetTypeRegID), NewTargetedMessage(msg2, regID2, TargetTypeRegID)}
	result, err := client.SendTargetMessageList(msgList)
	fmt.Println("result1, err1", result, err)
}

func TestMiPush_SendToAlias(t *testing.T) {
	result, err := client.SendToAlias(msg1, alias1)
	fmt.Println("result1, err1", result, err)
}

func TestMiPush_SendToAliasList(t *testing.T) {
	result, err := client.SendToAliasList(msg1, []string{alias1, alias2})
	fmt.Println("result1, err1", result, err)
}

func TestMiPush_SendToUserAccount(t *testing.T) {
	result, err := client.SendToUserAccount(msg1, account1)
	fmt.Println("result1, err1", result, err)
}

func TestMiPush_SendToUserAccountList(t *testing.T) {
	result, err := client.SendToUserAccountList(msg1, []string{account1, account2})
	fmt.Println("result1, err1", result, err)
}

// Not Finished
func TestMiPush_Broadcast(t *testing.T) {
	result, err := client.Broadcast(msg1, topic1)
	fmt.Println("result1, err1", result, err)
}

func TestMiPush_BroadcastAll(t *testing.T) {
	result, err := client.BroadcastAll(msg1)
	fmt.Println("result1, err1", result, err)
}

func TestMiPush_MultiTopicBroadcast(t *testing.T) {
	result, err := client.MultiTopicBroadcast(msg1, []string{topic1, topic2}, INTERSECTION)
	fmt.Println("result1, err1", result, err)
}
