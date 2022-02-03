package telego

import (
	"errors"
	"os"

	logger "github.com/SakoDroid/telego/logger"
	objs "github.com/SakoDroid/telego/objects"
)

//StickerSet is a set of stickers.
type StickerSet struct {
	bot        *Bot
	stickerSet *objs.StickerSet
	userId     int
}

/*update updates this sticker set*/
func (ss *StickerSet) update() {
	if ss != nil {
		res, err := ss.bot.apiInterface.GetStickerSet(ss.stickerSet.Name)
		if err != nil {
			logger.Logger.Println("Error while updating sticker set.", err.Error())
		} else {
			ss.stickerSet = res.Result
		}
	}
}

/*GetTitle returns the title of this sticker set*/
func (ss *StickerSet) GetTitle() string {
	if ss == nil {
		return ""
	}
	return ss.stickerSet.Title
}

/*GetName returnes the name of this sticker set*/
func (ss *StickerSet) GetName() string {
	if ss == nil {
		return ""
	}
	return ss.stickerSet.Name
}

/*GetStickers returns the sticker in this sticker set.*/
func (ss *StickerSet) GetStickers() []objs.Sticker {
	if ss == nil {
		return nil
	}
	ss.update()
	return ss.stickerSet.Stickers
}

/*GetThumb returns the thumbnail of this sticker set*/
func (ss *StickerSet) GetThumb() *objs.PhotoSize {
	if ss == nil {
		return nil
	}
	ss.update()
	return ss.stickerSet.Thumb
}

/*Deprecated: This function should no longer be used for adding stickers to a sticker set. It has been preserved for backward compatibility and will be removed in next versions. Use "AddPngSticker","AddAnimatedSticker" or "AddVideoSticker" methods instead.

Adds a sticker to the current set
Use this method to add a new sticker to a set created by the bot. You must use exactly one of the fields png_sticker or tgs_sticker. Animated stickers can be added to animated sticker sets and only to them. Animated sticker sets can have up to 50 stickers. Static sticker sets can have up to 120 stickers. Returns True on success.
png sticker can be passed as an file id or url (pngStickerFileIdOrUrl) or file(pngStickerFile).*/
func (ss *StickerSet) AddSticker(pngStickerFileIdOrUrl string, pngStickerFile *os.File, tgsSticker *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	if ss == nil {
		return nil, errors.New("sticker set is nil")
	}
	if tgsSticker == nil {
		if pngStickerFile == nil {
			if pngStickerFileIdOrUrl == "" {
				return nil, errors.New("wrong file id or url")
			}
			return ss.bot.apiInterface.AddStickerToSet(
				ss.userId, ss.stickerSet.Name, pngStickerFileIdOrUrl, "", "", emojies, maskPosition, nil,
			)
		} else {
			stat, er := pngStickerFile.Stat()
			if er != nil {
				return nil, er
			}
			return ss.bot.apiInterface.AddStickerToSet(
				ss.userId, ss.stickerSet.Name, "attach://"+stat.Name(), "", "", emojies, maskPosition, pngStickerFile,
			)
		}
	} else {
		stat, er := tgsSticker.Stat()
		if er != nil {
			return nil, er
		}
		return ss.bot.apiInterface.AddStickerToSet(
			ss.userId, ss.stickerSet.Name, "", "attach://"+stat.Name(), "", emojies, maskPosition, tgsSticker,
		)
	}
}

//AddPngSticker adds a new PNG picture to the sticker set. This method should be used when the PNG file in stored in telegram servers or it's an HTTP URL. If the file is stored in your computer, use "AddPngStickerByFile" method.
func (ss *StickerSet) AddPngSticker(pngPicFileIdOrUrl, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	return ss.bot.apiInterface.AddStickerToSet(
		ss.userId, ss.stickerSet.Name, pngPicFileIdOrUrl, "", "", emojies, maskPosition, nil,
	)
}

//AddPngStickerByFile adds a new PNG picture to the sticker set. This method should be used when the PNG file in stored in your computer.
func (ss *StickerSet) AddPngStickerByFile(pngPicFile *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	if pngPicFile == nil {
		return nil, errors.New("pngPicFile cannot be nil")
	}
	stat, er := pngPicFile.Stat()
	if er != nil {
		return nil, er
	}
	return ss.bot.apiInterface.AddStickerToSet(
		ss.userId, ss.stickerSet.Name, "attach://"+stat.Name(), "", "", emojies, maskPosition, pngPicFile,
	)
}

//AddAnimatedSticker adds a new TGS sticker (animated sticker) to the sticker set.
func (ss *StickerSet) AddAnimatedSticker(tgsFile *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	if tgsFile == nil {
		return nil, errors.New("tgsFile cannot be nil")
	}
	stat, er := tgsFile.Stat()
	if er != nil {
		return nil, er
	}
	return ss.bot.apiInterface.AddStickerToSet(
		ss.userId, ss.stickerSet.Name, "", "attach://"+stat.Name(), "", emojies, maskPosition, tgsFile,
	)
}

//AddVideoSticker adds a new WEBM sticker (video sticker) to the sticker set.
func (ss *StickerSet) AddVideoSticker(webmFile *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	if webmFile == nil {
		return nil, errors.New("webmFile cannot be nil")
	}
	stat, er := webmFile.Stat()
	if er != nil {
		return nil, er
	}
	return ss.bot.apiInterface.AddStickerToSet(
		ss.userId, ss.stickerSet.Name, "", "", "attach://"+stat.Name(), emojies, maskPosition, webmFile,
	)
}

/*SetStickerPosition can be used to move a sticker in a set created by the bot to a specific position. Returns True on success.

"sticker" is file identifier of the sticker and "position" is new sticker position in the set, zero-based*/
func (ss *StickerSet) SetStickerPosition(sticker string, position int) (*objs.LogicalResult, error) {
	if ss == nil {
		return nil, errors.New("sticker set is nil")
	}
	return ss.bot.apiInterface.SetStickerPositionInSet(sticker, position)
}

/*DeleteStickerFromSet can be used to delete a sticker from a set created by the bot. Returns True on success.

"sticker" is file identifier of the sticker.*/
func (ss *StickerSet) DeleteStickerFromSet(sticker string) (*objs.LogicalResult, error) {
	if ss == nil {
		return nil, errors.New("sticker set is nil")
	}
	return ss.bot.apiInterface.DeleteStickerFromSet(sticker)
}

/*SetThumb can be used to set the thumbnail of a sticker set using url or file id. Animated thumbnails can be set for animated sticker sets only. Returns True on success.*/
func (ss *StickerSet) SetThumb(userId int, thumb string) (*objs.LogicalResult, error) {
	if ss == nil {
		return nil, errors.New("sticker set is nil")
	}
	return ss.bot.apiInterface.SetStickerSetThumb(ss.stickerSet.Name, thumb, userId, nil)
}

/*SetThumbByFile can be used to set the thumbnail of a sticker set using a file on the computer. Animated thumbnails can be set for animated sticker sets only. Returns True on success.*/
func (ss *StickerSet) SetThumbByFile(userId int, thumb *os.File) (*objs.LogicalResult, error) {
	if ss == nil {
		return nil, errors.New("sticker set is nil")
	}
	stats, err := thumb.Stat()
	if err != nil {
		return nil, err
	}
	return ss.bot.apiInterface.SetStickerSetThumb(ss.stickerSet.Name, "attach://"+stats.Name(), userId, thumb)
}
