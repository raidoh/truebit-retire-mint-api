package models

type TruebitInfo struct {
	Display string      `json:"display"`
	Data    TruebitData `json:"data"`
}

type TruebitData struct {
	EthPrice    float64 `json:"ethPrice"`
	MintPrice   float64 `json:"mintPrice"`
	RetirePrice float64 `json:"retirePrice"`
	TotalSupply string  `json:"totalSupply"`
	Reserve     Reserve `json:"reserve"`
	Uniswap     Uniswap `json:"uniswap"`
}

type Reserve struct {
	ETH string `json:"eth"`
	USD string `json:"usd"`
}

type Uniswap struct {
	ETH    string `json:"eth"`
	TRU    string `json:"tru"`
	ETHUSD string `json:"ethUSD"`
}

type EthPrice struct {
	Ethereum struct {
		USD float64 `json:"usd"`
	} `json:"ethereum"`
}
