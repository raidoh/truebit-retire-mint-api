package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"truebit-api/internal/client"
	"truebit-api/internal/models"
	"truebit-api/pkg/utils"
)

type TruebitService struct {
	ethClient *client.EthereumClient
}

func NewTruebitService(client *client.EthereumClient) *TruebitService {
	return &TruebitService{
		ethClient: client,
	}
}

func (s *TruebitService) getETHPrice() (float64, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch ETH price: %v", err)
	}
	defer resp.Body.Close()

	var price models.EthPrice
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return 0, fmt.Errorf("failed to decode price data: %v", err)
	}

	return price.Ethereum.USD, nil
}

func (s *TruebitService) GetInfo(ctx context.Context) (*models.TruebitInfo, error) {
	purchase, token, uniswap := s.ethClient.GetContracts()

	// Get ETH price
	ethPrice, err := s.getETHPrice()
	if err != nil {
		return nil, err
	}

	// Get mint price
	var out []interface{}
	err = purchase.Call(&bind.CallOpts{Context: ctx}, &out, "getPurchasePrice", big.NewInt(1e18))
	if err != nil {
		return nil, fmt.Errorf("failed to get mint price: %v", err)
	}
	mintPrice := out[0].(*big.Int)
	mintPriceETH := utils.WeiToEth(mintPrice)

	// Get retire price
	out = make([]interface{}, 0)
	err = purchase.Call(&bind.CallOpts{Context: ctx}, &out, "getRetirePrice", big.NewInt(1e18))
	if err != nil {
		return nil, fmt.Errorf("failed to get retire price: %v", err)
	}
	retirePrice := out[0].(*big.Int)
	retirePriceETH := utils.WeiToEth(retirePrice)

	// Get total supply
	out = make([]interface{}, 0)
	err = token.Call(&bind.CallOpts{Context: ctx}, &out, "totalSupply")
	if err != nil {
		return nil, fmt.Errorf("failed to get total supply: %v", err)
	}
	totalSupply := out[0].(*big.Int)
	formattedSupply := utils.FormatTokenSupply(totalSupply)

	// Get reserve
	out = make([]interface{}, 0)
	err = purchase.Call(&bind.CallOpts{Context: ctx}, &out, "reserve")
	if err != nil {
		return nil, fmt.Errorf("failed to get reserve: %v", err)
	}
	reserve := out[0].(*big.Int)
	reserveETH := utils.FormatEthValue(reserve)
	reserveFloat, _ := new(big.Float).SetString(reserveETH)
	reserveF64, _ := reserveFloat.Float64()

	// Get Uniswap reserves
	out = make([]interface{}, 0)
	err = uniswap.Call(&bind.CallOpts{Context: ctx}, &out, "getReserves")
	if err != nil {
		return nil, fmt.Errorf("failed to get reserves: %v", err)
	}
	poolETH := out[0].(*big.Int)
	poolTRU := out[1].(*big.Int)

	poolETHStr := utils.FormatEthValue(poolETH)
	poolTRUStr := utils.FormatEthValue(poolTRU)
	poolETHFloat, _ := new(big.Float).SetString(poolETHStr)
	poolTRUFloat, _ := new(big.Float).SetString(poolTRUStr)
	poolETHF64, _ := poolETHFloat.Float64()
	poolTRUF64, _ := poolTRUFloat.Float64()

	// Format display output
	display := strings.Builder{}
	display.WriteString(fmt.Sprintf("1 ETH = $%.2f\n\n", ethPrice))
	display.WriteString(fmt.Sprintf("Truebit OS mint price: %.6f ETH ($%.2f)\n", mintPriceETH, mintPriceETH*ethPrice))
	display.WriteString(fmt.Sprintf("Truebit OS retire price: %.6f ETH ($%.2f)\n\n", retirePriceETH, retirePriceETH*ethPrice))
	display.WriteString(fmt.Sprintf("Total supply: %s TRU\n", formattedSupply))
	display.WriteString(fmt.Sprintf("Reserve: %s ETH ($%.2f)\n\n", reserveETH, reserveF64*ethPrice))
	display.WriteString(fmt.Sprintf("Uniswap v2 ETH: %s ($%.2f)\n", poolETHStr, poolETHF64*ethPrice))
	display.WriteString(fmt.Sprintf("Uniswap v2 TRU: %s\n", poolTRUStr))

	if poolTRUF64 > 0 {
		truPriceETH := poolETHF64 / poolTRUF64
		display.WriteString(fmt.Sprintf("Uniswap v2: 1 TRU = %.6f ETH ($%.2f)\n", truPriceETH, truPriceETH*ethPrice))
	}

	// Prepare response
	info := &models.TruebitInfo{
		Display: display.String(),
		Data: models.TruebitData{
			EthPrice:    ethPrice,
			MintPrice:   mintPriceETH,
			RetirePrice: retirePriceETH,
			TotalSupply: formattedSupply,
			Reserve: models.Reserve{
				ETH: reserveETH,
				USD: fmt.Sprintf("$%.2f", reserveF64*ethPrice),
			},
			Uniswap: models.Uniswap{
				ETH:    poolETHStr,
				TRU:    poolTRUStr,
				ETHUSD: fmt.Sprintf("$%.2f", poolETHF64*ethPrice),
			},
		},
	}

	return info, nil
}
