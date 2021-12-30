package objects

type InlineQueryResult interface {
	GetResultType() string
}

/*This object represents an incoming inline query. When the user sends an empty query, your bot could return some default or trending results.*/
type InlineQuery struct {
	/*Unique identifier for this query*/
	Id string `json:"id"`
	/*Sender*/
	From *User `json:"from"`
	/*Text of the query (up to 256 characters)*/
	Query string `json:"query"`
	/*Offset of the results to be returned, can be controlled by the bot*/
	Offset string `json:"offset"`
	/*Optional. Type of the chat, from which the inline query was sent.Can be either “sender” for a private chat with the inline query sender, “private”, “group”, “supergroup”, or “channel”. The chat type should be always known for requests sent from official clients and most third-party clients, unless the request was sent from a secret chat*/
	ChatType string `json:"chat_type,omitempty"`
	/*Optional. Sender location, only for bots that request user location*/
	Location *Location `json:"location,omitempty"`
}

/*This object should not be used at all.*/
type InlineQueryResultDefault struct {
	/*Type of the result*/
	Type string `json:"type"`
	/*Unique identifier for this result, 1-64 Bytes*/
	Id string `json:"id"`
	/*Title of the result*/
	Title string `json:"title"`
}

/*Returns the type of the result. The returned value is one of the following :
CachedAudio
CachedDocument
CachedGif
CachedMpeg4Gif
CachedPhoto
CachedSticker
CachedVideo
CachedVoice
Article
Audio
Contact
Game
Document
Gif
Location
Mpeg4Gif
Photo
Venue
Video
Voice*/
func (i *InlineQueryResultDefault) GetResultType() string {
	return i.Type
}

/*Represents a link to an article or web page.*/
type InlineQueryResultArticle struct {
	InlineQueryResultDefault
	/*Content of the message to be sent*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. URL of the result*/
	URL string `json:"url,omitempty"`
	/*Optional. Pass True, if you don't want the URL to be shown in the message*/
	HideURL bool `json:"hide_url,omitempty"`
	/*Optional. Short description of the result*/
	Description string `json:"description,omitempty"`
	/*Optional. Url of the thumbnail for the result*/
	ThumbURL string `json:"thumb_url,omitempty"`
	/*Optional. Thumbnail width*/
	ThumbWidth int `json:"thumb_width,omitempty"`
	/*Optional. Thumbnail height*/
	ThumbHeight int `json:"thumb_height,omitempty"`
}

/*Represents a link to a photo. By default, this photo will be sent by the user with optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.*/
type InlineQueryResultPhoto struct {
	InlineQueryResultDefault
	/*A valid URL of the photo. Photo must be in JPEG format. Photo size must not exceed 5MB*/
	PhotoURL string `json:"photo_url"`
	/*URL of the thumbnail for the photo*/
	ThumbURL string `json:"thumb_url"`
	/*Optional. Width of the photo*/
	PhotoWidth int `json:"photo_width,omitempty"`
	/*Optional. Height of the photo*/
	PhotoHeight int `json:"photo_height,omitempty"`
	/*Optional. Short description of the result*/
	Description string `json:"description,omitempty"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the photo*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to an animated GIF file. By default, this animated GIF file will be sent by the user with optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.*/
type InlineQueryResultGif struct {
	InlineQueryResultDefault
	/*A valid URL for the GIF file. File size must not exceed 1MB*/
	GifURL string `json:"gif_url"`
	/*Optional. Width of the GIF*/
	GifWidth int `json:"gif_width,omitempty"`
	/*Optional. Height of the GIF*/
	GifHeight int `json:"gif_height,omitempty"`
	/*Optional. Duration of the GIF in seconds*/
	GifDuration int `json:"gif_duration,omitempty"`
	/*URL of the static (JPEG or GIF) or animated (MPEG4) thumbnail for the result*/
	ThumbURL string `json:"thumb_url"`
	/*Optional. MIME type of the thumbnail, must be one of “image/jpeg”, “image/gif”, or “video/mp4”. Defaults to “image/jpeg”*/
	ThumbMIMEType string `json:"thumb_mime_type"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the GIF animation*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to a video animation (H.264/MPEG-4 AVC video without sound). By default, this animated MPEG-4 file will be sent by the user with optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.*/
type InlineQueryResultMpeg4Gif struct {
	InlineQueryResultDefault
	/*A valid URL for the MP4 file. File size must not exceed 1MB*/
	Mpeg4URL string `json:"mpeg4_url"`
	/*Optional. Video width*/
	Mpeg4Width int `json:"mpeg4_width,omitempty"`
	/*Optional. Video height*/
	Mpeg4Height int `json:"mpeg4_height,omitempty"`
	/*Optional. Video duration in seconds*/
	Mpeg4Duration int `json:"mpeg4_duration,omitempty"`
	/*URL of the static (JPEG or GIF) or animated (MPEG4) thumbnail for the result*/
	ThumbURL string `json:"thumb_url"`
	/*Optional. MIME type of the thumbnail, must be one of “image/jpeg”, “image/gif”, or “video/mp4”. Defaults to “image/jpeg”*/
	ThumbMIMEType string `json:"thumb_mime_type"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the animation*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to a page containing an embedded video player or a video file. By default, this video file will be sent by the user with an optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the video.

If an InlineQueryResultVideo message contains an embedded video (e.g., YouTube), you must replace its content using input_message_content.*/
type InlineQueryResultVideo struct {
	InlineQueryResultDefault
	/*A valid URL for the embedded video player or video file*/
	VideoURL string `json:"video_url"`
	/*Mime type of the content of video url, “text/html” or “video/mp4”*/
	MIMEType string `json:"mime_type"`
	/*URL of the thumbnail (JPEG only) for the video*/
	ThumbURL string `json:"thumb_url"`
	/*Optional. Video width*/
	VideoWidth int `json:"video_width,omitempty"`
	/*Optional. Video height*/
	VideoHeight int `json:"video_height,omitempty"`
	/*Optional. Video duration in seconds*/
	VideoDuration int `json:"video_duration,omitempty"`
	/*Optional. Short description of the result*/
	Description string `json:"description,omitempty"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the video. This field is required if InlineQueryResultVideo is used to send an HTML-page as a result (e.g., a YouTube video).*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to an MP3 audio file. By default, this audio file will be sent by the user. Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
Note: This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.
*/
type InlineQueryResultAudio struct {
	InlineQueryResultDefault
	/*A valid URL for the audio file*/
	AudioURL string `json:"audio_url"`
	/*Optional. Performer*/
	Performer string `json:"performer,omitempty"`
	/*Optional. Audio duration in seconds*/
	AudioDuration int `json:"audio_duration,omitempty"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the audio*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to a voice recording in an .OGG container encoded with OPUS. By default, this voice recording will be sent by the user. Alternatively, you can use input_message_content to send a message with the specified content instead of the the voice message.
Note: This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.*/
type InlineQueryResultVocie struct {
	InlineQueryResultDefault
	/*A valid URL for the voice recording*/
	VoiceURL string `json:"voice_url"`
	/*Optional. Recording duration in seconds*/
	VoiceDuration int `json:"voice_duration,omitempty"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the voice recording*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to a file. By default, this file will be sent by the user with an optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the file. Currently, only .PDF and .ZIP files can be sent using this method.
Note: This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.*/
type InlineQueryResultDocument struct {
	InlineQueryResultDefault
	/*A valid URL for the file*/
	DocumentURL string `json:"document_url"`
	/*Mime type of the content of the file, either “application/pdf” or “application/zip”*/
	MIMEType string `json:"mime_type"`
	/*Optional. Short description of the result*/
	Description string `json:"description,omitempty"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the file*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	/*Optional. Url of the thumbnail for the result*/
	ThumbURL string `json:"thumb_url,omitempty"`
	/*Optional. Thumbnail width*/
	ThumbWidth int `json:"thumb_width,omitempty"`
	/*Optional. Thumbnail height*/
	ThumbHeight int `json:"thumb_height,omitempty"`
}

/*Represents a location on a map. By default, the location will be sent by the user. Alternatively, you can use input_message_content to send a message with the specified content instead of the location.
Note: This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.*/
type InlineQueryResultLocation struct {
	InlineQueryResultDefault
	/*Location latitude in degrees*/
	Latitude float32 `json:"latitude"`
	/*Location longitude in degrees*/
	Longitude float32 `json:"longitude"`
	/*Optional. The radius of uncertainty for the location, measured in meters; 0-1500*/
	HorizontalAccuracy float32 `json:"horizontal_accuracy,omitempty"`
	/*Optional. Period in seconds for which the location can be updated, should be between 60 and 86400.*/
	LivePeriod int `json:"live_peroid,omitempty"`
	/*Optional. For live locations, a direction in which the user is moving, in degrees. Must be between 1 and 360 if specified.*/
	Heading int `json:"heading,omitempty"`
	/*Optional. For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters. Must be between 1 and 100000 if specified.*/
	ProximityAlertRadius int `json:"proximity_alert_radius,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the location*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	/*Optional. Url of the thumbnail for the result*/
	ThumbURL string `json:"thumb_url,omitempty"`
	/*Optional. Thumbnail width*/
	ThumbWidth int `json:"thumb_width,omitempty"`
	/*Optional. Thumbnail height*/
	ThumbHeight int `json:"thumb_height,omitempty"`
}

/*Represents a venue. By default, the venue will be sent by the user. Alternatively, you can use input_message_content to send a message with the specified content instead of the venue.
Note: This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.*/
type InlineQueryResultVenu struct {
	InlineQueryResultDefault
	/*Latitude of the venue location in degrees*/
	Latitude float32 `json:"latitude"`
	/*Longitude of the venue location in degrees*/
	Longitude float32 `json:"longitude"`
	/*Address of the venue*/
	Address string `json:"address"`
	/*Optional. Foursquare identifier of the venue if known*/
	FourquareId string `json:"fourquare_id,omitempty"`
	/*Optional. Foursquare type of the venue, if known. (For example, “arts_entertainment/default”, “arts_entertainment/aquarium” or “food/icecream”.)*/
	FoursquareType string `json:"fourquare_type,omitempty"`
	/*Optional. Google Places identifier of the venue*/
	GooglePlaceId string `json:"google_place_id,omitempty"`
	/*Optional. Google Places type of the venue.*/
	GoogleplaceType string `json:"google_place_type,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the venue*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	/*Optional. Url of the thumbnail for the result*/
	ThumbURL string `json:"thumb_url,omitempty"`
	/*Optional. Thumbnail width*/
	ThumbWidth int `json:"thumb_width,omitempty"`
	/*Optional. Thumbnail height*/
	ThumbHeight int `json:"thumb_height,omitempty"`
}

/*Represents a contact with a phone number. By default, this contact will be sent by the user. Alternatively, you can use input_message_content to send a message with the specified content instead of the contact.
Note: This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.*/
type InlineQueryResultContact struct {
	InlineQueryResultDefault
	/*Contact's phone number*/
	PhoneNumber string `json:"phone_number"`
	/*Optional. Contact's first name*/
	Firstname string `json:"first_name"`
	/*Optional. Contact's last name*/
	LastName string `json:"last_name,omitempty"`
	/*Optional. Additional data about the contact in the form of a vCard, 0-2048 bytes*/
	Vcard string `json:"vcard,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the contact*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	/*Optional. Url of the thumbnail for the result*/
	ThumbURL string `json:"thumb_url,omitempty"`
	/*Optional. Thumbnail width*/
	ThumbWidth int `json:"thumb_width,omitempty"`
	/*Optional. Thumbnail height*/
	ThumbHeight int `json:"thumb_height,omitempty"`
}

/*Represents a Game.
Note: This will only work in Telegram versions released after October 1, 2016. Older clients will not display any inline results if a game result is among them.*/
type InlineQueryResultGame struct {
	/*Type of the result*/
	Type string `json:"type"`
	/*Unique identifier for this result, 1-64 Bytes*/
	Id string `json:"id"`
	/*Type of the result, must be game*/
	GameShortName string `json:"game_short_name"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (i *InlineQueryResultGame) GetResultType() string {
	return i.Type
}

/*Represents a link to a photo stored on the Telegram servers. By default, this photo will be sent by the user with an optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.*/
type InlineQueryResultCachedPhoto struct {
	InlineQueryResultDefault
	/*A valid file identifier for the photo*/
	PhotoFileId string `json:"photo_file_id"`
	/*Optional. Short description of the result*/
	Description string `json:"description,omitempty"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the photo*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to an animated GIF file stored on the Telegram servers. By default, this animated GIF file will be sent by the user with an optional caption. Alternatively, you can use input_message_content to send a message with specified content instead of the animation.*/
type InlineQueryResultCachedGif struct {
	InlineQueryResultDefault
	/*A valid file identifier for the gif*/
	GifFileId string `json:"gif_file_id"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the GIF animation*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to a video animation (H.264/MPEG-4 AVC video without sound) stored on the Telegram servers. By default, this animated MPEG-4 file will be sent by the user with an optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.*/
type InlineQueryResultCachedMpeg4Gif struct {
	InlineQueryResultDefault
	/*A valid file identifier for the animation*/
	Mpeg4FileId string `json:"mpeg4_file_id"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the animation*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to a sticker stored on the Telegram servers. By default, this sticker will be sent by the user. Alternatively, you can use input_message_content to send a message with the specified content instead of the sticker.
This will only work in Telegram versions released after 9 April, 2016 for static stickers and after 06 July, 2019 for animated stickers. Older clients will ignore them.*/
type InlineQueryResultCachedSticker struct {
	/*Type of the result*/
	Type string `json:"type"`
	/*Unique identifier for this result, 1-64 Bytes*/
	Id string `json:"id"`
	/*A valid file identifier for the sticker*/
	StickerFileId string `json:"sticker_file_id"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the animation*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (i *InlineQueryResultCachedSticker) GetResultType() string {
	return i.Type
}

/*Represents a link to a video file stored on the Telegram servers. By default, this video file will be sent by the user with an optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the video.*/
type InlineQueryResultCachedVideo struct {
	InlineQueryResultDefault
	/*A valid file identifier for the video*/
	VideoFileId string `json:"video_file_id"`
	/*Optional. Short description of the result*/
	Description string `json:"description,omitempty"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the video. This field is required if InlineQueryResultVideo is used to send an HTML-page as a result (e.g., a YouTube video).*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to an MP3 audio file stored on the Telegram servers. By default, this audio file will be sent by the user. Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.*/
type InlineQueryResultCachedAudio struct {
	InlineQueryResultDefault
	/*A valid file identifier for the audio*/
	AudioFileId string `json:"audio_file_id"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the audio*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to a voice message stored on the Telegram servers. By default, this voice message will be sent by the user. Alternatively, you can use input_message_content to send a message with the specified content instead of the voice message.
This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.*/
type InlineQueryResultCachedVocie struct {
	InlineQueryResultDefault
	/*A valid file identifier for the video*/
	VoiceFileId string `json:"voice_file_id"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the voice recording*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

/*Represents a link to a file stored on the Telegram servers. By default, this file will be sent by the user with an optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the file.
Note: This will only work in Telegram versions released after 9 April, 2016. Older clients will ignore them.*/
type InlineQueryResultCachedDocument struct {
	InlineQueryResultDefault
	/*A valid file identifier for the file*/
	DocumentFileId string `json:"document_file_id"`
	/*Optional. Short description of the result*/
	Description string `json:"description,omitempty"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the caption. See formatting options for more details.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Inline keyboard attached to the message*/
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	/*Optional. Content of the message to be sent instead of the file*/
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}
