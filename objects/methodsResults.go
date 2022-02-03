package objects

import "encoding/json"

/*UpdateResult represents the response of getUpdates method*/
type UpdateResult struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}

//FailureResult represents a failure response that has "ok : false" field.
type FailureResult struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

/*SendMethodsResult represents the response of the methods which send a message. (sendMessage,sendPhoto,...)*/
type SendMethodsResult struct {
	Ok     bool     `json:"ok"`
	Result *Message `json:"result"`
}

/*SendMediaGroupMethodResult represents the response of "sendMediaGroup" method*/
type SendMediaGroupMethodResult struct {
	Ok     bool      `json:"ok"`
	Result []Message `json:"result"`
}

//DefaultResult represents the unparsed response of each method.
type DefaultResult struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

//StringResult represents a response that contains string reponse.
type StringResult struct {
	Ok     bool   `json:"ok"`
	Result string `json:"result"`
}

//IntResult represents a response that contains integer reponse.
type IntResult struct {
	Ok     bool `json:"ok"`
	Result int  `json:"result"`
}

//ChatInviteLinkResult represents a response that contains ChatInviteLink object.
type ChatInviteLinkResult struct {
	Ok     bool            `json:"ok"`
	Result *ChatInviteLink `json:"result"`
}

//ChatResult represents a response that contains Chat object.
type ChatResult struct {
	Ok     bool  `json:"ok"`
	Result *Chat `json:"result"`
}

//ChatAdministratorsResult represents a response that contains []ChatMemberOwner object.
type ChatAdministratorsResult struct {
	Ok     bool              `json:"ok"`
	Result []ChatMemberOwner `json:"result"`
}

//LogicalResult represents a response that contains boolean response.
type LogicalResult struct {
	Ok     bool `json:"ok"`
	Result bool `json:"result"`
}

//ProfilePhototsResult represents a response that contains UserProfilePhotos object.
type ProfilePhototsResult struct {
	Ok     bool               `json:"ok"`
	Result *UserProfilePhotos `json:"result"`
}

//GetFileResult represents a response that contains File object.
type GetFileResult struct {
	Ok     bool  `json:"ok"`
	Result *File `json:"result"`
}

//GetCommandsResult represents a response that contains []BotCommand object.
type GetCommandsResult struct {
	Ok     bool         `json:"ok"`
	Result []BotCommand `json:"result"`
}

//PollResult represents a response that contains Poll object.
type PollResult struct {
	Ok     bool  `json:"ok"`
	Result *Poll `json:"result"`
}

//StickerSetResult represents a response that contains StickerSet object.
type StickerSetResult struct {
	Ok     bool        `json:"ok"`
	Result *StickerSet `json:"result"`
}

//GameHighScoresResult represents a response that contains GameHighScores object.
type GameHighScoresResult struct {
	Ok     bool            `json:"ok"`
	Result []GameHighScore `json:"result"`
}

//UserResult represents a response that contains User object.
type UserResult struct {
	Ok     bool  `json:"ok"`
	Result *User `json:"result"`
}

//WebhookInfoResult represents a response that contains WebhookInfo object.
type WebhookInfoResult struct {
	Ok     bool         `json:"ok"`
	Result *WebhookInfo `json:"result"`
}
