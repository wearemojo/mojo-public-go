package stripeconst

import (
	mapset "github.com/deckarep/golang-set/v2"
)

// EuropeanCountries is the set of countries that Stripe considers to be part of
// Europe.
//
// https://stripe.com/docs/currencies#european-credit-cards
var EuropeanCountries = mapset.NewThreadUnsafeSet(
	"AD",
	"AT",
	"BE",
	"BG",
	"CY",
	"CZ",
	"DE",
	"DK",
	"EE",
	"ES",
	"FI",
	"FO",
	"FR",
	"GG",
	"GI",
	"GL",
	"GR",
	// Stripe has removed Croatia from the list, but
	// it seems like this might have been a mistake
	//
	// TODO: Check with Stripe if this is intentional
	//
	// "HR",
	"HU",
	"IE",
	"IL",
	"IM",
	"IS",
	"IT",
	"JE",
	"LI",
	"LT",
	"LU",
	"LV",
	"MC",
	"MT",
	"NL",
	"NO",
	"PL",
	"PM",
	"PT",
	"RO",
	"SE",
	"SI",
	"SK",
	"SM",
	"TR",
	"VA",
	"GB",
)
