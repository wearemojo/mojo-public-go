package stripeconst

// MinimumAmounts is a map of currencies to the minimum amount Stripe will
// allow for a charge.
//
// This is based on the settlement currency. It's unclear what exchange rates
// Stripe uses to map to the settlement currency.
//
// LM notes that this page is particularly useful:
// https://dashboard.stripe.com/currency_conversion
// It shows the currency conversion for Stripe within that 12 hour window.
//
// https://stripe.com/docs/currencies#minimum-and-maximum-charge-amounts
var MinimumAmounts = map[string]int64{
	"AED": 200,   // 2.00 د.إ
	"AUD": 60,    // $0.60
	"BGN": 100,   // лв1.00
	"BRL": 200,   // R$2.00
	"CAD": 60,    // $0.60
	"CHF": 50,    // 0.50 Fr
	"CZK": 1500,  // 15.00Kč
	"DKK": 300,   // 3.00-kr.
	"EUR": 50,    // €0.50
	"GBP": 30,    // £0.30
	"HKD": 400,   // $4.00
	"HRK": 300,   // 3.00 kn
	"HUF": 17500, // 175.00 Ft
	"INR": 4000,  // ₹40.00
	"JPY": 80,    // ¥80 (zero-decimal)
	"MXN": 1000,  // $10
	"MYR": 200,   // RM 2
	"NOK": 500,   // 5.00-kr.
	"NZD": 80,    // $0.80
	"PLN": 250,   // 2.50 zł
	"RON": 200,   // lei2.00
	"SEK": 500,   // 5.00-kr.
	"SGD": 60,    // $0.60
	"THB": 1500,  // ฿15
	"USD": 50,    // $0.50
}
