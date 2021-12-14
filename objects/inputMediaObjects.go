package objects

type InputMedia interface {
	blah()
}

/*This should not be used at all*/
type InputMediaDefault struct {
	/*Type of the result, can be photo,video,animation,audio or document.*/
	Type string `json:"type"`
	/*File to send. Pass a file_id to send a file that exists on the Telegram servers (recommended), pass an HTTP URL for Telegram to get a file from the Internet, or pass “attach://<file_attach_name>” to upload a new one using multipart/form-data under <file_attach_name> name.*/
	Media string `json:"media"`
	/*Optional. Caption of the file to be sent, 0-1024 characters after entities parsing*/
	Caption string `json:"caption,omitempty"`
	/*Optional. Mode for parsing entities in the file caption.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
}

/*Represents a photo to be sent.*/
type InputMediaPhoto struct {
	InputMediaDefault
}

func (is *InputMediaPhoto) blah() {}

/*Represents a video to be sent.*/
type InputMediaVideo struct {
	InputMediaDefault
	/*Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb string `json:"thumb,omitempty"`
	/*Optional. Video width*/
	Width int `json:"width,omitempty"`
	/*Optional. Video height*/
	Height int `json:"height,omitempty"`
	/*Optional. Video duration in seconds*/
	Duration int `json:"duration,omitempty"`
	/*Optional. Pass True, if the uploaded video is suitable for streaming*/
	SupportsStreaming bool `json:"supports_streaming,omitempty"`
}

func (is *InputMediaVideo) blah() {}

/*Represents an animation file (GIF or H.264/MPEG-4 AVC video without sound) to be sent.*/
type InputMediaAnimation struct {
	InputMediaDefault
	/*Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb string `json:"thumb,omitempty"`
	/*Optional. Animation width*/
	Width int `json:"width,omitempty"`
	/*Optional. Animation height*/
	Height int `json:"height,omitempty"`
	/*Optional. Animation duration in seconds*/
	Duration int `json:"duration,omitempty"`
}

func (is *InputMediaAnimation) blah() {}

/*Represents an audio file to be treated as music to be sent.*/
type InputMediaAudio struct {
	InputMediaDefault
	/*Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb string `json:"thumb,omitempty"`
	/*Optional. Animation duration in seconds*/
	Duration int `json:"duration,omitempty"`
	/*Optional. Performer of the audio*/
	Performer string `json:"performer,omitempty"`
	/*Optional. Title of the audio*/
	Title string `json:"title,omitempty"`
}

func (is *InputMediaAudio) blah() {}

/*Represents a general file to be sent.*/
type InputMediaDocument struct {
	InputMediaDefault
	/*Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side. The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320. Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file, so you can pass “attach://<file_attach_name>” if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.*/
	Thumb string `json:"thumb,omitempty"`
	/*Optional. Disables automatic server-side content type detection for files uploaded using multipart/form-data. Always True, if the document is sent as part of an album.*/
	DisableContentTypeDetection bool `json:"disable_content_type_detection,omitempty"`
}

func (is *InputMediaDocument) blah() {}
