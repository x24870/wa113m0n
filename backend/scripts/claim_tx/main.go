package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"wallemon/pkg/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	parsedABI  *abi.ABI
	client     *ethclient.Client
	chainID    *big.Int
	privateKey *ecdsa.PrivateKey
	// publicKey  *ecdsa.PublicKey
	address common.Address
)

func init() {
	var err error
	parsedABI, err = utils.GetContractABI("../config/abi.json")
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	client, err = ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	chainID, err = client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err = crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA)
}

func main() {
	from := address
	nftAddr := common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")
	// to := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	// safeMint(client, privateKey, from, to, nftAddr)
	userMint(client, privateKey, from, nftAddr)
}

func userMint(client *ethclient.Client, privateKey *ecdsa.PrivateKey, from, contractAddr common.Address) {
	refCode := "wallemon"
	signature, err := hex.DecodeString("f4b6424aebb6e151136076cfe601f8c20cb91603e4bfef367c384f7f5de6fd287c707cdce8aa7757770f756e0b6a4351998861bfb9417de77ab365181625e7af1b")
	if err != nil {
		log.Fatalf("Failed to decode signature: %v", err)
	}

	// Prepare the method input parameters.
	params, err := parsedABI.Pack("userMint", refCode, signature)
	if err != nil {
		log.Fatalf("Failed to pack ABI call: %v", err)
	}

	// New transaction
	tx, err := utils.NewTransaction(client, from, contractAddr, nil, params)
	if err != nil {
		log.Fatalf("Failed to create transaction: %v", err)
	}

	// Sign the transaction.
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Broadcast the transaction.
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Sent Transaction: %s\n", signedTx.Hash().Hex())
}

func safeMint(client *ethclient.Client, privateKey *ecdsa.PrivateKey, from, to, contractAddr common.Address) {
	parsedABI, err := utils.GetContractABI("../config/abi.json")

	// Prepare the method input parameters.
	params, err := parsedABI.Pack("safeMint", to)
	if err != nil {
		log.Fatalf("Failed to pack ABI call: %v", err)
	}

	// Estimate gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	// Get chainID
	cid, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Get the nonce
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		log.Fatal(err)
	}

	// Estimate gas limit.
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  from,
		To:    &contractAddr,
		Gas:   0,
		Value: nil,
		Data:  params,
	})
	if err != nil {
		log.Fatalf("Failed to estimate gas limit: %v", err)
	}

	// Create the transaction.
	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), gasLimit, gasPrice, params)

	// Sign the transaction.
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(cid), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Broadcast the transaction.
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Sent Transaction: %s\n", signedTx.Hash().Hex())
}
