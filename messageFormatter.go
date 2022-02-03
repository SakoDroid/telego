package telego

import objs "github.com/SakoDroid/telego/objects"

//TextFormatter is tool for creating formatted texts.
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

/*AddNormal adds a normal text to the original text*/
func (mf *TextFormatter) AddNormal(text string) {
	mf.text += text + " "
}

/*AddMention adds a mention to the original text. example : @username*/
func (mf *TextFormatter) AddMention(text string) {
	mf.addEntity(text, "mention", "", "", nil)
}

/*AddHashtag adds a hashtag to the original text. example : #hashtag*/
func (mf *TextFormatter) AddHashtag(text string) {
	mf.addEntity(text, "hashtag", "", "", nil)
}

/*AddCashtag adds a cashtag to the original text. exzmple : $USD*/
func (mf *TextFormatter) AddCashtag(text string) {
	mf.addEntity(text, "cashtag", "", "", nil)
}

/*AddBotCommand adds a bot command to the original text. example : /start@jobs_bot*/
func (mf *TextFormatter) AddBotCommand(text string) {
	mf.addEntity(text, "bot_command", "", "", nil)
}

/*AddURL adds a url (not a clickable text url) to the original text. example : https://telegram.org*/
func (mf *TextFormatter) AddURL(text string) {
	mf.addEntity(text, "url", "", "", nil)
}

/*AddEmail adds an email to the original text. example : do-not-reply@telegram.org*/
func (mf *TextFormatter) AddEmail(text string) {
	mf.addEntity(text, "email", "", "", nil)
}

/*AddPhoneNumber adds a phone number to the original text. example : +1-212-555-0123*/
func (mf *TextFormatter) AddPhoneNumber(text string) {
	mf.addEntity(text, "phone_number", "", "", nil)
}

/*AddBold adds a bold text to the original text.*/
func (mf *TextFormatter) AddBold(text string) {
	mf.addEntity(text, "bold", "", "", nil)
}

/*AddItalic adds an italic text to the original text.*/
func (mf *TextFormatter) AddItalic(text string) {
	mf.addEntity(text, "italic", "", "", nil)
}

/*AddUnderline adds an underlined text to the original text.*/
func (mf *TextFormatter) AddUnderline(text string) {
	mf.addEntity(text, "underline", "", "", nil)
}

/*AddStrike adds a strikethrough text to the original text.*/
func (mf *TextFormatter) AddStrike(text string) {
	mf.addEntity(text, "strikethrough", "", "", nil)
}

/*AddSpoiler adds a spoiler text to the original text. This text is hidden until user clicks on it.*/
func (mf *TextFormatter) AddSpoiler(text string) {
	mf.addEntity(text, "spoiler", "", "", nil)
}

/*AddCode adds a piece of code (programming code) to the original text.*/
func (mf *TextFormatter) AddCode(text, language string) {
	mf.addEntity(text, "pre", "", language, nil)
}

/*AddTextLink adds a text link (clickable text which opens a URL) to the original text.*/
func (mf *TextFormatter) AddTextLink(text, url string) {
	mf.addEntity(text, "text_link", url, "", nil)
}

/*AddTextMention adds a mention (for users without username) to the original text.*/
func (mf *TextFormatter) AddTextMention(text string, user *objs.User) {
	mf.addEntity(text, "text_mention", "", "", user)
}

/*GetText returns the original text*/
func (mf *TextFormatter) GetText() string {
	return mf.text
}

/*GetEntities returnss the entities array*/
func (mf *TextFormatter) GetEntities() []objs.MessageEntity {
	return mf.entites
}
