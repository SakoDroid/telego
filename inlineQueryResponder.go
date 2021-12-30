package telego

import (
	"errors"

	objs "github.com/SakoDroid/telego/objects"
)

type InlineQueryResponder struct {
	bot                                             *Bot
	id, nextOffset, switchPmText, switchPmParameter string
	results                                         []objs.InlineQueryResult
	cacheTime                                       int
	isPersonal                                      bool
}

/*Adds an article to the result. No more than 50 results are allowed.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing".You can use "Create***Message" method to create this messages.
*/
func (iqs *InlineQueryResponder) AddArticle(id, title, url, description, thumbUrl string, thumbWidth, thumbHeight int, hideUrl bool, message objs.InputMessageContent, keyboard *inlineKeyboard) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultArticle{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "article",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		URL:                 url,
		HideURL:             hideUrl,
		Description:         description,
		ThumbURL:            thumbUrl,
		ThumbWidth:          thumbWidth,
		ThumbHeight:         thumbHeight,
	})
	return nil
}

/*Adds a photo to the result. No more than 50 results are allowed

Represents a link to a photo. By default, this photo will be sent by the user with optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing".You can use "Create***Message" method to create this messages.
*/
func (iqs *InlineQueryResponder) AddPhoto(id, title, photoURL, description, caption, parseMode, thumbUrl string, photoWidth, photoHeight int, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultPhoto{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "photo",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		PhotoURL:            photoURL,
		Description:         description,
		ThumbURL:            thumbUrl,
		PhotoWidth:          photoWidth,
		PhotoHeight:         photoHeight,
	})
	return nil
}

/*Adds a gif to the result. No more than 50 results are allowed


Represents a link to an animated GIF file. By default, this animated GIF file will be sent by the user with optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing".You can use "Create***Message" method to create this messages.
*/
func (iqs *InlineQueryResponder) AddGif(id, title, gifURL, caption, parseMode, thumbUrl, thumbMIMEType string, gifWidth, gifHeight, gifDuration int, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultGif{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "gif",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		ThumbURL:            thumbUrl,
		ThumbMIMEType:       thumbMIMEType,
		GifURL:              gifURL,
		GifWidth:            gifWidth,
		GifHeight:           gifHeight,
		GifDuration:         gifDuration,
	})
	return nil
}

/*Adds a mpeg4 to the result. No more than 50 results are allowed


Represents a link to a video animation (H.264/MPEG-4 AVC video without sound). By default, this animated MPEG-4 file will be sent by the user with optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing".You can use "Create***Message" method to create this messages.
*/
func (iqs *InlineQueryResponder) AddMpeg4Gif(id, title, mpeg4URL, caption, parseMode, thumbUrl, thumbMIMEType string, mpeg4Width, mpeg4Height, mpeg4Duration int, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultMpeg4Gif{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "mpeg4_gif",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		ThumbURL:            thumbUrl,
		ThumbMIMEType:       thumbMIMEType,
		Mpeg4URL:            mpeg4URL,
		Mpeg4Width:          mpeg4Width,
		Mpeg4Height:         mpeg4Height,
		Mpeg4Duration:       mpeg4Duration,
	})
	return nil
}

/*Adds a video to the result. No more than 50 results are allowed


Represents a link to a page containing an embedded video player or a video file. By default, this video file will be sent by the user with an optional caption. Alternatively, you can use input_message_content to send a message with the specified content instead of the video.

If an InlineQueryResultVideo message contains an embedded video (e.g., YouTube), you must replace its content using message field.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing".You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddVideo(id, title, videoURL, mimeType, caption, description, parseMode, thumbUrl string, videoWidth, videoHeight, videoDuration int, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultVideo{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "video",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		ThumbURL:            thumbUrl,
		VideoURL:            videoURL,
		MIMEType:            mimeType,
		VideoWidth:          videoHeight,
		VideoHeight:         videoHeight,
		VideoDuration:       videoDuration,
		Description:         description,
	})
	return nil
}

/*Adds an audio to the result. No more than 50 results are allowed

Represents a link to an MP3 audio file. By default, this audio file will be sent by the user. Alternatively, you can use message field to send a message with the specified content instead of the audio.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing".You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddAudio(id, title, audioURL, caption, parseMode, performer string, audioDuration int, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultAudio{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "audio",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		AudioURL:            audioURL,
		Performer:           performer,
		AudioDuration:       audioDuration,
	})
	return nil
}

/*Adds a voice to the result. No more than 50 results are allowed


Represents a link to a voice recording in an .OGG container encoded with OPUS. By default, this voice recording will be sent by the user. Alternatively, you can use message field to send a message with the specified content instead of the the voice message.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing".You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddVoice(id, title, voiceURL, caption, parseMode string, voiceDuration int, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultVocie{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "voice",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		VoiceURL:            voiceURL,
		VoiceDuration:       voiceDuration,
	})
	return nil
}

/*Adds a document to the result. No more than 50 results are allowed


Represents a link to a file. By default, this file will be sent by the user with an optional caption. Alternatively, you can use message field to send a message with the specified content instead of the file. Currently, only .PDF and .ZIP files can be sent using this method.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddDocument(id, title, documentURL, mimeType, description, thumbUrl, caption, parseMode string, thumbWidth, thumbHeight int, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultDocument{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "document",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		DocumentURL:         documentURL,
		MIMEType:            mimeType,
		Description:         description,
		ThumbURL:            thumbUrl,
		ThumbWidth:          thumbWidth,
		ThumbHeight:         thumbHeight,
	})
	return nil
}

/*Adds a location to the result. No more than 50 results are allowed

Represents a location on a map. By default, the location will be sent by the user. Alternatively, you can use message field to send a message with the specified content instead of the location.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddLocation(id, title, thumbUrl string, latitude, longitude, horizontalAccuracy float32, livePeriod, heading, proximityAlertRadius, thumbWidth, thumbHeight int, message objs.InputMessageContent, keyboard *inlineKeyboard) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultLocation{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "location",
			Id:    id,
			Title: title,
		},
		InputMessageContent:  message,
		ReplyMarkup:          iqs.fixTheKeyboard(keyboard),
		Latitude:             latitude,
		Longitude:            longitude,
		HorizontalAccuracy:   horizontalAccuracy,
		LivePeriod:           livePeriod,
		Heading:              heading,
		ProximityAlertRadius: proximityAlertRadius,
		ThumbURL:             thumbUrl,
		ThumbWidth:           thumbWidth,
		ThumbHeight:          thumbHeight,
	})
	return nil
}

/*Adds a venue to the result. No more than 50 results are allowed


Represents a venue. By default, the venue will be sent by the user. Alternatively, you can use message field to send a message with the specified content instead of the venue.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddVenue(id, title, thumbUrl string, latitude, longitude float32, address, foursquareId, foursquareType, googlePlaceId, googlePlaceType string, thumbWidth, thumbHeight int, message objs.InputMessageContent, keyboard *inlineKeyboard) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultVenu{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "venue",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Latitude:            latitude,
		Longitude:           longitude,
		Address:             address,
		FourquareId:         foursquareId,
		FoursquareType:      foursquareType,
		GooglePlaceId:       googlePlaceId,
		GoogleplaceType:     googlePlaceType,
		ThumbURL:            thumbUrl,
		ThumbWidth:          thumbWidth,
		ThumbHeight:         thumbHeight,
	})
	return nil
}

/*Adds a contact to the result. No more than 50 results are allowed


Represents a contact with a phone number. By default, this contact will be sent by the user. Alternatively, you can use message field to send a message with the specified content instead of the contact.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddContact(id, title, thumbUrl, phoneNumber, firstName, lastName, vCard string, thumbWidth, thumbHeight int, message objs.InputMessageContent, keyboard *inlineKeyboard) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultContact{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "contact",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		PhoneNumber:         phoneNumber,
		Firstname:           firstName,
		LastName:            lastName,
		Vcard:               vCard,
		ThumbURL:            thumbUrl,
		ThumbWidth:          thumbWidth,
		ThumbHeight:         thumbHeight,
	})
	return nil
}

/*Adds a game to the result. No more than 50 results are allowed

Represents a game*/
func (iqs *InlineQueryResponder) AddGame(id, gameShortName string, keyboard *inlineKeyboard) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultGame{
		Type:          "game",
		Id:            id,
		GameShortName: gameShortName,
		ReplyMarkup:   iqs.fixTheKeyboard(keyboard),
	})
	return nil
}

/*Adds a cached photo to the result. No more than 50 results are allowed

Represents a link to a photo stored on the Telegram servers. By default, this photo will be sent by the user with an optional caption. Alternatively, you can use message field to send a message with the specified content instead of the photo.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddCachedPhoto(id, title, photoFileId, description, caption, parseMode string, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultCachedPhoto{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "photo",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		PhotoFileId:         photoFileId,
		Description:         description,
	})
	return nil
}

/*Adds a cached gif to the result. No more than 50 results are allowed

Represents a link to an animated GIF file stored on the Telegram servers. By default, this animated GIF file will be sent by the user with an optional caption. Alternatively, you can use message field to send a message with specified content instead of the animation.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddCachedGif(id, title, gifFileId, caption, parseMode string, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultCachedGif{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "gif",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		GifFileId:           gifFileId,
	})
	return nil
}

/*Adds a cached mpeg4 to the result. No more than 50 results are allowed.

Represents a link to a video animation (H.264/MPEG-4 AVC video without sound) stored on the Telegram servers. By default, this animated MPEG-4 file will be sent by the user with an optional caption. Alternatively, you can use message field to send a message with the specified content instead of the animation.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddCachedMpeg4Gif(id, title, mpeg4FileId, caption, parseMode string, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultCachedMpeg4Gif{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "mpeg4_gif",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		Mpeg4FileId:         mpeg4FileId,
	})
	return nil
}

/*Adds a cached mpeg4 to the result. No more than 50 results are allowed.

Represents a link to a sticker stored on the Telegram servers. By default, this sticker will be sent by the user. Alternatively, you can use message field to send a message with the specified content instead of the sticker.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddCachedSticker(id, stickerFileId string, message objs.InputMessageContent, keyboard *inlineKeyboard) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultCachedSticker{
		Type:                "sticker",
		Id:                  id,
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		StickerFileId:       stickerFileId,
	})
	return nil
}

/*Adds a cached document to the result. No more than 50 results are allowed

Represents a link to a file stored on the Telegram servers. By default, this file will be sent by the user with an optional caption. Alternatively, you can use message field to send a message with the specified content instead of the file.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddCachedDocument(id, title, documentFileId, description, caption, parseMode string, messsage objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultCachedDocument{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "document",
			Id:    id,
			Title: title,
		},
		InputMessageContent: messsage,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		DocumentFileId:      documentFileId,
		Description:         description,
	})
	return nil
}

/*Adds a cached video to the result. No more than 50 results are allowed

Represents a link to a video file stored on the Telegram servers. By default, this video file will be sent by the user with an optional caption. Alternatively, you can use message field to send a message with the specified content instead of the video.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddCachedVideo(id, title, videoFileId, caption, description, parseMode string, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultCachedVideo{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "video",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		VideoFileId:         videoFileId,
		Description:         description,
	})
	return nil
}

/*Adds an audio to the result. No more than 50 results are allowed

Represents a link to an MP3 audio file stored on the Telegram servers. By default, this audio file will be sent by the user. Alternatively, you can use message field to send a message with the specified content instead of the audio.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddCachedAudio(id, title, audioFileId, caption, parseMode string, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultCachedAudio{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "audio",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		AudioFileId:         audioFileId,
	})
	return nil
}

/*Adds a voice to the result. No more than 50 results are allowed

Represents a link to a voice message stored on the Telegram servers. By default, this voice message will be sent by the user. Alternatively, you can use message field to send a message with the specified content instead of the voice message.

"message" argument is the message that will be sent if this result is pressed. It is necessary and if it's not passed the API server will return "400 bad request, message_text is missing". You can use "Create***Message" method to create this messages.*/
func (iqs *InlineQueryResponder) AddCachedVoice(id, title, voiceFileId, caption, parseMode string, message objs.InputMessageContent, keyboard *inlineKeyboard, captionEntities []objs.MessageEntity) error {
	if len(iqs.results) >= 50 {
		return errors.New("cant add more than 50 results")
	}
	iqs.results = append(iqs.results, &objs.InlineQueryResultCachedVocie{
		InlineQueryResultDefault: objs.InlineQueryResultDefault{
			Type:  "voice",
			Id:    id,
			Title: title,
		},
		InputMessageContent: message,
		ReplyMarkup:         iqs.fixTheKeyboard(keyboard),
		Caption:             caption,
		CaptionEntities:     captionEntities,
		ParseMode:           parseMode,
		VoiceFileId:         voiceFileId,
	})
	return nil
}

/*Use this function to create a text message to be passed as "message" argument in "Add" methods.*/
func (iqs *InlineQueryResponder) CreateTextMessage(text, parseMode string, entities []objs.MessageEntity, disabelWebPagePreview bool) objs.InputMessageContent {
	return &objs.InputTextMessageContent{
		MessageText:           text,
		ParseMode:             parseMode,
		Entities:              entities,
		DisableWebPagePreview: disabelWebPagePreview,
	}
}

/*Use this function to create a location message to be passed as "message" argument in "Add" methods.*/
func (iqs *InlineQueryResponder) CreateLocationMessage(latitude, longitude, accuracy float32, liveperiod, heading, proximityAlertRadius int) objs.InputMessageContent {
	return &objs.InputLocationMessageContent{
		Latitude:             latitude,
		Longitude:            longitude,
		HorizontalAccuracy:   accuracy,
		LivePeriod:           liveperiod,
		Heading:              heading,
		ProximityAlertRadius: proximityAlertRadius,
	}
}

/*Use this function to create a venue message to be passed as "message" argument in "Add" methods.*/
func (iqs *InlineQueryResponder) CreateVenueMessage(latitude, longitude float32, title, address, foursquareId, foursquareType, googlePlaceId, googlePlaceType string) objs.InputMessageContent {
	return &objs.InputVenueMessageContent{
		Latitude:        latitude,
		Longitude:       longitude,
		Title:           title,
		Address:         address,
		FoursquareId:    foursquareId,
		FoursquareType:  foursquareType,
		GooglePlaceId:   googlePlaceId,
		GooglePlaceType: googlePlaceType,
	}
}

/*Use this function to create a contact message to be passed as "message" argument in "Add" methods.*/
func (iqs *InlineQueryResponder) CreateContactMessage(phoneNumber, firstName, lastName, vCard string) objs.InputMessageContent {
	return &objs.InputContactMessageContent{
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
		LastName:    lastName,
		Vcard:       vCard,
	}
}

/*Use this function to create an invoice message to be passed as "message" argument in "Add" methods.*/
func (iqs *InlineQueryResponder) CreateInvoiceMessage(invoice *Invoice) objs.InputMessageContent {
	return &objs.InputInvoiceMessageContent{
		Title:                     invoice.title,
		Description:               invoice.description,
		Payload:                   invoice.payload,
		ProviderToken:             invoice.providerToken,
		Currency:                  invoice.currency,
		Prices:                    invoice.prices,
		MaxTipAmount:              invoice.maxTipAmount,
		SuggestedTipAmounts:       invoice.suggestedTipAmounts,
		ProviderData:              invoice.providerData,
		PhotoURL:                  invoice.photoURL,
		PhotoSize:                 invoice.photoSize,
		PhotoWidth:                invoice.photoWidth,
		PhotoHeight:               invoice.photoHeight,
		NeedName:                  invoice.needName,
		NeedPhoneNumber:           invoice.needPhoneNumber,
		NeedEmail:                 invoice.needEmail,
		NeedShippingAddress:       invoice.needShippingAddress,
		SendPhoneNumberToProvider: invoice.sendPhoneNumberToProvider,
		SendEmailToProvider:       invoice.sendEmailToProvider,
		IsFlexible:                invoice.isFlexible,
	}
}

/*Sends this answer to the user*/
func (iqs *InlineQueryResponder) Send() (*objs.LogicalResult, error) {
	return iqs.bot.apiInterface.AnswerInlineQuery(
		iqs.id, iqs.results, iqs.cacheTime, iqs.isPersonal, iqs.nextOffset,
		iqs.switchPmText, iqs.switchPmParameter,
	)
}

func (iqs *InlineQueryResponder) fixTheKeyboard(keyboard *inlineKeyboard) *objs.InlineKeyboardMarkup {
	var replyMarkup *objs.InlineKeyboardMarkup
	if keyboard != nil {
		kb := keyboard.toInlineKeyboardMarkup()
		replyMarkup = &kb
	}
	return replyMarkup
}
