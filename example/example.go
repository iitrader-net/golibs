package main

import (
	"log"
	"time"

	"github.com/iitrader-net/golibs"
	"github.com/shopspring/decimal"
)

var (
	RestApi *golibs.RestApi
)

func getQuote(symbol string, timestamp int64) {
	log.Println("===== Test Quote =====")
	quote, err := RestApi.Quote(symbol, timestamp)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("Symbol: ", symbol, " Price: ", quote.Price, " Timestamp:", quote.Timestamp)
}

func getQuotePeriod(symbol string, ts1, ts2 int64) {
	log.Println("===== Test Quote Period =====")
	quote, err := RestApi.QuotePeriod(symbol, ts1, ts2)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("Symbol: ", symbol, " Start: ", quote.StartTimestamp, " End: ", quote.EndTimestamp)
	for _, c := range quote.Candles {
		log.Println("Open: ", c.Open, " High: ", c.High, " Low: ", c.Low, " Close: ", c.Close, " Time: ", c.RawTimestamp)
	}
}

func setOrder(symbol string, volume, price decimal.Decimal, callback string, otype int, tag string) string {
	log.Println("===== Test Order =====")
	order, err := RestApi.Order(symbol, volume, price, callback, otype, tag)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}
	log.Println("Order Id: ", order.OrderId)
	return order.OrderId
}

func cancelOrder(orderId string) {
	log.Println("===== Test Cancel Order =====")
	err := RestApi.CancelOrder(orderId)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("OK")
}

func getOpenOrders() {
	log.Println("===== Test Open Orders =====")
	orders, err := RestApi.OpenOrders()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for _, o := range orders.Orders {
		log.Println("Symbol: ", o.Symbol, " Volume: ", o.Volume, " Price: ", o.Price, " Date: ", o.Datetime,
			" Type: ", o.Type, " Tag: ", o.Tag, " Id: ", o.OrderId)
	}
}

func getHistoricalOrders(page int32) {
	log.Println("===== Test Historical Orders =====")
	orders, err := RestApi.HistoricalOrders(page)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for _, o := range orders.Orders {
		log.Println("Symbol: ", o.Symbol, " Volume: ", o.Volume, " Price: ", o.Price, " Date: ", o.Datetime,
			" Type: ", o.Type, " Tag: ", o.Tag, " Id: ", o.OrderId)
	}
}

func getHistoricalDeals(page int32) {
	log.Println("===== Test Historical Deals =====")
	deals, err := RestApi.HistoricalDeals(page)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for _, o := range deals.Deals {
		log.Println("Symbol: ", o.Symbol, " Volume: ", o.Volume, " Price: ", o.Price, " Date: ", o.Datetime,
			" Tag: ", o.Tag, " Id: ", o.OrderId, " Usd: ", o.Usd, " Ntd: ", o.Ntd)
	}
}

func getPosition() {
	log.Println("===== Test Position =====")
	position, err := RestApi.Position()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for i := 0; i < len(position.Symbols); i += 1 {
		log.Println("Symbol: ", position.Symbols[i], " Volume: ", position.Volumes[i], " Price: ", position.Prices[i])
	}
}

func getRight() {
	log.Println("===== Test Right =====")
	right, err := RestApi.Right()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("Right: ", right.Right)
}

func getDocId() {
	log.Println("===== Test Doc Id=====")
	docId, err := RestApi.DocId()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("Doc id: ", docId.DocId)
}

func setWatch(symbol string) {
	log.Println("===== Test Watch =====")
	err := RestApi.Watch(symbol)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("OK")
}

func delWatch(symbol string) {
	log.Println("===== Test Watch Delete =====")
	err := RestApi.WatchDel(symbol)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("OK")
}

func getWatchList() {
	log.Println("===== Test WatchList =====")
	watchList, err := RestApi.WatchList()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for _, w := range watchList.WatchList {
		log.Println("Symbol: ", w.Symbol, " Price: ", w.Price, " Change: ", w.Change)
	}
}

func getRank() {
	log.Println("===== Test Rank =====")
	ranks, err := RestApi.Rank()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for _, r := range ranks.Ranks {
		log.Println("Hash: ", r.Hash, " Performance: ", r.Performance, " Name: ", r.Name)
	}
}

func getRanks() string {
	log.Println("===== Test Ranks =====")
	ranks, err := RestApi.Ranks()
	if err != nil {
		log.Println("Error:", err)
		return ""
	}
	for _, r := range ranks.Ranks {
		log.Println("Hash: ", r.Hash, " Performance: ", r.Performance, " Name: ", r.Name)
	}
	return ranks.Ranks[0].Hash
}

func setSub(hash string) {
	log.Println("===== Test Sub =====")
	err := RestApi.Sub(hash)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("OK")
}

func getSubList() {
	log.Println("===== Test Sub List =====")
	list, err := RestApi.SubList()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for _, s := range list.Sub {
		log.Println("Hash: ", s.Hash, " Performance: ", s.Performance, " Name: ", s.Name, " Expire: ", s.Expire)
	}
}

func getAllTags() {
	log.Println("===== Test All Tags =====")
	tags, err := RestApi.AllTags()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for _, t := range tags.Tags {
		log.Println("Tag: ", t)
	}
}

func getNetValue() {
	log.Println("===== Test Net Value =====")
	netvalue, err := RestApi.NetValue()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	for _, v := range netvalue.NetValues {
		log.Println("Timestamp: ", v.Timestamp, " Balance: ", v.Balance)
	}
}

func main() {
	// Please fill your own token as the parameter.
	RestApi = golibs.NewRestApi("")

	getQuote("2454.TW", 0)

	startTime := time.Now().Unix() - 24*60*60*7
	getQuotePeriod("2454.TW", startTime, time.Now().Unix())

	orderId := setOrder("2454.TW", decimal.NewFromFloat(1000), decimal.NewFromFloat(200), "", 0, "test")

	getOpenOrders()

	cancelOrder(orderId)

	getHistoricalOrders(0)

	getHistoricalDeals(0)

	getPosition()

	getRight()

	getDocId()

	setWatch("2454.TW")

	delWatch("2454.TW")

	getRank()

	topRankHash := getRanks()

	setSub(topRankHash)

	getSubList()

	getWatchList()

	getAllTags()

	getNetValue()
}
