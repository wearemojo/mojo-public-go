package datahappy

type Integrations struct {
	Meta      *Meta      `json:"meta,omitempty"`
	GoogleAds *GoogleAds `json:"google_ads,omitempty"`
	// GA4       *GA4       `json:"ga4,omitempty"`
	// AppsFlyer *AppsFlyer `json:"appsflyer,omitempty"`
	// Adjust    *Adjust    `json:"adjust,omitempty"`
	// HubSpot   *HubSpot   `json:"hubspot,omitempty"`
}

type Meta struct {
	FBP string `json:"fbp,omitempty"`
	FBC string `json:"fbc,omitempty"`
}

type GoogleAds struct {
	GCLID  string `json:"gclid,omitempty"`
	GBRAID string `json:"gbraid,omitempty"`
	WBRAID string `json:"wbraid,omitempty"`
}
