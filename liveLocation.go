package telego

import (
	errs "github.com/SakoDroid/telego/errors"
	objs "github.com/SakoDroid/telego/objects"
)

//LiveLocation is a live location that can be sent to a user.
type LiveLocation struct {
	bot                                       *Bot
	chatIdInt                                 int
	chatIdString                              string
	messageId                                 int
	replyTo                                   int
	allowSendingWihoutReply                   bool
	replyMarkUp                               objs.ReplyMarkup
	latitude, longitude, horizontalAccuracy   float32
	livePeriod, heading, proximityAlertRadius int
}

/*Send sends this live location to all types of chats but channels. To send it to a channel use "SendToChannelMethod".

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.*/
func (ll *LiveLocation) Send(chatId int, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	ll.chatIdInt = chatId
	res, err := ll.bot.apiInterface.SendLocation(
		chatId, "", ll.latitude, ll.longitude, ll.horizontalAccuracy, ll.livePeriod,
		ll.heading, ll.proximityAlertRadius, ll.replyTo, silent, ll.allowSendingWihoutReply, protectContent,
		ll.replyMarkUp,
	)
	if err == nil {
		ll.messageId = res.Result.MessageId
	}
	return res, err
}

/*SendToChannel sends this live location to a channel. Chat id should be the username of the channel.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.*/
func (ll *LiveLocation) SendToChannel(chatId string, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	ll.chatIdString = chatId
	res, err := ll.bot.apiInterface.SendLocation(
		0, chatId, ll.latitude, ll.longitude, ll.horizontalAccuracy, ll.livePeriod,
		ll.heading, ll.proximityAlertRadius, ll.replyTo, silent, ll.allowSendingWihoutReply, protectContent,
		ll.replyMarkUp,
	)
	if err == nil {
		ll.messageId = res.Result.MessageId
	}
	return res, err
}

/*Edit edits the live location.

------------------------

Official telegram doc :

Use this method to edit live location messages. A location can be edited until its live_period expires or editing is explicitly disabled by a call to stopMessageLiveLocation. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.*/
func (ll *LiveLocation) Edit(latitude, langitude, horizontalAccuracy float32, heading, proximtyAlertRadius int, replyMarkUp *objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	if ll.messageId != 0 {
		ll.latitude = latitude
		ll.longitude = langitude
		ll.horizontalAccuracy = horizontalAccuracy
		ll.heading = heading
		ll.proximityAlertRadius = proximtyAlertRadius
		ll.replyMarkUp = replyMarkUp
		return ll.bot.apiInterface.EditMessageLiveLocation(
			ll.chatIdInt, ll.chatIdString, "", ll.messageId, ll.latitude, ll.longitude,
			ll.horizontalAccuracy, ll.heading, ll.proximityAlertRadius, replyMarkUp,
		)
	}
	return nil, &errs.LiveLocationNotStarted{}
}

/*Stop stops the live location.

------------------------

Official telegram doc :

Use this method to stop updating a live location message before live_period expires. On success, if the message is not an inline message, the edited Message is returned, otherwise True is returned.*/
func (ll *LiveLocation) Stop(replyMarkrup objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	if ll.messageId != 0 {
		return ll.bot.apiInterface.StopMessageLiveLocation(
			ll.chatIdInt, ll.chatIdString, "", ll.messageId, &replyMarkrup,
		)
	} else {
		return nil, &errs.LiveLocationNotStarted{}
	}
}
