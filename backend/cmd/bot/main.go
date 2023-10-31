package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os/signal"

	"os"
	"wallemon/pkg/database"
	"wallemon/pkg/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robfig/cron/v3"
)

var (
	env string
)

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
	err := utils.LoadEnvConfig("config/.env")
	if err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	rpc := os.Getenv("RPC")
	client, err = ethclient.Dial(rpc)
	if err != nil {
		panic(fmt.Errorf("failed to connect to the Ethereum client: %v", err))
	}

	chainID, err = client.NetworkID(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to get network ID: %v", err))
	}

	parsedABI, err = utils.GetContractABI("./config/abi.json")
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

	wallemonAddr = common.HexToAddress(os.Getenv("WALLEMON_ADDRESS"))

	env = os.Getenv("ENV")
	if env != "local" {
		env = "cloud"
	}
	fmt.Println("wallemon-bot ENV: ", env)
}

func main() {
	// Create root context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database.Initialize(ctx)
	defer database.Finalize()

	c := cron.New()
	c.AddFunc("@every 3s", func() {
		fmt.Println("Gogo")
	})

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
