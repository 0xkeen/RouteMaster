package quoter

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"route_master/commons"
	"route_master/db"
)

type QuoteCacheEntry struct {
	Quote      *Quote
	Expiration time.Time
}

type QuoteResult struct {
	Quoter        Quoter
	Quote         *Quote
	SegmentOutput *big.Float
	Err           error
}

type Quote struct {
	RequestId     string                 `json:"RequestId"`
	FromToken     common.Address         `json:"FromToken"`
	ToToken       common.Address         `json:"ToToken"`
	UserAddress   common.Address         `json:"UserAddress"`
	To            common.Address         `json:"To"`
	AmountIn      *big.Int               `json:"AmountIn"`
	AmountOut     *big.Int               `json:"AmountOut"`
	MiniAmountOut *big.Int               `json:"miniAmountOut"`
	Prices        map[string]string      `json:"Prices"`
	Routes        map[string]interface{} `json:"Routes"`
	Protocol      string                 `json:"Protocol"`
	Path          []common.Address       `json:"Path"`
	Tx            Transaction            `json:"tx"`
	Calldata      []byte                 `json:"Calldata"`
	Timestamp     int64                  `json:"Timestamp"`
}

func (q *Quote) FormatQuoteResponse() interface{} {
	return map[string]interface{}{
		"RequestId":     q.RequestId,
		"FromToken":     q.FromToken.String(),
		"ToToken":       q.ToToken.String(),
		"UserAddress":   q.UserAddress.String(),
		"AmountIn":      q.AmountIn.String(),
		"AmountOut":     q.AmountOut.String(),
		"MiniAmountOut": q.MiniAmountOut.String(),
		"Prices":        q.Prices,
		"Routes":        q.Routes,
		"Protocol":      q.Protocol,
		"Path":          q.Path,
		"Tx":            q.Tx,
		"Calldata":      q.Calldata,
		"Timestamp":     q.Timestamp,
	}
}

func (q *Quote) ToQuoteRecord() *db.QuoteRecord {
	route, _ := json.Marshal(q.Routes)
	return &db.QuoteRecord{
		FromToken:   q.FromToken.String(),
		ToToken:     q.ToToken.String(),
		AmountIn:    q.AmountIn.String(),
		UserAddress: q.UserAddress.String(),
		Routes:      string(route),
	}
}

// Quoter defines the interface that all quoters must implement
type Quoter interface {
	String() string
	// ExactIn gets quotes for exact input amounts
	ExactIn(
		fromToken commons.Token,
		toToken commons.Token,
		amountIn string,
		requestId string,
		route interface{},
		filler common.Address,
	) (*Quote, error)

	// ExactOut gets quotes for exact output amounts
	ExactOut(
		fromToken commons.Token,
		toToken commons.Token,
		amountOut string,
		requestId string,
		route interface{},
		filler common.Address,
	) (*Quote, error)
}

type Transaction struct {
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	From     string `json:"from"`
	Gas      int    `json:"gas"`
	GasPrice int64  `json:"gasPrice"`
}

type QuoterState struct {
	AccumulatedAmountIn  *big.Int
	AccumulatedAmountOut *big.Float
	LastQuote            *Quote
	LastQuoteAmount      *big.Int
	Settlements          []map[string]interface{}
}
