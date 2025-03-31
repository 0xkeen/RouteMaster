package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"github.com/urfave/cli"

	aggregator2 "route_master/aggregator"
	"route_master/aggregator/quoter"
	"route_master/aggregator/quoter/amm"
	"route_master/api"
	"route_master/db"
	"route_master/metrics"
	"route_master/networking"
)

var logger = log15.New("module", "main")

func main() {
	gin.SetMode(gin.ReleaseMode)

	app := cli.NewApp()
	app.Name = "route-master"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "listen", Value: ":80", Usage: "listen address and port", EnvVar: "LISTEN"},
		cli.Int64Flag{Name: "chain_id", Value: 1, Usage: "ethereum chain id", EnvVar: "CHAIN_ID"},
		cli.StringFlag{Name: "debug_listen", Value: ":8080", Usage: "debug server listen", EnvVar: "DEBUG_LISTEN"},
		cli.StringFlag{Name: "database_dsn", Value: "root@tcp(localhost:3306)/test?parseTime=true", Usage: "Data Source Name", EnvVar: "DATABASE_DSN"},
		cli.StringFlag{Name: "private_node_url", Value: "", Usage: "web3 rpc url", EnvVar: "PRIVATE_NODE_URL"},
		cli.StringFlag{Name: "redis_host", Value: "redis://localhost:6379/0", Usage: "balance redis", EnvVar: "REDIS_HOST"},
		cli.Int64Flag{Name: "segment_percentage", Value: 50, Usage: "split amount", EnvVar: "SEGMENT_PERCENTAGE"},
		cli.IntFlag{Name: "private_node_enlarge_factor", Usage: "enlarge private node call times", Value: 1, EnvVar: "PRIVATE_NODE_ENLARGE_FACTOR"},
	}

	if err := app.Run(os.Args); err != nil {
		panic(fmt.Sprintf("Failed to start routerMaster, err=%+v", err))
	}
}

func run(ctx *cli.Context) {
	chainId := big.NewInt(ctx.Int64("chain_id"))
	backupClients := getBackupNodeClients(chainId.Int64())
	logger.Info("init backup node rpc client", "clients", len(backupClients))

	//  1:1 enlarge the request volume of own node
	privateNodeClient, err := networking.Dial(ctx.String("private_node_url"))
	if err != nil {
		panic(err)
	}
	enlargeFactor := ctx.Int("private_node_enlarge_factor")
	if enlargeFactor > 0 {
		n := len(backupClients) * enlargeFactor
		for i := 0; i < n; i++ {
			backupClients = append(backupClients, privateNodeClient)
		}
	}
	solidClient := networking.NewSolidEthClient(privateNodeClient, backupClients...)

	var quoters []quoter.Quoter
	fees := []int64{100, 500, 1000, 3000, 10000}
	for _, fee := range fees {
		quoters = append(quoters, amm.NewUniswapV3(solidClient, big.NewInt(fee)))
		quoters = append(quoters, amm.NewSushiSwapV3(solidClient, big.NewInt(fee)))
	}
	aggregator := aggregator2.NewAggregator(quoters, ctx.Int64("segment_percentage"))

	quoteDB := db.NewQuote(ctx.String("database_dsn"))
	routeAPI := api.NewRouteAPI(chainId, aggregator, solidClient, quoteDB)

	apiRouter := gin.New()
	apiGroup := apiRouter.Group("/v1")
	registerRouter(chainId.Int64(), apiGroup, routeAPI)

	go func() {
		if err := metrics.NewMetricServer(ctx.String("debug_listen")); err != nil {
			logger.Error("metrics server starting failed", "err", err)
			panic(err)
		}
	}()

	logger.Info("starting server", "listen", ctx.String("listen"))
	_ = apiRouter.Run(ctx.String("listen"))
}

func registerRouter(chainId int64, router *gin.RouterGroup, routeAPI *api.RouteAPI) *gin.RouterGroup {
	router.GET(fmt.Sprintf("/quote/%d", chainId), routeAPI.Quote)
	router.GET(fmt.Sprintf("/quotes/history/%d", chainId), routeAPI.QuotesHistory)

	return router
}

func getBackupNodeClients(chainId int64) []*networking.StatefulEthClient {
	var clients []*networking.StatefulEthClient
	if chainId == 1 {
		freeRpcs := []string{
			"https://eth.merkle.io",
			"https://ethereum.publicnode.com",
			"https://eth.llamarpc.com",
			"https://ethereum.blockpi.network/v1/rpc/public",
			"https://rpc.ankr.com/eth",
		}
		for _, rpcUrl := range freeRpcs {
			if freeNodeClient, err := networking.Dial(rpcUrl); err == nil {
				clients = append(clients, freeNodeClient)
			}
		}
	}
	return clients
}
