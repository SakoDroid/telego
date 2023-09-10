package telego

import objs "github.com/SakoDroid/telego/objects"

// Invoice is an invoice that can be modified and sent to the user.
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

/*
AddPrice adds a new price label to this invoice.

"amount" is the price of the product in the smallest units of the currency (integer, not float/double). For example, for a price of US$ 1.45 pass amount = 145.
*/
func (is *Invoice) AddPrice(label string, amount int) {
	is.prices = append(is.prices, objs.LabeledPrice{Label: label, Amount: amount})
}

/*
Send sends this invoice.

-------------------------------

Official telegram doc :

Use this method to send invoices. On success, the sent Message is returned.
*/
func (is *Invoice) Send(replyTo, messageThreadId int, silent bool) (*objs.Result[*objs.Message], error) {
	return is.bot.apiInterface.SendInvoice(
		is.chatIdInt, is.chatIdString, is.title, is.description, is.payload, is.providerToken,
		is.currency, is.prices, is.maxTipAmount, is.suggestedTipAmounts, is.startParameter, is.providerData,
		is.photoURL, is.photoSize, is.photoWidth, is.photoHeight, is.needName, is.needPhoneNumber, is.needEmail, is.needShippingAddress,
		is.sendPhoneNumberToProvider, is.sendEmailToProvider, is.isFlexible, silent, replyTo, messageThreadId, is.allowSendingWithoutReply, is.replyMarkup,
	)
}

/*
CreateLink creates a link for the invoice and returnes the link.

-------------------------------

Official telegram doc :

Use this method to create a link for an invoice. Returns the created invoice link as String on success.
*/
func (is *Invoice) CreateLink() (*objs.Result[string], error) {
	return is.bot.apiInterface.CreateInvoiceLink(is.title, is.description, is.payload, is.providerToken,
		is.currency, is.prices, is.maxTipAmount, is.suggestedTipAmounts, is.providerData,
		is.photoURL, is.photoSize, is.photoWidth, is.photoHeight, is.needName, is.needPhoneNumber, is.needEmail, is.needShippingAddress,
		is.sendPhoneNumberToProvider, is.sendEmailToProvider, is.isFlexible)
}
