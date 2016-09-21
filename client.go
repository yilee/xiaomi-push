package xiaomipush

import (
	"encoding/json"
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

func (m *MiPush) Broadcast(msg *Message, topic string) (*Result, error) {
	params := m.assembleBroadcastParams(msg, topic)
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

type TopicOP string

const (
	UNION        TopicOP = "UNION"
	INTERSECTION TopicOP = "INTERSECTION"
	EXCEPT       TopicOP = "EXCEPT"
)

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

func (m *MiPush) BroadcastAll(msg *Message) (*Result, error) {
	params := m.assembleBroadcastAllParams(msg)
	bytes, err := m.doPost(m.host+MessageAllURL, params)
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
	form := url.Values{}
	form.Add("topic_op", string(topicOP))
	form.Add("topics", strings.Join(topics, ";$;"))
	form.Add("restricted_package_name", m.packageName)
	if msg.timeToLive > 0 {
		form.Add("time_to_live", strconv.FormatInt(msg.timeToLive, 10))
	}
	if len(msg.payload) > 0 {
		form.Add("payload", msg.payload)
	}
	if len(msg.title) > 0 {
		form.Add("title", msg.title)
	}
	if len(msg.description) > 0 {
		form.Add("description", msg.description)
	}
	form.Add("notify_type", strconv.FormatInt(int64(msg.notifyType), 10))
	form.Add("pass_through", strconv.FormatInt(int64(msg.passThrough), 10))
	if msg.notifyID > 0 {
		form.Add("notify_id", strconv.FormatInt(int64(msg.notifyID), 10))
	}
	if msg.timeToSend > 0 {
		form.Add("time_to_send", strconv.FormatInt(int64(msg.timeToSend), 10))
	}
	if msg.extra != nil && len(msg.extra) > 0 {
		for k, v := range msg.extra {
			form.Add("extra."+k, v)
		}
	}
	return form
}

func (m *MiPush) assembleBroadcastAllParams(msg *Message) url.Values {
	form := url.Values{}
	form.Add("restricted_package_name", m.packageName)
	if msg.timeToLive > 0 {
		form.Add("time_to_live", strconv.FormatInt(msg.timeToLive, 10))
	}
	if len(msg.payload) > 0 {
		form.Add("payload", msg.payload)
	}
	if len(msg.title) > 0 {
		form.Add("title", msg.title)
	}
	if len(msg.description) > 0 {
		form.Add("description", msg.description)
	}
	form.Add("notify_type", strconv.FormatInt(int64(msg.notifyType), 10))
	form.Add("pass_through", strconv.FormatInt(int64(msg.passThrough), 10))
	if msg.notifyID > 0 {
		form.Add("notify_id", strconv.FormatInt(int64(msg.notifyID), 10))
	}
	if msg.timeToSend > 0 {
		form.Add("time_to_send", strconv.FormatInt(int64(msg.timeToSend), 10))
	}
	if msg.extra != nil && len(msg.extra) > 0 {
		for k, v := range msg.extra {
			form.Add("extra."+k, v)
		}
	}
	return form
}

func (m *MiPush) assembleBroadcastParams(msg *Message, topic string) url.Values {
	form := url.Values{}
	form.Add("topic", topic)
	form.Add("restricted_package_name", m.packageName)
	if msg.timeToLive > 0 {
		form.Add("time_to_live", strconv.FormatInt(msg.timeToLive, 10))
	}
	if len(msg.payload) > 0 {
		form.Add("payload", msg.payload)
	}
	if len(msg.title) > 0 {
		form.Add("title", msg.title)
	}
	if len(msg.description) > 0 {
		form.Add("description", msg.description)
	}
	form.Add("notify_type", strconv.FormatInt(int64(msg.notifyType), 10))
	form.Add("pass_through", strconv.FormatInt(int64(msg.passThrough), 10))
	if msg.notifyID > 0 {
		form.Add("notify_id", strconv.FormatInt(int64(msg.notifyID), 10))
	}
	if msg.timeToSend > 0 {
		form.Add("time_to_send", strconv.FormatInt(int64(msg.timeToSend), 10))
	}
	if msg.extra != nil && len(msg.extra) > 0 {
		for k, v := range msg.extra {
			form.Add("extra."+k, v)
		}
	}
	return form
}

func (m *MiPush) assembleSendParams(msg *Message, regID string) url.Values {
	form := url.Values{}
	form.Add("registration_id", regID)
	form.Add("restricted_package_name", m.packageName)
	if msg.timeToLive > 0 {
		form.Add("time_to_live", strconv.FormatInt(msg.timeToLive, 10))
	}
	if len(msg.payload) > 0 {
		form.Add("payload", msg.payload)
	}
	if len(msg.title) > 0 {
		form.Add("title", msg.title)
	}
	if len(msg.description) > 0 {
		form.Add("description", msg.description)
	}
	form.Add("notify_type", strconv.FormatInt(int64(msg.notifyType), 10))
	form.Add("pass_through", strconv.FormatInt(int64(msg.passThrough), 10))
	if msg.notifyID > 0 {
		form.Add("notify_id", strconv.FormatInt(int64(msg.notifyID), 10))
	}
	if msg.timeToSend > 0 {
		form.Add("time_to_send", strconv.FormatInt(int64(msg.timeToSend), 10))
	}
	if msg.extra != nil && len(msg.extra) > 0 {
		for k, v := range msg.extra {
			form.Add("extra."+k, v)
		}
	}
	return form
}

func (m *MiPush) assembleStatsParams(start, end string) string {
	form := url.Values{}
	form.Add("start_date", start)
	form.Add("end_date", end)
	form.Add("restricted_package_name", m.packageName)
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
