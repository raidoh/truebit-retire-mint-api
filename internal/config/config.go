package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	InfuraProjectID string
	ContractAddresses
}

type ContractAddresses struct {
	PurchaseAddress string
	TokenAddress    string
	UniswapAddress  string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	infuraID := os.Getenv("INFURA_PROJECT_ID")
	if infuraID == "" {
		return nil, fmt.Errorf("INFURA_PROJECT_ID not set in environment")
	}

	return &Config{
		InfuraProjectID: infuraID,
		ContractAddresses: ContractAddresses{
			PurchaseAddress: "0x764C64b2A09b09Acb100B80d8c505Aa6a0302EF2",
			TokenAddress:    "0xf65B5C5104c4faFD4b709d9D60a185eAE063276c",
			UniswapAddress:  "0x80b4d4e9d88D9f78198c56c5A27F3BACB9A685C5",
		},
	}, nil
}
