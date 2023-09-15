package telego

import objs "github.com/SakoDroid/telego/v2/objects"

type StickerEditor struct {
	bot       *Bot
	stickerId string
}

// SetEmojiList sets a new emoji list for this sticker.
func (se *StickerEditor) SetEmojiList(emojiList []string) (*objs.Result[bool], error) {
	return se.bot.apiInterface.SetStickerEmojiList(se.stickerId, emojiList)
}

// SetKeyword sets a new keyword list for this sticker. Maximum length is 20.
func (se *StickerEditor) SetKeywords(keywords []string) (*objs.Result[bool], error) {
	return se.bot.apiInterface.SetStickerKeywords(se.stickerId, keywords)
}

// SetMaskPosition sets a new mask position for this sticker.
func (se *StickerEditor) SetMaskPosition(maskPosition *objs.MaskPosition) (*objs.Result[bool], error) {
	return se.bot.apiInterface.SetStickerMaskPosition(se.stickerId, maskPosition)
}

// Delete deletes this sticker from the set it belongs to.
func (se *StickerEditor) Delete() (*objs.Result[bool], error) {
	return se.bot.apiInterface.DeleteStickerFromSet(se.stickerId)
}
