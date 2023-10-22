package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
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
	client       *ethclient.Client
	chainID      *big.Int
	parsedABI    *abi.ABI
	ownerKey     *ecdsa.PrivateKey
	ownerAddr    common.Address
	wallemonAddr common.Address
)

func init() {
	var err error
	err = utils.LoadSecrets("config/.secrets")
	if err != nil {
		panic(fmt.Errorf("failed to load secrets: %v", err))
	}

	// client, err = ethclient.Dial("http://127.0.0.1:8545")
	client, err = ethclient.Dial("http://host.docker.internal:8545")
	if err != nil {
		panic(fmt.Errorf("failed to connect to the Ethereum client: %v", err))
	}

	chainID, err = client.NetworkID(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to get network ID: %v", err))
	}

	parsedABI, err = utils.GetContractABI("./config/abi.json")
	// parsedABI, err = utils.GetContractABI("../../config/abi.json")
	if err != nil {
		panic(fmt.Errorf("failed to parse contract ABI: %v", err))
	}

	k := os.Getenv("SIGNER_KEY")
	if k[:2] == "0x" {
		k = k[2:]
	}
	ownerKey, err = crypto.HexToECDSA(k)
	if err != nil {
		panic(fmt.Errorf("failed to get signer: %v", err))
	}
	ownerAddr = crypto.PubkeyToAddress(ownerKey.PublicKey)

	wallemonAddr = common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")
}

func main() {
	// cfg, err := loadConfig("config/env.yaml")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	c := cron.New()
	// c.AddFunc("@every 6h", func() {
	// fmt.Println(cfg.Database.Host)
	// 	sendEthTransaction()
	// })

	c.AddFunc("@every 3s", func() {
		sickBot()
	})

	c.AddFunc("@every 3s", func() {
		killBot()
	})

	c.Start()
	defer c.Stop()

	// set a channel to receive OS signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}

func sickBot() {
	ret, err := toBeSickList(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret)
	if len(ret) == 0 {
		fmt.Println("No one is sick.")
		return
	}

	err = batchSick(ret)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func killBot() {
	ret, err := toBeDeadkList(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret)
	if len(ret) == 0 {
		fmt.Println("No one is dead.")
		return
	}

	err = batachKill(ret)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func toBeSickList(client *ethclient.Client) (ToBeSickListRet, error) {
	data, err := parsedABI.Pack("toBeSickList")
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI call: %v", err)
	}

	callArgs := ethereum.CallMsg{
		From: ownerAddr,
		To:   &wallemonAddr,
		Data: data,
	}

	res, err := client.CallContract(context.Background(), callArgs, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}
	if len(res) == 0 {
		return nil, nil
	}

	var ret ToBeSickListRet
	err = parsedABI.UnpackIntoInterface(&ret, "toBeSickList", res)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack ABI call: %v", err)
	}

	return ret, nil
}

func batchSick(list ToBeSickListRet) error {
	// Prepare the method input parameters.
	params, err := parsedABI.Pack("batchSick", list)
	if err != nil {
		return fmt.Errorf("failed to pack ABI call: %v", err)
	}

	// New transaction
	tx, err := utils.NewTransaction(client, ownerAddr, wallemonAddr, nil, params)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign the transaction.
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), ownerKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Broadcast the transaction.
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Sent Transaction: %s\n", signedTx.Hash().Hex())

	return nil
}

func toBeDeadkList(client *ethclient.Client) (ToBeSickListRet, error) {
	data, err := parsedABI.Pack("toBeDeadList")
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI call: %v", err)
	}

	callArgs := ethereum.CallMsg{
		From: ownerAddr,
		To:   &wallemonAddr,
		Data: data,
	}

	res, err := client.CallContract(context.Background(), callArgs, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}
	if len(res) == 0 {
		return nil, nil
	}

	var ret ToBeSickListRet
	err = parsedABI.UnpackIntoInterface(&ret, "toBeDeadList", res)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack ABI call: %v", err)
	}

	return ret, nil
}

func batachKill(list ToBeSickListRet) error {
	// Prepare the method input parameters.
	params, err := parsedABI.Pack("batachKill", list)
	if err != nil {
		return fmt.Errorf("failed to pack ABI call: %v", err)
	}

	// New transaction
	tx, err := utils.NewTransaction(client, ownerAddr, wallemonAddr, nil, params)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign the transaction.
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), ownerKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Broadcast the transaction.
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Sent Transaction: %s\n", signedTx.Hash().Hex())

	return nil
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
