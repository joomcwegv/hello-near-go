package main

import (
	"strconv"

	"encoding/json"
	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/types"
)

type AuctionState struct {
	Data *collections.LookupMap[string, string]
}

func GetState() AuctionState {
	return AuctionState{
		Data: collections.NewLookupMap[string, string]("a_"),
	}
}

//go:export init
func Init() {
	options := types.ContractInputOptions{IsRawBytes: false}
	contractInput, _, _ := env.ContractInput(options)
	var args struct {
		EndTime    string `json:"end_time"`
		Auctioneer string `json:"auctioneer"`
	}
	if err := json.Unmarshal(contractInput, &args); err != nil {
		env.PanicStr("Failed to parse input")
	}

	endTimeStr := args.EndTime
	if endTimeStr == "" {
		env.PanicStr("Missing end_time")
	}
	
	auctioneer := args.Auctioneer
	if auctioneer == "" {
		env.PanicStr("Missing auctioneer")
	}

	state := GetState()
	state.Data.Insert("end_time", endTimeStr)
	state.Data.Insert("auctioneer", auctioneer)
	state.Data.Insert("highest_bid", "0")
	state.Data.Insert("highest_bidder", "")
	state.Data.Insert("claimed", "false")

	env.LogString("Auction initialized. Auctioneer: " + auctioneer + ", End Time: " + endTimeStr)
}

//go:export bid
func Bid() {
	state := GetState()

	// 1. Check if auction is still active
	endTimeStr, _ := state.Data.Get("end_time")
	endTime, _ := strconv.ParseUint(endTimeStr, 10, 64)
	
	currentTime := env.GetBlockTimeMs() * 1000000
	if currentTime >= endTime {
		env.PanicStr("Auction has already ended")
	}

	// 2. Check if bid is high enough
	attachedDeposit, _ := env.GetAttachedDeposit()

	highestBidStr, _ := state.Data.Get("highest_bid")
	highestBid, _ := types.U128FromString(highestBidStr)

	if attachedDeposit.Cmp(highestBid) <= 0 {
		env.PanicStr("Deposit must be higher than current highest bid")
	}

	// 3. Refund the previous highest bidder
	previousBidder, _ := state.Data.Get("highest_bidder")
	if previousBidder != "" {
		promiseIndex := env.PromiseBatchCreate([]byte(previousBidder))
		env.PromiseBatchActionTransfer(promiseIndex, highestBid)
	}

	// 4. Update state with new highest bidder
	accountId, _ := env.GetPredecessorAccountID()
	state.Data.Insert("highest_bid", attachedDeposit.String())
	state.Data.Insert("highest_bidder", accountId)

	env.LogString("New highest bid: " + attachedDeposit.String() + " by " + accountId)
}

//go:export claim
func Claim() {
	state := GetState()

	// 1. Check if auction has ended
	endTimeStr, _ := state.Data.Get("end_time")
	endTime, _ := strconv.ParseUint(endTimeStr, 10, 64)
	
	currentTime := env.GetBlockTimeMs() * 1000000
	if currentTime < endTime {
		env.PanicStr("Auction is still active")
	}

	// 2. Check if already claimed
	claimed, _ := state.Data.Get("claimed")
	if claimed == "true" {
		env.PanicStr("Proceeds have already been claimed")
	}

	// 3. Transfer highest bid to auctioneer
	auctioneer, _ := state.Data.Get("auctioneer")
	highestBidStr, _ := state.Data.Get("highest_bid")
	highestBid, _ := types.U128FromString(highestBidStr)

	if highestBidStr != "0" {
		promiseIndex := env.PromiseBatchCreate([]byte(auctioneer))
		env.PromiseBatchActionTransfer(promiseIndex, highestBid)
	}

	// 4. Mark as claimed
	state.Data.Insert("claimed", "true")
	env.LogString("Proceeds claimed by auctioneer")
}

//go:export get_highest_bid
func GetHighestBid() {
	state := GetState()
	bidStr, _ := state.Data.Get("highest_bid")
	bidderStr, _ := state.Data.Get("highest_bidder")

	response := `{"bidder": "` + bidderStr + `", "amount": "` + bidStr + `"}`
	env.ContractValueReturn([]byte(response))
}

//go:export get_auction_end_time
func GetAuctionEndTime() {
	state := GetState()
	endTime, _ := state.Data.Get("end_time")
	env.ContractValueReturn([]byte(endTime))
}

//go:export get_auctioneer
func GetAuctioneer() {
	state := GetState()
	auctioneer, _ := state.Data.Get("auctioneer")
	env.ContractValueReturn([]byte(auctioneer))
}

//go:export get_claimed
func GetClaimed() {
	state := GetState()
	claimed, _ := state.Data.Get("claimed")
	env.ContractValueReturn([]byte(claimed))
}

func main() {}
