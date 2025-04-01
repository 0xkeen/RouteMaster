package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"route_master/aggregator/quoter"
)

type BebopQuoteRequest struct {
	SellTokens     []string
	BuyTokens      []string
	SellAmounts    []string
	TakerAddress   string
	SkipValidation bool
	ApprovalType   string
	Gasless        bool
	Receiver       string
	Source         string
}

type BebopQuoteResponse struct {
	Type         string  `json:"type"`
	Status       string  `json:"status"`
	QuoteId      string  `json:"quoteId"`
	ChainId      int     `json:"chainId"`
	ApprovalType string  `json:"approvalType"`
	NativeToken  string  `json:"nativeToken"`
	Taker        string  `json:"taker"`
	Receiver     string  `json:"receiver"`
	Expiry       int     `json:"expiry"`
	Slippage     float64 `json:"slippage"`
	GasFee       struct {
		Native string  `json:"native"`
		Usd    float64 `json:"usd"`
	} `json:"gasFee"`
	BuyTokens map[string]struct {
		Amount            string  `json:"amount"`
		Decimals          int     `json:"decimals"`
		PriceUsd          float64 `json:"priceUsd"`
		Symbol            string  `json:"symbol"`
		MinimumAmount     string  `json:"minimumAmount"`
		Price             float64 `json:"price"`
		PriceBeforeFee    float64 `json:"priceBeforeFee"`
		AmountBeforeFee   string  `json:"amountBeforeFee"`
		DeltaFromExpected float64 `json:"deltaFromExpected"`
	} `json:"buyTokens"`
	SellTokens map[string]struct {
		Amount         string  `json:"amount"`
		Decimals       int     `json:"decimals"`
		PriceUsd       float64 `json:"priceUsd"`
		Symbol         string  `json:"symbol"`
		Price          float64 `json:"price"`
		PriceBeforeFee float64 `json:"priceBeforeFee"`
	} `json:"sellTokens"`
	SettlementAddress  string             `json:"settlementAddress"`
	ApprovalTarget     string             `json:"approvalTarget"`
	RequiredSignatures []interface{}      `json:"requiredSignatures"`
	PriceImpact        float64            `json:"priceImpact"`
	PartnerFeeNative   string             `json:"partnerFeeNative"`
	Warnings           []interface{}      `json:"warnings"`
	Tx                 quoter.Transaction `json:"tx"`
	Makers             []string           `json:"makers"`
	ToSign             struct {
		PartnerId      int    `json:"partner_id"`
		Expiry         int    `json:"expiry"`
		TakerAddress   string `json:"taker_address"`
		MakerAddress   string `json:"maker_address"`
		MakerNonce     string `json:"maker_nonce"`
		TakerToken     string `json:"taker_token"`
		MakerToken     string `json:"maker_token"`
		TakerAmount    string `json:"taker_amount"`
		MakerAmount    string `json:"maker_amount"`
		Receiver       string `json:"receiver"`
		PackedCommands string `json:"packed_commands"`
	} `json:"toSign"`
	OnchainOrderType  string `json:"onchainOrderType"`
	PartialFillOffset int    `json:"partialFillOffset"`
}

type BebopClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewBebopClient() *BebopClient {
	return &BebopClient{
		baseURL: "https://api.bebop.xyz/pmm/ethereum/v3/quote",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (b *BebopClient) GetQuote(params BebopQuoteRequest) (*BebopQuoteResponse, error) {
	queryParams := url.Values{}

	// Required parameters
	queryParams.Add("sell_tokens", params.SellTokens[0])   // Assuming single token for simplicity
	queryParams.Add("buy_tokens", params.BuyTokens[0])     // Assuming single token for simplicity
	queryParams.Add("sell_amounts", params.SellAmounts[0]) // Assuming single amount for simplicity

	// Optional parameters with defaults
	queryParams.Add("skip_validation", fmt.Sprintf("%t", params.SkipValidation))
	queryParams.Add("approval_type", params.ApprovalType)
	queryParams.Add("gasless", fmt.Sprintf("%t", params.Gasless))

	// Optional parameters
	if params.TakerAddress != "" {
		queryParams.Add("taker_address", params.TakerAddress)
	}
	if params.Receiver != "" {
		queryParams.Add("receiver", params.Receiver)
	}
	if params.Source != "" {
		queryParams.Add("source", params.Source)
	}

	// Make the HTTP request
	quoteUrl := fmt.Sprintf("%s?%s", b.baseURL, queryParams.Encode())
	resp, err := b.httpClient.Get(quoteUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response BebopQuoteResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
