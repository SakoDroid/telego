package telego

import (
	"errors"
	"os"

	errs "github.com/SakoDroid/telego/errors"
	logger "github.com/SakoDroid/telego/logger"
	objs "github.com/SakoDroid/telego/objects"
)

// StickerSet is a set of stickers.
type StickerSet struct {
	bot                                     *Bot
	stickerSet                              *objs.StickerSet
	initStickers                            []*objs.InputSticker
	initFiles                               []*os.File
	userId                                  int
	name, title, stickerFormat, stickerType string
	needsRepainting                         bool
	created                                 bool
}

/*update updates this sticker set*/
func (ss *StickerSet) update() {
	if ss.created && ss != nil {
		res, err := ss.bot.apiInterface.GetStickerSet(ss.name)
		if err != nil {
			logger.Logger.Println("Error while updating sticker set.", err.Error())
		} else {
			ss.stickerSet = res.Result
		}
	}
}

//Create is used for creating this sticker set if it has not been created before.
func (ss *StickerSet) Create() (bool, error) {
	if ss.created {
		return false, errors.New("sticker set already created")
	}
	res, err := ss.bot.apiInterface.CreateNewStickerSet(
		ss.userId, ss.name, ss.title, ss.stickerFormat, ss.stickerType, ss.needsRepainting, ss.initStickers, ss.initFiles...,
	)
	if err != nil || !res.Ok {
		return false, err
	}
	ss.created = true
	ss.update()
	return true, nil
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

/*
AddNewSticker adds a new sticker to this set. If the set has been created before then the sticker will be added to sticker set on telegram servers.

If not created yet, then the sticker is added to an internal array. The "Create" function should be called after all stickers have been added.

userId is the user id of the owner.
*/
func (ss *StickerSet) AddNewSticker(fileIdOrURL string, userId int, emojiList, keywords []string, maskPosition *objs.MaskPosition) (bool, error) {
	inputSticker := &objs.InputSticker{
		Sticker:      fileIdOrURL,
		EmojiList:    emojiList,
		MaskPosition: maskPosition,
		KeyWords:     keywords,
	}
	if ss.created {
		res, err := ss.bot.apiInterface.AddStickerToSet(
			userId,
			ss.name,
			inputSticker,
			nil,
		)
		defer ss.update()
		return res.Ok, err
	}
	ss.initStickers = append(ss.initStickers, inputSticker)
	return true, nil
}

/*
AddNewStickerByFile adds a new sticker to this set. If the set has been created before then the sticker will be added to sticker set on telegram servers.

If not created yet, then the sticker is added to an internal array. The "Create" function should be called after all stickers have been added.

userId is the user id of the owner.
*/
func (ss *StickerSet) AddNewStickerByFile(file *os.File, userId int, emojiList, keywords []string, maskPosition *objs.MaskPosition) (bool, error) {
	stat, err := file.Stat()
	if err != nil {
		return false, err
	}
	inputSticker := &objs.InputSticker{
		Sticker:      "attach://" + stat.Name(),
		EmojiList:    emojiList,
		MaskPosition: maskPosition,
		KeyWords:     keywords,
	}
	if ss.created {
		res, err := ss.bot.apiInterface.AddStickerToSet(
			userId,
			ss.name,
			inputSticker,
			nil,
		)
		defer ss.update()
		return res.Ok, err
	}
	ss.initStickers = append(ss.initStickers, inputSticker)
	ss.initFiles = append(ss.initFiles, file)
	return true, nil
}

// Deprecated: This function should no longer be used for adding stickers to a sticker set. Use "AddNewSticker" method instead.
func (ss *StickerSet) AddSticker(pngStickerFileIdOrUrl string, pngStickerFile *os.File, tgsSticker *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	// if ss == nil {
	// 	return nil, errors.New("sticker set is nil")
	// }
	// if tgsSticker == nil {
	// 	if pngStickerFile == nil {
	// 		if pngStickerFileIdOrUrl == "" {
	// 			return nil, errors.New("wrong file id or url")
	// 		}
	// 		return ss.bot.apiInterface.AddStickerToSet(
	// 			ss.userId, ss.stickerSet.Name, pngStickerFileIdOrUrl, "", "", emojies, maskPosition, nil,
	// 		)
	// 	} else {
	// 		stat, er := pngStickerFile.Stat()
	// 		if er != nil {
	// 			return nil, er
	// 		}
	// 		return ss.bot.apiInterface.AddStickerToSet(
	// 			ss.userId, ss.stickerSet.Name, "attach://"+stat.Name(), "", "", emojies, maskPosition, pngStickerFile,
	// 		)
	// 	}
	// } else {
	// 	stat, er := tgsSticker.Stat()
	// 	if er != nil {
	// 		return nil, er
	// 	}
	// 	return ss.bot.apiInterface.AddStickerToSet(
	// 		ss.userId, ss.stickerSet.Name, "", "attach://"+stat.Name(), "", emojies, maskPosition, tgsSticker,
	// 	)
	// }
	return nil, &errs.MethodDeprecated{MethodName: "AddSticker", Replacement: "AddNewSticker"}
}

// Deprecated: This function should no longer be used for adding stickers to a sticker set. Use "AddNewSticker" method instead.
func (ss *StickerSet) AddPngSticker(pngPicFileIdOrUrl, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	return nil, &errs.MethodDeprecated{MethodName: "AddPngSticker", Replacement: "AddNewSticker"}
	// return ss.bot.apiInterface.AddStickerToSet(
	// 	ss.userId, ss.stickerSet.Name, pngPicFileIdOrUrl, "", "", emojies, maskPosition, nil,
	// )
}

// Deprecated: This function should no longer be used for adding stickers to a sticker set. Use "AddNewSticker" method instead.
func (ss *StickerSet) AddPngStickerByFile(pngPicFile *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	return nil, &errs.MethodDeprecated{MethodName: "AddPngStickerByFile", Replacement: "AddNewStickerByFile"}
	// if pngPicFile == nil {
	// 	return nil, errors.New("pngPicFile cannot be nil")
	// }
	// stat, er := pngPicFile.Stat()
	// if er != nil {
	// 	return nil, er
	// }
	// return ss.bot.apiInterface.AddStickerToSet(
	// 	ss.userId, ss.stickerSet.Name, "attach://"+stat.Name(), "", "", emojies, maskPosition, pngPicFile,
	// )
}

// Deprecated: This function should no longer be used for adding stickers to a sticker set. Use "AddNewSticker" method instead.
func (ss *StickerSet) AddAnimatedSticker(tgsFile *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	return nil, &errs.MethodDeprecated{MethodName: "AddAnimatedSticker", Replacement: "AddNewSticker"}
	// if tgsFile == nil {
	// 	return nil, errors.New("tgsFile cannot be nil")
	// }
	// stat, er := tgsFile.Stat()
	// if er != nil {
	// 	return nil, er
	// }
	// return ss.bot.apiInterface.AddStickerToSet(
	// 	ss.userId, ss.stickerSet.Name, "", "attach://"+stat.Name(), "", emojies, maskPosition, tgsFile,
	// )
}

// Deprecated: This function should no longer be used for adding stickers to a sticker set. Use "AddNewSticker" method instead.
func (ss *StickerSet) AddVideoSticker(webmFile *os.File, emojies string, maskPosition *objs.MaskPosition) (*objs.LogicalResult, error) {
	return nil, &errs.MethodDeprecated{MethodName: "AddVideoSticker", Replacement: "AddNewSticker"}
	// if webmFile == nil {
	// 	return nil, errors.New("webmFile cannot be nil")
	// }
	// stat, er := webmFile.Stat()
	// if er != nil {
	// 	return nil, er
	// }
	// return ss.bot.apiInterface.AddStickerToSet(
	// 	ss.userId, ss.stickerSet.Name, "", "", "attach://"+stat.Name(), emojies, maskPosition, webmFile,
	// )
}

/*
SetStickerPosition can be used to move a sticker in a set created by the bot to a specific position. Returns True on success.

"sticker" is file identifier of the sticker and "position" is new sticker position in the set, zero-based
*/
func (ss *StickerSet) SetStickerPosition(sticker string, position int) (*objs.LogicalResult, error) {
	if ss == nil {
		return nil, errors.New("sticker set is nil")
	}
	return ss.bot.apiInterface.SetStickerPositionInSet(sticker, position)
}

/*
DeleteStickerFromSet can be used to delete a sticker from a set created by the bot. Returns True on success.

"sticker" is file identifier of the sticker.
*/
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
