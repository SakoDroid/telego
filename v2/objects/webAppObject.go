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

/*This object represents a service message about a user allowing a bot to write messages after adding the bot to the attachment menu or launching a Web App from a link.*/
type WriteAccessAllowed struct {
	/*Optional. True, if the access was granted after the user accepted an explicit request from a Web App sent by the method requestWriteAccess*/
	FromRequest bool `json:"from_request"`
	/*Optional. Name of the Web App which was launched from a link*/
	WebAppName string `json:"web_app_name"`
	/*Optional. True, if the access was granted when the bot was added to the attachment or side menu*/
	FromAttachmentMenu bool `json:"from_attachment_menu"`
}
