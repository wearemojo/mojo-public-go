package northbeam

import (
	"context"
	"time"

	"github.com/igrmk/decimal"
)

type Order struct {
	OrderID                 string           `json:"order_id"`
	CustomerID              string           `json:"customer_id"`
	TimeOfPurchase          time.Time        `json:"time_of_purchase"`
	CustomerEmail           *string          `json:"customer_email"`
	CustomerPhoneNumber     *string          `json:"customer_phone_number"`
	CustomerName            *string          `json:"customer_name"`
	CustomerIP              *string          `json:"customer_ip_address"`
	DiscountCodes           []string         `json:"discount_codes"`
	DiscountAmount          *decimal.Decimal `json:"discount_amount"`
	OrderTags               []string         `json:"order_tags"`
	Tax                     decimal.Decimal  `json:"tax"`
	IsRecurringOrder        *bool            `json:"is_recurring_order"`
	Currency                string           `json:"currency"`
	PurchaseTotal           decimal.Decimal  `json:"purchase_total"`
	Products                []Product        `json:"products"`
	Refunds                 []Refund         `json:"refunds"`
	CustomerShippingAddress *ShippingAddress `json:"customer_shipping_address"`
	ShippingCost            decimal.Decimal  `json:"shipping_cost"`
}

type Product struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Quantity    decimal.Decimal `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	VariantID   *string         `json:"variant_id"`
	VariantName *string         `json:"variant_name"`
}

type Refund struct {
	ProductID    string          `json:"product_id"`
	Quantity     decimal.Decimal `json:"quantity"`
	RefundAmount decimal.Decimal `json:"refund_amount"`
	RefundCost   decimal.Decimal `json:"refund_cost"`
	RefundMadeAt string          `json:"refund_made_at"`
	VariantID    *string         `json:"variant_id"`
}

type ShippingAddress struct {
	Address1    string  `json:"address1"`
	Address2    *string `json:"address2"`
	City        string  `json:"city"`
	State       *string `json:"state"`
	Zip         *string `json:"zip"`
	CountryCode string  `json:"country_code"`
}

func (c *Client) ReportOrders(ctx context.Context, req []Order) error {
	return c.client.Do(ctx, "POST", "/v2/orders", nil, req, nil)
}
