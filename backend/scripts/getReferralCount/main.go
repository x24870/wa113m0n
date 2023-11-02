package main

import (
	"context"
	"log"
	"math/big"
	"wallemon/pkg/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL       = "https://arbitrum-mainnet.infura.io/v3/55b0bda240044d149de944863519ae3e"
	contractAddress = "0xCd3933B359A73DEBbDcfd0E6C3a6BB36f114E187"
)

var referralCodes = []string{
	"wH7w3jpl", "QAh3KVJi", "vDMj6kEN", "csxGU2Uv", "ZH9HicFK", "duFZUYZM", "bR8kfYcp",
	"Qvj007Fa", "NelvR8Qh", "Ttv9xsEA", "fcD0Pygl", "vhqx7m5S", "19zywS9i", "FNTawp9U",
	"Pr18AjLi", "75xkFJfU", "ft3rVwg7", "qRHImAfH", "FPbeLToM", "uOK1Psm7", "85APKVFW",
	"4gtFeTkE", "0rRxOklt", "JzA1OUKx",
}

func main() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	parsedABI, err := utils.GetContractABI("../../config/abiReferral.json")
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	contract := common.HexToAddress(contractAddress)

	for _, code := range referralCodes {
		data, err := parsedABI.Pack("getReferralCount", code)
		if err != nil {
			log.Fatalf("Failed to pack data: %v", err)
		}
		callMsg := ethereum.CallMsg{
			To:   &contract,
			Data: data,
		}

		result, err := client.CallContract(context.Background(), callMsg, nil)
		if err != nil {
			log.Printf("Failed to call contract for referral code %s: %v", code, err)
			continue
		}

		count := new(big.Int).SetBytes(result)
		log.Printf("Referral code %s has count: %d", code, count.Uint64())
	}
}
