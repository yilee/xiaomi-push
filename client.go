package xiaomipush

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type MiPush struct {
	packageName []string
	host        string
	appSecret   string
}

func NewClient(appSecret string, packageName []string) *MiPush {
	return &MiPush{
		packageName: packageName,
		host:        ProductionHost,
		appSecret:   appSecret,
	}
}

//----------------------------------------Sender----------------------------------------//
// 根据registrationId，发送消息到指定设备上
func (m *MiPush) Send(ctx context.Context, msg *Message, regID string) (*SendResult, error) {
	params := m.assembleSendParams(msg, regID)
	bytes, err := m.doPost(ctx, m.host+RegURL, params)
	if err != nil {
		return nil, err
	}
	var result SendResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 根据regIds，发送消息到指定的一组设备上
// regIds的个数不得超过1000个。
func (m *MiPush) SendToList(ctx context.Context, msg *Message, regIDList []string) (*SendResult, error) {
	if len(regIDList) == 0 || len(regIDList) > 1000 {
		panic("wrong number regIDList")
	}
	return m.Send(ctx, msg, strings.Join(regIDList, ","))
}

// 发送一组消息。其中TargetedMessage类中封装了Message对象和该Message所要发送的目标。注意：messages内所有TargetedMessage对象的targetType必须相同，
// 不支持在一个调用中同时给regid和alias发送消息。
// 如果是定时消息, 所有消息的time_to_send必须相同
// 消息必须设置packagename, 见client_test TestMiPush_SendTargetMessageList
func (m *MiPush) SendTargetMessageList(ctx context.Context, msgList []*TargetedMessage) (*SendResult, error) {
	if len(msgList) == 0 {
		return nil, errors.New("empty msg")
	}
	if len(msgList) == 1 {
		return m.Send(ctx, msgList[0].message, msgList[0].target)
	}
	params := m.assembleTargetMessageListParams(msgList)
	var bytes []byte
	var err error
	if msgList[0].targetType == TargetTypeRegID {
		bytes, err = m.doPost(ctx, m.host+MultiMessagesRegIDURL, params)
	} else if msgList[0].targetType == TargetTypeReAlias {
		bytes, err = m.doPost(ctx, m.host+MultiMessagesAliasURL, params)
	} else if msgList[0].targetType == TargetTypeAccount {
		bytes, err = m.doPost(ctx, m.host+MultiMessagesUserAccountURL, params)
	} else {
		panic("bad targetType")
	}

	if err != nil {
		return nil, err
	}
	var result SendResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 根据alias，发送消息到指定设备上
func (m *MiPush) SendToAlias(ctx context.Context, msg *Message, alias string) (*SendResult, error) {
	params := m.assembleSendToAlisaParams(msg, alias)
	bytes, err := m.doPost(ctx, m.host+MessageAlisaURL, params)
	if err != nil {
		return nil, err
	}
	var result SendResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 根据aliasList，发送消息到指定的一组设备上
// 元素的个数不得超过1000个。
func (m *MiPush) SendToAliasList(ctx context.Context, msg *Message, aliasList []string) (*SendResult, error) {
	if len(aliasList) == 0 || len(aliasList) > 1000 {
		panic("wrong number aliasList")
	}
	return m.SendToAlias(ctx, msg, strings.Join(aliasList, ","))
}

// 根据account，发送消息到指定account上
func (m *MiPush) SendToUserAccount(ctx context.Context, msg *Message, userAccount string) (*SendResult, error) {
	params := m.assembleSendToUserAccountParams(msg, userAccount)
	bytes, err := m.doPost(ctx, m.host+MessageUserAccountURL, params)
	if err != nil {
		return nil, err
	}
	var result SendResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 根据accountList，发送消息到指定的一组设备上
// 元素的个数不得超过1000个。
func (m *MiPush) SendToUserAccountList(ctx context.Context, msg *Message, accountList []string) (*SendResult, error) {
	if len(accountList) == 0 || len(accountList) > 1000 {
		panic("wrong number accountList")
	}
	return m.SendToUserAccount(ctx, msg, strings.Join(accountList, ","))
}

// 根据topic，发送消息到指定一组设备上
func (m *MiPush) Broadcast(ctx context.Context, msg *Message, topic string) (*SendResult, error) {
	params := m.assembleBroadcastParams(msg, topic)
	var bytes []byte
	var err error
	if len(m.packageName) > 1 {
		bytes, err = m.doPost(ctx, m.host+MultiPackageNameMessageMultiTopicURL, params)
	} else {
		bytes, err = m.doPost(ctx, m.host+MessageMultiTopicURL, params)
	}
	if err != nil {
		return nil, err
	}
	var result SendResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 向所有设备发送消息
func (m *MiPush) BroadcastAll(ctx context.Context, msg *Message) (*SendResult, error) {
	params := m.assembleBroadcastAllParams(msg)
	var bytes []byte
	var err error
	if len(m.packageName) > 1 {
		bytes, err = m.doPost(ctx, m.host+MultiPackageNameMessageAllURL, params)
	} else {
		bytes, err = m.doPost(ctx, m.host+MessageAllURL, params)
	}
	if err != nil {
		return nil, err
	}
	var result SendResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type TopicOP string

const (
	UNION        TopicOP = "UNION"        // 并集
	INTERSECTION TopicOP = "INTERSECTION" // 交集
	EXCEPT       TopicOP = "EXCEPT"       // 差集
)

// 向多个topic广播消息，支持topic间的交集、并集或差集（如果只有一个topic请用单topic版本）
// TOPIC_OP是一个枚举类型，指定了发送广播消息时多个topic之间的运算关系。
// 例如：topics的列表元素是[A, B, C, D]，则并集结果是A∪B∪C∪D，交集的结果是A∩B∩C∩D，差集的结果是A-B-C-D
func (m *MiPush) MultiTopicBroadcast(ctx context.Context, msg *Message, topics []string, topicOP TopicOP) (*SendResult, error) {
	if len(topics) > 5 || len(topics) == 0 {
		panic("topics size invalid")
	}
	if len(topics) == 1 {
		return m.Broadcast(ctx, msg, topics[0])
	}
	params := m.assembleMultiTopicBroadcastParams(msg, topics, topicOP)
	bytes, err := m.doPost(ctx, m.host+MultiTopicURL, params)
	if err != nil {
		return nil, err
	}
	var result SendResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 检测定时消息的任务是否存在。
// result.code = 0 为任务存在, 否则不存在
func (m *MiPush) CheckScheduleJobExist(ctx context.Context, msgID string) (*Result, error) {
	params := m.assembleCheckScheduleJobParams(msgID)
	bytes, err := m.doPost(ctx, m.host+ScheduleJobExistURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 删除指定的定时消息
func (m *MiPush) DeleteScheduleJob(ctx context.Context, msgID string) (*Result, error) {
	params := m.assembleDeleteScheduleJobParams(msgID)
	bytes, err := m.doPost(ctx, m.host+ScheduleJobDeleteURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 删除指定的定时消息
func (m *MiPush) DeleteScheduleJobByJobKey(ctx context.Context, jobKey string) (*Result, error) {
	params := m.assembleDeleteScheduleJobByJobKeyParams(jobKey)
	bytes, err := m.doPost(ctx, m.host+ScheduleJobDeleteByJobKeyURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//----------------------------------------Stats----------------------------------------//
// 获取指定日期范围内的日统计数据（如果日期范围包含今日，则今日数据为从今天00：00开始到现在的累计量)。
// packageName:
// Android设备，传入App的包名
// IOS设备，传入App的Bundle Id
func (m *MiPush) Stats(ctx context.Context, start, end, packageName string) (*StatsResult, error) {
	params := m.assembleStatsParams(start, end, packageName)
	bytes, err := m.doGet(ctx, m.host+StatsURL, params)
	if err != nil {
		return nil, err
	}
	var result StatsResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//----------------------------------------Tracer----------------------------------------//
// 获取指定ID的消息状态
func (m *MiPush) GetMessageStatusByMsgID(ctx context.Context, msgID string) (*SingleStatusResult, error) {
	params := m.assembleStatusParams(msgID)
	bytes, err := m.doGet(ctx, m.host+MessageStatusURL, params)
	if err != nil {
		return nil, err
	}
	var result SingleStatusResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 获取某个时间间隔内所有消息的状态。
func (m *MiPush) GetMessageStatusByJobKey(ctx context.Context, jobKey string) (*BatchStatusResult, error) {
	params := m.assembleStatusByJobKeyParams(jobKey)
	bytes, err := m.doGet(ctx, m.host+MessagesStatusURL, params)
	if err != nil {
		return nil, err
	}
	var result BatchStatusResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 获取某个时间间隔内所有消息的状态。
func (m *MiPush) GetMessageStatusPeriod(ctx context.Context, beginTime, endTime int64) (*BatchStatusResult, error) {
	params := m.assembleStatusPeriodParams(beginTime, endTime)
	bytes, err := m.doGet(ctx, m.host+MessagesStatusURL, params)
	if err != nil {
		return nil, err
	}
	var result BatchStatusResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//----------------------------------------Subscription----------------------------------------//

// 给某个regid订阅标签
func (m *MiPush) SubscribeTopicForRegID(ctx context.Context, regID, topic, category string) (*Result, error) {
	params := m.assembleSubscribeTopicForRegIDParams(regID, topic, category)
	bytes, err := m.doPost(ctx, m.host+TopicSubscribeURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 给一组regid列表订阅标签
func (m *MiPush) SubscribeTopicForRegIDList(ctx context.Context, regIDList []string, topic, category string) (*Result, error) {
	return m.SubscribeTopicForRegID(ctx, strings.Join(regIDList, ","), topic, category)
}

// 取消某个regid的标签。
func (m *MiPush) UnSubscribeTopicForRegID(ctx context.Context, regID, topic, category string) (*Result, error) {
	params := m.assembleUnSubscribeTopicForRegIDParams(regID, topic, category)
	bytes, err := m.doPost(ctx, m.host+TopicUnSubscribeURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 取消一组regid列表的标签
func (m *MiPush) UnSubscribeTopicForRegIDList(ctx context.Context, regIDList []string, topic, category string) (*Result, error) {
	return m.UnSubscribeTopicForRegID(ctx, strings.Join(regIDList, ","), topic, category)
}

// 给一组alias列表订阅标签
func (m *MiPush) SubscribeTopicByAlias(ctx context.Context, aliases []string, topic, category string) (*Result, error) {
	params := m.assembleSubscribeTopicByAliasParams(aliases, topic, category)
	bytes, err := m.doPost(ctx, m.host+TopicSubscribeByAliasURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 取消一组alias列表的标签
func (m *MiPush) UnSubscribeTopicByAlias(ctx context.Context, aliases []string, topic, category string) (*Result, error) {
	params := m.assembleUnSubscribeTopicByAliasParams(aliases, topic, category)
	bytes, err := m.doPost(ctx, m.host+TopicUnSubscribeByAliasURL, params)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//----------------------------------------Feedback----------------------------------------//

// 获取失效的regId列表
// 获取失效的regId列表，每次请求最多返回1000个regId。
// 每次请求之后，成功返回的失效的regId将会从MiPush数据库删除。
func (m *MiPush) GetInvalidRegIDs(ctx context.Context) (*InvalidRegIDsResult, error) {
	params := m.assembleGetInvalidRegIDsParams()
	bytes, err := m.doGet(ctx, InvalidRegIDsURL, params)
	if err != nil {
		return nil, err
	}
	var result InvalidRegIDsResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//----------------------------------------DevTools----------------------------------------//

// 获取一个应用的某个用户目前设置的所有Alias
func (m *MiPush) GetAliasesOfRegID(ctx context.Context, regID string) (*AliasesOfRegIDResult, error) {
	params := m.assembleGetAliasesOfParams(regID)
	bytes, err := m.doGet(ctx, m.host+AliasAllURL, params)
	if err != nil {
		return nil, err
	}
	var result AliasesOfRegIDResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 	获取一个应用的某个用户的目前订阅的所有Topic
func (m *MiPush) GetTopicsOfRegID(ctx context.Context, regID string) (*TopicsOfRegIDResult, error) {
	params := m.assembleGetTopicsOfParams(regID)
	bytes, err := m.doGet(ctx, m.host+TopicsAllURL, params)
	if err != nil {
		return nil, err
	}
	var result TopicsOfRegIDResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) assembleSendParams(msg *Message, regID string) url.Values {
	form := m.defaultForm(msg)
	form.Add("registration_id", regID)
	return form
}

func (m *MiPush) assembleTargetMessageListParams(msgList []*TargetedMessage) url.Values {
	form := url.Values{}
	type OneMsg struct {
		Target  string   `json:"target"`
		Message *Message `json:"message"`
	}
	var messages []*OneMsg

	for _, m := range msgList {
		messages = append(messages, &OneMsg{
			Target:  m.target,
			Message: m.message,
		})
	}
	bytes, err := json.Marshal(messages)
	if err != nil {
		panic(err)
	}
	form.Add("messages", string(bytes))
	form.Add("time_to_send", strconv.FormatInt(msgList[0].message.TimeToSend, 10))
	return form
}

func (m *MiPush) assembleSendToAlisaParams(msg *Message, alias string) url.Values {
	form := m.defaultForm(msg)
	form.Add("alias", alias)
	return form
}

func (m *MiPush) assembleSendToUserAccountParams(msg *Message, userAccount string) url.Values {
	form := m.defaultForm(msg)
	form.Add("user_account", userAccount)
	return form
}

func (m *MiPush) assembleBroadcastParams(msg *Message, topic string) url.Values {
	form := m.defaultForm(msg)
	form.Add("topic", topic)
	return form
}

func (m *MiPush) assembleBroadcastAllParams(msg *Message) url.Values {
	form := m.defaultForm(msg)
	return form
}

func (m *MiPush) assembleMultiTopicBroadcastParams(msg *Message, topics []string, topicOP TopicOP) url.Values {
	form := m.defaultForm(msg)
	form.Add("topic_op", string(topicOP))
	form.Add("topics", strings.Join(topics, ";$;"))
	return form
}

func (m *MiPush) assembleCheckScheduleJobParams(msgID string) url.Values {
	form := url.Values{}
	form.Add("job_id", msgID)
	return form
}

func (m *MiPush) assembleDeleteScheduleJobParams(msgID string) url.Values {
	form := url.Values{}
	form.Add("job_id", msgID)
	return form
}

func (m *MiPush) assembleDeleteScheduleJobByJobKeyParams(jobKey string) url.Values {
	form := url.Values{}
	form.Add("job_key", jobKey)
	return form
}

func (m *MiPush) assembleStatsParams(start, end, packageName string) string {
	form := url.Values{}
	form.Add("start_date", start)
	form.Add("end_date", end)
	form.Add("restricted_package_name", packageName)
	return "?" + form.Encode()
}

func (m *MiPush) assembleStatusParams(msgID string) string {
	form := url.Values{}
	form.Add("msg_id", msgID)
	return "?" + form.Encode()
}

func (m *MiPush) assembleStatusByJobKeyParams(jobKey string) string {
	form := url.Values{}
	form.Add("job_key", jobKey)
	return "?" + form.Encode()
}

func (m *MiPush) assembleStatusPeriodParams(beginTime, endTime int64) string {
	form := url.Values{}
	form.Add("begin_time", strconv.FormatInt(int64(beginTime), 10))
	form.Add("end_time", strconv.FormatInt(int64(endTime), 10))
	return "?" + form.Encode()
}

func (m *MiPush) assembleSubscribeTopicForRegIDParams(regID, topic, category string) url.Values {
	form := url.Values{}
	form.Add("registration_id", regID)
	form.Add("topic", topic)
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	if category != "" {
		form.Add("category", category)
	}
	return form
}

func (m *MiPush) assembleUnSubscribeTopicForRegIDParams(regID, topic, category string) url.Values {
	form := url.Values{}
	form.Add("registration_id", regID)
	form.Add("topic", topic)
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	if category != "" {
		form.Add("category", category)
	}
	return form
}

func (m *MiPush) assembleSubscribeTopicByAliasParams(aliases []string, topic, category string) url.Values {
	form := url.Values{}
	form.Add("aliases", strings.Join(aliases, ","))
	form.Add("topic", topic)
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	if category != "" {
		form.Add("category", category)
	}
	return form
}

func (m *MiPush) assembleUnSubscribeTopicByAliasParams(aliases []string, topic, category string) url.Values {
	form := url.Values{}
	form.Add("aliases", strings.Join(aliases, ","))
	form.Add("topic", topic)
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	if category != "" {
		form.Add("category", category)
	}
	return form
}

func (m *MiPush) assembleGetInvalidRegIDsParams() string {
	form := url.Values{}
	return "?" + form.Encode()
}

func (m *MiPush) assembleGetAliasesOfParams(regID string) string {
	form := url.Values{}
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	form.Add("registration_id", regID)
	return "?" + form.Encode()
}

func (m *MiPush) assembleGetTopicsOfParams(regID string) string {
	form := url.Values{}
	form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	form.Add("registration_id", regID)
	return "?" + form.Encode()
}

func (m *MiPush) doPost(ctx context.Context, url string, form url.Values) ([]byte, error) {
	var result []byte
	var req *http.Request
	var res *http.Response
	var err error
	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+m.appSecret)
	client := &http.Client{}
	tryTime := 0
tryAgain:
	res, err = ctxhttp.Do(ctx, client, req)
	if err != nil {
		fmt.Println("xiaomi push post err:", err, tryTime)
		select {
		case <-ctx.Done():
			return nil, err
		default:
		}
		tryTime += 1
		if tryTime < PostRetryTimes {
			goto tryAgain
		}
		return nil, err
	}
	if res.Body == nil {
		panic("xiaomi response is nil")
	}
	defer res.Body.Close()
	fmt.Println("res.StatusCode=", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("network error")
	}
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiPush) doGet(ctx context.Context, url string, params string) ([]byte, error) {
	var result []byte
	var req *http.Request
	var res *http.Response
	var err error
	req, err = http.NewRequest("GET", url+params, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+m.appSecret)

	client := &http.Client{}
	res, err = ctxhttp.Do(ctx, client, req)
	if res.Body == nil {
		panic("xiaomi response is nil")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("network error")
	}
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiPush) defaultForm(msg *Message) url.Values {
	form := url.Values{}
	if len(m.packageName) > 0 {
		form.Add("restricted_package_name", strings.Join(m.packageName, ","))
	}
	if msg.TimeToLive > 0 {
		form.Add("time_to_live", strconv.FormatInt(msg.TimeToLive, 10))
	}
	if len(msg.Payload) > 0 {
		form.Add("payload", msg.Payload)
	}
	if len(msg.Title) > 0 {
		form.Add("title", msg.Title)
	}
	if len(msg.Description) > 0 {
		form.Add("description", msg.Description)
	}
	form.Add("notify_type", strconv.FormatInt(int64(msg.NotifyType), 10))
	form.Add("pass_through", strconv.FormatInt(int64(msg.PassThrough), 10))
	if msg.NotifyID != 0 {
		form.Add("notify_id", strconv.FormatInt(int64(msg.NotifyID), 10))
	}
	if msg.TimeToSend > 0 {
		form.Add("time_to_send", strconv.FormatInt(int64(msg.TimeToSend), 10))
	}
	if len(msg.Extra) > 0 {
		for k, v := range msg.Extra {
			form.Add("extra."+k, v)
		}
	}
	return form
}
