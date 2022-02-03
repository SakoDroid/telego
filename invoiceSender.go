package telego

import objs "github.com/SakoDroid/telego/objects"

//Invoice is an invoice that can be modified and sent to the user.
type Invoice struct {
	bot                                                                                                                                             *Bot
	chatIdInt                                                                                                                                       int
	chatIdString                                                                                                                                    string
	replyMarkup                                                                                                                                     objs.InlineKeyboardMarkup
	prices                                                                                                                                          []objs.LabeledPrice
	suggestedTipAmounts                                                                                                                             []int
	photoURL, startParameter, providerData, title, description, payload, providerToken, currency                                                    string
	photoSize, photoWidth, photoHeight, maxTipAmount                                                                                                int
	allowSendingWithoutReply, needName, needPhoneNumber, needEmail, needShippingAddress, sendPhoneNumberToProvider, sendEmailToProvider, isFlexible bool
}

/*AddPrice adds a new price label to this invoice.

"amount" is the price of the product in the smallest units of the currency (integer, not float/double). For example, for a price of US$ 1.45 pass amount = 145.*/
func (is *Invoice) AddPrice(label string, amount int) {
	is.prices = append(is.prices, objs.LabeledPrice{Label: label, Amount: amount})
}

/*Send sends this invoice.

-------------------------------

Official telegram doc :

Use this method to send invoices. On success, the sent Message is returned.*/
func (is *Invoice) Send(replyTo int, silent bool) (*objs.SendMethodsResult, error) {
	return is.bot.apiInterface.SendInvoice(
		is.chatIdInt, is.chatIdString, is.title, is.description, is.payload, is.providerToken,
		is.currency, is.prices, is.maxTipAmount, is.suggestedTipAmounts, is.startParameter, is.providerData,
		is.photoURL, is.photoSize, is.photoWidth, is.photoHeight, is.needName, is.needPhoneNumber, is.needEmail, is.needShippingAddress,
		is.sendPhoneNumberToProvider, is.sendEmailToProvider, is.isFlexible, silent, replyTo, is.allowSendingWithoutReply, is.replyMarkup,
	)
}
