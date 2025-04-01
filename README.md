# RouteMaster

RouteMaster is a high-performance cross-chain route aggregator that finds optimal trading paths across multiple DEXs and cross-chain protocols. It features intelligent segmentation and caching mechanisms for efficient handling of large transactions.

## Features

- **Multi-Protocol Support**: Supports multiple DEXs and cross-chain protocols
- **Smart Segmentation**: Automatically splits large transactions into smaller segments for optimal pricing
- **Caching Mechanism**: Implements efficient quote caching to reduce redundant requests
- **Rate Limiting**: Built-in request rate limiting to prevent excessive protocol load
- **Data Persistence**: Supports saving quote history to MySQL database
- **Concurrent Processing**: Uses goroutines for concurrent quote queries
- **Error Handling**: Comprehensive error handling and logging mechanisms

## Installation

1. Clone the repository:
```bash
git clone https://github.com/0xkeen/RouteMaster.git
cd RouteMaster
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables, including:

- `SEGMENT_PERCENT`: Transaction segmentation percentage
- `CACHE_EXPIRATION`: Cache expiration time
- `DATABASE_DSN`: Database connection string
- `PRIVATE_NODE_URL`: Private node URL
- `CHAIN_ID`: Chain ID
- `PRIVATE_NODE_ENLARGE_FACTOR`: Private node enlargement factor

## Usage

1. Start the service:
```bash
# Start services using docker-compose
docker-compose up -d

# View service logs
docker-compose logs -f app

# Stop services
docker-compose down
```

2. Get quote example:
```bash
# Get quote with default 1% slippage
curl "http://localhost:8081/v1/quote/1?fromToken=0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2&toToken=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48&amount=100000000000000000000000&userAddress=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"

# Get quote with custom 0.5% slippage (50 bps)
curl "http://localhost:8081/v1/quote/1?fromToken=0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2&toToken=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48&amount=100000000000000000000000&userAddress=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48&slippage=50"
```

Parameters:
- `fromToken`: Source token address
- `toToken`: Target token address
- `amount`: Amount to swap (in source token's decimals)
- `userAddress`: User's wallet address
- `slippage`: Slippage tolerance in basis points (bps), 1 bps = 0.01%. Default is 100 bps (1%). Range: 0-10000 bps.

Response includes:
- `requestId`: Unique request identifier
- `fromToken`: Source token address
- `toToken`: Target token address
- `amountIn`: Input amount
- `amountOut`: Expected output amount
- `minAmountOut`: Minimum output amount based on slippage
- `userAddress`: User's wallet address
- `routes`: Detailed routing information
- `timestamp`: Quote timestamp

3. Get quote history example:
```bash
curl "http://localhost:8081/v1/quotes/history/1?fromToken=0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2&toToken=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48&amount=100000000000000000000000&limit=10&offset=0"
```

Parameters:
- `fromToken`: Source token address
- `toToken`: Target token address
- `amount`: Amount to swap (in source token's decimals)
- `limit`: Maximum number of records to return (default: 10)
- `offset`: Number of records to skip (default: 0)

Response includes:
- `id`: Record ID
- `fromToken`: Source token address
- `toToken`: Target token address
- `amountIn`: Input amount
- `userAddress`: User's wallet address
- `routes`: Detailed routing information
- `createdAt`: Record creation timestamp
- `updatedAt`: Record last update timestamp

## Testing

Run tests:
```bash
go test ./...
```