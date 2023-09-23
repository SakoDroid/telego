package objects

import (
	"bytes"
	"encoding/json"
	"io"
	mp "mime/multipart"
	"strconv"
	"strings"
)

type MethodArguments interface {
	ToJson() []byte
	ToMultiPart(wr *mp.Writer)
}

type SetWebhookArgs struct {
	/*HTTPS url to send updates to. Use an empty string to remove webhook integration*/
	URL string
	/*public key certificate */
	Certificate string
	/*The fixed IP address which will be used to send webhook requests instead of the IP address resolved through DNS*/
	IPAddress string
	/*Maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery, 1-100. Defaults to 40. Use lower values to limit the load on your bot's server, and higher values to increase your bot's throughput.*/
	MaxConnections int
	/*A JSON-serialized list of the update types you want your bot to receive. For example, specify [“message”, “edited_channel_post”, “callback_query”] to only receive updates of these types. See Update for a complete list of available update types. Specify an empty list to receive all update types except chat_member (default). If not specified, the previous setting will be used.
	Please note that this parameter doesn't affect updates created before the call to the setWebhook, so unwanted updates may be received for a short period of time.*/
	AllowedUpdates []string
	/*Pass True to drop all pending updates*/
	DropPendingUpdates bool
	/*A secret token to be sent in a header “X-Telegram-Bot-Api-Secret-Token” in every webhook request, 1-256 characters. Only characters A-Z, a-z, 0-9, _ and - are allowed. The header is useful to ensure that the request comes from a webhook set by you.*/
	SecretToken string `json:"secret_token,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetWebhookArgs) ToJson() []byte {
	//Arguments of this method are never passed as json.
	return nil
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetWebhookArgs) ToMultiPart(wr *mp.Writer) {
	fr, _ := wr.CreateFormField("url")
	_, _ = io.Copy(fr, strings.NewReader(args.URL))
	fr, _ = wr.CreateFormField("certificate")
	_, _ = io.Copy(fr, strings.NewReader(args.Certificate))
	if args.IPAddress != "" {
		fr, _ = wr.CreateFormField("ip_address")
		_, _ = io.Copy(fr, strings.NewReader(args.IPAddress))
	}
	if args.MaxConnections != 0 {
		fr, _ = wr.CreateFormField("max_connections")
		_, _ = io.Copy(fr, strings.NewReader(strconv.Itoa(args.MaxConnections)))
	}
	if args.AllowedUpdates != nil {
		fr, _ = wr.CreateFormField("allowed_updates")
		bts, _ := json.Marshal(args.AllowedUpdates)
		_, _ = io.Copy(fr, bytes.NewReader(bts))
	}
	fr, _ = wr.CreateFormField("drop_pending_updates")
	_, _ = io.Copy(fr, strings.NewReader(strconv.FormatBool(args.DropPendingUpdates)))
}

type DeleteWebhookArgs struct {
	DropPendingUpdates bool `json:"drop_pending_updates"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *DeleteWebhookArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *DeleteWebhookArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this method are never passed as multipart.
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

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetUpdatesArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetUpdatesArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this method are never passed as multipart.
}

type DefaultSendMethodsArguments struct {
	/*Unique identifier for the target chat or Username of the target channel (in the format @channelusername).*/
	ChatId json.RawMessage `json:"chat_id"`
	/*Sends the message silently. Users will receive a notification with no sound.*/
	DisableNotification bool `json:"disable_notification"`
	/*If the message is a reply, ID of the original message*/
	ReplyToMessageId int `json:"reply_to_message_id,omitempty"`
	/*Pass True, if the message should be sent even if the specified replied-to message is not found*/
	AllowSendingWithoutReply bool `json:"allow_sending_without_reply"`
	/*Protects the contents of sent messages from forwarding and saving*/
	ProtectContent bool `json:"protect_content"`
	/*Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove reply keyboard or to force a reply from the user.*/
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
	/*Unique identifier for the target message thread (topic) of the forum; for forum supergroups only*/
	MessageThreadId int `json:"message_thread_id,omitempty"`
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (df *DefaultSendMethodsArguments) toMultiPart(wr *mp.Writer) {
	fw, _ := wr.CreateFormField("chat_id")
	_, _ = io.Copy(fw, strings.NewReader(string(df.ChatId)))
	fw, _ = wr.CreateFormField("disable_notification")
	_, _ = io.Copy(fw, strings.NewReader(strconv.FormatBool(df.DisableNotification)))
	if df.ReplyToMessageId != 0 {
		fw, _ = wr.CreateFormField("reply_to_message_id")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(df.ReplyToMessageId)))
	}
	fw, _ = wr.CreateFormField("allow_sending_without_reply")
	_, _ = io.Copy(fw, strings.NewReader(strconv.FormatBool(df.AllowSendingWithoutReply)))
	if df.ReplyMarkup != nil {
		fw, _ = wr.CreateFormField("reply_markup")
		bt, _ := json.Marshal(df.ReplyMarkup)
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
}

type SendMessageArgs struct {
	DefaultSendMethodsArguments
	/*Text of the message to be sent, 1-4096 characters after entities parsing*/
	Text string `json:"text"`
	/*Mode for parsing entities in the message text. */
	ParseMode string `json:"parse_mode,omitempty"`
	/*A JSON-serialized list of special entities that appear in message text, which can be specified instead of parse_mode*/
	Entities []MessageEntity `json:"entities,omitempty"`
	/*Disables link previews for links in this message*/
	DisableWebPagePreview bool `json:"disable_web_page_preview,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendMessageArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this method are never passed as multipart.
}

type ForwardMessageArgs struct {
	/*Unique identifier for the target chat or Username of the target channel (in the format @channelusername).*/
	ChatId json.RawMessage `json:"chat_id"`
	/*Unique identifier for the chat where the original message was sent or Channel username in the format @channelusername*/
	FromChatId json.RawMessage `json:"from_chat_id"`
	/*Sends the message silently. Users will receive a notification with no sound.*/
	DisableNotification bool `json:"disable_notification,omitempty"`
	/*Protects the contents of sent messages from forwarding and saving*/
	ProtectContent bool `json:"protect_content"`
	/*Message identifier in the chat specified in from_chat_id*/
	MessageId int `json:"message_id"`
	/*Unique identifier for the target message thread (topic) of the forum; for forum supergroups only*/
	MessageThreadId int `json:"message_thread_id,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *ForwardMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *ForwardMessageArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this method are never passed as multipart.
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

// ToJson converts this strcut into json to be sent to the API server.
func (args *CopyMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *CopyMessageArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this method are never passed as multipart.
}

type SendPhotoArgs struct {
	DefaultSendMethodsArguments
	/*Photo to send. Pass a file_id as String to send a photo that exists on the Telegram servers (recommended), pass an HTTP URL as a String for Telegram to get a photo from the Internet, or upload a new photo using multipart/form-data. The photo must be at most 10 MB in size. The photo's width and height must not exceed 10000 in total. Width and height ratio must be at most 20.*/
	Photo string `json:"photo"`
	/*Photo caption (may also be used when resending photos by file_id), 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Mode for parsing entities in the photo caption. */
	ParseMode string `json:"parse_mode,omitempty"`
	/*A JSON-serialized list of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Pass True if the photo needs to be covered with a spoiler animation*/
	HasSpoiler bool `json:"has_spoiler"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendPhotoArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendPhotoArgs) ToMultiPart(wr *mp.Writer) {
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("photo")
	_, _ = io.Copy(fw, strings.NewReader(args.Photo))
	if args.Caption != "" {
		fw, _ = wr.CreateFormField("caption")
		_, _ = io.Copy(fw, strings.NewReader(args.Caption))
	}
	if args.ParseMode != "" {
		fw, _ = wr.CreateFormField("parse_mode")
		_, _ = io.Copy(fw, strings.NewReader(args.ParseMode))
	}
	if args.CaptionEntities != nil {
		fw, _ = wr.CreateFormField("caption_entities")
		bt, _ := json.Marshal(args.CaptionEntities)
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
}

type SendStickerArgs struct {
	DefaultSendMethodsArguments
	Sticker string `json:"sticker"`
	//Emoji associated with the sticker; only for just uploaded stickers
	Emoji string `json:"emoji"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendStickerArgs) ToJson() []byte {
	//The arguments of this method are never passed as json.
	return nil
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendStickerArgs) ToMultiPart(wr *mp.Writer) {
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("sticker")
	_, _ = io.Copy(fw, strings.NewReader(args.Sticker))
}

type GetStickerSetArgs struct {
	Name string `json:"name"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetStickerSetArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetStickerSetArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this methods are never passed as multipart.
}

type UploadStickerFileArgs struct {
	UserId        int           `json:"user_id"`
	Sticker       *InputSticker `json:"sticker"`
	StickerFormat string        `json:"sticker_format"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *UploadStickerFileArgs) ToJson() []byte {
	//The arguments of this method are neverr passed as json.
	return nil
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *UploadStickerFileArgs) ToMultiPart(wr *mp.Writer) {
	fw, _ := wr.CreateFormField("user_id")
	io.Copy(fw, strings.NewReader(strconv.Itoa(args.UserId)))
	fw, _ = wr.CreateFormField("sticker")
	jsn, _ := json.Marshal(args.Sticker)
	io.Copy(fw, bytes.NewReader(jsn))
	fw, _ = wr.CreateFormField("sticker_format")
	io.Copy(fw, strings.NewReader(args.StickerFormat))
}

type CreateNewStickerSetArgs struct {
	UserId          int             `json:"user_id"`
	Name            string          `json:"name"`
	Title           string          `json:"title"`
	Stickers        []*InputSticker `json:"stickers"`
	StickerFormat   string          `json:"sticker_format"`
	StickerType     string          `json:"sticker_type"`
	NeedsRepainting bool            `json:"needs_repainting"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *CreateNewStickerSetArgs) ToJson() []byte {
	//The arguments of this methos is never passed as json.
	return nil
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *CreateNewStickerSetArgs) ToMultiPart(wr *mp.Writer) {
	fw, _ := wr.CreateFormField("user_id")
	_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.UserId)))
	fw, _ = wr.CreateFormField("name")
	_, _ = io.Copy(fw, strings.NewReader(args.Name))
	fw, _ = wr.CreateFormField("title")
	_, _ = io.Copy(fw, strings.NewReader(args.Title))
	fw, _ = wr.CreateFormField("stickers")
	jsn, _ := json.Marshal(args.Stickers)
	io.Copy(fw, bytes.NewReader(jsn))
	fw, _ = wr.CreateFormField("sticker_format")
	io.Copy(fw, strings.NewReader(args.StickerFormat))
	fw, _ = wr.CreateFormField("sticker_type")
	io.Copy(fw, strings.NewReader(args.StickerType))
	fw, _ = wr.CreateFormField("needs_repainting")
	io.Copy(fw, strings.NewReader(strconv.FormatBool(args.NeedsRepainting)))
}

type AddStickerSetArgs struct {
	UserId  int           `json:"user_id"`
	Name    string        `json:"name"`
	Sticker *InputSticker `json:"sticker"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *AddStickerSetArgs) ToJson() []byte {
	//The arguments of this methos is never passed as json.
	return nil
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *AddStickerSetArgs) ToMultiPart(wr *mp.Writer) {
	fw, _ := wr.CreateFormField("user_id")
	_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.UserId)))
	fw, _ = wr.CreateFormField("name")
	_, _ = io.Copy(fw, strings.NewReader(args.Name))
	fw, _ = wr.CreateFormField("sticker")
	jsn, _ := json.Marshal(args.Sticker)
	io.Copy(fw, bytes.NewReader(jsn))
}

type SetStickerPositionInSetArgs struct {
	Sticker  string `json:"sticker"`
	Position int    `json:"position"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetStickerPositionInSetArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetStickerPositionInSetArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this meethod are never passed as multipart.
}

type DeleteStickerFromSetArgs struct {
	Sticker string `json:"sticker"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *DeleteStickerFromSetArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *DeleteStickerFromSetArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this meethod are never passed as multipart.
}

type SetStickerEmojiListArgs struct {
	Sticker   string   `json:"sticker"`
	EmojiList []string `json:"emoji_list"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetStickerEmojiListArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetStickerEmojiListArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this meethod are never passed as multipart.
}

type SetStickerKeywordsArgs struct {
	Sticker   string   `json:"sticker"`
	Keywoards []string `json:"keywords"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetStickerKeywordsArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetStickerKeywordsArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this meethod are never passed as multipart.
}

type SetStickerMaskPositionArgs struct {
	Sticker      string        `json:"sticker"`
	MaskPosition *MaskPosition `json:"mask_position"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetStickerMaskPositionArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetStickerMaskPositionArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this meethod are never passed as multipart.
}

type SetStickerSetTitleArgs struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetStickerSetTitleArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetStickerSetTitleArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this meethod are never passed as multipart.
}

type DeleteStickerSetArgs struct {
	Name string `json:"name"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *DeleteStickerSetArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *DeleteStickerSetArgs) ToMultiPart(wr *mp.Writer) {
	//The arguments of this meethod are never passed as multipart.
}

type SetStickerSetThumbnailArgs struct {
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
	Thumb  string `json:"thumb"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetStickerSetThumbnailArgs) ToJson() []byte {
	//The arguments of this methos is never passed as json.
	return nil
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetStickerSetThumbnailArgs) ToMultiPart(wr *mp.Writer) {
	fw, _ := wr.CreateFormField("user_id")
	_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.UserId)))
	fw, _ = wr.CreateFormField("name")
	_, _ = io.Copy(fw, strings.NewReader(args.Name))
	fw, _ = wr.CreateFormField("thumbnail")
	_, _ = io.Copy(fw, strings.NewReader(args.Thumb))
}

type SendAudioArgs struct {
	DefaultSendMethodsArguments
	/*Audio file to send. Pass a file_id as String to send an audio file that exists on the Telegram servers (recommended), pass an HTTP URL as a String for Telegram to get an audio file from the Internet, or upload a new one using multipart/form-data.*/
	Audio           string          `json:"audio"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Duration of the audio in secconds*/
	Duration  int    `json:"duration,omitempty"`
	Performer string `json:"performer,omitempty"`
	/*Track name*/
	Title string `json:"title,omitempty"`
	/*Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb string `json:"thumbnail,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendAudioArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendAudioArgs) ToMultiPart(wr *mp.Writer) {
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("audio")
	_, _ = io.Copy(fw, strings.NewReader(args.Audio))
	if args.Duration != 0 {
		fw, _ = wr.CreateFormField("duration")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.Duration)))
	}
	if args.Thumb != "" {
		fw, _ = wr.CreateFormField("thumbnail")
		_, _ = io.Copy(fw, strings.NewReader(args.Thumb))
	}
	if args.Caption != "" {
		fw, _ = wr.CreateFormField("caption")
		_, _ = io.Copy(fw, strings.NewReader(args.Caption))
	}
	if args.ParseMode != "" {
		fw, _ = wr.CreateFormField("parse_mode")
		_, _ = io.Copy(fw, strings.NewReader(args.ParseMode))
	}
	if args.CaptionEntities != nil {
		fw, _ = wr.CreateFormField("caption_entities")
		bt, _ := json.Marshal(args.CaptionEntities)
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
	if args.Performer != "" {
		fw, _ = wr.CreateFormField("performer")
		_, _ = io.Copy(fw, strings.NewReader(args.Performer))
	}
	if args.Title != "" {
		fw, _ = wr.CreateFormField("title")
		_, _ = io.Copy(fw, strings.NewReader(args.Title))
	}
}

type SendDocumentArgs struct {
	DefaultSendMethodsArguments
	Document string `json:"document"`
	/*Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb           string          `json:"thumbnail,omitempty"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Disables automatic server-side content type detection for files uploaded using multipart/form-data*/
	DisableContentTypeDetection bool `json:"disable_content_type_detection,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendDocumentArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendDocumentArgs) ToMultiPart(wr *mp.Writer) {
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("document")
	_, _ = io.Copy(fw, strings.NewReader(args.Document))
	if args.Thumb != "" {
		fw, _ = wr.CreateFormField("thumbnail")
		_, _ = io.Copy(fw, strings.NewReader(args.Thumb))
	}
	if args.Caption != "" {
		fw, _ = wr.CreateFormField("caption")
		_, _ = io.Copy(fw, strings.NewReader(args.Caption))
	}
	if args.ParseMode != "" {
		fw, _ = wr.CreateFormField("parse_mode")
		_, _ = io.Copy(fw, strings.NewReader(args.ParseMode))
	}
	if args.CaptionEntities != nil {
		fw, _ = wr.CreateFormField("caption_entities")
		bt, _ := json.Marshal(args.CaptionEntities)
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
	fw, _ = wr.CreateFormField("disable_content_type_detection")
	_, _ = io.Copy(fw, strings.NewReader(strconv.FormatBool(args.DisableContentTypeDetection)))
}

type SendVideoArgs struct {
	DefaultSendMethodsArguments
	Video string `json:"video"`
	/*Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb           string          `json:"thumbnail,omitempty"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Pass True if the photo needs to be covered with a spoiler animation*/
	HasSpoiler bool `json:"has_spoiler"`
	/*Duration of sent video in seconds.*/
	Duration int `json:"duration,omitempty"`
	/*Pass True, if the uploaded video is suitable for streaming*/
	SupportsStreaming bool `json:"supports_streaming,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendVideoArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendVideoArgs) ToMultiPart(wr *mp.Writer) {
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("video")
	_, _ = io.Copy(fw, strings.NewReader(args.Video))
	fw, _ = wr.CreateFormField("supports_streaming")
	_, _ = io.Copy(fw, strings.NewReader(strconv.FormatBool(args.SupportsStreaming)))
	if args.Duration != 0 {
		fw, _ = wr.CreateFormField("duration")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.Duration)))
	}
	if args.Thumb != "" {
		fw, _ = wr.CreateFormField("thumbnail")
		_, _ = io.Copy(fw, strings.NewReader(args.Thumb))
	}
	if args.Caption != "" {
		fw, _ = wr.CreateFormField("caption")
		_, _ = io.Copy(fw, strings.NewReader(args.Caption))
	}
	if args.ParseMode != "" {
		fw, _ = wr.CreateFormField("parse_mode")
		_, _ = io.Copy(fw, strings.NewReader(args.ParseMode))
	}
	if args.CaptionEntities != nil {
		fw, _ = wr.CreateFormField("caption_entities")
		bt, _ := json.Marshal(args.CaptionEntities)
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
}

type SendAnimationArgs struct {
	DefaultSendMethodsArguments
	Animation string `json:"animation"`
	/*Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb           string          `json:"thumbnail,omitempty"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Pass True if the photo needs to be covered with a spoiler animation*/
	HasSpoiler bool `json:"has_spoiler"`
	Width      int  `json:"width,omitempty"`
	Height     int  `json:"height,omitempty"`
	/*Duration of sent video in seconds.*/
	Duration int `json:"duration,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendAnimationArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendAnimationArgs) ToMultiPart(wr *mp.Writer) {
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("animation")
	_, _ = io.Copy(fw, strings.NewReader(args.Animation))
	if args.Duration != 0 {
		fw, _ = wr.CreateFormField("duration")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.Duration)))
	}
	if args.Thumb != "" {
		fw, _ = wr.CreateFormField("thumbnail")
		_, _ = io.Copy(fw, strings.NewReader(args.Thumb))
	}
	if args.Caption != "" {
		fw, _ = wr.CreateFormField("caption")
		_, _ = io.Copy(fw, strings.NewReader(args.Caption))
	}
	if args.ParseMode != "" {
		fw, _ = wr.CreateFormField("parse_mode")
		_, _ = io.Copy(fw, strings.NewReader(args.ParseMode))
	}
	if args.CaptionEntities != nil {
		fw, _ = wr.CreateFormField("caption_entities")
		bt, _ := json.Marshal(args.CaptionEntities)
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
	if args.Width != 0 {
		fw, _ = wr.CreateFormField("width")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.Width)))
	}
	if args.Height != 0 {
		fw, _ = wr.CreateFormField("height")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.Height)))
	}
}

type SendVoiceArgs struct {
	DefaultSendMethodsArguments
	Voice           string          `json:"voice"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Duration of sent video in seconds.*/
	Duration int `json:"duration,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendVoiceArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendVoiceArgs) ToMultiPart(wr *mp.Writer) {
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("voice")
	_, _ = io.Copy(fw, strings.NewReader(args.Voice))
	if args.Duration != 0 {
		fw, _ = wr.CreateFormField("duration")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.Duration)))
	}
	if args.Caption != "" {
		fw, _ = wr.CreateFormField("caption")
		_, _ = io.Copy(fw, strings.NewReader(args.Caption))
	}
	if args.ParseMode != "" {
		fw, _ = wr.CreateFormField("parse_mode")
		_, _ = io.Copy(fw, strings.NewReader(args.ParseMode))
	}
	if args.CaptionEntities != nil {
		fw, _ = wr.CreateFormField("caption_entities")
		bt, _ := json.Marshal(args.CaptionEntities)
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
}

type SendVideoNoteArgs struct {
	DefaultSendMethodsArguments
	VideoNote string `json:"video_note"`
	/*Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb           string          `json:"thumbnail,omitempty"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Length          int             `json:"length,omitempty"`
	/*Duration of sent video in seconds.*/
	Duration int `json:"duration,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendVideoNoteArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendVideoNoteArgs) ToMultiPart(wr *mp.Writer) {
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("video_note")
	_, _ = io.Copy(fw, strings.NewReader(args.VideoNote))
	if args.Duration != 0 {
		fw, _ = wr.CreateFormField("duration")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.Duration)))
	}
	if args.Thumb != "" {
		fw, _ = wr.CreateFormField("thumbnail")
		_, _ = io.Copy(fw, strings.NewReader(args.Thumb))
	}
	if args.Caption != "" {
		fw, _ = wr.CreateFormField("caption")
		_, _ = io.Copy(fw, strings.NewReader(args.Caption))
	}
	if args.ParseMode != "" {
		fw, _ = wr.CreateFormField("parse_mode")
		_, _ = io.Copy(fw, strings.NewReader(args.ParseMode))
	}
	if args.CaptionEntities != nil {
		fw, _ = wr.CreateFormField("caption_entities")
		bt, _ := json.Marshal(args.CaptionEntities)
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
	if args.Length != 0 {
		fw, _ = wr.CreateFormField("lengt")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.Length)))
	}
}

type SendMediaGroupArgs struct {
	DefaultSendMethodsArguments
	Media []InputMedia `json:"media"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendMediaGroupArgs) ToJson() []byte {
	args.ReplyMarkup = nil
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendMediaGroupArgs) ToMultiPart(wr *mp.Writer) {
	args.ReplyMarkup = nil
	args.toMultiPart(wr)
	fw, _ := wr.CreateFormField("media")
	bt, _ := json.Marshal(args.Media)
	_, _ = io.Copy(fw, bytes.NewReader(bt))
}

type SendLocationArgs struct {
	DefaultSendMethodsArguments
	/*Latitude of the location*/
	Latitude float32 `json:"latitude"`
	/*Longitude of the location*/
	Longitude float32 `json:"longitude"`
	/*The radius of uncertainty for the location, measured in meters; 0-1500*/
	HorizontalAccuracy float32 `json:"horizontal_accuracy,omitempty"`
	/*Period in seconds for which the location will be updated (see Live Locations, should be between 60 and 86400.*/
	LivePeriod int `json:"live_period,omitempty"`
	/*	For live locations, a direction in which the user is moving, in degrees. Must be between 1 and 360 if specified.*/
	Heading int `json:"heading,omitempty"`
	/*For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters. Must be between 1 and 100000 if specified.*/
	ProximityAlertRadius int `json:"proximity_alert_radius,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendLocationArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendLocationArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type EditMessageLiveLocationArgs struct {
	/*Required if inline_message_id is not specified. Unique identifier for the target chat or username of the target channel (in the format @channelusername)*/
	ChatId json.RawMessage `json:"chat_id,omitempty"`
	/*Required if inline_message_id is not specified. Identifier of the message to edit*/
	MessageId int `json:"message_id,omitempty"`
	/*Required if chat_id and message_id are not specified. Identifier of the inline message*/
	InlineMessageId string `json:"inline_message_id,omitempty"`
	/*Latitude of the location*/
	Latitude float32 `json:"latitude"`
	/*Longitude of the location*/
	Longitude float32 `json:"longitude"`
	/*The radius of uncertainty for the location, measured in meters; 0-1500*/
	HorizontalAccuracy float32 `json:"horizontal_accuracy,omitempty"`
	/*	For live locations, a direction in which the user is moving, in degrees. Must be between 1 and 360 if specified.*/
	Heading int `json:"heading,omitempty"`
	/*For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters. Must be between 1 and 100000 if specified.*/
	ProximityAlertRadius int `json:"proximity_alert_radius,omitempty"`
	/*Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove reply keyboard or to force a reply from the user.*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *EditMessageLiveLocationArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *EditMessageLiveLocationArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type StopMessageLiveLocationArgs struct {
	/*Required if inline_message_id is not specified. Unique identifier for the target chat or username of the target channel (in the format @channelusername)*/
	ChatId json.RawMessage `json:"chat_id,omitempty"`
	/*Required if inline_message_id is not specified. Identifier of the message to edit*/
	MessageId int `json:"message_id,omitempty"`
	/*Required if chat_id and message_id are not specified. Identifier of the inline message*/
	InlineMessageId string `json:"inline_message_id,omitempty"`
	/*Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove reply keyboard or to force a reply from the user.*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *StopMessageLiveLocationArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *StopMessageLiveLocationArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SendVenueArgs struct {
	DefaultSendMethodsArguments
	/*Latitude of the location*/
	Latitude float32 `json:"latitude"`
	/*Longitude of the location*/
	Longitude float32 `json:"longitude"`
	Title     string  `json:"title"`
	Address   string  `json:"address"`
	/*Foursquare identifier of the venue*/
	FoursquareId string `json:"foursquare_id,omitempty"`
	/*Foursquare type of the venue, if known. (For example, “arts_entertainment/default”, “arts_entertainment/aquarium” or “food/icecream”.)*/
	FoursquareType  string `json:"foursquare_type,omitempty"`
	GooglePlaceId   string `json:"google_place_id,omitempty"`
	GooglePlaceType string `json:"google_place_type,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendVenueArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendVenueArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SendContactArgs struct {
	DefaultSendMethodsArguments
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	Vcard       string `json:"vcard,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendContactArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendContactArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SendPollArgs struct {
	DefaultSendMethodsArguments
	/*Poll question, 1-300 characters*/
	Question string `json:"question"`
	/*A JSON-serialized list of answer options, 2-10 strings 1-100 characters each*/
	Options []string `json:"options"`
	/*Pass True, if the poll needs to be immediately closed. This can be useful for poll preview.*/
	IsClosed bool `json:"is_closed,omitempty"`
	/*True, if the poll needs to be anonymous, defaults to True*/
	IsAnonymous bool `json:"is_anonymous,omitempty"`
	/*Poll type, “quiz” or “regular”, defaults to “regular”*/
	Type string `json:"type,omitempty"`
	/*	True, if the poll allows multiple answers, ignored for polls in quiz mode, defaults to False*/
	AllowsMultipleAnswers bool `json:"allows_multiple_answers,omitempty"`
	/*0-based identifier of the correct answer option, required for polls in quiz mode*/
	CorrectOptionId int `json:"correct_option_id"`
	/*Text that is shown when a user chooses an incorrect answer or taps on the lamp icon in a quiz-style poll, 0-200 characters with at most 2 line feeds after entities parsing*/
	Explanation string `json:"explanation,omitempty"`
	/*Mode for parsing entities in the explanation. See formatting options for more details.*/
	ExplanationParseMode string `json:"explanation_parse_mode,omitempty"`
	/*	A JSON-serialized list of special entities that appear in the poll explanation, which can be specified instead of parse_mode*/
	ExplanationEntities []MessageEntity `json:"explanation_entities,omitempty"`
	/*	Amount of time in seconds the poll will be active after creation, 5-600. Can't be used together with close_date.*/
	OpenPeriod int `json:"open_period,omitempty"`
	/*the poll will be automatically closed. Must be at least 5 and no more than 600 seconds in the future. Can't be used together with open_period.*/
	CloseDate int `json:"close_date,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendPollArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendPollArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SendDiceArgs struct {
	DefaultSendMethodsArguments
	/*Emoji on which the dice throw animation is based. Currently, must be one of “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”. Dice can have values 1-6 for “🎲”, “🎯” and “🎳”, values 1-5 for “🏀” and “⚽”, and values 1-64 for “🎰”. Defaults to “🎲”*/
	Emoji string `json:"emoji,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendDiceArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendDiceArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SendChatActionArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
	/*Type of action to broadcast. Choose one, depending on what the user is about to receive: typing for text messages, upload_photo for photos, record_video or upload_video for videos, record_voice or upload_voice for voice notes, upload_document for general files, choose_sticker for stickers, find_location for location data, record_video_note or upload_video_note for video notes.*/
	Action string `json:"action"`
	/*Unique identifier for the target message thread; supergroups only*/
	MessageThreaddId int `json:"message_thread_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendChatActionArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendChatActionArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type GetUserProfilePhototsArgs struct {
	/*Unique identifier of the target user*/
	UserId int `json:"user_id"`
	/*Sequential number of the first photo to be returned. By default, all photos are returned.*/
	Offset int `json:"offset,omitempty"`
	/*Limits the number of photos to be retrieved. Values between 1-100 are accepted. Defaults to 100.*/
	Limit int `json:"limit,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetUserProfilePhototsArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetUserProfilePhototsArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type GetFileArgs struct {
	FileId string `json:"file_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetFileArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetFileArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type DefaultChatArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *DefaultChatArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *DefaultChatArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type BanChatMemberArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
	UserId int             `json:"user_id"`
	/*Date when the user will be unbanned, unix time. If user is banned for more than 366 days or less than 30 seconds from the current time they are considered to be banned forever. Applied for supergroups and channels only.*/
	UntilDate int `json:"until_date,omitempty"`
	/*Pass True to delete all messages from the chat for the user that is being removed. If False, the user will be able to see messages in the group that were sent before the user was removed. Always True for supergroups and channels.*/
	RevokeMessages bool `json:"revoke_messages,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *BanChatMemberArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *BanChatMemberArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type UnbanChatMemberArgsArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
	UserId int             `json:"user_id"`
	/*Do nothing if the user is not banned*/
	OnlyIfBanned bool `json:"only_if_banned,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *UnbanChatMemberArgsArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *UnbanChatMemberArgsArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type RestrictChatMemberArgs struct {
	ChatId                        json.RawMessage `json:"chat_id"`
	UserId                        int             `json:"user_id"`
	Permission                    ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions"`
	UntilDate                     int             `json:"until_date,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *RestrictChatMemberArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *RestrictChatMemberArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type PromoteChatMemberArgs struct {
	ChatId              json.RawMessage `json:"chat_id"`
	UserId              int             `json:"user_id"`
	IsAnonymous         bool            `json:"is_anonymous"`
	CanManageChat       bool            `json:"can_manage_chat"`
	CanPostMessages     bool            `json:"can_post_messages"`
	CanEditMessages     bool            `json:"can_edit_messages"`
	CanDeleteMessages   bool            `json:"can_delete_messages"`
	CanPostStories      bool            `json:"can_post_stories"`
	CanEditStories      bool            `json:"can_edit_stories"`
	CanDeleteStories    bool            `json:"can_delete_stories"`
	CanManageVideoChats bool            `json:"can_manage_video_chats"`
	CanRestrictMembers  bool            `json:"can_restrict_members"`
	CanPromoteMembers   bool            `json:"can_promote_members"`
	CanChangeInfo       bool            `json:"can_change_info"`
	CanInviteUsers      bool            `json:"can_invite_users"`
	CanPinMessages      bool            `json:"can_pin_messages"`
	CanManageTopics     bool            `json:"can_manage_topics"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *PromoteChatMemberArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *PromoteChatMemberArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetChatAdministratorCustomTitleArgs struct {
	ChatId      json.RawMessage `json:"chat_id"`
	UserId      int             `json:"user_id"`
	CustomTitle string          `json:"custom_title"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetChatAdministratorCustomTitleArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetChatAdministratorCustomTitleArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type BanChatSenderChatArgs struct {
	ChatId       json.RawMessage `json:"chat_id"`
	SenderChatId int             `json:"sender_chat_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *BanChatSenderChatArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *BanChatSenderChatArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type UnbanChatSenderChatArgs struct {
	ChatId       json.RawMessage `json:"chat_id"`
	SenderChatId int             `json:"sender_chat_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *UnbanChatSenderChatArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *UnbanChatSenderChatArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetChatPermissionsArgs struct {
	ChatId                        json.RawMessage `json:"chat_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetChatPermissionsArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetChatPermissionsArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type CreateChatInviteLinkArgs struct {
	ChatId             json.RawMessage `json:"chat_id"`
	Name               string          `json:"name,omitempty"`
	ExpireDate         int             `json:"expire_date,omitempty"`
	MemberLimit        int             `json:"member_limit,omitempty"`
	CreatesjoinRequest bool            `json:"creates_join_request,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *CreateChatInviteLinkArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *CreateChatInviteLinkArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type EditChatInviteLinkArgs struct {
	ChatId             json.RawMessage `json:"chat_id"`
	InviteLink         string          `json:"invite_link"`
	Name               string          `json:"name,omitempty"`
	ExpireDate         int             `json:"expire_date,omitempty"`
	MemberLimit        int             `json:"member_limit,omitempty"`
	CreatesjoinRequest bool            `json:"creates_join_request,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *EditChatInviteLinkArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *EditChatInviteLinkArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type RevokeChatInviteLinkArgs struct {
	ChatId     json.RawMessage `json:"chat_id"`
	InviteLink string          `json:"invite_link"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *RevokeChatInviteLinkArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *RevokeChatInviteLinkArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type ApproveChatJoinRequestArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
	UserId int             `json:"user_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *ApproveChatJoinRequestArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *ApproveChatJoinRequestArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type DeclineChatJoinRequestArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
	UserId int             `json:"user_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *DeclineChatJoinRequestArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *DeclineChatJoinRequestArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetChatPhotoArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
	Photo  string          `json:"photo"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetChatPhotoArgs) ToJson() []byte {
	//The arguments of this method are never passed as json.
	return nil
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetChatPhotoArgs) ToMultiPart(wr *mp.Writer) {
	fw, _ := wr.CreateFormField("chat_id")
	_, _ = io.Copy(fw, bytes.NewReader(args.ChatId))
	fw, _ = wr.CreateFormField("photo")
	_, _ = io.Copy(fw, strings.NewReader(args.Photo))
}

type SetChatTitleArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
	Title  string          `json:"title"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetChatTitleArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetChatTitleArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetChatDescriptionArgs struct {
	ChatId      json.RawMessage `json:"chat_id"`
	Description string          `json:"description"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetChatDescriptionArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetChatDescriptionArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type PinChatMessageArgs struct {
	ChatId              json.RawMessage `json:"chat_id"`
	MessageId           int             `json:"message_id"`
	DisableNotification bool            `json:"disable_notification"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *PinChatMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *PinChatMessageArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type UnpinChatMessageArgs struct {
	ChatId    json.RawMessage `json:"chat_id"`
	MessageId int             `json:"message_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *UnpinChatMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *UnpinChatMessageArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type GetChatMemberArgs struct {
	ChatId json.RawMessage `json:"chat_id"`
	UserId int             `json:"user_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetChatMemberArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetChatMemberArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetChatStcikerSet struct {
	ChatId         json.RawMessage `json:"chat_id"`
	StickerSetName string          `json:"sticker_set_name"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetChatStcikerSet) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetChatStcikerSet) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type AnswerCallbackQueryArgs struct {
	CallbackQueyId string `json:"callback_query_id"`
	Text           string `json:"text,omitempty"`
	ShowAlert      bool   `json:"show_alert"`
	URL            string `json:"url,omitempty"`
	CacheTime      int    `json:"cache_time,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *AnswerCallbackQueryArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *AnswerCallbackQueryArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type MyCommandsDefault struct {
	Scope        BotCommandScope `json:"scope,omitempty"`
	LanguageCode string          `json:"language_code,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *MyCommandsDefault) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *MyCommandsDefault) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetMyCommandsArgs struct {
	Commands []BotCommand `json:"commands"`
	MyCommandsDefault
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetMyCommandsArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetMyCommandsArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type EditMessageDefaultArgs struct {
	ChatId          json.RawMessage       `json:"chat_id,omitempty"`
	MessageId       int                   `json:"message_id,omitempty"`
	InlineMessageId string                `json:"inline_message_id,omitempty"`
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type EditMessageTextArgs struct {
	EditMessageDefaultArgs
	Text                  string          `json:"text"`
	ParseMode             string          `json:"parse_mode,omitempty"`
	Entities              []MessageEntity `json:"entities,omitempty"`
	DisablewebpagePreview bool            `json:"disable_web_page_preview"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *EditMessageTextArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *EditMessageTextArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type EditMessageCaptionArgs struct {
	EditMessageDefaultArgs
	Caption               string          `json:"caption,omitempty"`
	ParseMode             string          `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	DisablewebpagePreview bool            `json:"disable_web_page_preview"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *EditMessageCaptionArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *EditMessageCaptionArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type EditMessageMediaArgs struct {
	EditMessageDefaultArgs
	Media InputMedia
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *EditMessageMediaArgs) ToJson() []byte {
	//The arguments of this method are never passed as json.
	return nil
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *EditMessageMediaArgs) ToMultiPart(wr *mp.Writer) {
	if args.ChatId != nil {
		fw, _ := wr.CreateFormField("chat_id")
		_, _ = io.Copy(fw, strings.NewReader(string(args.ChatId)))
	}
	if args.MessageId != 0 {
		fw, _ := wr.CreateFormField("message_id")
		_, _ = io.Copy(fw, strings.NewReader(strconv.Itoa(args.MessageId)))
	}
	if args.InlineMessageId != "" {
		fw, _ := wr.CreateFormField("inline_message_id")
		_, _ = io.Copy(fw, strings.NewReader(args.InlineMessageId))
	}
	bt, _ := json.Marshal(args.Media)
	fw, _ := wr.CreateFormField("media")
	_, _ = io.Copy(fw, bytes.NewReader(bt))
	if args.ReplyMarkup != nil {
		bt, _ = json.Marshal(args.ReplyMarkup)
		fw, _ = wr.CreateFormField("reply_markup")
		_, _ = io.Copy(fw, bytes.NewReader(bt))
	}
}

type EditMessageReplyMakrupArgs struct {
	EditMessageDefaultArgs
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *EditMessageReplyMakrupArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *EditMessageReplyMakrupArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type DeleteMessageArgs struct {
	ChatId    json.RawMessage `json:"chat_id"`
	MessageId int             `json:"message_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *DeleteMessageArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *DeleteMessageArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type StopPollArgs struct {
	ChatId      json.RawMessage       `json:"chat_id"`
	MessageId   int                   `json:"message_id"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *StopPollArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *StopPollArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type AnswerInlineQueryArgs struct {
	InlineQueryId string                    `json:"inline_query_id"`
	Results       []InlineQueryResult       `json:"results"`
	CacheTime     int                       `json:"cache_time,omitempty"`
	IsPersonal    bool                      `json:"is_personal"`
	NextOffset    string                    `json:"next_offset,omitempty"`
	Button        *InlineQueryResultsButton `json:"button,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *AnswerInlineQueryArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *AnswerInlineQueryArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SendInvoiceArgs struct {
	DefaultSendMethodsArguments
	Title                     string         `json:"title"`
	Description               string         `json:"description"`
	Payload                   string         `json:"payload"`
	ProviderToken             string         `json:"provider_token"`
	Currency                  string         `json:"currency"`
	Prices                    []LabeledPrice `json:"prices"`
	MaxTipAmount              int            `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts       []int          `json:"suggested_tip_amounts,omitempty"`
	StartParameter            string         `json:"start_parameter,omitempty"`
	ProviderData              string         `json:"provider_data,omitempty"`
	PhotoURL                  string         `json:"photo_url,omitempty"`
	PhotoSize                 int            `json:"photo_size,omitempty"`
	PhotoWidth                int            `json:"photo_width,omitempty"`
	PhotoHeight               int            `json:"photo_height,omitempty"`
	NeedName                  bool           `json:"need_name"`
	NeedPhoneNumber           bool           `json:"need_phone_number"`
	NeedEmail                 bool           `json:"need_email"`
	NeedShippingAddress       bool           `json:"need_shipping_address"`
	SendPhoneNumberToProvider bool           `json:"send_phone_number_to_provider"`
	SendEmailToProvider       bool           `json:"send_email_to_provider"`
	IsFlexible                bool           `json:"is_flexible"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendInvoiceArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendInvoiceArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type AnswerShippingQueryArgs struct {
	ShippingQueryId string           `json:"shipping_query_id"`
	OK              bool             `json:"ok"`
	ShippingOptions []ShippingOption `json:"shipping_options,omitempty"`
	ErrorMessage    string           `json:"error_message,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *AnswerShippingQueryArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *AnswerShippingQueryArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type AnswerPreCheckoutQueryArgs struct {
	PreCheckoutQueryId string `json:"pre_checkout_query_id"`
	Ok                 bool   `json:"ok"`
	ErrorMessage       string `json:"error_message,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *AnswerPreCheckoutQueryArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *AnswerPreCheckoutQueryArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetPassportDataErrorsArgs struct {
	UserId int                    `json:"user_id"`
	Errors []PassportElementError `json:"errors"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetPassportDataErrorsArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetPassportDataErrorsArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SendGameArgs struct {
	DefaultSendMethodsArguments
	GameShortName string `json:"game_short_name"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SendGameArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SendGameArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetGameScoreArgs struct {
	UserId             int    `json:"user_id"`
	Score              int    `json:"score"`
	Force              bool   `json:"force"`
	DisableEditMessage bool   `json:"disable_edit_message"`
	ChatId             int    `json:"chat_id,omitempty"`
	MessageId          int    `json:"message_id,omitempty"`
	InlineMessageId    string `json:"inline_message_id,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetGameScoreArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetGameScoreArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type GetGameHighScoresArgs struct {
	UserId          int    `json:"user_id"`
	ChatId          int    `json:"chat_id,omitempty"`
	MessageId       int    `json:"message_id,omitempty"`
	InlineMessageId string `json:"inline_message_id,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetGameHighScoresArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetGameHighScoresArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type AnswerWebAppQueryArgs struct {
	WebAppQueryId string            `json:"web_app_query_id"`
	Result        InlineQueryResult `json:"result,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *AnswerWebAppQueryArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *AnswerWebAppQueryArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type ChatMenuButtonArgs struct {
	ChatId     int64       `json:"chat_id"`
	MenuButton *MenuButton `json:"menu_button,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *ChatMenuButtonArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *ChatMenuButtonArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type MyDefaultAdministratorRightsArgs struct {
	/*Pass True to change the default administrator rights of the bot in channels. Otherwise, the default administrator rights of the bot for groups and supergroups will be changed.*/
	ForChannels bool `json:"for_channels"`
	/*A JSON-serialized object describing new default administrator rights. If not specified, the default administrator rights will be cleared.*/
	Rights *ChatAdministratorRights `json:"rights,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *MyDefaultAdministratorRightsArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *MyDefaultAdministratorRightsArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type GetCustomEmojiStickersArgs struct {
	CustomEmojiIds []string `json:"custom_emoji_ids"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetCustomEmojiStickersArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetCustomEmojiStickersArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type CreateForumTopicArgs struct {
	ChatId            json.RawMessage `json:"chat_id"`
	Name              string          `json:"name"`
	IconColor         int             `json:"icon_color,omitempty"`
	IconCustomEmojiId string          `json:"icon_custom_emoji_id,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *CreateForumTopicArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *CreateForumTopicArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type EditForumTopicArgs struct {
	ChatId            json.RawMessage `json:"chat_id"`
	MessageThreadId   int             `json:"message_thread_id"`
	Name              string          `json:"name"`
	IconCustomEmojiId string          `json:"icon_custom_emoji_id,omitempty"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *EditForumTopicArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *EditForumTopicArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type CloseForumTopicArgs struct {
	ChatId          json.RawMessage `json:"chat_id"`
	MessageThreadId int             `json:"message_thread_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *CloseForumTopicArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *CloseForumTopicArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type ReopenForumTopicArgs struct {
	*CloseForumTopicArgs
}

type DeleteForumTopicArgs struct {
	*CloseForumTopicArgs
}

type UnpinAllForumTopicMessages struct {
	*CloseForumTopicArgs
}

type EditGeneralForumTopic struct {
	ChatId json.RawMessage `json:"chat_id"`
	Name   string          `json:"name"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *EditGeneralForumTopic) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *EditGeneralForumTopic) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type CloseGeneralForumTopic struct {
	ChatId json.RawMessage `json:"chat_id"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *CloseGeneralForumTopic) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *CloseGeneralForumTopic) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type ReopenGeneralForumTopic struct {
	*CloseGeneralForumTopic
}

type HideGeneralForumTopic struct {
	*CloseGeneralForumTopic
}

type UnhideGeneralForumTopic struct {
	*CloseGeneralForumTopic
}

type UnpinAllGeneralForumTopicMessages struct {
	*CloseGeneralForumTopic
}

type SetMyDescriptionArgs struct {
	/*New bot description; 0-512 characters. Pass an empty string to remove the dedicated description for the given language.*/
	Description string `json:"description"`
	/*A two-letter ISO 639-1 language code. If empty, the description will be applied to all users for whose language there is no dedicated description.*/
	LanguageCode string `json:"language_code"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetMyDescriptionArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetMyDescriptionArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type GetMyDescriptionArgs struct {
	//A two-letter ISO 639-1 language code or an empty string
	LanguageCode string `json:"language_code"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetMyDescriptionArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetMyDescriptionArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetMyShortDescriptionArgs struct {
	/*New short description for the bot; 0-120 characters. Pass an empty string to remove the dedicated short description for the given language.*/
	ShortDescription string `json:"short_description"`
	/*A two-letter ISO 639-1 language code. If empty, the description will be applied to all users for whose language there is no dedicated description.*/
	LanguageCode string `json:"language_code"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetMyShortDescriptionArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetMyShortDescriptionArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type SetMyNameArgs struct {
	Name         string `json:"name"`
	LanguageCode string `json:"language_code"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *SetMyNameArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *SetMyNameArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}

type GetMyNameArgs struct {
	LanguageCode string `json:"language_code"`
}

// ToJson converts this strcut into json to be sent to the API server.
func (args *GetMyNameArgs) ToJson() []byte {
	bt, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	return bt
}

// ToMultiPart converts this strcut into HTTP multipart form to be sent to the API server.
func (args *GetMyNameArgs) ToMultiPart(wr *mp.Writer) {
	//This method arguments are never passed as multipart
}
