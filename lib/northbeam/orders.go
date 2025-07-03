package northbeam

import (
	"context"
	"time"

	"github.com/igrmk/decimal"
)

type Decimal decimal.Decimal

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(decimal.Decimal(d).String()), nil
}

func (d *Decimal) UnmarshalJSON(data []byte) (err error) {
	var dec decimal.Decimal
	if err = dec.UnmarshalJSON(data); err == nil {
		*d = Decimal(dec)
	}
	return
}

type Order struct {
	OrderID                 string           `json:"order_id"`
	CustomerID              string           `json:"customer_id"`
	TimeOfPurchase          time.Time        `json:"time_of_purchase"`
	CustomerEmail           *string          `json:"customer_email,omitempty"`
	CustomerPhoneNumber     *string          `json:"customer_phone_number,omitempty"`
	CustomerName            *string          `json:"customer_name,omitempty"`
	CustomerIP              *string          `json:"customer_ip_address,omitempty"`
	DiscountCodes           []string         `json:"discount_codes,omitempty"`
	DiscountAmount          *Decimal         `json:"discount_amount,omitempty"`
	OrderTags               []string         `json:"order_tags,omitempty"`
	Tax                     Decimal          `json:"tax"`
	IsRecurringOrder        *bool            `json:"is_recurring_order,omitempty"`
	Currency                string           `json:"currency"`
	PurchaseTotal           Decimal          `json:"purchase_total"`
	Products                []Product        `json:"products"`
	Refunds                 []Refund         `json:"refunds,omitempty"`
	CustomerShippingAddress *ShippingAddress `json:"customer_shipping_address,omitempty"`
	ShippingCost            Decimal          `json:"shipping_cost"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Quantity    Decimal `json:"quantity"`
	Price       Decimal `json:"price"`
	VariantID   *string `json:"variant_id,omitempty"`
	VariantName *string `json:"variant_name,omitempty"`
}

type Refund struct {
	ProductID    string  `json:"product_id"`
	Quantity     Decimal `json:"quantity"`
	RefundAmount Decimal `json:"refund_amount"`
	RefundCost   Decimal `json:"refund_cost"`
	RefundMadeAt string  `json:"refund_made_at"`
	VariantID    *string `json:"variant_id,omitempty"`
}

type ShippingAddress struct {
	Address1    string  `json:"address1"`
	Address2    *string `json:"address2,omitempty"`
	City        string  `json:"city"`
	State       *string `json:"state,omitempty"`
	Zip         *string `json:"zip,omitempty"`
	CountryCode string  `json:"country_code"`
}

func (c *Client) ReportOrders(ctx context.Context, req []Order) error {
	return c.client.Do(ctx, "POST", "/v2/orders", nil, req, nil)
}
