package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os/signal"

	"os"
	"wallemon/pkg/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Ethereum struct {
		RPCURL string `yaml:"rpcURL"`
	} `yaml:"ethereum"`
}

type ToBeSickListRet = []*big.Int

var (
	rpcClient    *rpc.Client
	client       *ethclient.Client
	chainID      *big.Int
	parsedABI    *abi.ABI
	wallemonAddr common.Address
)

func init() {
	// err := utils.LoadSecrets("config/.secrets")
	// if err != nil {
	// 	panic(fmt.Errorf("failed to load secrets: %v", err))
	// }

	var err error
	rpcClient, err = rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(fmt.Errorf("failed to connect to the rpc Ethereum client: %v", err))
	}

	client, err = ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(fmt.Errorf("failed to connect to the Ethereum client: %v", err))
	}

	chainID, err = client.NetworkID(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to get network ID: %v", err))
	}

	// parsedABI, err = utils.GetContractABI("../config/abi.json")
	parsedABI, err = utils.GetContractABI("../../config/abi.json")
	if err != nil {
		panic(fmt.Errorf("failed to parse contract ABI: %v", err))
	}

	wallemonAddr = common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")
}

func main() {
	// cfg, err := loadConfig("config/env.yaml")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	c := cron.New()
	// c.AddFunc("@every 6h", func() {
	// 	sendEthTransaction()
	// })
	c.AddFunc("@every 3s", func() {
		// fmt.Println(cfg.Database.Host)
		fmt.Println("get to be sick list")
		// ret, err := getToBeSickList(rpcClient)
		ret, err := getToBeSickList(client)
		if err != nil {
			log.Fatal(err)

		}
		fmt.Println(ret)
	})
	c.Start()
	defer c.Stop()

	// set a channel to receive OS signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}

// func getToBeSickList(client *rpc.Client) (ToBeSickListRet, error) {
func getToBeSickList(client *ethclient.Client) (ToBeSickListRet, error) {
	data, err := parsedABI.Pack("toBeSickList")
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI call: %v", err)
	}

	callArgs := ethereum.CallMsg{
		From: common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
		To:   &wallemonAddr,
		Data: data,
	}

	res, err := client.CallContract(context.Background(), callArgs, nil)
	fmt.Println(res)
	return nil, nil

	// result := []byte{}
	// err = rpcClient.Call(&result, "eth_call", callArgs, "latest")
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to call contract: %v", err)
	// }

	// var ret ToBeSickListRet
	// err = parsedABI.UnpackIntoInterface(&ret, "toBeSickList", result)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to unpack result: %v", err)
	// }

	// return ret, nil
}

func loadConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func sendEthTransaction() {
	client, err := ethclient.Dial("YOUR_ETHEREUM_NODE_URL")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("RECIPIENT_ETHEREUM_ADDRESS")
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

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

	log.Printf("Sent transaction: %s", signedTx.Hash().Hex())
}
