package stripeconst

// MinimumAmounts is a map of currencies to the minimum amount Stripe will
// allow for a charge.
//
// This is based on the settlement currency. It's unclear what exchange rates
// Stripe uses to map to the settlement currency.
//
// https://stripe.com/docs/currencies#minimum-and-maximum-charge-amounts
var MinimumAmounts = map[string]int64{
	"AED": 200,   // 2.00 د.إ
	"AUD": 50,    // $0.50
	"BGN": 100,   // лв1.00
	"BRL": 50,    // R$0.50
	"CAD": 50,    // $0.50
	"CHF": 50,    // 0.50 Fr
	"CZK": 1500,  // 15.00Kč
	"DKK": 250,   // 2.50-kr.
	"EUR": 50,    // €0.50
	"GBP": 30,    // £0.30
	"HKD": 400,   // $4.00
	"HRK": 50,    // 0.50 kn
	"HUF": 17500, // 175.00 Ft
	"INR": 50,    // ₹0.50
	"JPY": 50,    // ¥50 (zero-decimal)
	"MXN": 1000,  // $10
	"MYR": 200,   // RM 2
	"NOK": 300,   // 3.00-kr.
	"NZD": 50,    // $0.50
	"PLN": 200,   // 2.00 zł
	"RON": 200,   // lei2.00
	"SEK": 300,   // 3.00-kr.
	"SGD": 50,    // $0.50
	"THB": 1000,  // ฿10
	"USD": 50,    // $0.50
}
