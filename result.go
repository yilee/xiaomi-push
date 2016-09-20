package xiaomipush

type Result struct {
	MessageID   string `json:"trace_id"`
	Code        int64  `json:"code"`
	Reason      string `json:"reason"`
	Description string `json:"description"`
}

type StatsResult struct {
	Result
	Data struct {
		Data []struct {
			Date                  string `json:"date"`
			AliasRecipients       int64  `json:"alias_recipients"`
			UserAccountRecipients int64  `json:"useraccount_recipients"`
			RegidRecipients       int64  `json:"regid_recipients"`
			Received              int64  `json:"received"`
			BroadcastRecipients   int64  `json:"broadcast_recipients"`
			Click                 int64  `json:"click"`
			SingleRecipients      int64  `json:"single_recipients"`
		} `json:"data,omitempty"`
	} `json:"data,omitempty"`
}

type StatusResult struct {
	Result
	Data struct {
		ID           string
		DeliveryRate string
		Delivered    string
		Resolved     string
		MsgType      string
		CreateTime   string
		TimeToLive   string
		ClickRate    string
	} `json:"data,omitempty"`
}
