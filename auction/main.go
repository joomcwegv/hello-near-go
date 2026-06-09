package main

import (
	"strconv"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/types"
)

// Helper function to read from storage
func storageRead(key string) string {
	data, err := env.StorageRead([]byte(key))
	if err != nil {
		return ""
	}
	return string(data)
}

// Helper function to write to storage
func storageWrite(key string, value string) {
	env.StorageWrite([]byte(key), []byte(value))
}

//go:export init
func Init() {
	// Simple manual JSON parsing to avoid heavy reflection in tinygo
	options := types.ContractInputOptions{IsRawBytes: false}
	contractInput, _, _ := env.ContractInput(options)
	inputStr := string(contractInput)

	// Naive manual extraction for vibecoder demo
	// We expect input like: {"end_time":"123","auctioneer":"abc"}
	endTime := extractJsonValue(inputStr, "end_time")
	auctioneer := extractJsonValue(inputStr, "auctioneer")

	if endTime == "" {
		env.PanicStr("Missing end_time")
	}
	if auctioneer == "" {
		env.PanicStr("Missing auctioneer")
	}

	storageWrite("end_time", endTime)
	storageWrite("auctioneer", auctioneer)
	storageWrite("highest_bid", "0")
	storageWrite("highest_bidder", "")
	storageWrite("claimed", "false")

	env.LogString("Auction initialized. Auctioneer: " + auctioneer + ", End Time: " + endTime)
}

// Naive JSON value extractor to avoid encoding/json reflection
func extractJsonValue(jsonStr string, key string) string {
	searchKey := `"` + key + `":"`
	startIdx := -1
	for i := 0; i < len(jsonStr)-len(searchKey); i++ {
		if jsonStr[i:i+len(searchKey)] == searchKey {
			startIdx = i + len(searchKey)
			break
		}
	}
	if startIdx == -1 {
		return ""
	}
	endIdx := startIdx
	for i := startIdx; i < len(jsonStr); i++ {
		if jsonStr[i] == '"' {
			endIdx = i
			break
		}
	}
	return jsonStr[startIdx:endIdx]
}

//go:export bid
func Bid() {
	endTimeStr := storageRead("end_time")
	endTime, _ := strconv.ParseUint(endTimeStr, 10, 64)

	currentTime := env.GetBlockTimeMs() * 1000000
	if currentTime >= endTime {
		env.PanicStr("Auction has already ended")
	}

	attachedDeposit, _ := env.GetAttachedDeposit()

	highestBidStr := storageRead("highest_bid")
	if highestBidStr == "" {
		highestBidStr = "0"
	}
	highestBid, _ := types.U128FromString(highestBidStr)

	if attachedDeposit.Cmp(highestBid) <= 0 {
		env.PanicStr("Deposit must be higher than current highest bid")
	}

	previousBidder := storageRead("highest_bidder")
	if previousBidder != "" {
		promiseIndex := env.PromiseBatchCreate([]byte(previousBidder))
		env.PromiseBatchActionTransfer(promiseIndex, highestBid)
	}

	accountId, _ := env.GetPredecessorAccountID()
	storageWrite("highest_bid", attachedDeposit.String())
	storageWrite("highest_bidder", accountId)

	env.LogString("New highest bid: " + attachedDeposit.String() + " by " + accountId)
}

//go:export claim
func Claim() {
	endTimeStr := storageRead("end_time")
	endTime, _ := strconv.ParseUint(endTimeStr, 10, 64)

	currentTime := env.GetBlockTimeMs() * 1000000
	if currentTime < endTime {
		env.PanicStr("Auction is still active")
	}

	claimed := storageRead("claimed")
	if claimed == "true" {
		env.PanicStr("Proceeds have already been claimed")
	}

	auctioneer := storageRead("auctioneer")
	highestBidStr := storageRead("highest_bid")
	if highestBidStr == "" {
		highestBidStr = "0"
	}
	highestBid, _ := types.U128FromString(highestBidStr)

	if highestBidStr != "0" {
		promiseIndex := env.PromiseBatchCreate([]byte(auctioneer))
		env.PromiseBatchActionTransfer(promiseIndex, highestBid)
	}

	storageWrite("claimed", "true")
	env.LogString("Proceeds claimed by auctioneer")
}

//go:export get_highest_bid
func GetHighestBid() {
	bidStr := storageRead("highest_bid")
	bidderStr := storageRead("highest_bidder")

	response := `{"bidder": "` + bidderStr + `", "amount": "` + bidStr + `"}`
	env.ContractValueReturn([]byte(response))
}

//go:export get_auction_end_time
func GetAuctionEndTime() {
	env.ContractValueReturn([]byte(storageRead("end_time")))
}

//go:export get_auctioneer
func GetAuctioneer() {
	env.ContractValueReturn([]byte(storageRead("auctioneer")))
}

//go:export get_claimed
func GetClaimed() {
	env.ContractValueReturn([]byte(storageRead("claimed")))
}

func main() {}
