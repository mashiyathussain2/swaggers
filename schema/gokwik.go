package schema

import "time"

//Define name of the cart collection
const (
	RTOColl string = "rto_data"
)

type CheckCODEligiblityOpts struct {
	Customer GoKwikCustomer `json:"customer"`
	Order    GoKwikOrder    `json:"order"`
}
type GoKwikCustomer struct {
	Age           int     `json:"age"`
	Gender        string  `json:"gender"`
	WalletBalance float64 `json:"wallet_balance"`
	CustomerSince int     `json:"customer_since"`
}
type GoKwikLineItems struct {
	ProductID           string  `json:"product_id"`
	LineItemID          string  `json:"line_item_id"`
	ItemBrand           string  `json:"item_brand"`
	ItemRating          float64 `json:"item_rating"`
	ItemSize            string  `json:"item_size"`
	ItemColor           string  `json:"item_color"`
	IsExclusive         bool    `json:"is_exclusive"`
	ItemWeight          int     `json:"item_weight"`
	ItemLength          int     `json:"item_length"`
	ItemBreadth         int     `json:"item_breadth"`
	ItemHeight          int     `json:"item_height"`
	ItemDiscount        int     `json:"item_discount"`
	VariantID           string  `json:"variant_id"`
	Name                string  `json:"name"`
	Sku                 string  `json:"sku"`
	Price               float64 `json:"price"`
	Quantity            int     `json:"quantity"`
	Subtotal            float64 `json:"subtotal"`
	Total               int     `json:"total"`
	Tax                 int     `json:"tax"`
	ProductURL          string  `json:"product_url"`
	ProductThumbnailURL string  `json:"product_thumbnail_url"`
	ArticlePrice        float64 `json:"article_price"`
	TargetGroup         string  `json:"target_group"`
	SubCategory         string  `json:"sub_category"`
	MajorCategory       string  `json:"major_category"`
}
type GoKwikShippingAddress struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Type      string `json:"type"`
}
type GoKwikBillingAddress struct {
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Postcode string `json:"postcode"`
	Type     string `json:"type"`
}
type GoKwikSession struct {
	Source            string   `json:"source"`
	SessionHistory    []string `json:"session_history"`
	SessionLength     int      `json:"session_length"`
	TotalPagesViewed  int      `json:"total_pages_viewed"`
	CustomerIP        string   `json:"customer_ip"`
	CustomerUserAgent string   `json:"customer_user_agent"`
}
type GoKwikOrder struct {
	OrderDate              time.Time             `json:"order_date"`
	Subtotal               int                   `json:"subtotal"`
	TotalLineItems         int                   `json:"total_line_items"`
	TotalLineItemsQuantity int                   `json:"total_line_items_quantity"`
	TotalTax               int                   `json:"total_tax"`
	TotalShipping          int                   `json:"total_shipping"`
	TotalDiscount          int                   `json:"total_discount"`
	Total                  int                   `json:"total"`
	PromoCode              string                `json:"promo_code"`
	LineItems              []GoKwikLineItems     `json:"line_items"`
	ShippingAddress        GoKwikShippingAddress `json:"shipping_address"`
	BillingAddress         GoKwikBillingAddress  `json:"billing_address"`
	Session                GoKwikSession         `json:"session"`
}

type RTOResp struct {
	StatusCode int         `json:"status_code"`
	Msg        string      `json:"msg"`
	Error      interface{} `json:"error"`
	Data       Data        `json:"data"`
}
type Data struct {
	RequestID string      `json:"request_id"`
	Score     int         `json:"score"`
	RiskFlag  interface{} `json:"risk_flag"`
	Reason    interface{} `json:"reason"`
}
