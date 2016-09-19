package xiaomipush

const (
	ProductionHost = "https://api.xmpush.xiaomi.com"
	SandboxHost    = "https://sandbox.xmpush.xiaomi.com" // iOS supported only
)

const (
	RegURL           = "/v3/message/regid"          // 向某个regid或一组regid列表推送某条消息
	MessageAllURL    = "/v3/message/all"            // 向所有设备推送某条消息
	MultiMessagesURL = "/v2/multi_messages/regids"  // 针对不同的regid推送不同的消息
	StatsURL         = "/v1/stats/message/counters" // 统计push
	MessageStatusURL = "/v1/trace/message/status"   // 获取指定ID的消息状态
)

var (
	ErrorCodeMap = map[string]string{
		"-1":    "未知错误",
		"0":     "成功",
		"100":   "参数不能为空",
		"10001": "系统错误",
		"10002": "服务暂停",
		"10003": "服务暂停",
		"10004": "服务暂停",
		"10005": "服务暂停",
		"10006": "服务暂停",
		"10007": "服务暂停",
		"10008": "服务暂停",
		"10009": "服务暂停",
		"10010": "服务暂停",
		"10011": "服务暂停",
		"10012": "服务暂停",
		"10013": "服务暂停",
		"10014": "服务暂停",
		"10016": "服务暂停",
		"10017": "服务暂停",
		"10018": "服务暂停",
		"10020": "服务暂停",
		"10021": "服务暂停",
		"10022": "服务暂停",
		"10023": "服务暂停",
		"10024": "服务暂停",
		"10025": "服务暂停",
		"10026": "服务暂停",
		"10027": "服务暂停",
		"10028": "服务暂停",
		"10029": "服务暂停",

		"22000": "服务暂停",
		"22001": "服务暂停",
		"22002": "服务暂停",
		"22003": "服务暂停",
		"22004": "服务暂停",
		"22005": "服务暂停",
		"22006": "服务暂停",
		"22007": "服务暂停",
		"22008": "服务暂停",
		"22020": "服务暂停",
		"22021": "服务暂停",
		"22022": "服务暂停",
		"22100": "服务暂停",
		"22101": "服务暂停",
		"22102": "服务暂停",
		"22103": "服务暂停",
	}
)

// for targeted push
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
