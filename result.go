package xiaomipush

import "encoding/json"

type Result struct {
	MessageID   string `json:"trace_id"`
	Code        int64  `json:"code"`
	Reason      string `json:"reason"`
	Description string `json:"description"`
}

type StatsResult struct {
	Result
	Data []struct {
		Date                string `json:"date"`
		SingleRecipients    int64  `json:"single_recipients"`
		BroadcastRecipients int64  `json:"broadcast_recipients"`
		Received            int64  `json:"received"`
		Click               int64  `json:"click"`
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

func getResultFromJSON(data []byte) (*Result, error) {
	var result Result
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
