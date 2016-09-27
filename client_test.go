package xiaomipush

import (
	"testing"
	"time"
)

var packageName string = "sbkssbkssbkssbkssbkssbkssbkssbks"

var client = NewClient("sbkssbkssbkssbkssbkssbkssbkssbks", []string{packageName})

var msg1 *Message = NewAndroidMessage("hi baby1", "hi1").SetPayload("开心么1").SetPassThrough(0)
var msg2 *Message = NewAndroidMessage("hi baby2", "hi2 ").SetPayload("开心么2").SetPassThrough(1)

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
	if err != nil {
		t.Errorf("TestMiPush_Send failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToList(t *testing.T) {
	result, err := client.SendToList(msg1, []string{regID1, regID2})
	if err != nil {
		t.Errorf("TestMiPush_SendToList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendTargetMessageList(t *testing.T) {
	msgList := []*TargetedMessage{NewTargetedMessage(msg1.SetRestrictedPackageName(client.packageName), regID1, TargetTypeRegID), NewTargetedMessage(msg2.SetRestrictedPackageName(client.packageName), regID2, TargetTypeRegID)}
	result, err := client.SendTargetMessageList(msgList)
	if err != nil {
		t.Errorf("TestMiPush_SendTargetMessageList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToAlias(t *testing.T) {
	result, err := client.SendToAlias(msg1, alias1)
	if err != nil {
		t.Errorf("TestMiPush_SendToAlias failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToAliasList(t *testing.T) {
	result, err := client.SendToAliasList(msg1, []string{alias1, alias2})
	if err != nil {
		t.Errorf("TestMiPush_SendToAliasList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToUserAccount(t *testing.T) {
	result, err := client.SendToUserAccount(msg1, account1)
	if err != nil {
		t.Errorf("TestMiPush_SendToUserAccount failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToUserAccountList(t *testing.T) {
	result, err := client.SendToUserAccountList(msg1, []string{account1, account2})
	if err != nil {
		t.Errorf("TestMiPush_SendToUserAccountList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_Broadcast(t *testing.T) {
	result, err := client.Broadcast(msg1, topic1)
	if err != nil {
		t.Errorf("TestMiPush_Broadcast failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_BroadcastAll(t *testing.T) {
	result, err := client.BroadcastAll(msg1)
	if err != nil {
		t.Errorf("TestMiPush_BroadcastAll failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_MultiTopicBroadcast(t *testing.T) {
	result, err := client.MultiTopicBroadcast(msg1, []string{topic1, topic2}, INTERSECTION)
	if err != nil {
		t.Errorf("TestMiPush_MultiTopicBroadcast failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_CheckScheduleJobExist(t *testing.T) {
	result, err := client.CheckScheduleJobExist("slm30b80474526081454i5")
	if err != nil {
		t.Errorf("TestMiPush_CheckScheduleJobExist failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_DeleteScheduleJob(t *testing.T) {
	result, err := client.DeleteScheduleJob("slm30b80474526081454i5")
	if err != nil {
		t.Errorf("TestMiPush_DeleteScheduleJob failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_DeleteScheduleJobByJobKey(t *testing.T) {
	result, err := client.DeleteScheduleJobByJobKey("Xcm45b21474513716292EL")
	if err != nil {
		t.Errorf("TestMiPush_DeleteScheduleJobByJobKey failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_Stats(t *testing.T) {
	result, err := client.Stats("20160922", "20160922", packageName)
	if err != nil {
		t.Errorf("TestMiPush_Stats failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_GetMessageStatusByMsgID(t *testing.T) {
	result, err := client.GetMessageStatusByMsgID("scm23b964745244861922w")
	if err != nil {
		t.Errorf("TestMiPush_GetMessageStatusByMsgID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_GetMessageStatusByJobKey(t *testing.T) {
	result, err := client.GetMessageStatusByJobKey("key111")
	if err != nil {
		t.Errorf("TestMiPush_GetMessageStatusByJobKey failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_GetMessageStatusPeriod(t *testing.T) {
	result, err := client.GetMessageStatusPeriod(time.Now().Add(-time.Hour*24).Unix()*1000, time.Now().Unix()*1000)
	if err != nil {
		t.Errorf("TestMiPush_GetMessageStatusPeriod failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

//----------------------------------------Subscription----------------------------------------//

func TestMiPush_SubscribeTopicForRegID(t *testing.T) {
	result, err := client.SubscribeTopicForRegID(regID1, "topic3", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicForRegID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SubscribeTopicForRegIDList(t *testing.T) {
	result, err := client.SubscribeTopicForRegIDList([]string{regID1, regID2}, "topic5", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicForRegIDList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_UnSubscribeTopicForRegID(t *testing.T) {
	result, err := client.UnSubscribeTopicForRegID(regID1, "topic3", "")
	if err != nil {
		t.Errorf("TestMiPush_UnSubscribeTopicForRegID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_UnSubscribeTopicForRegIDList(t *testing.T) {
	result, err := client.UnSubscribeTopicForRegIDList([]string{regID1, regID2}, "topic5", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicForRegIDList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SubscribeTopicByAlias(t *testing.T) {
	result, err := client.SubscribeTopicByAlias([]string{alias1, alias2}, "topic5", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicByAlias failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_UnSubscribeTopicByAlias(t *testing.T) {
	result, err := client.UnSubscribeTopicByAlias([]string{alias1, alias2}, "topic5", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicByAlias failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

//----------------------------------------Feedback----------------------------------------//

func TestMiPush_GetInvalidRegIDs(t *testing.T) {
	result, err := client.GetInvalidRegIDs()
	if err != nil {
		t.Errorf("TestMiPush_GetInvalidRegIDs failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

//----------------------------------------DevTools----------------------------------------//

func TestMiPush_GetAliasesOfRegID(t *testing.T) {
	result, err := client.GetAliasesOfRegID(regID1)
	if err != nil {
		t.Errorf("TestMiPush_GetAliasesOfRegID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_GetTopicsOfRegID(t *testing.T) {
	result, err := client.GetTopicsOfRegID(regID2)
	if err != nil {
		t.Errorf("TestMiPush_GetTopicsOfRegID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}
