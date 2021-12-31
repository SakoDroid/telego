package telego

import objs "github.com/SakoDroid/telego/objects"

type TextFormatter struct {
	text    string
	entites []objs.MessageEntity
}

func (mf *TextFormatter) addEntity(text, tp, url, lang string, user *objs.User) {
	length := len(text)
	if length != 0 {
		me := objs.MessageEntity{Type: tp, Length: length, URL: url, Language: lang, User: user, Offset: len(mf.text)}
		mf.entites = append(mf.entites, me)
		mf.text += text + " "
	}
}

/*Adds a normal text to the original text*/
func (mf *TextFormatter) AddNormal(text string) {
	mf.text += text + " "
}

/*Adds a mention to the original text. example : @username*/
func (mf *TextFormatter) AddMention(text string) {
	mf.addEntity(text, "mention", "", "", nil)
}

/*Adds a hashtag to the original text. example : #hashtag*/
func (mf *TextFormatter) AddHashtag(text string) {
	mf.addEntity(text, "hashtag", "", "", nil)
}

/*Adds a cashtag to the original text. exzmple : $USD*/
func (mf *TextFormatter) AddCashtag(text string) {
	mf.addEntity(text, "cashtag", "", "", nil)
}

/*Adds a bot command to the original text. example : /start@jobs_bot*/
func (mf *TextFormatter) AddBotCommand(text string) {
	mf.addEntity(text, "bot_command", "", "", nil)
}

/*Adds a url (not a clickable text url) to the original text. example : https://telegram.org*/
func (mf *TextFormatter) AddURL(text string) {
	mf.addEntity(text, "url", "", "", nil)
}

/*Adds an email to the original text. example : do-not-reply@telegram.org*/
func (mf *TextFormatter) AddEmail(text string) {
	mf.addEntity(text, "email", "", "", nil)
}

/*Adds a phone number to the original text. example : +1-212-555-0123*/
func (mf *TextFormatter) AddPhoneNumber(text string) {
	mf.addEntity(text, "phone_number", "", "", nil)
}

/*Adds a bold text to the original text.*/
func (mf *TextFormatter) AddBold(text string) {
	mf.addEntity(text, "bold", "", "", nil)
}

/*Adds an italic text to the original text.*/
func (mf *TextFormatter) AddItalic(text string) {
	mf.addEntity(text, "italic", "", "", nil)
}

/*Adds an underlined text to the original text.*/
func (mf *TextFormatter) AddUnderline(text string) {
	mf.addEntity(text, "underline", "", "", nil)
}

/*Adds a strikethrough text to the original text.*/
func (mf *TextFormatter) AddStrike(text string) {
	mf.addEntity(text, "strikethrough", "", "", nil)
}

/*Adds a spoiler text to the original text. This text is hidden until user clicks on it.*/
func (mf *TextFormatter) AddSpoiler(text string) {
	mf.addEntity(text, "spoiler", "", "", nil)
}

/*Adds a piece of code (programming code) to the original text.*/
func (mf *TextFormatter) AddCode(text, language string) {
	mf.addEntity(text, "pre", "", language, nil)
}

/*Adds a text link (clickable text which opens a URL) to the original text.*/
func (mf *TextFormatter) AddTextLink(text, url string) {
	mf.addEntity(text, "text_link", url, "", nil)
}

/*Adds a mention (for users without username) to the original text.*/
func (mf *TextFormatter) AddTextMention(text string, user *objs.User) {
	mf.addEntity(text, "text_mention", "", "", user)
}

/*Returnes the original text*/
func (mf *TextFormatter) GetText() string {
	return mf.text
}

/*Returnes the entities array*/
func (mf *TextFormatter) GetEntities() []objs.MessageEntity {
	return mf.entites
}
