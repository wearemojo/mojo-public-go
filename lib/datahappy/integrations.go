package datahappy

type Integrations struct {
	CAPI *CAPI `json:"capi,omitempty"`
	GAds *GAds `json:"gads,omitempty"`
	// GA4       *GA4       `json:"ga4,omitempty"`
	// AppsFlyer *AppsFlyer `json:"appsflyer,omitempty"`
	// Adjust    *Adjust    `json:"adjust,omitempty"`
	// HubSpot   *HubSpot   `json:"hubspot,omitempty"`
}

type CAPI struct {
	FBP string `json:"fbp,omitempty"`
	FBC string `json:"fbc,omitempty"`
}

type GAds struct {
	GCLID  string `json:"gclid,omitempty"`
	GBRAID string `json:"gbraid,omitempty"`
	WBRAID string `json:"wbraid,omitempty"`
}
