package golibs

import (
	"time"

	"github.com/shopspring/decimal"
)

type RestReplier interface {
	IsSuccess() bool
	GetError() string
}

type RestReply struct {
	Ret string `json:"ret"`
}

func (reply *RestReply) IsSuccess() bool {
	return reply.Ret == "OK"
}

func (reply *RestReply) GetError() string {
	return reply.Ret
}

type Candle struct {
	Open         decimal.Decimal `json:"o"`
	High         decimal.Decimal `json:"h"`
	Low          decimal.Decimal `json:"l"`
	Close        decimal.Decimal `json:"c"`
	RawTimestamp int64           `json:"ts,string"`
	Timestamp    time.Time
}

type SymbolData struct {
	Symbol       string
	CurrentPrice float64
	Candles      []*Candle
}

type Order struct {
	Symbol   string          `json:"sym"`
	Volume   decimal.Decimal `json:"vol"`
	Price    decimal.Decimal `json:"pri"`
	Datetime string          `json:"date"`
	Type     int             `json:"type"`
	Tag      string          `json:"tag"`
	OrderId  string          `json:"oid"`
}

type Deal struct {
	Symbol   string          `json:"sym"`
	Volume   decimal.Decimal `json:"vol"`
	Price    decimal.Decimal `json:"pri"`
	Datetime string          `json:"date"`
	Tag      string          `json:"tag"`
	OrderId  string          `json:"oid"`
	Usd      decimal.Decimal `json:"usd"`
	Ntd      decimal.Decimal `json:"ntd"`
}

type WatchSymbol struct {
	Symbol string          `json:"sym"`
	Price  decimal.Decimal `json:"pri"`
	Change decimal.Decimal `json:"change"`
}

type Rank struct {
	Hash        string          `json:"hash"`
	Performance decimal.Decimal `json:"perf"`
	Name        string          `json:"name"`
	Tag         string          `json:"tag"`
	Count       int             `json:"cnt"`
	Expire      string          `json:"expire"`
}

type NetValue struct {
	Timestamp int64           `json:"ts,string"`
	Balance   decimal.Decimal `json:"balance"`
}

type QuoteReply struct {
	RestReply
	Price     decimal.Decimal `json:"v"`
	Timestamp int64           `json:"ts,string"`
}

type QuotePeriodReply struct {
	RestReply
	StartTimestamp int64     `json:"ts1,string"`
	EndTimestamp   int64     `json:"ts2,string"`
	Candles        []*Candle `json:"v"`
}

type OrderReply struct {
	RestReply
	OrderId string `json:"roid"`
}

type OrdersReply struct {
	RestReply
	Orders []*Order `json:"orders"`
}

type DealsReply struct {
	RestReply
	Deals []*Deal `json:"deals"`
}

type PositionReply struct {
	RestReply
	Symbols []string          `json:"sym"`
	Volumes []decimal.Decimal `json:"vol"`
	Prices  []decimal.Decimal `json:"pri"`
}

type RightReply struct {
	RestReply
	Right string `json:"right"`
}

type DocIdReply struct {
	RestReply
	DocId string `json:"doc"`
}

type RanksReply struct {
	RestReply
	Ranks []*Rank `json:"rank"`
}

type SubListReply struct {
	RestReply
	Sub []*Rank `json:"sub"`
}

type AllTagsReply struct {
	RestReply
	Tags []string `json:"tags"`
}

type NetValueReply struct {
	RestReply
	NetValues []*NetValue `json:"netvalue"`
}

type ApiTokenReply struct {
	RestReply
	Token string `json:"token"`
}

type WatchListReply struct {
	RestReply
	WatchList []*WatchSymbol `json:"watches"`
}
