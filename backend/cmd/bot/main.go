package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os/signal"
	"time"

	"os"
	"wallemon/pkg/database"
	"wallemon/pkg/models"
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
type HeathListRet = []uint8

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

	wallemonAddr = common.HexToAddress(os.Getenv("WALLEMON_ADDR"))
	fmt.Println("wallemon-bot WALLEMON_ADDR: ", wallemonAddr.Hex())

	env = os.Getenv("ENV")
	if env != "local" {
		env = "cloud"
	}
	fmt.Println("wallemon-bot ENV: ", env)
}

func main() {
	// sleep 5s to wait for service and db to be ready
	time.Sleep(5 * time.Second)

	// Create root context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database.Initialize(ctx)
	defer database.Finalize()

	c := cron.New()
	c.AddFunc("@every 20s", func() {
		sickBot()
	})

	c.AddFunc("@every 30s", func() {
		killBot()
	})

	c.AddFunc(fmt.Sprintf("@every %ds", models.PoopDuration), func() {
		// c.AddFunc("@every 5s", func() {
		poopBot()
	})

	c.AddFunc("@every 1m", func() {
		healthBot()
	})

	c.Start()
	defer c.Stop()

	// set a channel to receive OS signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}

// if on-chain state is health, update tokenID to be healthy in DB
func healthBot() {
	ret, err := healthList(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := database.GetSQL()
	for tokenID, state := range ret {
		t := models.NewToken(uint(tokenID))
		t, err := models.Token.GetByTokenIDAndLock(db)
		if err != nil {
			fmt.Println(fmt.Errorf("healthBot: failed to get tokenID[%d] from DB: %v", tokenID, err))
			continue
		}

		// if on-chain state is health, update tokenID to be healthy in DB
		if state == 0 && t.GetState() != 0 {
			if err := t.Update(db, map[string]interface{}{
				"state": 0,
			}); err != nil {
				fmt.Println(fmt.Errorf("healthBot: failed to update tokenID[%d] to health: %v", tokenID, err))
			}
		}
	}
}

func poopBot() {
	db := database.GetSQL()
	poops, err := models.Poop.List(db)
	if err != nil {
		fmt.Println(fmt.Errorf("poopBot: failed to list poops from DB: %v", err))
		return
	}
	fmt.Println("poopBot: poops: ", poops)

	for _, p := range poops {
		// check if tokenID is already dead
		t, err := models.Token.GetByTokenID(db)
		if err != nil {
			fmt.Println(fmt.Errorf("poopBot: failed to get tokenID[%d] from DB: %v", p.GetTokenID(), err))
			continue
		}
		if t.GetState() == 2 {
			continue
		}

		a := p.GetAmount()
		if a >= models.PoopMaxAmount {
			continue
		}

		// increase amount by 1
		if err := p.Update(db, map[string]interface{}{
			"amount": a + 1,
		}); err != nil {
			fmt.Println(fmt.Errorf("poopBot: failed to update tokenID[%d] poop: %v", p.GetTokenID(), err))
		}

	}
}

func sickBot() {
	// get sick list if exeeded last meal time duration
	ret, err := toBeSickList(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("toBeSickList onchain: ", ret)

	// get sick list if poop amount >= 6
	db := database.GetSQL()
	poopSickList, err := models.Poop.ListShouldSick(db)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get poop sick list from DB: %v", err))

	}
	fmt.Println("poop sick list from DB: ", poopSickList)

	// update state in DB based on poop sick list
	for _, tokenID := range poopSickList {
		t := models.NewToken(uint(tokenID))
		t, err := t.GetByTokenID(db)
		if err != nil {
			fmt.Println(fmt.Errorf("sickBot: failed to get tokenID[%d] from DB: %v", tokenID, err))
			continue
		}

		//update tokenID to be sick in DB
		if err := t.Update(db, map[string]interface{}{
			"state": 1,
		}); err != nil {
			fmt.Println(fmt.Errorf("sickBot: failed to update tokenID[%d] to sick: %v", tokenID, err))
		}

	}

	// merge two list
	for _, p := range poopSickList {
		ret = append(ret, big.NewInt(int64(p)))
	}
	ret = uniqueBigInts(ret)

	// filter the tokenIDs that are already sick
	lst, err := healthList(client)
	if err != nil {
		fmt.Println(err)
	}
	sickTokenIDs := HealthListToSickTokenIDs(lst)
	ret = removeCommon(ret, sickTokenIDs)

	if len(ret) == 0 {
		fmt.Println("No one is sick.")
		return
	} else {
		fmt.Println("Sick list: ", ret)
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

	if len(ret) == 0 {
		fmt.Println("No one is dead.")
		return
	} else {
		fmt.Println("Dead list: ", ret)
	}

	err = batachKill(ret)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update state in DB
	db := database.GetSQL()
	for _, tokenID := range ret {
		t := models.NewToken(uint(tokenID.Int64()))
		t, err := t.GetByTokenIDAndLock(db)
		if err != nil {
			fmt.Println(fmt.Errorf("killBot: failed to get tokenID[%d] from DB: %v", tokenID, err))
			continue
		}

		// if on-chain state is dead, update tokenID to be dead in DB
		if t.GetState() != 2 {
			if err := t.Update(db, map[string]interface{}{
				"state": 2,
			}); err != nil {
				fmt.Println(fmt.Errorf("killBot: failed to update tokenID[%d] to dead: %v", tokenID, err))
			}
		}
	}
}

func healthList(client *ethclient.Client) (HeathListRet, error) {
	data, err := parsedABI.Pack("healthList")
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

	var ret HeathListRet
	err = parsedABI.UnpackIntoInterface(&ret, "healthList", res)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack ABI call: %v", err)
	}

	return ret, nil
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

func uniqueBigInts(ints []*big.Int) []*big.Int {
	seen := make(map[string]struct{})
	result := []*big.Int{}

	for _, n := range ints {
		// Use the string representation of the number as the key because big.Int can't be a key.
		str := n.String()
		if _, found := seen[str]; !found {
			seen[str] = struct{}{}
			result = append(result, n)
		}
	}

	return result
}

// removeCommon removes elements from slice1 that are present in slice2
func removeCommon(slice1 []*big.Int, slice2 []uint) []*big.Int {
	// Create a map to store the occurrences of elements in slice2
	occurrences := make(map[uint]struct{})
	for _, item := range slice2 {
		occurrences[item] = struct{}{}
	}

	// Create a new slice to hold the result
	var result []*big.Int

	// Add only elements that are not present in slice2
	for _, item := range slice1 {
		// Convert *big.Int to int for comparison
		// Note: This assumes the values in big.Int can fit into an int, which may not always be true
		// You should add additional checks if your big.Int values are larger than int can handle
		if _, exists := occurrences[uint(item.Int64())]; !exists {
			result = append(result, item)
		}
	}

	return result
}

func HealthListToSickTokenIDs(healthListRet HeathListRet) []uint {
	var ret []uint
	for tokenID, state := range healthListRet {
		if state == 1 {
			ret = append(ret, uint(tokenID))
		}
	}
	return ret
}
