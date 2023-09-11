# Telego

[![Go Reference](https://pkg.go.dev/badge/github.com/SakoDroid/telego/v2.svg)](https://pkg.go.dev/github.com/SakoDroid/telego/v2)
![example workflow](https://github.com/SakoDroid/telego/v2/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/SakoDroid/telego/v2)](https://goreportcard.com/report/github.com/SakoDroid/telego/v2)
![Version](https://img.shields.io/badge/%20%20Version%20%20-%20%201.8.0%20-brightgreen)
![Development status](https://img.shields.io/badge/%20%20Development%20%20-%20%20Active%20%20-blueviolet)

A Go library for creating telegram bots.

![telego logo inspired by Golang logo](https://i.ibb.co/zr9tNJQ/telego-logo-removebg.png)

---------------------------------
## ⚠️ Deprecation notice ⚠️
All versions and releases before v2.0.0 (including v1.8.0) have been deprecated and are considered out of date. Please consider updating to [v2](https://github.com/SakoDroid/telego/tree/master/v2) ASAP.

---------------------------------
## ✅ Upgrading to v2
Telego [v2](https://github.com/SakoDroid/telego/tree/master/v2) is a ground breaking change in Telego. Many new features have been added and many methods declarations have been changed. Some methods have been completely deprecated. Please read the change log before upgrading to [v2](https://github.com/SakoDroid/telego/tree/master/v2) and take extra caution while upgrading.

---------------------------------
## Supporting Telego
If you like Telego and it has helped you to build your telegram bot easily 🫶🏼, then you can support it to get even better!
You can support Telego by donating to the project or even giving it a star! Any kind of support is appreciated ❤️

If you want to donate to telego project and keep it running, [click here](https://github.com/SakoDroid/telego/v2/blob/master/donate.md)


---------------------------------
## License

Telego is licensed under [MIT lisence](https://en.wikipedia.org/wiki/MIT_License). Which means it can be used for commercial and private apps and can be modified.

---------------------------

## Change logs

### v2.0.0
**Improvements** :
* Added full support for Telegram bot API 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8 .
* Run method of the bot now has an auto pause option. If true is passed, the Run method will lock the go routine.
* Many many classes have been upgraded to newest version and many fields have been added. (Naming all of them is not possible due to the huge amount of change)
* Implemented a thread safe map in callback handlers to avoid race condition.
* Golang version (in go.mod) has been upgraded to **1.18**.
* Bug fixes.

**New features** :
* Added a new `stickerEditor` which has several methods for editing a sticker. Such as editing the emoji list, keywordas and etc. Accessible via `GetStickerEditor` method.
* Added `forumTopicManager` and `generalForumTopicManager` for managing topics. They both provide a set of tools that can be leveraged for manaing topics.
* Added spoiler support for media (photo,video,animation) messages.
* Added two new keys to the keyboard. `RequestUser` and `ReqestChat`.
* Added a set of new handlers used for handling messages triggered by `RequestUser` and `ReqestChat` buttons.
* Added `SetButton` methods to `inlineQueryResponder` for setting a button that will be shown above inline query results.
* Added `AddSwitchInlineQueryChoseChatButton` method to the inline keyboard to support the new `SwitchInlineQueryChosenChat` option in inline keyboards.
* Introduced bot manager tool, a tool for managing bot's personal information such as name and description.
* Removed row argument from `AddPayButton` of inline keybpard and prechecking invoice message keyboard. (Regarding ISSUE #15).
* Methods `ACreateInvoice` and `ACreateInvoiceUN` now return error along side with the invoiceSender.

**Deprecations** : 

1. CreateNewStickerSet method of the bot. Use CreateStickerSet instead.
2. AddSticker, AddPngSticker, AddPngStickerByFile, AddAnimatedSticker,AddVideoSticker methods of stickerSet. Use AddNewSticker and AddNewStickerByFile instead.

### v1.8.0
* Added full support for telegram bot API 6.0 and 6.1
* Added support for web apps
* Fixed some code errors
* Fixed a bug in webhook

### v1.7.0
* Added config updating option while bot is running.
* Added block option.
* Added `VerifyJoin` method for checking if a user has joined a channel or supergroup or not.
* Added file based configs.
* Improved logging system.
* Improved documentation

### v1.6.7
* Added support for telegram bot API 5.7
* Improved sticker creation experience by adding new separate methods.
* Correct syntax errors by @ityulkanov
* Bug fixes

### v1.5.7
* Bug fixes

### v1.5.5
* Added webhook support
* Improved handlers and regex bug fixed.
* Some other bug fixes.

### v1.4.5
* Added TextFormatter tool for formatting texts.
* Bug fixes

### v1.3.5
* Added support for telegram bot API 5.6 .
* Improved documentation.

### v1.3.4
* Custom keyboard button handler
* Major bug fixes

### v1.3.3
* Callback handlers
* keyboard creation tool

---------------------------
![telego logo inspired by Golang logo](https://i.ibb.co/zr9tNJQ/telego-logo-removebg.png)
