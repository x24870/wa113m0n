package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"wallemon/abi"
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

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8") // Replace '0xTO_ADDRESS' with actual destination address
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sent Transaction: %s\n", signedTx.Hash().Hex())

	nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	nftAddr := common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F")
	toAddr := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	mintNFT(client, privateKey, nonce, nftAddr, toAddr)
}

func mintNFT(client *ethclient.Client, privateKey *ecdsa.PrivateKey, nonce uint64, contractAddr, toMintAddress common.Address) {

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	cid, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, cid)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	// Get new contract instance
	contract, err := abi.NewWalleMon(contractAddr, client)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	fmt.Println("contract", contract)
	// call function of contract instance and get count value
	tx, err := contract.SafeMint(auth, toMintAddress)
	if err != nil {
		log.Fatalf("Failed to call SafeMint function: %v", err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex())
	// Read the ABI
	// parsedABI, err := abi.JSON(strings.NewReader(""))
	// if err != nil {
	// 	log.Fatalf("Failed to parse contract ABI: %v", err)
	// }

	// // Pack ABI with the method and arguments
	// data, err := parsedABI.Pack("safeMint", toMintAddress)
	// if err != nil {
	// 	log.Fatalf("Failed to pack ABI with method and arguments: %v", err)
	// }

	// // Create a transaction
	// tx := types.NewTransaction(
	// 	nonce, // nonce for the sender's account
	// 	contractAddress,
	// 	big.NewInt(0), // amount to transfer
	// 	gasLimit,
	// 	gasPrice,
	// 	data, // contract method input parameters
	// )

	// // Sign the transaction
	// chainID := big.NewInt(1) // Mainnet ID. For Rinkeby Testnet use 4, for Ropsten Testnet use 3, etc.
	// signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	// if err != nil {
	// 	log.Fatalf("Failed to sign transaction: %v", err)
	// }

	// // Broadcast the transaction
	// err = client.SendTransaction(context.Background(), signedTx)
	// if err != nil {
	// 	log.Fatalf("Failed to send transaction: %v", err)
	// }

	// fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
}
