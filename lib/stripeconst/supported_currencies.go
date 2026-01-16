package stripeconst

// SupportedCurrenciesPaypal is a list of currencies supported for presentment
// when using PayPal as a payment method.
//
// https://stripe.com/docs/payments/paypal
var SupportedCurrenciesPaypal = []string{
	"AUD",
	"CAD",
	"CHF",
	"CZK",
	"DKK",
	"EUR",
	"GBP",
	"HKD",
	"NOK",
	"NZD",
	"PLN",
	"SEK",
	"SGD",
	"USD",
}

// SupportedCurrenciesKlarna is a list of currencies supported for presentment
// when using Klarna as a payment method.
//
// https://stripe.com/docs/payments/klarna
var SupportedCurrenciesKlarna = []string{
	"UR",
	"DKK",
	"GBP",
	"NOK",
	"SEK",
	"CZK",
	"RON",
	"PLN",
	"CHF",
}
