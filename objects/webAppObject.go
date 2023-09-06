package objects

type WebAppInfo struct {
	URL string `json:"url"`
}
type SentWebAppMessage struct {
	InlineMessageId string `json:"inline_message_id"`
}

type WebAppData struct {
	/*The data. Be aware that a bad client can send arbitrary data in this field.*/
	Data string `json:"data"`
	/*Text of the web_app keyboard bu	tton from which the Web App was opened. Be aware that a bad client can send arbitrary data in this field.*/
	ButtonText string `json:"button_text"`
}
