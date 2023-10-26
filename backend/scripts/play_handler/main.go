package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	hdlr "wallemon/pkg/api/handlers"
	"wallemon/pkg/utils"

	"github.com/ethereum/go-ethereum/crypto"
)

const baseURL = "http://localhost:8080" // Modify this according to your actual server's address and port.

func main() {
	// GetGem request
	getGemResp, err := getGemRequest("0") // Provide the appropriate token_id value.
	if err != nil {
		fmt.Println("Error in GetGem request:", err)
		return
	}
	fmt.Println("GetGem response:", getGemResp)

	// GetPlay request
	getPlayResp, err := getPlayRequest("0x70997970C51812dc3A010C7d01b50e0d17dc79C8") // Provide the appropriate address value.
	if err != nil {
		fmt.Println("Error in GetPlay request:", err)
		return
	}
	fmt.Println("GetPlay response:", getPlayResp)

	var resp hdlr.GetPlayResp
	err = json.Unmarshal([]byte(getPlayResp), &resp)
	if err != nil {
		fmt.Println("Error in unmarshalling GetPlay response:", err)
		return
	}

	// Play request
	pk := "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	addr := "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	privateKey, err := utils.KeyStringToPrivateKey(pk)
	if err != nil {
		fmt.Println("Error in converting private key:", err)
		return
	}
	tokenID := 0
	msgHash := crypto.Keccak256([]byte(resp.Message))
	sig, err := crypto.Sign(msgHash, privateKey)
	if err != nil {
		fmt.Println("Error in signing:", err)
		return
	}

	// modify the signature to match the format in the backend
	sig[64] += 27

	playResp, err := playRequest(uint(tokenID), addr, hex.EncodeToString(sig)) // Provide appropriate values.
	if err != nil {
		fmt.Println("Error in Play request:", err)
		return
	}
	fmt.Println("Play response:", playResp)

	// GetPoop request
	getPoopResp, err := getPoopRequest("0") // Provide the appropriate token_id value.
	if err != nil {
		fmt.Println("Error in GetPoop request:", err)
		return
	}
	fmt.Println("GetPoop response:", getPoopResp)

	// GetClean request
	getCleanResp, err := getCleanRequest("0x70997970C51812dc3A010C7d01b50e0d17dc79C8") // Provide the appropriate address value.
	if err != nil {
		fmt.Println("Error in GetClean request:", err)
		return
	}
	fmt.Println("GetClean response:", getCleanResp)

	var resp2 hdlr.GetCleanResp
	err = json.Unmarshal([]byte(getCleanResp), &resp2)
	if err != nil {
		fmt.Println("Error in unmarshalling GetClean response:", err)
		return
	}

	// Clean request
	msgHash = crypto.Keccak256([]byte(resp2.Message))
	sig, err = crypto.Sign(msgHash, privateKey)
	if err != nil {
		fmt.Println("Error in signing:", err)
		return
	}

	// modify the signature to match the format in the backend
	sig[64] += 27

	cleanResp, err := cleanRequest(uint(tokenID), addr, hex.EncodeToString(sig)) // Provide appropriate values.
	if err != nil {
		fmt.Println("Error in Clean request:", err)
		return
	}
	fmt.Println("Clean response:", cleanResp)

}

func getGemRequest(tokenID string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/gem", nil) // Adjust the endpoint as needed.
	if err != nil {
		return "", err
	}
	req.Header.Set("token_id", tokenID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getPlayRequest(address string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/play", nil) // Adjust the endpoint as needed.
	if err != nil {
		return "", err
	}
	req.Header.Set("address", address)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func playRequest(tokenID uint, address, signature string) (string, error) {
	data := map[string]interface{}{
		"token_id":  tokenID,
		"address":   address,
		"signature": signature,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+"/play", "application/json", bytes.NewBuffer(jsonData)) // Adjust the endpoint as needed.
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getPoopRequest(tokenID string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/poop", nil) // Adjust the endpoint as needed.
	if err != nil {
		return "", err
	}
	req.Header.Set("token_id", tokenID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getCleanRequest(address string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/clean", nil) // Adjust the endpoint as needed.
	if err != nil {
		return "", err
	}
	req.Header.Set("address", address)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func cleanRequest(tokenID uint, address, signature string) (string, error) {
	data := map[string]interface{}{
		"token_id":  tokenID,
		"address":   address,
		"signature": signature,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+"/clean", "application/json", bytes.NewBuffer(jsonData)) // Adjust the endpoint as needed.
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
