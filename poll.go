package telebot

import (
	"errors"

	objs "github.com/SakoDroid/telebot/objects"
)

var Polls = make(map[string]*Poll)

type Poll struct {
	bot                                                                     *Bot
	chatIdInt, messageId, totalVoterCount                                   int
	id, question, pollType, explanation, explanationParseMode, chatIdString string
	options                                                                 []string
	result                                                                  []objs.PollOption
	isClosed, isAnonymouse, allowMultipleAnswers                            bool
	correctOptionId, openPeriod, closeDate                                  int
	updateChannel                                                           *chan bool
	explanationEntities                                                     []objs.MessageEntity
}

/*Adds an option to poll.

-If the poll has been sent this method will do nothing.*/
func (p *Poll) AddOption(option string) {
	if p.id == "" {
		p.options = append(p.options, option)
	}
}

/*Set explanation if this poll.

If the poll has been sent this method will do nothing.*/
func (p *Poll) SetExplanation(explanation, explanationParseMode string, explanationEntities []objs.MessageEntity) {
	if p.id == "" {
		p.explanation = explanation
		p.explanationParseMode = explanationParseMode
		p.explanationEntities = explanationEntities
	}
}

/*Sets the correct option id for the poll. The correct option is the index of the the true option in the options array.(0 based)

-If the type of this bot is "regular" this method will do nothing.

-If the poll has been sent this method will do nothing.*/
func (p *Poll) SetCorrectOption(co int) {
	if p.id == "" && p.pollType == "quiz" {
		p.correctOptionId = co
	}
}

/*Sets the flags of this poll. Flags are "isClosed","isAnonymous" and "allowMultipleAnswers".

-If the poll has been sent this method will do nothing.*/
func (p *Poll) SetFlags(isClosed, isAnonymous, allowMA bool) {
	if p.id == "" {
		p.isClosed = isClosed
		p.isAnonymouse = isAnonymous
		p.allowMultipleAnswers = allowMA
	}
}

/*Sets open period of this poll. According to official telegram doc, open period is amount of time in seconds the poll will be active after creation, 5-600. Can't be used together with close_date.

-If close date has been specified, this method wont set open period.

-If the poll has been sent this method will do nothing.*/
func (p *Poll) SetOpenPeriod(op int) {
	if p.id == "" && p.closeDate == 0 {
		p.openPeriod = op
	}
}

/*Sets the close date of this poll.

-If open period has been specified, this method wont set close date for the poll.

-If the poll has been sent this method will do nothing.*/
func (p *Poll) SetCloseDate(cd int) {
	if p.id == "" && p.openPeriod == 0 {
		p.closeDate = cd
	}
}

/*This method takes an update object and extracts the poll update from it.

May return error if the update does not contain any poll or the poll in the update is not this poll*/
func (p *Poll) Update(poll *objs.Poll) error {
	if poll.Id == "" {
		return errors.New("no poll found in the update")
	}
	if poll.Id != p.id {
		return errors.New("this update dos not belong to this poll")
	}
	p.result = poll.Options
	Polls[p.id] = p
	*p.updateChannel <- true
	return nil
}

/*Returns the poll type. Its either "regular" or "quiz".*/
func (p *Poll) GetType() string {
	return p.pollType
}

/*Returns the id of this poll*/
func (p *Poll) GetId() string {
	return p.id
}

/*Returns the question of this poll*/
func (p *Poll) GetQuestion() string {
	return p.question
}

/*Returns the explanation of this poll*/
func (p *Poll) GetExplanation() string {
	return p.explanation
}

/*Returns the options of this poll*/
func (p *Poll) GetOptions() []string {
	return p.options
}

/*Returns the correct option id. returnes 0 if type of the poll is "regular"*/
func (p *Poll) GetCorrectOption() int {
	return p.correctOptionId
}

/*Returns the update channel for this poll. Everytime an update is received which contains update for this poll, true is passed into the channel.*/
func (p *Poll) GetUpdateChannel() *chan bool {
	return p.updateChannel
}

/*Returns the up to date result of the poll*/
func (p *Poll) GetResult() []objs.PollOption {
	return p.result
}

/*Returns the up to date total number of voters for this poll*/
func (p *Poll) GetTotalVoters() int {
	return p.totalVoterCount
}

/*Sends the poll. If you want more options foe sending the bot, use "SendAdvanced" method.*/
func (p *Poll) Send(silent bool, replyTo int) error {
	res, err := p.bot.apiInterface.SendPoll(
		p.chatIdInt, p.chatIdString, p.question, p.options, p.isClosed, p.isAnonymouse,
		p.pollType, p.allowMultipleAnswers, p.correctOptionId, p.explanation, p.explanationParseMode,
		p.explanationEntities, p.openPeriod, p.closeDate, replyTo, silent, false, nil,
	)
	if err != nil {
		return err
	}
	p.messageId = res.Result.MessageId
	p.id = res.Result.Poll.Id
	p.result = res.Result.Poll.Options
	ch := make(chan bool)
	p.updateChannel = &ch
	Polls[p.id] = p
	return nil
}

/*Sends the poll. This method has more options than "Send" method*/
func (p *Poll) SendAdvanced(replyTo int, silent bool, allowSendingWithOutReply bool, replyMarkup objs.ReplyMarkup) error {
	res, err := p.bot.apiInterface.SendPoll(
		p.chatIdInt, p.chatIdString, p.question, p.options, p.isClosed, p.isAnonymouse,
		p.pollType, p.allowMultipleAnswers, p.correctOptionId, p.explanation, p.explanationParseMode,
		p.explanationEntities, p.openPeriod, p.closeDate, replyTo, silent, allowSendingWithOutReply, replyMarkup,
	)
	if err != nil {
		return err
	}
	p.messageId = res.Result.MessageId
	p.id = res.Result.Poll.Id
	p.result = res.Result.Poll.Options
	ch := make(chan bool)
	p.updateChannel = &ch
	Polls[p.id] = p
	return nil
}

/*stops the poll*/
func (p *Poll) Stop() error {
	_, err := p.bot.apiInterface.StopPoll(
		p.chatIdInt, p.chatIdString, p.messageId, objs.InlineKeyboardMarkup{},
	)
	return err
}
