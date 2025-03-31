package types

import (
	"math/big"
	"time"
)

// Quote represents a quote response
type Quote struct {
	RequestId string
	Timestamp int64
	Prices    map[string]string
	Routes    map[string]interface{}
	Protocol  string
	Path      []string
	AmountIn  *big.Int
	AmountOut *big.Int
	Calldata  []byte
	To        string
}

// QuoteCacheEntry represents a cached quote with expiration time
type QuoteCacheEntry struct {
	Quote      *Quote
	Expiration time.Time
}

// QuoteResult represents the result of a quote operation
type QuoteResult struct {
	Quoter        interface{}
	Quote         *Quote
	SegmentOutput *big.Float
	Err           error
}

// QuoterState represents the state of a quoter
type QuoterState struct {
	AccumulatedAmountIn  *big.Int
	AccumulatedAmountOut *big.Float
	LastQuote            *Quote
	LastQuoteAmount      *big.Int
	Settlements          []map[string]interface{}
}

// SegmentDetail represents a segment of the quote
type SegmentDetail struct {
	Protocol       string
	AmountIn       string
	AmountOut      string
	TotalAmountIn  string
	TotalAmountOut string
	Calldata       []byte
	Routes         map[string]interface{}
	To             string
}

// Settlement represents a settlement in the quote
type Settlement struct {
	Settler string
	Value   string
	Data    []byte
}

// Route represents a route in the quote
type Route struct {
	Protocol       string
	TotalAmountIn  string
	TotalAmountOut string
}
