# xiaomi-push
小米推送服务 Golang SDK

Production ready, full golang implementation of Xiaomi Push API (http://dev.xiaomi.com/console/?page=appservice&mod=push)

```Go
var client = xiaomipush.NewClient("yourappSecret", []string{"packageName"})

func main() {
    var msg1 *Message = xiaomipush.NewAndroidMessage("title", "body").SetPayload("this is payload1")
    client.Send(context.Background(), msg1, regID1)
}

```

### Sender APIs

- [x] Send(msg *Message, regID string)
- [x] SendToList(msg *Message, regIDList []string)
- [x] SendTargetMessageList(msgList []*TargetedMessage)
- [x] SendToAlias(msg *Message, alias string)
- [x] SendToAliasList(msg *Message, aliasList []string)
- [x] SendToUserAccount(msg *Message, userAccount string) 
- [x] SendToUserAccountList(msg *Message, accountList []string)
- [x] Broadcast(msg *Message, topic string)
- [x] BroadcastAll(msg *Message) (*SendResult, error)
- [x] MultiTopicBroadcast(msg *Message, topics []string, topicOP TopicOP)
- [x] CheckScheduleJobExist(msgID string)
- [x] DeleteScheduleJob(msgID string) (*Result, error)
- [x] DeleteScheduleJobByJobKey(jobKey string) (*Result, error) 

### Stats APIs

- [x] Stats(start, end, packageName string)
- [x] GetMessageStatusByMsgID(msgID string) (*SingleStatusResult, error)
- [x] GetMessageStatusByJobKey(jobKey string) (*BatchStatusResult, error) 
- [x] GetMessageStatusPeriod(beginTime, endTime int64) (*BatchStatusResult, error) 

### Subscription APIs

- [x] SubscribeTopicForRegID(regID, topic, category string) (*Result, error)
- [x] SubscribeTopicForRegIDList(regIDList []string, topic, category string) (*Result, error)
- [x] UnSubscribeTopicForRegID(regID, topic, category string) (*Result, error)
- [x] UnSubscribeTopicForRegIDList(regIDList []string, topic, category string) (*Result, error)
- [x] SubscribeTopicByAlias(aliases []string, topic, category string) (*Result, error)
- [x] UnSubscribeTopicByAlias(aliases []string, topic, category string) (*Result, error)

### Feedback APIs

- [x] GetInvalidRegIDs() (*InvalidRegIDsResult, error)

### DevTools APIs

- [x] GetAliasesOfRegID(regID string) (*AliasesOfRegIDResult, error)
- [x] GetTopicsOfRegID(regID string) (*TopicsOfRegIDResult, error)
