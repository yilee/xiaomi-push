package xiaomipush

const (
	ProductionHost = "https://api.xmpush.xiaomi.com"
)

const (
	RegURL                               = "/v3/message/regid"                // 向某个regid或一组regid列表推送某条消息
	MultiMessagesRegIDURL                = "/v2/multi_messages/regids"        // 针对不同的regid推送不同的消息
	MultiMessagesAliasURL                = "/v2/multi_messages/aliases"       // 针对不同的aliases推送不同的消息
	MultiMessagesUserAccountURL          = "/v2/multi_messages/user_accounts" // 针对不同的accounts推送不同的消息
	MessageAlisaURL                      = "/v3/message/alias"                // 根据alias，发送消息到指定设备上
	MessageUserAccountURL                = "/v2/message/user_account"         // 根据account，发送消息到指定account上
	MultiPackageNameMessageMultiTopicURL = "/v3/message/multi_topic"          // 根据topic，发送消息到指定一组设备上
	MessageMultiTopicURL                 = "/v2/message/topic"                // 根据topic，发送消息到指定一组设备上
	MultiPackageNameMessageAllURL        = "/v3/message/all"                  // 向所有设备推送某条消息
	MessageAllURL                        = "/v2/message/all"                  // 向所有设备推送某条消息
	MultiTopicURL                        = "/v3/message/multi_topic"          // 向多个topic广播消息
	ScheduleJobExistURL                  = "/v2/schedule_job/exist"           // 检测定时消息的任务是否存在。
	ScheduleJobDeleteURL                 = "/v2/schedule_job/delete"          // 删除指定的定时消息。
	ScheduleJobDeleteByJobKeyURL         = "/v3/schedule_job/delete"          // 删除指定的定时消息。

)

const (
	StatsURL          = "/v1/stats/message/counters" // 统计push
	MessageStatusURL  = "/v1/trace/message/status"   // 获取指定ID的消息状态
	MessagesStatusURL = "/v1/trace/messages/status"  // 获取某个时间间隔内所有消息的状态
)

const (
	TopicSubscribeURL          = "/v2/topic/subscribe"         // 给某个regid订阅标签。
	TopicUnSubscribeURL        = "/v2/topic/unsubscribe"       // 取消某个regid的标签。
	TopicSubscribeByAliasURL   = "/v2/topic/subscribe/alias"   // 给一组alias列表订阅标签
	TopicUnSubscribeByAliasURL = "/v2/topic/unsubscribe/alias" // 取消一组alias列表的标签
)

const (
	InvalidRegIDsURL = "https://feedback.xmpush.xiaomi.com/v1/feedback/fetch_invalid_regids"
)

const (
	AliasAllURL  = "/v1/alias/all" // 获取一个应用的某个用户目前设置的所有Alias
	TopicsAllURL = "/v1/topic/all" // 获取一个应用的某个用户的目前订阅的所有Topic
)

var (
	PostRetryTimes = 3
)

// for future targeted push
var (
	BrandsMap = map[string]string{
		"品牌":    "MODEL",
		"小米":    "xiaomi",
		"三星":    "samsung",
		"华为":    "huawei",
		"中兴":    "zte",
		"中兴努比亚": "nubia",
		"酷派":    "coolpad",
		"联想":    "lenovo",
		"魅族":    "meizu",
		"HTC":   "htc",
		"OPPO":  "oppo",
		"VIVO":  "vivo",
		"摩托罗拉":  "motorola",
		"索尼":    "sony",
		"LG":    "lg",
		"金立":    "jinli",
		"天语":    "tianyu",
		"诺基亚":   "nokia",
		"美图秀秀":  "meitu",
		"谷歌":    "google",
		"TCL":   "tcl",
		"锤子手机":  "chuizi",
		"一加手机":  "1+",
		"中国移动":  "chinamobile",
		"昂达":    "angda",
		"邦华":    "banghua",
		"波导":    "bird",
		"长虹":    "changhong",
		"大可乐":   "dakele",
		"朵唯":    "doov",
		"海尔":    "haier",
		"海信":    "hisense",
		"康佳":    "konka",
		"酷比魔方":  "kubimofang",
		"米歌":    "mige",
		"欧博信":   "ouboxin",
		"欧新":    "ouxin",
		"飞利浦":   "philip",
		"维图":    "voto",
		"小辣椒":   "xiaolajiao",
		"夏新":    "xiaxin",
		"亿通":    "yitong",
		"语信":    "yuxin",
	}

	PriceMap = map[string]string{
		"0-999":     "0-999",
		"1000-1999": "1000-1999",
		"2000-3999": "2000-3999",
		"4000+":     "4000+",
	}
)
