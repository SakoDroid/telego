package telego

import (
	objs "github.com/SakoDroid/telego/objects"
	upp "github.com/SakoDroid/telego/parser"
)

//MarkUps is the interface used for creating normal keyboards and inline keyboards.
type MarkUps interface {
	toMarkUp() objs.ReplyMarkup
}

//keyboard is a normal keyboard.
type keyboard struct {
	keys                                       [][]*objs.KeyboardButton
	resizeKeyBoard, oneTimeKeyboard, selective bool
	inputFieldPlaceHolder                      string
}

func (kb *keyboard) fixRows(row int) {
	dif := (row) - len(kb.keys)
	for i := 0; i < dif; i++ {
		kb.keys = append(kb.keys, make([]*objs.KeyboardButton, 0))
	}
}

/*AddButton adds a new button holding the given text to the specified row. According to telegram bot api if this button is pressed the text inside the button will be sent to the bot as a message.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added*/
func (kb *keyboard) AddButton(text string, row int) {
	kb.addButton(text, row, false, false, nil)
}

/*AddButtonHandler adds a new button holding the given text to the specified row. This method also adds a handler for that button so everytime this button is pressed the handler will be called. You can read the documentation of "AddHandler" for better understanding on handlers.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added*/
func (kb *keyboard) AddButtonHandler(text string, row int, handler func(*objs.Update), chatTypes ...string) {
	kb.addButton(text, row, false, false, nil)
	upp.AddHandler(text, handler, chatTypes...)
}

/*AddContactButton adds a new contact button. According to telegram bot api when this button is pressed,the user's phone number will be sent as a contact. Available in private chats only.

Note: ContactButtons and LocationButtons will only work in Telegram versions released after 9 April, 2016. Older clients will display unsupported message.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added*/
func (kb *keyboard) AddContactButton(text string, row int) {
	kb.addButton(text, row, true, false, nil)
}

/*AddLocationButton adds a new location button. According to telegram bot api when this button is pressed,the user's location will be sent. Available in private chats only.

Note: ContactButtons and LocationButtons will only work in Telegram versions released after 9 April, 2016. Older clients will display unsupported message.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added*/
func (kb *keyboard) AddLocationButton(text string, row int) {
	kb.addButton(text, row, false, true, nil)
}

/*AddPollButton adds a new poll button. According to telegram bot api, the user will be asked to create a poll and send it to the bot when this button is pressed. Available in private chats only.

Note: PollButton will only work in Telegram versions released after 23 January, 2020. Older clients will display unsupported message.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added.

Note : poll type can be "regular" or "quiz". Any other value will cause the button not to be added.*/
func (kb *keyboard) AddPollButton(text string, row int, pollType string) {
	if pollType == "regular" || pollType == "quiz" {
		kb.addButton(text, row, false, false, &objs.KeyboardButtonPollType{Type: pollType})
	}
}

func (kb *keyboard) addButton(text string, row int, contact, location bool, poll *objs.KeyboardButtonPollType) {
	if row >= 1 {
		kb.fixRows(row)
		kb.keys[row-1] = append(kb.keys[row-1], &objs.KeyboardButton{
			Text:            text,
			RequestContact:  contact,
			RequestLocation: location,
			RequestPoll:     poll,
		})
	}
}

func (kb *keyboard) toMarkUp() objs.ReplyMarkup {
	return &objs.ReplyKeyboardMarkup{
		Keyboard:              kb.keys,
		ResizeKeyboard:        kb.resizeKeyBoard,
		OneTimeKeyboard:       kb.oneTimeKeyboard,
		InputFieldPlaceholder: kb.inputFieldPlaceHolder,
		Selective:             kb.selective,
	}
}

type inlineKeyboard struct {
	keys [][]*objs.InlineKeyboardButton
}

/*AddURLButton adds a button that will open an url when pressed.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added*/
func (in *inlineKeyboard) AddURLButton(text, url string, row int) {
	in.addButton(text, url, "", "", "", nil, nil, false, row)
}

/*AddLoginURLButton adds a button that will be used for automatic authorization. According to telegram bot api, login url is an HTTP URL used to automatically authorize the user. Can be used as a replacement for the Telegram Login Widget.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added.

Arguments :

1. url : An HTTP URL to be opened with user authorization data added to the query string when the button is pressed. If the user refuses to provide authorization data, the original URL without information about the user will be opened. The data added is the same as described in Receiving authorization data. NOTE: You must always check the hash of the received data to verify the authentication and the integrity of the data as described in Checking authorization.

2. forwardText : New text of the button in forwarded messages.

3. botUsername : Username of a bot, which will be used for user authorization. See Setting up a bot for more details. If not specified, the current bot's username will be assumed. The url's domain must be the same as the domain linked with the bot. See Linking your domain to the bot for more details.

4. requestWriteAccess : Pass True to request the permission for your bot to send messages to the user.
*/
func (in *inlineKeyboard) AddLoginURLButton(text, url, forwardText, botUsername string, requestWriteAccess bool, row int) {
	in.addButton(text, "", "", "", "", &objs.LoginUrl{
		URL:                url,
		ForwardText:        forwardText,
		BotUsername:        botUsername,
		RequestWriteAccess: requestWriteAccess,
	}, nil, false, row)
}

/*AddCallbackButton adds a button that when its pressed, a call back query with the given data is sen to the bot

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added.
*/
func (in *inlineKeyboard) AddCallbackButton(text, callbackData string, row int) {
	in.addButton(text, "", callbackData, "", "", nil, nil, false, row)
}

/*AddCallbackButtonHandler adds a button that when its pressed, a call back query with the given data is sen to the bot. A handler is also added which will be called everytime a call back query is received for this button.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added.
*/
func (in *inlineKeyboard) AddCallbackButtonHandler(text, callbackData string, row int, handler func(*objs.Update)) {
	in.addButton(text, "", callbackData, "", "", nil, nil, false, row)
	upp.AddCallbackHandler(callbackData, handler)
}

/*AddSwitchInlineQueryButton adds a switch inline query button. According to tlegram bot api, pressing the button will prompt the user to select one of their chats, open that chat and insert the bot's username and the specified inline query in the input field. Can be empty, in which case just the bot's username will be inserted. Note: This offers an easy way for users to start using your bot in inline mode when they are currently in a private chat with it. Especially useful when combined with switch_pm… actions – in this case the user will be automatically returned to the chat they switched from, skipping the chat selection screen.

Note : If "currentChat" option is true, the inline query will be inserted in the current chat's input field.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added.
*/
func (in *inlineKeyboard) AddSwitchInlineQueryButton(text, inlineQuery string, row int, currenChat bool) {
	if currenChat {
		in.addButton(text, "", "", "", inlineQuery, nil, nil, false, row)
	} else {
		in.addButton(text, "", "", inlineQuery, "", nil, nil, false, row)
	}
}

/*AddGameButton adds a game button. Everytime a user presses this button a game will be launched. Use botfather to set up a game.
NOTE: This type of button must always be the first button in the first row.*/
func (in *inlineKeyboard) AddGameButton(text string, row int) {
	in.addButton(text, "", "", "", "", nil, &objs.CallbackGame{}, false, row)
}

/*AddPayButton adds a pay button.

NOTE: This type of button must always be the first button in the first row.

Note : row number starts from 1. (it's not zero based). If any number lower than 1 is passed, no button will be added.
*/
func (in *inlineKeyboard) AddPayButton(text string, row int) {
	in.addButton(text, "", "", "", "", nil, nil, true, row)
}

func (in *inlineKeyboard) addButton(text, url, callbackData, switchInlineQuery, switchInlineQueryCurrentChat string, loginUrl *objs.LoginUrl, callbackGame *objs.CallbackGame, pay bool, row int) {
	if row >= 1 {
		in.fixRows(row)
		in.keys[row-1] = append(in.keys[row-1], &objs.InlineKeyboardButton{
			Text:                         text,
			URL:                          url,
			LoginURL:                     loginUrl,
			CallbackData:                 callbackData,
			SwitchInlineQuery:            switchInlineQuery,
			SwitchInlineQueryCurrentChat: switchInlineQueryCurrentChat,
			CallbackGame:                 callbackGame,
			Pay:                          pay,
		})
	}
}

func (in *inlineKeyboard) fixRows(row int) {
	dif := (row) - len(in.keys)
	for i := 0; i < dif; i++ {
		in.keys = append(in.keys, make([]*objs.InlineKeyboardButton, 0))
	}
}

func (in *inlineKeyboard) toInlineKeyboardMarkup() objs.InlineKeyboardMarkup {
	return objs.InlineKeyboardMarkup{
		InlineKeyboard: in.keys,
	}
}

func (in *inlineKeyboard) toMarkUp() objs.ReplyMarkup {
	return &objs.InlineKeyboardMarkup{
		InlineKeyboard: in.keys,
	}
}
