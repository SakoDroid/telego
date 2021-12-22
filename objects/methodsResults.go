package objects

import "encoding/json"

/*This object represents the response of getUpdates method*/
type UpdateResult struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}

type FailureResult struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

/*This object represents the reponse of the methods which send a message. (sendMessage,sendPhoto,...)*/
type SendMethodsResult struct {
	Ok     bool     `json:"ok"`
	Result *Message `json:"result"`
}

/*This object represents the reponse of "sendMediaGroup" method*/
type SendMediaGroupMethodResult struct {
	Ok     bool      `json:"ok"`
	Result []Message `json:"result"`
}

type DefaultResult struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

type StringResult struct {
	Ok     bool   `json:"ok"`
	Result string `json:"result"`
}

type IntResult struct {
	Ok     bool `json:"ok"`
	Result int  `json:"result"`
}

type ChatInviteLinkResult struct {
	Ok     bool            `json:"ok"`
	Result *ChatInviteLink `json:"result"`
}

type ChatResult struct {
	Ok     bool  `json:"ok"`
	Result *Chat `json:"result"`
}

type ChatAdministratorsResult struct {
	Ok     bool              `json:"ok"`
	Result []ChatMemberOwner `json:"result"`
}

type LogicalResult struct {
	Ok     bool `json:"ok"`
	Result bool `json:"result"`
}

type ProfilePhototsResult struct {
	Ok     bool               `json:"ok"`
	Result *UserProfilePhotos `json:"result"`
}

type GetFileResult struct {
	Ok     bool  `json:"ok"`
	Result *File `json:"result"`
}

type GetCommandsResult struct {
	Ok     bool         `json:"ok"`
	Result []BotCommand `json:"result"`
}

type PollResult struct {
	Ok     bool  `json:"ok"`
	Result *Poll `json:"result"`
}

type StickerSetResult struct {
	Ok     bool        `json:"ok"`
	Result *StickerSet `json:"result"`
}

type GameHighScoresResult struct {
	Ok     bool            `json:"ok"`
	Result []GameHighScore `json:"result"`
}

type UserResult struct {
	Ok     bool  `json:"ok"`
	Result *User `json:"result"`
}
