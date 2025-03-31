package api

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"route_master/aggregator"
	"route_master/commons"
	"route_master/config/types"
	"route_master/db"
	"route_master/metrics"
	"route_master/networking"
	"route_master/utils"
)

type RouteAPI struct {
	utils.Loggable
	chainId     *big.Int
	quoteDB     *db.Quote
	aggregator  *aggregator.Aggregator
	solidClient *networking.SolidEthClient
}

func NewRouteAPI(chainId *big.Int, aggregator *aggregator.Aggregator, solidClient *networking.SolidEthClient, quoteDB *db.Quote) *RouteAPI {
	return &RouteAPI{
		chainId:     chainId,
		quoteDB:     quoteDB,
		aggregator:  aggregator,
		solidClient: solidClient,
		Loggable:    utils.Loggable{ModuleName: "route"},
	}
}

func (api *RouteAPI) Quote(c *gin.Context) {
	fromToken := c.Query("fromToken")
	if !common.IsHexAddress(fromToken) {
		api.Info("invalid parameter fromToken", "val", fromToken)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter fromToken=%s", fromToken)})
		return
	}

	toToken := c.Query("toToken")
	if !common.IsHexAddress(toToken) {
		api.Info("invalid parameter toToken", "val", toToken)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter toToken=%s", toToken)})
		return
	}

	fromTokenAddr := common.HexToAddress(fromToken)
	toTokenAddr := common.HexToAddress(toToken)

	if types.IsETHAddress(fromTokenAddr) {
		fromToken = types.WETHAddress.String()
		api.Info("replaced ETH address with WETH", "fromToken", fromToken)
	}

	if types.IsETHAddress(toTokenAddr) {
		toToken = types.WETHAddress.String()
		api.Info("replaced ETH address with WETH", "toToken", toToken)
	}

	//query token info from cache/ethclient
	fromTokenInfo, err := api.GetTokenInfo(fromToken)
	if err != nil {
		api.Info("fromToken not found", "val", fromToken)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("fromToken not found, fromToken=%s", fromToken)})
		return
	}
	toTokenInfo, err := api.GetTokenInfo(toToken)
	if err != nil {
		api.Info("toToken not found", "val", fromToken)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("toToken not found, toToken=%s", fromToken)})
		return
	}

	amount := c.Query("amount")
	amountNum, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		api.Info("invalid parameter amount", "val", amount)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter amount=%s", amount)})
		return
	}
	if amountNum.Cmp(new(big.Int).SetInt64(0)) <= 0 {
		api.Info("invalid parameter amount should be an positive value", "val", amount)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter amount=%s, should be an positive value", amount)})
		return
	}

	userAddress := c.Query("userAddress")
	if !common.IsHexAddress(userAddress) {
		api.Info("invalid parameter userAddress", "val", userAddress)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter userAddress=%s", userAddress)})
		return
	}
	slippageStr := c.DefaultQuery("slippage", "100")
	slippage, err := strconv.ParseInt(slippageStr, 10, 64)
	if err != nil {
		api.Error("failed to parse slippage", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse slippage: %v", err)})
		return
	}
	if slippage < 0 || slippage > 10000 {
		api.Error("slippage must be between 0 and 10000 bps", "slippage", slippage)
		c.JSON(http.StatusBadRequest, gin.H{"error": "slippage must be between 0 and 10000 bps"})
		return
	}

	requestId := uuid.NewString()
	response := api.aggregator.GetQuote(fromTokenInfo, toTokenInfo, amountNum, requestId, common.HexToAddress(userAddress), slippage)
	if response == nil {
		metrics.QuoteStatusCounter(fromToken, toToken, "failed")
		api.Info("failed to retrieve quote", "fromToken", fromToken, "toToken", toToken, "amount", amountNum.String(), "requestId", requestId)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no quote"})
		return
	}

	go func() {
		if err := api.quoteDB.CreateQuote(response.ToQuoteRecord()); err != nil {
			api.Error("creat quote record error", "err", err)
		}
	}()

	metrics.QuoteStatusCounter(fromToken, toToken, "success")
	c.JSON(http.StatusOK, response.FormatQuoteResponse())
}

func (api *RouteAPI) QuotesHistory(c *gin.Context) {
	fromToken := c.Query("fromToken")
	if fromToken != "" && !common.IsHexAddress(fromToken) {
		api.Info("invalid parameter fromToken", "val", fromToken)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter fromToken=%s", fromToken)})
		return
	}
	toToken := c.Query("toToken")
	if toToken != "" && !common.IsHexAddress(toToken) {
		api.Info("invalid parameter toToken", "val", toToken)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter toToken=%s", toToken)})
		return
	}
	amount := c.Query("amount")
	amountNum, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		api.Info("invalid parameter amount", "val", amount)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter amount=%s", amount)})
		return
	}
	if amountNum.Cmp(new(big.Int).SetInt64(0)) <= 0 {
		api.Info("invalid parameter amount should be an positive value", "val", amount)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter amount=%s, should be an positive value", amount)})
		return
	}
	userAddress := c.Query("userAddress")

	limit := c.GetInt("limit")
	offset := c.GetInt("offset")
	if limit == 0 {
		limit = 10
	}
	if offset == 0 {
		offset = 0
	}

	list, err := api.quoteDB.GetQuoteList(fromToken, toToken, amount, userAddress, limit, offset)
	if err != nil {
		api.Info("get quote list error", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (api *RouteAPI) GetTokenInfo(tokenAddress string) (commons.Token, error) {
	client, err := api.solidClient.Client()
	if err != nil {
		return commons.Token{}, err
	}
	tokens := commons.NewTokens(client)
	return tokens.GetToken(common.HexToAddress(tokenAddress)), nil
}
