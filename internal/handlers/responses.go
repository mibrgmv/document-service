package handlers

type Response struct {
	Error    interface{} `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type Error struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
