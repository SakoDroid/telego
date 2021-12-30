package objects

/*This object represents the content of a message to be sent as a result of an inline query. Telegram clients currently support the following 5 types:

InputTextMessageContent
InputLocationMessageContent
InputVenueMessageContent
InputContactMessageContent
InputInvoiceMessageContent*/
type InputMessageContent interface {
	/*Returns the full name of this object*/
	GetType() string
}

/*Represents the content of a text message to be sent as the result of an inline query.*/
type InputTextMessageContent struct {
	/*Text of the message to be sent, 1-4096 characters*/
	MessageText string `json:"message_text"`
	/*Optional. Mode for parsing entities in the message text.*/
	ParseMode string `json:"parse_mode,omitempty"`
	/*Optional. List of special entities that appear in message text, which can be specified instead of parse_mode*/
	Entities []MessageEntity `json:"entities,omitempty"`
	/*Optional. Disables link previews for links in the sent message*/
	DisableWebPagePreview bool `json:"disable_web_page_preview,omitempty"`
}

func (*InputTextMessageContent) GetType() string {
	return "InputTextMessageContent"
}

/*Represents the content of a location message to be sent as the result of an inline query.*/
type InputLocationMessageContent struct {
	/*Latitude of the location in degrees*/
	Latitude float32 `json:"latitude"`
	/*Longitude of the location in degrees*/
	Longitude float32 `json:"longitude"`
	/*Optional. The radius of uncertainty for the location, measured in meters; 0-1500*/
	HorizontalAccuracy float32 `json:"horizontal_accuracy,omitempty"`
	/*Optional. Period in seconds for which the location can be updated, should be between 60 and 86400.*/
	LivePeriod int `json:"live_period,omitempty"`
	/*Optional. For live locations, a direction in which the user is moving, in degrees. Must be between 1 and 360 if specified.*/
	Heading int `json:"heading,omitempty"`
	/*	Optional. For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters. Must be between 1 and 100000 if specified.*/
	ProximityAlertRadius int `json:"proximity_alert_radius,omitempty"`
}

func (*InputLocationMessageContent) GetType() string {
	return "InputLocationMessageContent"
}

/*Represents the content of a venue message to be sent as the result of an inline query.*/
type InputVenueMessageContent struct {
	/*Latitude of the venue in degrees*/
	Latitude float32 `json:"latitude"`
	/*Longitude of the venue in degrees*/
	Longitude float32 `json:"longitude"`
	/*Name of the venue*/
	Title string `json:"title"`
	/*Address of the venue*/
	Address string `json:"address"`
	/*Optional. Foursquare identifier of the venue, if known*/
	FoursquareId string `json:"fourquare_id,omitempty"`
	/*Optional. Foursquare type of the venue, if known. (For example, “arts_entertainment/default”, “arts_entertainment/aquarium” or “food/icecream”.)*/
	FoursquareType string `json:"foursquare_type,omitempty"`
	/*Optional. Google Places identifier of the venue*/
	GooglePlaceId string `json:"google_place_id,omitempty"`
	/*Optional. Google Places type of the venue.*/
	GooglePlaceType string `json:"google_place_type,omitempty"`
}

func (*InputVenueMessageContent) GetType() string {
	return "InputVenueMessageContent"
}

/*Represents the content of a contact message to be sent as the result of an inline query.*/
type InputContactMessageContent struct {
	/*Contact's phone number*/
	PhoneNumber string `json:"phone_number"`
	/*Contact's first name*/
	FirstName string `json:"first_name"`
	/*Optional. Contact's last name*/
	LastName string `json:"last_name,omitempty"`
	/*Optional. Additional data about the contact in the form of a vCard, 0-2048 bytes*/
	Vcard string `json:"vcard,omitempty"`
}

func (*InputContactMessageContent) GetType() string {
	return "InputContactMessageContent"
}

type InputInvoiceMessageContent struct {
	/*Product name, 1-32 characters*/
	Title string `json:"title"`
	/*Product description, 1-255 characters*/
	Description string `json:"description"`
	/*Bot-defined invoice payload, 1-128 bytes. This will not be displayed to the user, use for your internal processes.*/
	Payload string `json:"payload"`
	/*Payment provider token, obtained via Botfather*/
	ProviderToken string `json:"provider_token"`
	/*Three-letter ISO 4217 currency code*/
	Currency string `json:"currency"`
	/*Price breakdown, a JSON-serialized list of components (e.g. product price, tax, discount, delivery cost, delivery tax, bonus, etc.)*/
	Prices []LabeledPrice `json:"prices"`
	/*Optional. The maximum accepted amount for tips in the smallest units of the currency (integer, not float/double). For example, for a maximum tip of US$ 1.45 pass max_tip_amount = 145. See the exp parameter in currencies.json, it shows the number of digits past the decimal point for each currency (2 for the majority of currencies). Defaults to 0*/
	MaxTipAmount int `json:"max_tip_amount,omitempty"`
	/*Optional. A JSON-serialized array of suggested amounts of tip in the smallest units of the currency (integer, not float/double). At most 4 suggested tip amounts can be specified. The suggested tip amounts must be positive, passed in a strictly increased order and must not exceed max_tip_amount.*/
	SuggestedTipAmounts []int `json:"suggested_tip_amounts,omitempty"`
	/*Optional. A JSON-serialized object for data about the invoice, which will be shared with the payment provider. A detailed description of the required fields should be provided by the payment provider.*/
	ProviderData string `json:"provider_data,omitempty"`
	/*	Optional. URL of the product photo for the invoice. Can be a photo of the goods or a marketing image for a service. People like it better when they see what they are paying for.*/
	PhotoURL string `json:"photo_url,omitempty"`
	/*	Optional. Photo size*/
	PhotoSize int `json:"photo_size,omitempty"`
	/*	Optional. Photo width*/
	PhotoWidth int `json:"photo_width,omitempty"`
	/*	Optional. Photo height*/
	PhotoHeight int `json:"photo_height,omitempty"`
	/*Optional. Pass True, if you require the user's full name to complete the order*/
	NeedName bool `json:"need_name,omitempty"`
	/*Optional. Pass True, if you require the user's phone number to complete the order*/
	NeedPhoneNumber bool `json:"need_phone_number,omitempty"`
	/*Optional. Pass True, if you require the user's email address to complete the order*/
	NeedEmail bool `json:"need_email,omitempty"`
	/*Optional. Pass True, if you require the user's shipping address to complete the order*/
	NeedShippingAddress bool `json:"need_shipping_address,omitempty"`
	/*	Optional. Pass True, if user's phone number should be sent to provider*/
	SendPhoneNumberToProvider bool `json:"send_phone_number_to_provider,omitempty"`
	/*Optional. Pass True, if user's email address should be sent to provider*/
	SendEmailToProvider bool `json:"send_email_to_provider,omitempty"`
	/*Optional. Pass True, if the final price depends on the shipping method*/
	IsFlexible bool `json:"is_flexible,omitempty"`
}

func (*InputInvoiceMessageContent) GetType() string {
	return "InputVoiceMessageContent"
}

/*Represents a result of an inline query that was chosen by the user and sent to their chat partner.*/
type ChosenInlineResult struct {
	/*The unique identifier for the result that was chosen*/
	ResultId string `json:"result_id"`
	/*The user that chose the result*/
	From User `json:"user"`
	/*Optional. Sender location, only for bots that require user location*/
	Location Location `json:"location,omitempty"`
	/*Optional. Identifier of the sent inline message. Available only if there is an inline keyboard attached to the message. Will be also received in callback queries and can be used to edit the message.*/
	InlineMessageId string `json:"inline_message_id,omitempty"`
	/*The query that was used to obtain the result*/
	Query string `json:"query,omitempty"`
}

func (*ChosenInlineResult) GetType() string {
	return "ChosenInlineResult"
}
