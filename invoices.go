package stripe

import (
	"fmt"
	"net/url"
	"strconv"
)

type InvoiceLineType string

const (
	TypeInvoiceItem  InvoiceLineType = "invoiceitem"
	TypeSubscription InvoiceLineType = "subscription"
)

type InvoiceParams struct {
	Customer             string
	Desc, Statement, Sub string
	Fee                  uint64
	Meta                 map[string]string
	Closed               bool
}

type Invoice struct {
	Id          string            `json:"id"`
	Live        bool              `json:"livemode"`
	Amount      int64             `json:"amount_due"`
	Attempts    uint64            `json:"attempt_count"`
	Attempted   bool              `json:"attempted"`
	Closed      bool              `json:"closed"`
	Currency    Currency          `json:"currency"`
	Customer    string            `json:"customer"`
	Date        int64             `json:"date"`
	Lines       *InvoiceLineList  `json:lines"`
	Paid        bool              `json:"paid"`
	End         int64             `json:"period_end"`
	Start       int64             `json:"period_start"`
	Subtotal    int64             `json:"subtotal"`
	Total       int64             `json:"total"`
	Fee         uint64            `json:"application_fee"`
	Charge      string            `json:"charge"`
	Desc        string            `json:"description"`
	Discount    *Discount         `json:"discount"`
	Balance     int64             `json:"ending_balance"`
	NextAttempt int64             `json:"next_payment_attempt"`
	Statement   string            `json:"statement_description"`
	Sub         string            `json:"subscription"`
	Webhook     int64             `json:"webhooks_delivered_at"`
	Meta        map[string]string `json:"metadata"`
}

type InvoiceLine struct {
	Id        string            `json:"id"`
	Live      bool              `json:"live_mode"`
	Amount    int64             `json:"amount"`
	Currency  Currency          `json:"currency"`
	Period    *Period           `json:"period"`
	Proration bool              `json:"proration"`
	Type      InvoiceLineType   `json:"type"`
	Desc      string            `json:"description"`
	Meta      map[string]string `json:"metadata"`
	Plan      *Plan             `json:"plan"`
	Quantity  int64             `json:"quantity"`
}

type Period struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type InvoiceLineList struct {
	Count  uint16         `json:"total_count"`
	More   bool           `json:"has_more"`
	Url    string         `json:"url"`
	Values []*InvoiceLine `json:"data"`
}

type InvoiceClient struct {
	api   Api
	token string
}

func (c *InvoiceClient) Create(params *InvoiceParams) (*Invoice, error) {
	body := &url.Values{
		"customer": {params.Customer},
	}

	if len(params.Desc) > 0 {
		body.Add("description", params.Desc)
	}

	if len(params.Statement) > 0 {
		body.Add("statement_description", params.Statement)
	}

	if len(params.Sub) > 0 {
		body.Add("subscription", params.Sub)
	}

	if params.Fee > 0 {
		body.Add("application_fee", strconv.FormatUint(params.Fee, 10))
	}

	for k, v := range params.Meta {
		body.Add(fmt.Sprintf("metadata[%v]", k), v)
	}

	invoice := &Invoice{}
	err := c.api.Call("POST", "/invoices", c.token, body, invoice)

	return invoice, err
}

func (c *InvoiceClient) Get(id string) (*Invoice, error) {
	invoice := &Invoice{}
	err := c.api.Call("GET", "/invoices/"+id, c.token, nil, invoice)

	return invoice, err
}

func (c *InvoiceClient) Pay(id string) (*Invoice, error) {
	invoice := &Invoice{}
	err := c.api.Call("POST", fmt.Sprintf("/invoices/%v/pay", id), c.token, nil, invoice)

	return invoice, err
}

func (c *InvoiceClient) Update(id string, params *InvoiceParams) (*Invoice, error) {
	body := &url.Values{}

	if len(params.Desc) > 0 {
		body.Add("description", params.Desc)
	}

	if len(params.Statement) > 0 {
		body.Add("statement_description", params.Statement)
	}

	if len(params.Sub) > 0 {
		body.Add("subscription", params.Sub)
	}

	if params.Fee > 0 {
		body.Add("application_fee", strconv.FormatUint(params.Fee, 10))
	}

	if params.Closed {
		body.Add("closed", strconv.FormatBool(true))
	}

	for k, v := range params.Meta {
		body.Add(fmt.Sprintf("metadata[%v]", k), v)
	}

	invoice := &Invoice{}
	err := c.api.Call("POST", "/invoices/"+id, c.token, body, invoice)

	return invoice, err
}

func (c *InvoiceClient) GetNext(params *InvoiceParams) (*Invoice, error) {
	body := &url.Values{
		"customer": {params.Customer},
	}

	if len(params.Sub) > 0 {
		body.Add("subscription", params.Sub)
	}

	invoice := &Invoice{}
	err := c.api.Call("GET", "/invoices", c.token, body, invoice)

	return invoice, err
}