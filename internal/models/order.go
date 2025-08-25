package models

import "time"

// Order — основная сущность заказа
type Order struct {
	ID                uint   `json:"id"`
	OrderUID          string `json:"order_uid" gorm:"type:uuid;uniqueIndex;not null"`
	TrackNumber       string `json:"track_number" gorm:"size:64;uniqueIndex;not null"`
	Entry             string `json:"entry" gorm:"size:10;not null"`
	Locale            string `json:"locale" gorm:"size:10;not null"`
	InternalSignature string `json:"internal_signature"`
	CustomerID        string `json:"customer_id" gorm:"size:128;not null"`
	DeliveryService   string `json:"delivery_service" gorm:"size:64;not null"`
	Shardkey          string `json:"shardkey" gorm:"size:2;not null"`
	SMID              int    `json:"sm_id" gorm:"not null"`
	DateCreated       string `json:"date_created" gorm:"not null"`
	OofShard          string `json:"oof_shard" gorm:"size:2;not null"`

	// Ассоциации
	Delivery *Delivery `json:"delivery" gorm:"foreignKey:OrderID;references:OrderUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Payment  *Payment  `json:"payment" gorm:"foreignKey:Transaction;references:OrderUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Items    []Item    `json:"items" gorm:"foreignKey:OrderID;references:OrderUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// GORM timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// Delivery — данные доставки
type Delivery struct {
	OrderID string `json:"-" gorm:"type:uuid;primaryKey;not null"`
	Name    string `json:"name" gorm:"size:128;not null"`
	Phone   string `json:"phone" gorm:"size:15;not null"`
	Zip     string `json:"zip" gorm:"size:10;not null"`
	City    string `json:"city" gorm:"size:64;not null"`
	Address string `json:"address" gorm:"size:64;not null"`
	Region  string `json:"region" gorm:"size:64;not null"`
	Email   string `json:"email" gorm:"size:128"`
}

func (Delivery) TableName() string {
	return "delivery"
}

// Payment — данные оплаты
type Payment struct {
	Transaction   string    `json:"transaction" gorm:"type:uuid;primaryKey;not null"`
	OrderID       string    `json:"-" gorm:"type:uuid;uniqueIndex;not null"`
	RequestID     string    `json:"request_id" gorm:"size:64"`
	Currency      string    `json:"currency" gorm:"size:10;not null"`
	Provider      string    `json:"provider" gorm:"size:32;not null"`
	Amount        int       `json:"amount" gorm:"not null"`
	PaymentDt     int64     `json:"payment_dt" gorm:"-"`
	PaymentDtTime time.Time `gorm:"column:payment_dt;not null"`
	Bank          string    `json:"bank" gorm:"size:20;not null"`
	DeliveryCost  int       `json:"delivery_cost" gorm:"not null"`
	GoodsTotal    int       `json:"goods_total" gorm:"not null"`
	CustomFee     int       `json:"custom_fee" gorm:"not null"`
}

func (Payment) TableName() string {
	return "payment"
}

// Item — товар в заказе
type Item struct {
	ChrtID      int64  `json:"chrt_id" gorm:"primaryKey"`
	OrderID     string `json:"-" gorm:"type:uuid;index;not null"`
	TrackNumber string `json:"track_number" gorm:"size:64;not null"`
	Price       int    `json:"price" gorm:"not null"`
	RID         string `json:"rid" gorm:"column:rid;not null"`
	Name        string `json:"name" gorm:"size:128;not null"`
	Sale        int    `json:"sale" gorm:"not null;check:sale >= 0 AND sale <= 100"`
	Size        string `json:"size" gorm:"size:16;not null"`
	TotalPrice  int    `json:"total_price" gorm:"not null"`
	NMID        int64  `json:"nm_id" gorm:"not null"`
	Brand       string `json:"brand" gorm:"size:64;not null"`
	Status      int    `json:"status" gorm:"not null"`
}
