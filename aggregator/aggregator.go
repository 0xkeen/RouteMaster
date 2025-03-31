package aggregator

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/common"

	"route_master/aggregator/quoter"
	"route_master/commons"
	"route_master/utils"
)

type Aggregator struct {
	utils.Loggable
	quoters        []quoter.Quoter
	segmentPercent int64
	quoteCache     map[string]map[string]*quoter.QuoteCacheEntry
	cacheMutex     sync.Mutex
}

func NewAggregator(quoters []quoter.Quoter, segmentPercentage int64) *Aggregator {
	quoteCache := make(map[string]map[string]*quoter.QuoteCacheEntry)
	for _, q := range quoters {
		quoteCache[q.String()] = make(map[string]*quoter.QuoteCacheEntry)
	}

	return &Aggregator{
		quoters:        quoters,
		segmentPercent: segmentPercentage,
		quoteCache:     quoteCache,
		cacheMutex:     sync.Mutex{},
		Loggable:       utils.Loggable{ModuleName: "aggregator"},
	}
}

// splitAmount splits the input amount into segments based on segmentPercent
func (a *Aggregator) splitAmount(inputAmount *big.Int) []*big.Int {
	segmentAmount := new(big.Int).Div(
		new(big.Int).Mul(inputAmount, big.NewInt(a.segmentPercent)),
		big.NewInt(100),
	)

	totalRounds := new(big.Int).Div(inputAmount, segmentAmount)

	var baseSegments []*big.Int
	remainderAmount := new(big.Int).Mod(inputAmount, segmentAmount)

	for i := int64(0); i < totalRounds.Int64(); i++ {
		baseSegments = append(baseSegments, new(big.Int).Set(segmentAmount))
	}

	if remainderAmount.Cmp(big.NewInt(0)) > 0 {
		if len(baseSegments) > 0 {
			lastIndex := len(baseSegments) - 1
			baseSegments[lastIndex] = new(big.Int).Add(baseSegments[lastIndex], remainderAmount)
		} else {
			baseSegments = append(baseSegments, new(big.Int).Set(remainderAmount))
		}
	}

	a.Debug("cache: created base segments", "count", len(baseSegments))
	for i, segment := range baseSegments {
		a.Debug("cache: base segment", "index", i, "amount", segment.String())
	}

	return baseSegments
}

// generateAmountCombinations generates all possible combinations of amounts to query
func (a *Aggregator) generateAmountCombinations(baseSegments []*big.Int) []*big.Int {
	amountsToQuery := make(map[string]bool)
	for _, segment := range baseSegments {
		amountsToQuery[segment.String()] = true
	}

	cumulativeAmount := big.NewInt(0)
	for _, segment := range baseSegments {
		cumulativeAmount = new(big.Int).Add(cumulativeAmount, segment)
		amountsToQuery[cumulativeAmount.String()] = true
	}

	if len(baseSegments) >= 2 {
		lastSegment := baseSegments[len(baseSegments)-1]
		for i := 0; i < len(baseSegments)-1; i++ {
			combinedAmount := new(big.Int).Add(lastSegment, baseSegments[i])
			amountsToQuery[combinedAmount.String()] = true

			tempAmount := new(big.Int).Set(combinedAmount)
			for j := i + 1; j < len(baseSegments)-1; j++ {
				tempAmount = new(big.Int).Add(tempAmount, baseSegments[j])
				amountsToQuery[tempAmount.String()] = true
			}
		}

		if len(baseSegments) >= 4 {
			for step := 2; step < len(baseSegments)-1; step++ {
				for start := 0; start+step < len(baseSegments)-1; start++ {
					skipCombination := new(big.Int).Add(lastSegment, baseSegments[start])
					skipCombination = new(big.Int).Add(skipCombination, baseSegments[start+step])
					amountsToQuery[skipCombination.String()] = true
				}
			}
		}
	}

	var amountsList []*big.Int
	for amountStr := range amountsToQuery {
		amount, _ := new(big.Int).SetString(amountStr, 10)
		amountsList = append(amountsList, amount)
	}

	a.Info("cache: generated combinations to query", "count", len(amountsList))
	for i, amount := range amountsList {
		a.Info("cache: combination", "index", i, "amount", amount.String())
	}

	return amountsList
}

// preQuoteAndCache performs quote requests and caches the results
func (a *Aggregator) preQuoteAndCache(
	fromToken commons.Token,
	toToken commons.Token,
	amountsList []*big.Int,
	requestId string,
	userAddress common.Address,
) {
	var wg sync.WaitGroup
	for _, q := range a.quoters {
		quoterName := q.String()
		for _, amount := range amountsList {
			amountStr := amount.String()

			wg.Add(1)
			go func(quoter quoter.Quoter, amountToQuery string) {
				defer wg.Done()

				quote, err := quoter.ExactIn(fromToken, toToken, amountToQuery, requestId, nil, userAddress)
				if err == nil && quote != nil && quote.Prices[amountToQuery] != "" {
					a.addQuoteToCache(quoterName, amountToQuery, quote)
					a.Info("cache: preQuote success",
						"requestId", requestId,
						"quoter", quoterName,
						"amount", amountToQuery,
						"price", quote.Prices[amountToQuery])
				} else {
					a.Error("cache: preQuote failed",
						"requestId", requestId,
						"quoter", quoterName,
						"amount", amountToQuery,
						"error", err)
				}
			}(q, amountStr)
		}
	}
	wg.Wait()
}

// getQuoteFromCacheWithLogging retrieves a quote from cache and logs the result
func (a *Aggregator) getQuoteFromCacheWithLogging(protocol string, currentAmountInStr string) *quoter.Quote {
	quote := a.getQuoteFromCache(protocol, currentAmountInStr)
	if quote == nil {
		a.Warn("exactIn: cache miss, skipping quoter",
			"quoter", protocol,
			"amount", currentAmountInStr)
		return nil
	}

	a.Info("exactIn: using cached quote",
		"quoter", protocol,
		"amount", currentAmountInStr)
	return quote
}

// parseQuotePrice parses the quote price and returns the total output
func (a *Aggregator) parseQuotePrice(quote *quoter.Quote, protocol string, currentAmountInStr string) (*big.Float, error) {
	currentTotalOutput, ok := new(big.Float).SetString(quote.Prices[currentAmountInStr])
	if !ok {
		a.Error("exactIn: failed to parse price",
			"quoter", protocol,
			"price", quote.Prices[currentAmountInStr])
		return nil, fmt.Errorf("failed to parse price")
	}
	return currentTotalOutput, nil
}

// calculateSegmentOutput calculates the output for the current segment
func (a *Aggregator) calculateSegmentOutput(currentTotalOutput *big.Float, state *quoter.QuoterState) *big.Float {
	if state.LastQuote != nil && state.LastQuoteAmount != nil {
		lastTotalOutput, ok := new(big.Float).SetString(state.LastQuote.Prices[state.LastQuoteAmount.String()])
		if ok {
			return new(big.Float).Sub(currentTotalOutput, lastTotalOutput)
		}
	}
	return currentTotalOutput
}

// logQuoteResult logs the quote result details
func (a *Aggregator) logQuoteResult(protocol string, state *quoter.QuoterState, currentAmountInStr string, currentTotalOutput, currentSegmentOutput *big.Float) {
	a.Info("exactIn: quoted",
		"quoter", protocol,
		"accumulatedAmountIn", state.AccumulatedAmountIn.String(),
		"currentAmountIn", currentAmountInStr,
		"totalOutput", currentTotalOutput.Text('f', 0),
		"segmentOutput", currentSegmentOutput.Text('f', 0))
}

func (a *Aggregator) processQuoter(
	q quoter.Quoter,
	state *quoter.QuoterState,
	currentAmountIn *big.Int,
	resultChan chan quoter.QuoteResult,
) {
	protocol := q.String()
	currentAmountInStr := currentAmountIn.String()

	quote := a.getQuoteFromCacheWithLogging(protocol, currentAmountInStr)
	if quote == nil {
		resultChan <- quoter.QuoteResult{Quoter: q, Quote: nil, SegmentOutput: nil, Err: fmt.Errorf("cache miss")}
		return
	}

	state.LastQuote = quote
	state.LastQuoteAmount = currentAmountIn

	currentTotalOutput, err := a.parseQuotePrice(quote, protocol, currentAmountInStr)
	if err != nil {
		resultChan <- quoter.QuoteResult{Quoter: q, Quote: nil, SegmentOutput: nil, Err: err}
		return
	}

	currentSegmentOutput := a.calculateSegmentOutput(currentTotalOutput, state)
	a.logQuoteResult(protocol, state, currentAmountInStr, currentTotalOutput, currentSegmentOutput)

	resultChan <- quoter.QuoteResult{
		Quoter:        q,
		Quote:         quote,
		SegmentOutput: currentSegmentOutput,
		Err:           nil,
	}
}

func (a *Aggregator) GetQuote(
	fromToken commons.Token,
	toToken commons.Token,
	inputAmount *big.Int,
	requestId string,
	userAddress common.Address,
	slippage int64,
) *quoter.Quote {
	if inputAmount.Cmp(big.NewInt(0)) <= 0 {
		a.Error("getQuote: amountIn must be positive", "requestId", requestId, "inputAmount", inputAmount.String())
		return nil
	}

	baseSegments := a.splitAmount(inputAmount)
	amountsList := a.generateAmountCombinations(baseSegments)
	a.preQuoteAndCache(fromToken, toToken, amountsList, requestId, userAddress)

	quoterStates := make(map[string]*quoter.QuoterState)
	for _, q := range a.quoters {
		quoterStates[q.String()] = &quoter.QuoterState{
			AccumulatedAmountIn:  big.NewInt(0),
			AccumulatedAmountOut: new(big.Float).SetInt64(0),
			LastQuote:            nil,
			LastQuoteAmount:      big.NewInt(0),
			Settlements:          []map[string]interface{}{},
		}
	}

	activeQuoters := make([]quoter.Quoter, len(a.quoters))
	copy(activeQuoters, a.quoters)

	var finalSegments []map[string]interface{}
	remainingAmount := new(big.Int).Set(inputAmount)

	currentRound := int64(0)
	a.Info("getQuote: start aggregation", "requestId", requestId, "activeQuoters", len(activeQuoters), "remainingAmount", remainingAmount.String())

	for len(activeQuoters) > 0 && remainingAmount.Cmp(big.NewInt(0)) > 0 {
		currentRound++
		isLastRound := currentRound == int64(len(baseSegments))
		resultChan := make(chan quoter.QuoteResult, len(activeQuoters))

		var currentSegmentAmount *big.Int
		if isLastRound {
			currentSegmentAmount = new(big.Int).Set(remainingAmount)
		} else {
			currentSegmentAmount = new(big.Int).Set(baseSegments[currentRound-1])
		}

		a.Info("getQuote: current round",
			"requestId", requestId,
			"round", currentRound,
			"totalRounds", len(baseSegments),
			"isLastRound", isLastRound,
			"activeQuoters", len(activeQuoters),
			"remainingAmount", remainingAmount.String(),
			"currentSegmentAmount", currentSegmentAmount.String())

		for _, activeQuoter := range activeQuoters {
			go func(q quoter.Quoter) {
				state := quoterStates[q.String()]
				currentAmountIn := new(big.Int).Add(state.AccumulatedAmountIn, currentSegmentAmount)
				if currentAmountIn.Cmp(inputAmount) > 0 {
					currentAmountIn = new(big.Int).Set(inputAmount)
				}
				a.processQuoter(q, state, currentAmountIn, resultChan)
			}(activeQuoter)
		}

		var nextActiveQuoters []quoter.Quoter
		var bestQuoter quoter.Quoter
		var bestQuoteResult *quoter.Quote
		var bestSegmentOutput *big.Float

		for i := 0; i < len(activeQuoters); i++ {
			result := <-resultChan

			if result.Err != nil || result.Quote == nil || result.SegmentOutput == nil {
				continue
			}

			nextActiveQuoters = append(nextActiveQuoters, result.Quoter)

			if bestSegmentOutput == nil || result.SegmentOutput.Cmp(bestSegmentOutput) > 0 {
				bestQuoter = result.Quoter
				bestQuoteResult = result.Quote
				bestSegmentOutput = result.SegmentOutput
			}
		}

		a.Info("getQuote: next activeQuoters", "requestId", requestId, "count", len(nextActiveQuoters), "remainingAmount", remainingAmount.String())
		if bestQuoter == nil || bestSegmentOutput == nil {
			a.Warn("getQuote: no valid quotes for current segment, ending", "requestId", requestId)
			break
		}

		if bestQuoteResult == nil {
			a.Error("getQuote: best quote result is nil", "requestId", requestId, "activeQuoter", bestQuoter.String())
			continue
		}

		bestProtocol := bestQuoter.String()
		state := quoterStates[bestProtocol]

		actualSegmentAmount := new(big.Int).Set(currentSegmentAmount)
		currentAmountIn := new(big.Int).Add(state.AccumulatedAmountIn, actualSegmentAmount)
		currentAmountInStr := currentAmountIn.String()
		currentTotalOutput, ok := new(big.Float).SetString(bestQuoteResult.Prices[currentAmountInStr])
		if !ok {
			a.Error("getQuote: failed to parse total output price",
				"requestId", requestId,
				"activeQuoter", bestProtocol,
				"price", bestQuoteResult.Prices[currentAmountInStr],
				"amountKey", currentAmountInStr)
			continue
		}
		activeQuoters = nextActiveQuoters

		state.LastQuote = bestQuoteResult
		state.LastQuoteAmount = currentAmountIn
		state.AccumulatedAmountOut = currentTotalOutput
		state.AccumulatedAmountIn = new(big.Int).Add(state.AccumulatedAmountIn, actualSegmentAmount)

		remainingAmount = new(big.Int).Sub(remainingAmount, actualSegmentAmount)
		segmentDetail := map[string]interface{}{
			"protocol":       bestProtocol,
			"amountIn":       actualSegmentAmount.String(),
			"amountOut":      bestSegmentOutput.Text('f', 0),
			"totalAmountIn":  state.AccumulatedAmountIn.String(),
			"totalAmountOut": state.AccumulatedAmountOut.Text('f', 0),
			"calldata":       hexutil.Encode(bestQuoteResult.Calldata),
			"routes":         bestQuoteResult.Routes,
			"to":             bestQuoteResult.To,
		}

		state.Settlements = []map[string]interface{}{segmentDetail}
		finalSegments = []map[string]interface{}{}
		for _, quoterState := range quoterStates {
			finalSegments = append(finalSegments, quoterState.Settlements...)
		}

		a.Info("getQuote: selected best activeQuoter for segment",
			"requestId", requestId,
			"activeQuoter", bestProtocol,
			"segmentAmountIn", actualSegmentAmount.String(),
			"segmentAmountOut", bestSegmentOutput.Text('f', 0),
			"remainingAmount", remainingAmount.String())
	}

	if len(finalSegments) == 0 {
		return nil
	}

	totalProcessedAmount := big.NewInt(0)
	for _, segment := range finalSegments {
		segmentTotalAmountIn, ok := new(big.Int).SetString(segment["totalAmountIn"].(string), 10)
		if !ok {
			a.Error("getQuote: failed to parse segment totalAmountIn", "requestId", requestId, "totalAmountIn", segment["totalAmountIn"])
			return nil
		}
		totalProcessedAmount = new(big.Int).Add(totalProcessedAmount, segmentTotalAmountIn)
	}

	if totalProcessedAmount.Cmp(inputAmount) != 0 {
		a.Error("getQuote: total processed amount does not match input amount",
			"requestId", requestId,
			"totalProcessed", totalProcessedAmount.String(),
			"inputAmount", inputAmount.String())
		return nil
	}

	combinedQuote := &quoter.Quote{
		RequestId:   requestId,
		FromToken:   fromToken.Address,
		ToToken:     toToken.Address,
		UserAddress: userAddress,
		To:          common.Address{},
		Timestamp:   time.Now().Unix(),
		Prices:      make(map[string]string),
		Routes:      make(map[string]interface{}),
		Protocol:    "aggregated",
		Path:        nil,
	}

	totalAmountOut := new(big.Float)

	var settlements []map[string]interface{}
	var routes []interface{}
	for i, segment := range finalSegments {
		if i == 0 {
			protocol := segment["protocol"].(string)
			state := quoterStates[protocol]
			if state.LastQuote != nil {
				combinedQuote.Path = state.LastQuote.Path
			}
		}

		segmentAmountOut, _ := new(big.Float).SetString(segment["totalAmountOut"].(string))
		totalAmountOut = new(big.Float).Add(totalAmountOut, segmentAmountOut)

		settlement := map[string]interface{}{
			"settler": segment["to"].(common.Address),
			"value":   "0",
			"data":    segment["calldata"],
		}
		settlements = append(settlements, settlement)

		routes = append(routes, map[string]interface{}{
			"protocol":       segment["protocol"],
			"totalAmountIn":  segment["totalAmountIn"],
			"totalAmountOut": segment["totalAmountOut"],
		})
	}

	combinedQuote.Prices[inputAmount.String()] = totalAmountOut.Text('f', 0)
	combinedQuote.AmountIn = inputAmount
	combinedQuote.AmountOut, _ = new(big.Int).SetString(totalAmountOut.Text('f', 0), 10)

	amountOutFloat := new(big.Float)
	amountOutFloat.SetString(totalAmountOut.Text('f', 0))
	minAmountOutFloat := new(big.Float)
	minAmountOutFloat.Mul(amountOutFloat, big.NewFloat(1-float64(slippage)/10000))
	minAmountOutInt := new(big.Int)
	minAmountOutInt.SetString(minAmountOutFloat.Text('f', 0), 10)
	combinedQuote.MiniAmountOut = minAmountOutInt

	combinedQuote.Routes["settlements"] = settlements
	combinedQuote.Routes["routes"] = routes

	return combinedQuote
}

func (a *Aggregator) addQuoteToCache(quoterName string, amountInStr string, quote *quoter.Quote) {
	a.cacheMutex.Lock()
	defer a.cacheMutex.Unlock()

	quoterCache, exists := a.quoteCache[quoterName]
	if !exists {
		a.quoteCache[quoterName] = make(map[string]*quoter.QuoteCacheEntry)
		quoterCache = a.quoteCache[quoterName]
	}

	quoterCache[amountInStr] = &quoter.QuoteCacheEntry{
		Quote:      quote,
		Expiration: time.Now().Add(10 * time.Second),
	}

	a.Info("cache: added to cache",
		"quoter", quoterName,
		"amount", amountInStr)
}

func (a *Aggregator) getQuoteFromCache(quoterName string, amountInStr string) *quoter.Quote {
	quoterCache, exists := a.quoteCache[quoterName]
	if !exists {
		return nil
	}

	cacheEntry, exists := quoterCache[amountInStr]
	if !exists {
		return nil
	}

	if time.Now().After(cacheEntry.Expiration) {
		a.cacheMutex.Lock()
		delete(quoterCache, amountInStr)
		a.cacheMutex.Unlock()
		return nil
	}

	a.Info("cache: cache hit",
		"quoter", quoterName,
		"amount", amountInStr)

	return cacheEntry.Quote
}
