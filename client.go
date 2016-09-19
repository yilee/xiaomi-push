package xiaomipush

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
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

func (m *MiPush) Push(msg *Message, regIDList []string) (*Result, error) {
	params := m.assemblePushParams(msg, regIDList)
	return m.doPost(m.host+RegURL, params)
}

func (m *MiPush) Stats(start, end string) (*Result, error) {
	params := m.assembleStatsParams(start, end)
	return m.doGet(m.host+StatsURL, params)
}

func (m *MiPush) Status(msgID string) (*Result, error) {
	params := m.assembleStatusParams(msgID)
	return m.doGet(m.host+MessageStatusURL, params)
}

func (m *MiPush) assemblePushParams(msg *Message, regIDList []string) map[string]string {
	params := make(map[string]string)
	params["registration_id"] = strings.Join(regIDList, ",")
	params["restricted_package_name"] = m.packageName

	if msg.timeToLive > 0 {
		params["time_to_live"] = strconv.FormatInt(msg.timeToLive, 10)
	}
	if len(msg.payload) > 0 {
		params["payload"] = msg.payload
	}
	if len(msg.title) > 0 {
		params["title"] = msg.title
	}
	if len(msg.description) > 0 {
		params["description"] = msg.description
	}
	params["notify_type"] = strconv.FormatInt(int64(msg.notifyType), 10)
	params["pass_through"] = strconv.FormatInt(int64(msg.passThrough), 10)
	if msg.notifyID > 0 {
		params["notify_id"] = strconv.FormatInt(msg.notifyID, 10)
	}
	if msg.timeToSend > 0 {
		params["time_to_send"] = strconv.FormatInt(msg.timeToSend, 10)
	}
	if msg.extra != nil && len(msg.extra) > 0 {
		for k, v := range msg.extra {
			params["extra."+k] = v
		}
	}
	return params
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

func (m *MiPush) handlePushResponse(response *http.Response) (*Result, error) {
	defer func() {
		_ = response.Body.Close()
	}()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return getResultFromJSON(data)
}

func (m *MiPush) doPost(url string, params map[string]string) (*Result, error) {
	var result *Result
	var req *http.Request
	var resp *http.Response
	var err error
	data, _ := json.Marshal(params)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("X-PUSH-OS", runtime.GOOS)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorizatione", "key="+m.appSecret)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	result, err = m.handlePushResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiPush) doGet(url string, params string) (*Result, error) {
	var result *Result
	var req *http.Request
	var resp *http.Response
	var err error
	req, err = http.NewRequest("GET", url+params, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorizatione", "key="+m.appSecret)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	result, err = m.handlePushResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}
