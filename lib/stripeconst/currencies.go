package stripeconst

import (
	mapset "github.com/deckarep/golang-set/v2"
)

// SpecialCaseCurrencies is a set of currencies that have special rules.
//
// Practically, amounts in these currencies are required to be integer multiples
// of 100 in some situations.
//
// https://stripe.com/docs/currencies#special-cases
var SpecialCaseCurrencies = mapset.NewThreadUnsafeSet(
	"ISK",
	"HUF",
	"TWD",
	"UGX",
)

// ZeroDecimalCurrencies is a set of currencies that have zero decimal places.
//
// https://stripe.com/docs/currencies#zero-decimal
var ZeroDecimalCurrencies = mapset.NewThreadUnsafeSet(
	"BIF",
	"CLP",
	"DJF",
	"GNF",
	"JPY",
	"KMF",
	"KRW",
	"MGA",
	"PYG",
	"RWF",
	"UGX",
	"VND",
	"VUV",
	"XAF",
	"XOF",
	"XPF",
)

// ThreeDecimalCurrencies is a set of currencies that have three decimal places.
//
// https://stripe.com/docs/currencies#three-decimal
var ThreeDecimalCurrencies = mapset.NewThreadUnsafeSet(
	"BHD",
	"JOD",
	"KWD",
	"OMR",
	"TND",
)
