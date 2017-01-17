package xiaomipush

import (
	"testing"
	"time"

	"golang.org/x/net/context"
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
	result, err := client.Send(context.Background(), msg1, regID1)
	if err != nil {
		t.Errorf("TestMiPush_Send failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendWithTimeout(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	time.Sleep(time.Second)
	_, err := client.Send(ctx, msg1, regID1)
	if err == nil {
		t.Errorf("TestMiPush_Send failed :%v\n", err)
	} else {
		t.Logf("err=%v\n", err)
	}
}

func TestMiPush_SendWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.Send(ctx, msg1, regID1)
	if err == nil {
		t.Errorf("TestMiPush_Send failed :%v\n", err)
	} else {
		t.Logf("err=%v\n", err)
	}
}

func TestMiPush_SendToList(t *testing.T) {
	result, err := client.SendToList(context.TODO(), msg1, []string{regID1, regID2})
	if err != nil {
		t.Errorf("TestMiPush_SendToList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendTargetMessageList(t *testing.T) {
	msgList := []*TargetedMessage{NewTargetedMessage(msg1.SetRestrictedPackageName(client.packageName), regID1, TargetTypeRegID), NewTargetedMessage(msg2.SetRestrictedPackageName(client.packageName), regID2, TargetTypeRegID)}
	result, err := client.SendTargetMessageList(context.TODO(), msgList)
	if err != nil {
		t.Errorf("TestMiPush_SendTargetMessageList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToAlias(t *testing.T) {
	result, err := client.SendToAlias(context.TODO(), msg1, alias1)
	if err != nil {
		t.Errorf("TestMiPush_SendToAlias failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToAliasList(t *testing.T) {
	result, err := client.SendToAliasList(context.TODO(), msg1, []string{alias1, alias2})
	if err != nil {
		t.Errorf("TestMiPush_SendToAliasList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToUserAccount(t *testing.T) {
	result, err := client.SendToUserAccount(context.TODO(), msg1, account1)
	if err != nil {
		t.Errorf("TestMiPush_SendToUserAccount failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SendToUserAccountList(t *testing.T) {
	result, err := client.SendToUserAccountList(context.TODO(), msg1, []string{account1, account2})
	if err != nil {
		t.Errorf("TestMiPush_SendToUserAccountList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_Broadcast(t *testing.T) {
	result, err := client.Broadcast(context.TODO(), msg1, topic1)
	if err != nil {
		t.Errorf("TestMiPush_Broadcast failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_BroadcastAll(t *testing.T) {
	result, err := client.BroadcastAll(context.TODO(), msg1)
	if err != nil {
		t.Errorf("TestMiPush_BroadcastAll failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_MultiTopicBroadcast(t *testing.T) {
	result, err := client.MultiTopicBroadcast(context.TODO(), msg1, []string{topic1, topic2}, INTERSECTION)
	if err != nil {
		t.Errorf("TestMiPush_MultiTopicBroadcast failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_CheckScheduleJobExist(t *testing.T) {
	result, err := client.CheckScheduleJobExist(context.TODO(), "slm30b80474526081454i5")
	if err != nil {
		t.Errorf("TestMiPush_CheckScheduleJobExist failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_DeleteScheduleJob(t *testing.T) {
	result, err := client.DeleteScheduleJob(context.TODO(), "slm30b80474526081454i5")
	if err != nil {
		t.Errorf("TestMiPush_DeleteScheduleJob failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_DeleteScheduleJobByJobKey(t *testing.T) {
	result, err := client.DeleteScheduleJobByJobKey(context.TODO(), "Xcm45b21474513716292EL")
	if err != nil {
		t.Errorf("TestMiPush_DeleteScheduleJobByJobKey failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_Stats(t *testing.T) {
	result, err := client.Stats(context.TODO(), "20160922", "20160922", packageName)
	if err != nil {
		t.Errorf("TestMiPush_Stats failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_GetMessageStatusByMsgID(t *testing.T) {
	result, err := client.GetMessageStatusByMsgID(context.TODO(), "scm23b964745244861922w")
	if err != nil {
		t.Errorf("TestMiPush_GetMessageStatusByMsgID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_GetMessageStatusByJobKey(t *testing.T) {
	result, err := client.GetMessageStatusByJobKey(context.TODO(), "key111")
	if err != nil {
		t.Errorf("TestMiPush_GetMessageStatusByJobKey failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_GetMessageStatusPeriod(t *testing.T) {
	result, err := client.GetMessageStatusPeriod(context.TODO(), time.Now().Add(-time.Hour*24).Unix()*1000, time.Now().Unix()*1000)
	if err != nil {
		t.Errorf("TestMiPush_GetMessageStatusPeriod failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

//----------------------------------------Subscription----------------------------------------//

func TestMiPush_SubscribeTopicForRegID(t *testing.T) {
	result, err := client.SubscribeTopicForRegID(context.TODO(), regID1, "topic3", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicForRegID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SubscribeTopicForRegIDList(t *testing.T) {
	result, err := client.SubscribeTopicForRegIDList(context.TODO(), []string{regID1, regID2}, "topic5", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicForRegIDList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_UnSubscribeTopicForRegID(t *testing.T) {
	result, err := client.UnSubscribeTopicForRegID(context.TODO(), regID1, "topic3", "")
	if err != nil {
		t.Errorf("TestMiPush_UnSubscribeTopicForRegID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_UnSubscribeTopicForRegIDList(t *testing.T) {
	result, err := client.UnSubscribeTopicForRegIDList(context.TODO(), []string{regID1, regID2}, "topic5", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicForRegIDList failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_SubscribeTopicByAlias(t *testing.T) {
	result, err := client.SubscribeTopicByAlias(context.TODO(), []string{alias1, alias2}, "topic5", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicByAlias failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_UnSubscribeTopicByAlias(t *testing.T) {
	result, err := client.UnSubscribeTopicByAlias(context.TODO(), []string{alias1, alias2}, "topic5", "")
	if err != nil {
		t.Errorf("TestMiPush_SubscribeTopicByAlias failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

//----------------------------------------Feedback----------------------------------------//

func TestMiPush_GetInvalidRegIDs(t *testing.T) {
	result, err := client.GetInvalidRegIDs(context.TODO())
	if err != nil {
		t.Errorf("TestMiPush_GetInvalidRegIDs failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

//----------------------------------------DevTools----------------------------------------//

func TestMiPush_GetAliasesOfRegID(t *testing.T) {
	result, err := client.GetAliasesOfRegID(context.TODO(), regID1)
	if err != nil {
		t.Errorf("TestMiPush_GetAliasesOfRegID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}

func TestMiPush_GetTopicsOfRegID(t *testing.T) {
	result, err := client.GetTopicsOfRegID(context.TODO(), regID2)
	if err != nil {
		t.Errorf("TestMiPush_GetTopicsOfRegID failed :%v\n", err)
	}
	t.Logf("result=%#v\n", result)
}
