package objects

/*This object represents a portion of the price for goods or services.*/
type LabeledPrice struct {
	/*Portion label*/
	Label string `json:"label"`
	/*Price of the product in the smallest units of the currency (integer, not float/double). For example, for a price of US$ 1.45 pass amount = 145.*/
	Amount int `json:"amount"`
}

/*This object contains basic information about an invoice.*/
type Invoice struct {
	/*Product name*/
	Title string `json:"title"`
	/*Product description*/
	Description string `json:"description"`
	/*Unique bot deep-linking parameter that can be used to generate this invoice*/
	StartParameter string `json:"start_parameter"`
	/*Three-letter ISO 4217 currency code*/
	Currency string `json:"currency"`
	/*Total price in the smallest units of the currency (integer, not float/double). For example, for a price of US$ 1.45 pass amount = 145.*/
	TotalAmount string `json:"total_amount"`
}

/*This object represents a shipping address.*/
type ShippingAddress struct {
	/*ISO 3166-1 alpha-2 country code*/
	CountryCode string `json:"country_code"`
	/*State, if applicable*/
	State string `json:"state"`
	/*City*/
	City string `json:"city"`
	/*First line for the address*/
	StreetLine1 string `json:"street_line1"`
	/*Second line for the address*/
	StreetLine2 string `json:"street_line2"`
	/*Address post code*/
	PostCode string `json:"post_code"`
}

/*This object represents information about an order.*/
type OrderInfo struct {
	/*Optional. User name*/
	Name string `json:"name,omitempty"`
	/*Optional. User's phone number*/
	PhoneNumber string `json:"phone_number,omitempty"`
	/*Optional. User email*/
	Email string `json:"email,omitempty"`
	/*Optional. User shipping address*/
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

/*This object represents one shipping option.*/
type ShippingOption struct {
	/*Shipping option identifier*/
	Id string `json:"id"`
	/*Option title*/
	Title string `json:"title"`
	/*List of price portions*/
	Prices []LabeledPrice `json:"prices"`
}

/*This object contains basic information about a successful payment.*/
type SuccessfulPayment struct {
	/*Three-letter ISO 4217 currency code*/
	Currency string `json:"currency"`
	/*Total price in the smallest units of the currency (integer, not float/double). For example, for a price of US$ 1.45 pass amount = 145. */
	TotalAmount int `json:"total_amount"`
	/*Bot specified invoice payload*/
	InvoicePayload string `json:"invoice_payload"`
	/*Optional. Identifier of the shipping option chosen by the user*/
	ShippingOptionId string `json:"shipping_option_id,omitempty"`
	/*Optional. Order info provided by the user*/
	OrderInfo *OrderInfo `json:"order_info,omitempty"`
	/*Telegram payment identifier*/
	TelegramPaymentChargeId string `json:"telegram_payment_charge_id"`
	/*Provider payment identifier*/
	ProviderPaymentChargeId string `json:"provider_payment_charge_id"`
}

/*This object contains information about an incoming shipping query.*/
type ShippingQuery struct {
	/*Unique query identifier*/
	Id string `json:"id"`
	/*User who sent the query*/
	From *User `json:"from"`
	/*Bot specified invoice payload*/
	InvoicePayload string `json:"invoice_payload"`
	/*User specified shipping address*/
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

/*This object contains information about an incoming pre-checkout query.*/
type PreCheckoutQuery struct {
	/*Unique query identifier*/
	Id string `json:"id"`
	/*User who sent the query*/
	From *User `json:"from"`
	/*Three-letter ISO 4217 currency code*/
	Currency string `json:"currency"`
	/*Total price in the smallest units of the currency (integer, not float/double). For example, for a price of US$ 1.45 pass amount = 145.*/
	TotalAmount int `json:"total_amount"`
	/*Bot specified invoice payload*/
	InvoicePayload string `json:"invoice_payload"`
	/*Optional. Identifier of the shipping option chosen by the user*/
	ShippingOptionId string `json:"shipping_option_id,omitempty"`
	/*Optional. Order info provided by the user*/
	OrderInfo *OrderInfo `json:"order_info,omitempty"`
}
