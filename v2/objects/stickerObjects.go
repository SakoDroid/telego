package objects

/*Sticker represents a sticker.*/
type Sticker struct {
	/*Identifier for this file, which can be used to download or reuse the file*/
	FileId string `json:"file_id"`
	/*Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.*/
	FileUniqueId string `json:"file_unique_id"`
	/*Type of the sticker, currently one of “regular”, “mask”, “custom_emoji”. The type of the sticker is independent from its format, which is determined by the fields is_animated and is_video.*/
	Type string `json:"type"`
	/*Sticker width*/
	Width int `json:"width"`
	/*Sticker height*/
	Height int `json:"height"`
	/*True, if the sticker is animated*/
	IsAnimated bool `json:"is_animated"`
	/*True, if the sticker is a video sticker*/
	IsVideo bool `json:"is_video"`
	/*Optional. Sticker thumbnail in the .WEBP or .JPG format*/
	Thumb *PhotoSize `json:"thumbnail,omitempty"`
	/*Optional. Emoji associated with the sticker*/
	Emoji string `json:"emoji,omitempty"`
	/*Optional. Name of the sticker set to which the sticker belongs*/
	SetName string `json:"set_name,omitempty"`
	/*Optional. Premium animation for the sticker, if the sticker is premium*/
	PremiumAnimation *File `json:"premium_animation,omitempty"`
	/*Optional. For mask stickers, the position where the mask should be placed*/
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	/*Optional. For custom emoji stickers, unique identifier of the custom emoji*/
	CustomEmojiId string `json:"custom_emoji_id"`
	/*Optional. True, if the sticker must be repainted to a text color in messages, the color of the Telegram Premium badge in emoji status, white color on chat photos, or another appropriate color in other places*/
	NeedsRepainting bool `json:"needs_repainting"`
	/*Optional. File size in bytes*/
	FileSize int `json:"file_size,omitempty"`
}

/*StickerSet represents a sticker set.*/
type StickerSet struct {
	/*Sticker set name*/
	Name string `json:"name"`
	/*Sticker set title*/
	Title string `json:"title"`
	/*Type of stickers in the set, currently one of “regular”, “mask”, “custom_emoji”*/
	StickerType string `json:"sticker_type"`
	/*True, if the sticker set contains animated stickers*/
	IsAnimated bool `json:"is_animated"`
	/*True, if the sticker set contains video stickers*/
	IsVideo bool `json:"is_video"`
	/*True, if the sticker set contains masks
	Caution : The field contains_masks has been removed from the documentation of the class StickerSet. The field is still returned in the object for backward compatibility, but new bots should use the field sticker_type instead.
	*/
	ContainsMask bool `json:"contains_mask"`
	/*List of all set stickers*/
	Stickers []Sticker `json:"stickers"`
	/*Optional. Sticker set thumbnail in the .WEBP or .TGS format*/
	Thumb *PhotoSize `json:"thumbnail,omitempty"`
}

/*MaskPosition describes the position on faces where a mask should be placed by default.*/
type MaskPosition struct {
	/*The part of the face relative to which the mask should be placed. One of “forehead”, “eyes”, “mouth”, or “chin”.*/
	Point string `json:"point"`
	/*Shift by X-axis measured in widths of the mask scaled to the face size, from left to right. For example, choosing -1.0 will place mask just to the left of the default mask position.*/
	XShift float32 `json:"x_shift"`
	/*Shift by Y-axis measured in heights of the mask scaled to the face size, from top to bottom. For example, 1.0 will place the mask just below the default mask position.*/
	YShift float32 `json:"y_shift"`
	/*Mask scaling coefficient. For example, 2.0 means double size.*/
	Scale float32 `json:"scale"`
}

type InputSticker struct {
	Sticker string `json:"sticker"`
	/*List of 1-20 emoji associated with the sticker*/
	EmojiList []string `json:"emoji_list"`
	/*Optional. Position where the mask should be placed on faces. For “mask” stickers only.*/
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	/*Optional. List of 0-20 search keywords for the sticker with total length of up to 64 characters. For “regular” and “custom_emoji” stickers only.*/
	KeyWords []string `json:"keywords"`
}
