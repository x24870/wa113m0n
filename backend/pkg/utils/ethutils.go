package utils

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumClient struct {
	Client      *ethclient.Client
	PrivateKey  *ecdsa.PrivateKey
	FromAddress common.Address
	ChainID     *big.Int
}

func NewEthereumClient(rpcURL string, privateKeyHex string) *EthereumClient {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return &EthereumClient{
		Client:      client,
		PrivateKey:  privateKey,
		FromAddress: fromAddress,
		ChainID:     chainID,
	}
}

func (e *EthereumClient) EstimateGas(toAddress common.Address, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From: e.FromAddress,
		To:   &toAddress,
		Data: data,
	}
	gasLimit, err := e.Client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}

func (e *EthereumClient) SendTransaction(toAddress common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	nonce, err := e.Client.PendingNonceAt(context.Background(), e.FromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// Estimating gas limit for the transaction
	gasLimit, err := e.EstimateGas(toAddress, data)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(e.ChainID), e.PrivateKey)
	if err != nil {
		return nil, err
	}

	err = e.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func KeyStringToPrivateKey(key string) (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func EncodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}
