package telego

import (
	"os"

	logger "github.com/SakoDroid/telego/logger"
	objs "github.com/SakoDroid/telego/objects"
)

type StickerSet struct {
	bot        *Bot
	stickerSet *objs.StickerSet
}

/*Updates this sticker set*/
func (ss *StickerSet) update() {
	res, err := ss.bot.apiInterface.GetStickerSet(ss.stickerSet.Name)
	if err != nil {
		logger.Logger.Println("Error while updating sticker set.", err.Error())
	} else {
		ss.stickerSet = res.Result
	}
}

/*Returns the title of this sticker set*/
func (ss *StickerSet) GetTitle() string {
	return ss.stickerSet.Title
}

/*Returnes the name of this sticker set*/
func (ss *StickerSet) GetName() string {
	return ss.stickerSet.Name
}

/*Returns the sticker in this sticker set.*/
func (ss *StickerSet) GetStickers() []objs.Sticker {
	ss.update()
	return ss.stickerSet.Stickers
}

/*Returns the thumbnail of this sticker set*/
func (ss *StickerSet) GetThumb() *objs.PhotoSize {
	ss.update()
	return ss.stickerSet.Thumb
}

/*Adds a sticker to the current set

Use this method to add a new sticker to a set created by the bot. You must use exactly one of the fields png_sticker or tgs_sticker. Animated stickers can be added to animated sticker sets and only to them. Animated sticker sets can have up to 50 stickers. Static sticker sets can have up to 120 stickers. Returns True on success.

png sticker can be passed as an file id or url (pngStickerFileIdOrUrl) or file(pngStickerFile).*/
func (ss *StickerSet) AddSticker(userId int, pngStickerFileIdOrUrl string, pngStickerFile *os.File, tgsSticker *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	if tgsSticker == nil {
		if pngStickerFile == nil {
			return ss.bot.apiInterface.AddStickerToSet(
				userId, ss.stickerSet.Name, pngStickerFileIdOrUrl, "", emojies, maskPosition, nil,
			)
		} else {
			stat, er := pngStickerFile.Stat()
			if er != nil {
				return nil, er
			}
			return ss.bot.apiInterface.AddStickerToSet(
				userId, ss.stickerSet.Name, "attach://"+stat.Name(), "", emojies, maskPosition, pngStickerFile,
			)
		}
	} else {
		stat, er := tgsSticker.Stat()
		if er != nil {
			return nil, er
		}
		return ss.bot.apiInterface.AddStickerToSet(
			userId, ss.stickerSet.Name, "", "attach://"+stat.Name(), emojies, maskPosition, tgsSticker,
		)
	}
}

/*
Use this method to move a sticker in a set created by the bot to a specific position. Returns True on success.

"sticker" is file identifier of the sticker and "position" is new sticker position in the set, zero-based*/
func (ss *StickerSet) SetStickerPosition(sticker string, position int) (*objs.LogicalResult, error) {
	return ss.bot.apiInterface.SetStickerPositionInSet(sticker, position)
}

/*
Use this method to delete a sticker from a set created by the bot. Returns True on success.

"sticker" is file identifier of the sticker.*/
func (ss *StickerSet) DeleteStickerFromSet(sticker string) (*objs.LogicalResult, error) {
	return ss.bot.apiInterface.DeleteStickerFromSet(sticker)
}

/*Use this method to set the thumbnail of a sticker set. Animated thumbnails can be set for animated sticker sets only. Returns True on success.*/
func (ss *StickerSet) SetThumb(userId int, thumb string) (*objs.LogicalResult, error) {
	return ss.bot.apiInterface.SetStickerSetThumb(ss.stickerSet.Name, thumb, userId, nil)
}

/*Use this method to set the thumbnail of a sticker set. Animated thumbnails can be set for animated sticker sets only. Returns True on success.*/
func (ss *StickerSet) SetThumbByFile(userId int, thumb *os.File) (*objs.LogicalResult, error) {
	stats, err := thumb.Stat()
	if err != nil {
		return nil, err
	}
	return ss.bot.apiInterface.SetStickerSetThumb(ss.stickerSet.Name, "attach://"+stats.Name(), userId, thumb)
}
