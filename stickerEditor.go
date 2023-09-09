package telego

import objs "github.com/SakoDroid/telego/objects"

type stickerEditor struct {
	bot       *Bot
	stickerId string
}

// SetEmojiList sets a new emoji list for this sticker.
func (se *stickerEditor) SetEmojiList(emojiList []string) (*objs.LogicalResult, error) {
	return se.bot.apiInterface.SetStickerEmojiList(se.stickerId, emojiList)
}

// SetKeyword sets a new keyword list for this sticker. Maximum length is 20.
func (se *stickerEditor) SetKeywords(keywords []string) (*objs.LogicalResult, error) {
	return se.bot.apiInterface.SetStickerKeywords(se.stickerId, keywords)
}

// SetMaskPosition sets a new mask position for this sticker.
func (se *stickerEditor) SetMaskPosition(maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	return se.bot.apiInterface.SetStickerMaskPosition(se.stickerId, maskPosition)
}

// Delete deletes this sticker from the set it belongs to.
func (se *stickerEditor) Delete() (*objs.LogicalResult, error) {
	return se.bot.apiInterface.DeleteStickerFromSet(se.stickerId)
}
