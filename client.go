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
)

type MiPush struct {
	packageName string
	host        string
	appSecret   string
}

func NewClient(appSecret, packageName string) *MiPush {
	return &MiPush{
		packageName: packageName,
		host:        ProductionHost,
		appSecret:   appSecret,
	}
}

// 根据registrationId，发送消息到指定设备上
func (m *MiPush) Send(msg *Message, regID string) (*Result, error) {
	params := m.assembleSendParams(msg, regID)
	bytes, err := m.doPost(m.host+RegURL, params)
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

// 根据regIds，发送消息到指定的一组设备上
// regIds的个数不得超过1000个。
func (m *MiPush) SendToList(msg *Message, regIDList []string) (*Result, error) {
	return m.Send(msg, strings.Join(regIDList, ","))
}

// 发送一组消息。其中TargetedMessage类中封装了Message对象和该Message所要发送的目标。注意：messages内所有TargetedMessage对象的targetType必须相同，
// 不支持在一个调用中同时给regid和alias发送消息。
func (m *MiPush) SendTargetMessageList(msgList []*TargetedMessage) (*Result, error) {
	if len(msgList) == 0 {
		return nil, errors.New("empty msg")
	}
	if len(msgList) == 1 {
		return m.Send(msgList[0].message, msgList[0].target)
	}
	params := m.assembleTargetMessageListParams(msgList)
	var bytes []byte
	var err error
	if msgList[0].targetType == TargetTypeRegID {
		bytes, err = m.doPost2(m.host+MultiMessagesRegIDURL, params)
	} else if msgList[0].targetType == TargetTypeReAlias {
		bytes, err = m.doPost2(m.host+MultiMessagesAliasURL, params)
	} else if msgList[0].targetType == TargetTypeAccount {
		bytes, err = m.doPost2(m.host+MultiMessagesUserAccountURL, params)
	} else {
		panic("bad targetType")
	}

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

// 根据alias，发送消息到指定设备上
func (m *MiPush) SendToAlias(msg *Message, alias string) (*Result, error) {
	params := m.assembleSendToAlisaParams(msg, alias)
	bytes, err := m.doPost(m.host+MessageAlisaURL, params)
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

// 根据aliasList，发送消息到指定的一组设备上
// 元素的个数不得超过1000个。
func (m *MiPush) SendToAliasList(msg *Message, aliasList []string) (*Result, error) {
	return m.SendToAlias(msg, strings.Join(aliasList, ","))
}

// 根据account，发送消息到指定account上
func (m *MiPush) SendToUserAccount(msg *Message, userAccount string) (*Result, error) {
	params := m.assembleSendToUserAccountParams(msg, userAccount)
	bytes, err := m.doPost(m.host+MessageUserAccountURL, params)
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

// 根据accountList，发送消息到指定的一组设备上
// 元素的个数不得超过1000个。
func (m *MiPush) SendToUserAccountList(msg *Message, accountList []string) (*Result, error) {
	return m.SendToUserAccount(msg, strings.Join(accountList, ","))
}

// 根据topic，发送消息到指定一组设备上
func (m *MiPush) Broadcast(msg *Message, topic string) (*Result, error) {
	params := m.assembleBroadcastParams(msg, topic)
	var bytes []byte
	var err error
	if msg.IsMultiPackageName() {
		bytes, err = m.doPost(m.host+MultiPackageNameMessageMultiTopicURL, params)
	} else {
		bytes, err = m.doPost(m.host+MessageMultiTopicURL, params)
	}
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

// 向所有设备发送消息
func (m *MiPush) BroadcastAll(msg *Message) (*Result, error) {
	params := m.assembleBroadcastAllParams(msg)
	var bytes []byte
	var err error
	if msg.IsMultiPackageName() {
		bytes, err = m.doPost(m.host+MultiPackageNameMessageAllURL, params)
	} else {
		bytes, err = m.doPost(m.host+MessageAllURL, params)
	}
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

type TopicOP string

const (
	UNION        TopicOP = "UNION"        // 并集
	INTERSECTION TopicOP = "INTERSECTION" // 交集
	EXCEPT       TopicOP = "EXCEPT"       // 差集
)

// 向多个topic广播消息，支持topic间的交集、并集或差集（如果只有一个topic请用单topic版本）
// TOPIC_OP是一个枚举类型，指定了发送广播消息时多个topic之间的运算关系。
// 例如：topics的列表元素是[A, B, C, D]，则并集结果是A∪B∪C∪D，交集的结果是A∩B∩C∩D，差集的结果是A-B-C-D
func (m *MiPush) MultiTopicBroadcast(msg *Message, topics []string, topicOP TopicOP) (*Result, error) {
	if len(topics) > 5 || len(topics) == 0 {
		panic("topics size invalid")
	}
	if len(topics) == 1 {
		return m.Broadcast(msg, topics[0])
	}
	params := m.assembleMultiTopicBroadcastParams(msg, topics, topicOP)
	bytes, err := m.doPost(m.host+MultiTopicURL, params)
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

func (m *MiPush) Stats(start, end string) (*StatsResult, error) {
	params := m.assembleStatsParams(start, end)
	bytes, err := m.doGet(m.host+StatsURL, params)
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

func (m *MiPush) Status(msgID string) (*StatusResult, error) {
	params := m.assembleStatusParams(msgID)
	bytes, err := m.doGet(m.host+MessageStatusURL, params)
	if err != nil {
		return nil, err
	}
	var result StatusResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MiPush) assembleMultiTopicBroadcastParams(msg *Message, topics []string, topicOP TopicOP) url.Values {
	form := m.defaultForm(msg)
	form.Add("topic_op", string(topicOP))
	form.Add("topics", strings.Join(topics, ";$;"))
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

func (m *MiPush) assembleSendParams(msg *Message, regID string) url.Values {
	form := m.defaultForm(msg)
	form.Add("registration_id", regID)
	return form
}

func (m *MiPush) assembleTargetMessageListParams(msgList []*TargetedMessage) string {
	form := url.Values{}
	type OneMsg struct {
		Target  string `json:"target"`
		Message string `json:"message"`
	}
	messages := struct {
		Messages []*OneMsg `json:"messages"`
	}{}

	for _, m := range msgList {
		messages.Messages = append(messages.Messages, &OneMsg{
			Target:  m.target,
			Message: string(m.message.JSON()),
		})
	}
	bytes, err := json.Marshal(messages)
	if err != nil {
		panic(err)
	}
	fmt.Println("bytes", string(bytes))
	form.Add("messages", string(bytes))
	return string(bytes)
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

func (m *MiPush) assembleStatsParams(start, end string) string {
	form := url.Values{}
	form.Add("start_date", start)
	form.Add("end_date", end)
	return "?" + form.Encode()
}

func (m *MiPush) assembleStatusParams(msgID string) string {
	form := url.Values{}
	form.Add("msg_id", msgID)
	return "?" + form.Encode()
}

func (m *MiPush) handleResponse(response *http.Response) ([]byte, error) {
	defer func() {
		_ = response.Body.Close()
	}()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *MiPush) doPost(url string, form url.Values) ([]byte, error) {
	var result []byte
	var req *http.Request
	var resp *http.Response
	var err error
	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+m.appSecret)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	result, err = m.handleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiPush) doPost2(url string, jsonStr string) ([]byte, error) {
	var result []byte
	var req *http.Request
	var resp *http.Response
	var err error
	fmt.Println("jsonStr", jsonStr)
	req, err = http.NewRequest("POST", url, strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+m.appSecret)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	result, err = m.handleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiPush) doGet(url string, params string) ([]byte, error) {
	var result []byte
	var req *http.Request
	var resp *http.Response
	var err error
	req, err = http.NewRequest("GET", url+params, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+m.appSecret)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	result, err = m.handleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiPush) defaultForm(msg *Message) url.Values {
	form := url.Values{}
	if len(msg.RestrictedPackageNames) > 0 {
		form.Add("restricted_package_name", strings.Join(msg.RestrictedPackageNames, ","))
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
	if msg.NotifyID > 0 {
		form.Add("notify_id", strconv.FormatInt(int64(msg.NotifyID), 10))
	}
	if msg.TimeToSend > 0 {
		form.Add("time_to_send", strconv.FormatInt(int64(msg.TimeToSend), 10))
	}
	if msg.Extra != nil && len(msg.Extra) > 0 {
		for k, v := range msg.Extra {
			form.Add("extra."+k, v)
		}
	}
	return form
}
