package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Quote struct {
	db *gorm.DB
}

func NewQuote(dsn string) *Quote {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&QuoteRecord{})
	return &Quote{db: db}
}

type QuoteRecord struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FromToken   string    `gorm:"type:varchar(64);not null;index:idx_quotes_from_token_to_token_amount" json:"from_token"`
	ToToken     string    `gorm:"type:varchar(64);not null;index:idx_quotes_from_token_to_token_amount" json:"to_token"`
	AmountIn    string    `gorm:"type:varchar(100);not null;index:idx_quotes_from_token_to_token_amount" json:"amount_in"`
	UserAddress string    `gorm:"type:varchar(64);not null;index:idx_quotes_user_address" json:"user_address"`
	Routes      string    `gorm:"type:longtext;not null" json:"routes"`
	CreatedAt   time.Time `gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_quotes_created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (QuoteRecord) TableName() string {
	return "quotes"
}

func (q *Quote) CreateQuote(quote *QuoteRecord) error {
	return q.db.Create(quote).Error
}

func (q *Quote) GetQuoteList(fromToken, toToken, amount, userAddress string, limit, offset int) ([]*QuoteRecord, error) {
	query := q.db.Table("quotes")
	if fromToken != "" {
		query = query.Where("from_token = ?", fromToken)
	}
	if toToken != "" {
		query = query.Where("to_token = ?", toToken)
	}
	if amount != "" {
		query = query.Where("amount_in = ?", amount)
	}
	if userAddress != "" {
		query = query.Where("user_address = ?", userAddress)
	}

	var records []*QuoteRecord
	err := query.Order("id desc").Limit(limit).Offset(offset).Scan(&records).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(records)

	return records, nil
}
