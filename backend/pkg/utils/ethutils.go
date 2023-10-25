package utils

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
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

func NewTransaction(client *ethclient.Client, from, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	// Estimate gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	// Get the nonce
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Estimate gas limit.
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  from,
		To:    &to,
		Gas:   0,
		Value: value,
		Data:  data,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas limit: %v", err)
	}

	return types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data), nil
}

func VerifySignature(address common.Address, message string, signature []byte) (bool, error) {
	// hash the message
	hash := crypto.Keccak256Hash([]byte(message))

	// recover the public key
	pubkeyECDSA, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %v", err)
	}

	// []byte to ecdsa.PublicKey
	pubkey, err := crypto.UnmarshalPubkey(pubkeyECDSA)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal public key: %v", err)
	}

	// obtain the address from the public key
	addr := crypto.PubkeyToAddress(*pubkey).Hex()

	return addr == address.Hex(), nil
}

func SignatureToBytes(signature string) ([]byte, error) {
	// decompose the signature
	signatureBytes := common.FromHex(signature)
	if len(signatureBytes) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}

	// Ethereum uses `V` values of 27 or 28.
	// `go-ethereum` uses 0 or 1. Let's convert.
	signatureBytes[64] -= 27

	return signatureBytes, nil
}
