package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Replace 'YOUR_INFURA_URL' with your actual Infura URL
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	// Replace 'YOUR_PRIVATE_KEY' with your actual private key
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}

	from := crypto.PubkeyToAddress(*publicKeyECDSA)
	nftAddr := common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F")
	to := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	mintNFT(client, privateKey, from, to, nftAddr)
}

func mintNFT(client *ethclient.Client, privateKey *ecdsa.PrivateKey, from, to, contractAddr common.Address) {
	// Initialize contract's ABI.
	contractABI := getContractABI("w.json")
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

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

func getContractABI(path string) string {
	// Read the file content
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read ABI file: %v", err)
	}
	return string(data)
}
