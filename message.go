package xiaomipush

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Message struct {
	UniqueID               string            `json:"unique_id"`                // 消息唯一ID
	Payload                string            `json:"payload"`                  // 消息内容payload
	Title                  string            `json:"title"`                    // 通知栏展示的通知的标题
	Description            string            `json:"description"`              // 通知栏展示的通知的描述
	PassThrough            int32             `json:"pass_through"`             // 是否通过透传的方式送给app，1表示透传消息，0表示通知栏消息。
	NotifyType             int32             `json:"notify_type"`              // DEFAULT_ALL = -1; DEFAULT_SOUND  = 1;   // 使用默认提示音提示 DEFAULT_VIBRATE = 2;   // 使用默认震动提示 DEFAULT_LIGHTS = 4;    // 使用默认led灯光提示
	RestrictedPackageNames []string          `json:"restricted_package_names"` // 设置app的包名packageName。packageName必须和开发者网站上申请的结果一致。
	TimeToLive             int64             `json:"time_to_live"`             // 可选项。如果用户离线，设置消息在服务器保存的时间，单位：ms。服务器默认最长保留两周。
	TimeToSend             int64             `json:"time_to_send"`             // 可选项。定时发送消息。timeToSend是以毫秒为单位的时间戳。注：仅支持七天内的定时消息。
	NotifyID               int64             `json:"notify_id"`                // 可选项。默认情况下，通知栏只显示一条推送消息。如果通知栏要显示多条推送消息，需要针对不同的消息设置不同的notify_id（相同notify_id的通知栏消息会覆盖之前的）。
	Extra                  map[string]string `json:"unique_id"`                // 可选项，对app提供一些扩展的功能，请参考2.2。除了这些扩展功能，开发者还可以定义一些key和value来控制客户端的行为。注：key和value的字符数不能超过1024，至多可以设置10个key-value键值对。
}

const (
	MaxTimeToSend = time.Hour * 24 * 7
	MaxTimeToLive = time.Hour * 24 * 7 * 2
)

func NewAndroidMessage(title, description string) *Message {
	return &Message{
		UniqueID:               "",
		Payload:                "",
		Title:                  title,
		Description:            description,
		PassThrough:            0,
		NotifyType:             -1, // default notify type
		RestrictedPackageNames: nil,
		TimeToLive:             0,
		TimeToSend:             0,
		NotifyID:               0,
		Extra:                  make(map[string]string),
	}
}

func (m *Message) SetUniqueID(uniqueID string) *Message {
	m.UniqueID = uniqueID
	return m
}

func (m *Message) SetPayload(payload string) *Message {
	m.Payload = payload
	return m
}

func (m *Message) SetPassThrough(passThrough int32) *Message {
	m.PassThrough = passThrough
	return m
}

func (m *Message) SetNotifyType(notifyType int32) *Message {
	m.NotifyType = notifyType
	return m
}

func (m *Message) SetRestrictedPackageNames(restrictedPackageNames []string) *Message {
	m.RestrictedPackageNames = restrictedPackageNames
	return m
}

func (m *Message) IsMultiPackageName() bool {
	return len(m.RestrictedPackageNames) > 1
}

func (m *Message) SetTimeToSend(tts int64) *Message {
	if time.Since(time.Unix(0, tts*int64(time.Millisecond))) > MaxTimeToSend {
		m.TimeToSend = time.Now().Add(MaxTimeToSend).UnixNano() / 1e6
	} else {
		m.TimeToSend = tts
	}
	return m
}

func (m *Message) SetTimeToLive(ttl int64) *Message {
	if time.Since(time.Unix(0, ttl*int64(time.Millisecond))) > MaxTimeToLive {
		m.TimeToLive = time.Now().Add(MaxTimeToLive).UnixNano() / 1e6
	} else {
		m.TimeToLive = ttl
	}
	return m
}

func (m *Message) SetNotifyID(notifyID int64) *Message {
	m.NotifyID = notifyID
	return m
}

func (m *Message) EnableFlowControl() *Message {
	m.Extra["flow_control"] = "1"
	return m
}

func (m *Message) DisableFlowControl() *Message {
	delete(m.Extra, "flow_control")
	return m
}

// 小米推送服务器每隔1s将已送达或已点击的消息ID和对应设备的regid或alias通过调用第三方http接口传给开发者。
func (m *Message) SetCallback(callbackURL string) *Message {
	m.Extra["callback"] = callbackURL
	m.Extra["callback.param"] = m.UniqueID
	m.Extra["callback.type"] = "3" // 1:送达回执, 2:点击回执, 3:送达和点击回执
	return m
}

func (m *Message) JSON() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println("m", string(bytes))
	return bytes
}

//-----------------------------------------------------------------------------------//
// TargetedMessage封装了MiPush推送服务系统中的消息Message对象，和该Message对象所要发送到的目标。

type TargetType int32

const (
	TargetTypeRegID   TargetType = 1
	TargetTypeReAlias TargetType = 2
	TargetTypeAccount TargetType = 3
)

type TargetedMessage struct {
	message    *Message
	targetType TargetType
	target     string
}

func NewTargetedMessage(m *Message, target string, targetType TargetType) *TargetedMessage {
	return &TargetedMessage{
		message:    m,
		targetType: targetType,
		target:     target,
	}
}

func (tm *TargetedMessage) SetTargetType(targetType TargetType) *TargetedMessage {
	tm.targetType = targetType
	return tm
}

func (tm *TargetedMessage) SetTarget(target string) *TargetedMessage {
	tm.target = target
	return tm
}

func (tm *TargetedMessage) JSON() []byte {
	bytes, err := json.Marshal(tm)
	if err != nil {
		panic(err)
	}
	return bytes
}

//-----------------------------------------------------------------------------------//
// 发送给IOS设备的Message对象
func NewIOSMessage(description string) *Message {
	return &Message{
		UniqueID:               "",
		Payload:                "",
		Title:                  "",
		Description:            description,
		PassThrough:            0,
		NotifyType:             -1, // default notify type
		RestrictedPackageNames: nil,
		TimeToLive:             0,
		TimeToSend:             0,
		NotifyID:               0,
		Extra:                  make(map[string]string),
	}
}

// 可选项，自定义通知数字角标。
func (i *Message) SetBadge(badge int64) *Message {
	i.Extra["badge"] = strconv.FormatInt(badge, 10)
	return i
}

// 可选项，iOS8推送消息快速回复类别。
func (i *Message) SetCategory(category string) *Message {
	i.Extra["category"] = category
	return i
}

// 可选项，自定义消息铃声。
func (i *Message) SetSoundURL(soundURL string) *Message {
	i.Extra["sound_url"] = soundURL
	return i
}
