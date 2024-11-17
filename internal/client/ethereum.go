package client

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"truebit-api/internal/config"
)

const (
	purchaseABI = `[{"inputs":[{"internalType":"uint256","name":"numTRU","type":"uint256"}],"name":"getPurchasePrice","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"numTRU","type":"uint256"}],"name":"getRetirePrice","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"reserve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"opex","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"OPEX_COST","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"THETA","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	tokenABI    = `[{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	uniswapABI  = `[{"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"stateMutability":"view","type":"function"}]`
)

type EthereumClient struct {
	eth      *ethclient.Client
	purchase *bind.BoundContract
	token    *bind.BoundContract
	uniswap  *bind.BoundContract
}

func NewEthereumClient(cfg *config.Config) (*EthereumClient, error) {
	infuraURL := fmt.Sprintf("https://mainnet.infura.io/v3/%s", cfg.InfuraProjectID)
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %v", err)
	}

	purchaseABIParsed, err := abi.JSON(strings.NewReader(purchaseABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse purchase ABI: %v", err)
	}

	tokenABIParsed, err := abi.JSON(strings.NewReader(tokenABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token ABI: %v", err)
	}

	uniswapABIParsed, err := abi.JSON(strings.NewReader(uniswapABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse uniswap ABI: %v", err)
	}

	return &EthereumClient{
		eth: client,
		purchase: bind.NewBoundContract(
			common.HexToAddress(cfg.ContractAddresses.PurchaseAddress),
			purchaseABIParsed,
			client, client, client,
		),
		token: bind.NewBoundContract(
			common.HexToAddress(cfg.ContractAddresses.TokenAddress),
			tokenABIParsed,
			client, client, client,
		),
		uniswap: bind.NewBoundContract(
			common.HexToAddress(cfg.ContractAddresses.UniswapAddress),
			uniswapABIParsed,
			client, client, client,
		),
	}, nil
}

func (c *EthereumClient) GetContracts() (*bind.BoundContract, *bind.BoundContract, *bind.BoundContract) {
	return c.purchase, c.token, c.uniswap
}
