package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseURL = "http://localhost:8080" // Adjust this to your server address and port
)

type ClaimReq struct {
	Address string `json:"address"`
	RefCode string `json:"ref_code"`
}

type JoinWaitlistReq struct {
	Email string `json:"email"`
}

func claimWallemon(address, refCode string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/claim", baseURL)
	payload := ClaimReq{
		Address: address,
		RefCode: refCode,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func joinWaitlist(email string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/joinWaitlist", baseURL)
	payload := JoinWaitlistReq{
		Email: email,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func main() {
	// Testing claim wallemon
	claimResponse, err := claimWallemon("0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "wallemon")
	if err != nil {
		fmt.Printf("Error claiming wallemon: %v\n", err)
	} else {
		fmt.Printf("Claim Response: %+v\n", claimResponse)
	}

	// Testing join waitlist
	email := "test@example.com"
	joinResponse, err := joinWaitlist(email)
	if err != nil {
		fmt.Printf("Error joining waitlist: %v\n", err)
	} else {
		fmt.Printf("Join Waitlist Response: %+v\n", joinResponse)
	}
}
