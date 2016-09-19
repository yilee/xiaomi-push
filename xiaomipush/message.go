package xiaomipush

import "time"

type Message struct {
	uniqueID    string
	collapseKey string
	payload     string
	title       string
	description string
	passThrough int32
	notifyType  int32
	notifyID    int64
	timeToLive  int64
	timeToSend  int64
	extra       map[string]string
}

const (
	MaxTimeToSend = time.Hour * 24 * 7
	MaxTimeToLive = time.Hour * 24 * 7 * 2
)

type TargetType int32

const (
	TargetTypeRegID   TargetType = 1
	TargetTypeReAlias TargetType = 2
	TargetTypeAccount TargetType = 3
)

func NewMessage(payload string) *Message {
	return &Message{
		uniqueID:    "",
		payload:     payload,
		title:       "", // TBD
		description: "", // TBD
		passThrough: 1,
		notifyType:  -1, // default notify type
		extra:       make(map[string]string),
	}
}

func (m *Message) SetUniqueID(uniqueID string) *Message {
	m.uniqueID = uniqueID
	return m
}

func (m *Message) SetTimeToSend(tts int64) *Message {
	if time.Since(time.Unix(0, tts*int64(time.Millisecond))) > MaxTimeToSend {
		m.timeToSend = time.Now().Add(MaxTimeToSend).UnixNano() / 1e6
	} else {
		m.timeToSend = tts
	}
	return m
}

func (m *Message) SetTimeToLive(ttl int64) *Message {
	if time.Since(time.Unix(0, ttl*int64(time.Millisecond))) > MaxTimeToLive {
		m.timeToLive = time.Now().Add(MaxTimeToLive).UnixNano() / 1e6
	} else {
		m.timeToLive = ttl
	}
	return m
}

func (m *Message) EnableFlowControl() *Message {
	m.extra["flow_control"] = "1"
	return m
}

func (m *Message) DisableFlowControl() *Message {
	delete(m.extra, "flow_control")
	return m
}

// 小米推送服务器每隔1s将已送达或已点击的消息ID和对应设备的regid或alias通过调用第三方http接口传给开发者。
func (m *Message) SetCallback(callbackURL string) *Message {
	m.extra["callback"] = callbackURL
	m.extra["callback.param"] = m.uniqueID
	m.extra["callback.type"] = "3" // 1:送达回执, 2:点击回执, 3:送达和点击回执
	return m
}

// TargetedMessage封装了MiPush推送服务系统中的消息Message对象，和该Message对象所要发送到的目标。
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
