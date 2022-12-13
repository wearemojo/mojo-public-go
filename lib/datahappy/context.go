//nolint:tagliatelle // datahappy uses camel case
package datahappy

import (
	"time"

	"cloud.google.com/go/civil"
)

type Context struct {
	Library *Library `json:"library,omitempty"`
	// Campaign   *Campaign `json:"campaign,omitempty"`
	// Device     *Device   `json:"device,omitempty"`
	App    *App   `json:"app,omitempty"`
	IP     string `json:"ip,omitempty"`
	Locale string `json:"locale,omitempty"`
	// Location   *Location `json:"location,omitempty"`
	// Network    *Network  `json:"network,omitempty"`
	// OS         *OS       `json:"os,omitempty"`
	// Page       *Page     `json:"page,omitempty"`
	// Referrer   *Referrer `json:"referrer,omitempty"`
	// Screen     *Screen   `json:"screen,omitempty"`
	Timezone   string  `json:"timezone,omitempty"`
	GroupID    string  `json:"groupId,omitempty"`
	Traits     *Traits `json:"traits,omitempty"`
	UserAgent  string  `json:"userAgent,omitempty"`
	ExternalID string  `json:"externalId,omitempty"`
}

type Library struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type App struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Build   string `json:"build,omitempty"`
}

type Traits struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Title     string `json:"title,omitempty"`
	// Company     *Company    `json:"company,omitempty"`
	Website     string      `json:"website,omitempty"`
	Description string      `json:"description,omitempty"`
	Gender      string      `json:"gender,omitempty"`
	Birthday    *civil.Date `json:"birthday,omitempty"`
	// Address     *Address    `json:"address,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}
