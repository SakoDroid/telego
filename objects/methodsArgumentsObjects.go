package objects

import (
	"encoding/json"
	mp "mime/multipart"
)

type MethodArguments interface {
	ToJson() []byte
	ToMultiPart(wr *mp.Writer)
	GetMediaType() string
}

type GetUpdatesArgs struct {
	/*Identifier of the first update to be returned. Must be greater by one than the highest among the identifiers of previously received updates. By default, updates starting with the earliest unconfirmed update are returned. An update is considered confirmed as soon as getUpdates is called with an offset higher than its update_id. The negative offset can be specified to retrieve updates starting from -offset update from the end of the updates queue. All previous updates will forgotten.*/
	Offset int `json:"offset,omitempty"`
	/*Limits the number of updates to be retrieved. Values between 1-100 are accepted. Defaults to 100.*/
	Limit int `json:"limit,omitempty"`
	/*Timeout in seconds for long polling. Defaults to 0, i.e. usual short polling. Should be positive, short polling should be used for testing purposes only.*/
	Timeout int `json:"timeout,omitempty"`
	/*A JSON-serialized list of the update types you want your bot to receive. For example, specify [“message”, “edited_channel_post”, “callback_query”] to only receive updates of these types. See Update for a complete list of available update types. Specify an empty list to receive all update types except chat_member (default). If not specified, the previous setting will be used.
	Please note that this parameter doesnt affect updates created before the call to the getUpdates, so unwanted updates may be received for a short period of time.*/
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

func (args *GetUpdatesArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

func (args *GetUpdatesArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this method are never passed as multipart.
}

func (args *GetUpdatesArgs) GetMediaType() string {
	return "update"
}

type defaultSendMethodsArguments struct {
	/*Unique identifier for the target chat.(only one of ChatIdInt and ChatIdString should be present.)*/
	ChatIdInt int `json:"chat_id,omitempty"`
	/*Username of the target channel (in the format @channelusername). (only one of ChatIdInt and ChatIdString should be present.)*/
	ChatIdString string `json:"chat_id,omitempty"`
	/*Sends the message silently. Users will receive a notification with no sound.*/
	DisableNotification bool `json:"disable_notification,omitempty"`
	/*If the message is a reply, ID of the original message*/
	ReplyToMessageId int `json:"reply_to_message_id,omitempty"`
	/*Pass True, if the message should be sent even if the specified replied-to message is not found*/
	AllowSendingWithoutReply bool `json:"allow_sending_without_reply,omitempty"`
	/*Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove reply keyboard or to force a reply from the user.*/
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type SendMessageArgs struct {
	defaultSendMethodsArguments
	/*Text of the message to be sent, 1-4096 characters after entities parsing*/
	Text string `json:"text"`
	/*Mode for parsing entities in the message text. */
	ParseMode string `json:"parse_mode,omitempty"`
	/*A JSON-serialized list of special entities that appear in message text, which can be specified instead of parse_mode*/
	Entities []MessageEntity `json:"entities,omitempty"`
	/*Disables link previews for links in this message*/
	DisableWebPagePreview bool `json:"disable_web_page_preview,omitempty"`
}

func (args *SendMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

func (args *SendMessageArgs) ToMultiPart(wr *mp.Writer) {

}

func (args *SendMessageArgs) GetMediaType() string {
	return "text"
}

type ForwardMessageArgs struct {
	/*Unique identifier for the target chat. (only one of ChatIdInt and ChatIdString should be present.)*/
	ChatIdInt int `json:"chat_id,omitempty"`
	/*Username of the target channel (in the format @channelusername). (only one of ChatIdInt and ChatIdString should be present.)*/
	ChatIdString string `json:"chat_id,omitempty"`
	/*Unique identifier for the chat where the original message was sent*/
	FromChatIdInt int `json:"from_chat_id,omitempty"`
	/*Channel username in the format @channelusername*/
	FromChatIdString string `json:"from_chat_id,omitempty"`
	/*Sends the message silently. Users will receive a notification with no sound.*/
	DisableNotification bool `json:"disable_notification,omitempty"`
	/*Message identifier in the chat specified in from_chat_id*/
	MessageId int `json:"message_id"`
}

func (args *ForwardMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

func (args *ForwardMessageArgs) ToMultiPart(wr *mp.Writer) {

}

func (args *ForwardMessageArgs) GetMediaType() string {
	return "text"
}

type CopyMessageArgs struct {
	ForwardMessageArgs
	/*New caption for media, 0-1024 characters after entities parsing. If not specified, the original caption is kept*/
	Caption string `json:"caption,omitempty"`
	/*Mode for parsing entities in the message text. */
	ParseMode string `json:"parse_mode,omitempty"`
	/*A JSON-serialized list of special entities that appear in the new caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*If the message is a reply, ID of the original message*/
	ReplyToMessageId int `json:"reply_to_message_id,omitempty"`
	/*Pass True, if the message should be sent even if the specified replied-to message is not found*/
	AllowSendingWithoutReply bool `json:"allow_sending_without_reply,omitempty"`
	/*Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove reply keyboard or to force a reply from the user.*/
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

func (args *CopyMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

func (args *CopyMessageArgs) ToMultiPart(wr *mp.Writer) {

}

func (args *CopyMessageArgs) GetMediaType() string {
	return "text"
}

type SendPhotoArgs struct {
	defaultSendMethodsArguments
	/*Photo to send. Pass a file_id as String to send a photo that exists on the Telegram servers (recommended), pass an HTTP URL as a String for Telegram to get a photo from the Internet, or upload a new photo using multipart/form-data. The photo must be at most 10 MB in size. The photo's width and height must not exceed 10000 in total. Width and height ratio must be at most 20.*/
	Photo string `json:"photo"`
	/*Photo caption (may also be used when resending photos by file_id), 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Mode for parsing entities in the photo caption. */
	ParseMode string `json:"parse_mode,omitempty"`
	/*A JSON-serialized list of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
}

func (args *SendPhotoArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

func (args *SendPhotoArgs) ToMultiPart(wr *mp.Writer) {

}

func (args *SendPhotoArgs) GetMediaType() string {
	return "photo"
}

type SendAudioArgs struct {
	defaultSendMethodsArguments
	/*Audio file to send. Pass a file_id as String to send an audio file that exists on the Telegram servers (recommended), pass an HTTP URL as a String for Telegram to get an audio file from the Internet, or upload a new one using multipart/form-data.*/
	Audio           string          `json:"audio"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Duration of the audio in secconds"
	Duration int `json:"duration,omitempty"`
	Performer string `json:"performer,omitempty"`
	/*Track name*/
	Title string `json:"title,omitempty"`
	/*Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb string `json:"thumb,omitempty"`
}

func (args *SendAudioArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

func (args *SendAudioArgs) ToMultiPart(wr *mp.Writer) {

}

func (args *SendAudioArgs) GetMediaType() string {
	return "audio"
}

type SendDocumentArgs struct {
	defaultSendMethodsArguments
	Document string `json:"document"`
	/*Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb           string          `json:"thumb,omitempty"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Disables automatic server-side content type detection for files uploaded using multipart/form-data*/
	DisableContentTypeDetection bool `json:"disable_content_type_detection,omitempty`
}

func (args *SendDocumentArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

func (args *SendDocumentArgs) ToMultiPart(wr *mp.Writer) {

}

func (args *SendDocumentArgs) GetMediaType() string {
	return "document"
}

type SendVideoArgs struct {
	defaultSendMethodsArguments
	Video string `json:"video"`
	/*Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb           string          `json:"thumb,omitempty"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Duration of sent video in seconds.*/
	Duration int `json:"duration,omitempty"`
	/*Pass True, if the uploaded video is suitable for streaming*/
	SupportsStreaming bool `json:"supports_streaming,omitempty"`
}

func (args *SendVideoArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

func (args *SendVideoArgs) ToMultiPart(wr *mp.Writer) {

}

func (args *SendVideoArgs) GetMediaType() string {
	return "video"
}
