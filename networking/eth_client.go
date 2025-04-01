package networking

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/samber/lo"
	"github.com/sony/gobreaker"
)

const (
	Health State = iota // equal StateClosed
	Sick                // equal StateHalfOpen
	Dead                // equal StateOpen
)

type EthClient interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	Client() (*ethclient.Client, error)
}

type SolidEthClient struct {
	mutex         sync.Mutex
	clients       []*StatefulEthClient
	size          int64
	defaultClient EthClient
	counter       int64
}

func NewSolidEthClient(defaultClient *StatefulEthClient, client ...*StatefulEthClient) *SolidEthClient {
	var clients []*StatefulEthClient
	clients = append(clients, client...)
	clients = lo.Shuffle(clients)

	return &SolidEthClient{
		clients:       clients,
		defaultClient: defaultClient,
		size:          int64(len(clients)),
	}
}

func (c *SolidEthClient) pick() (EthClient, error) {
	if len(c.clients) == 0 {
		return c.defaultClient, nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()

	index := c.counter % c.size
	client := c.clients[index]
	if client.State() == Health || client.State() == Sick {
		c.counter++
		return client, nil
	}

	var healthOne *StatefulEthClient
	for _, client := range c.clients {
		if client.State() == Health || client.State() == Sick {
			healthOne = client
		}
	}
	if healthOne != nil {
		c.counter++
		return healthOne, nil
	}

	return c.defaultClient, nil
}

func (c *SolidEthClient) Client() (*ethclient.Client, error) {
	client, err := c.pick()
	if err != nil {
		return nil, err
	}
	return client.Client()
}

func (c *SolidEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	client, err := c.pick()
	if err != nil {
		return nil, err
	}

	return client.SuggestGasPrice(ctx)
}

func (c *SolidEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	client, err := c.pick()
	if err != nil {
		return 0, err
	}

	return client.BlockNumber(ctx)
}

func (c *SolidEthClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	client, err := c.pick()
	if err != nil {
		return 0, err
	}

	return client.NonceAt(ctx, account, blockNumber)
}

func (c *SolidEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	client, err := c.pick()
	if err != nil {
		return nil, err
	}
	return client.BalanceAt(ctx, account, blockNumber)
}

type StatefulEthClient struct {
	client  *ethclient.Client
	breaker *gobreaker.CircuitBreaker
}

func Dial(rawurl string) (*StatefulEthClient, error) {
	client, err := ethclient.DialContext(context.Background(), rawurl)
	if err != nil {
		return nil, err
	}

	var setting gobreaker.Settings
	// url had parsed in ethclient.DialContext
	_url, _ := url.Parse(rawurl)
	setting.Name = fmt.Sprintf("breaker:%s", _url.Host)
	setting.Timeout = 2 * time.Second
	setting.Interval = 2 * time.Second
	setting.ReadyToTrip = func(counts gobreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 3

	}

	breaker := gobreaker.NewCircuitBreaker(setting)

	return &StatefulEthClient{
		client:  client,
		breaker: breaker,
	}, nil
}

func (c *StatefulEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.client.SuggestGasPrice(ctx)
	})
	if err != nil {
		return nil, err
	}
	return result.(*big.Int), nil
}

func (c *StatefulEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.client.BlockNumber(ctx)
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %s", c.breaker.Name(), err.Error())
	}
	return result.(uint64), nil
}

func (c *StatefulEthClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.client.NonceAt(ctx, account, blockNumber)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

func (c *StatefulEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.client.BalanceAt(ctx, account, blockNumber)
	})
	if err != nil {
		return nil, err
	}
	return result.(*big.Int), nil
}

func (c *StatefulEthClient) Client() (*ethclient.Client, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.client, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*ethclient.Client), nil
}

func (c *StatefulEthClient) State() State {
	state := c.breaker.State()
	return State(state)
}

type State int

func (s State) String() string {
	return [...]string{"Health", "Sick", "Dead"}[s]
}
