package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SMID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OOFShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (o *Order) Validate() error {
	if o == nil {
		return errors.New("order is nil")
	}

	if strings.TrimSpace(o.OrderUID) == "" {
		return fmt.Errorf("order: empty order_uid")
	}
	if strings.TrimSpace(o.TrackNumber) == "" {
		return fmt.Errorf("order: empty track_number")
	}
	if strings.TrimSpace(o.CustomerID) == "" {
		return fmt.Errorf("order: empty customer_id")
	}

	if o.DateCreated.IsZero() {
		return fmt.Errorf("order: date_created is zero")
	}
	now := time.Now()
	tenYearsAgo := now.AddDate(-10, 0, 0)
	if o.DateCreated.Before(tenYearsAgo) {
		return fmt.Errorf("order: date_created too old: %v", o.DateCreated)
	}
	if o.DateCreated.After(now.Add(1 * time.Hour)) {
		return fmt.Errorf("order: date_created is in the future: %v", o.DateCreated)
	}

	if err := o.Delivery.Validate(); err != nil {
		return fmt.Errorf("order: invalid delivery: %w", err)
	}

	if err := o.Payment.Validate(); err != nil {
		return fmt.Errorf("order: invalid payment: %w", err)
	}

	if len(o.Items) == 0 {
		return errors.New("order: must contain at least one item")
	}
	for i, item := range o.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("order: invalid item[%d]: %w", i, err)
		}
	}

	return nil
}

func (d *Delivery) Validate() error {
	if d == nil {
		return errors.New("delivery is nil")
	}
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("delivery: name is empty")
	}
	if strings.TrimSpace(d.Phone) == "" {
		return errors.New("delivery: phone is empty")
	}
	if strings.TrimSpace(d.City) == "" {
		return errors.New("delivery: city is empty")
	}
	if strings.TrimSpace(d.Address) == "" {
		return errors.New("delivery: address is empty")
	}
	if !strings.Contains(d.Email, "@") {
		return fmt.Errorf("delivery: invalid email %q", d.Email)
	}
	return nil
}

func (p *Payment) Validate() error {
	if p == nil {
		return errors.New("payment is nil")
	}
	if p.Amount <= 0 {
		return fmt.Errorf("payment: invalid amount %d", p.Amount)
	}
	if strings.TrimSpace(p.Currency) == "" {
		return errors.New("payment: currency is empty")
	}
	if strings.TrimSpace(p.Provider) == "" {
		return errors.New("payment: provider is empty")
	}
	if strings.TrimSpace(p.Transaction) == "" {
		return errors.New("payment: transaction is empty")
	}
	return nil
}

func (i *Item) Validate() error {
	if i == nil {
		return errors.New("item is nil")
	}
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("item: name is empty")
	}
	if i.Price <= 0 {
		return fmt.Errorf("item: invalid price %d", i.Price)
	}
	if i.TotalPrice < i.Price-i.Sale {
		return fmt.Errorf("item: total price seems inconsistent (price=%d, sale=%d, total=%d)",
			i.Price, i.Sale, i.TotalPrice)
	}
	if strings.TrimSpace(i.Brand) == "" {
		return errors.New("item: brand is empty")
	}
	return nil
}
