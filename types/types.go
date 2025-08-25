package types

import "time"

type TzktDelegationsResponse struct {
	Level     int32     `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Sender    Address   `json:"sender"`
	Amount    int64     `json:"amount"`
}

type Address struct {
	Address string `json:"address"`
}

type Delegation struct {
	Delegator string    `json:"delegator"`
	Timestamp time.Time `json:"timestamp"`
	Block     int32     `json:"block"`
	Amount    int64     `json:"amount"`
}

type DelegationsResponse struct {
	Delegations []Delegation `json:"data"`
}
