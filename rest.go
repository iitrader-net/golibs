package golibs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/daveoxley/dnscache"
	"github.com/shopspring/decimal"
)

const (
	MaxRetry     = 5
	RestEndPoint = "http://50.18.230.41:5691"
)

type RestApi struct {
	client *http.Client
	token  string
}

// Customized http.Client, won't kill idle connection.
func NewHttpClient() *http.Client {
	resolver := &dnscache.Resolver{}
	go func() {
		for {
			time.Sleep(time.Hour)
			resolver.Refresh(false)
		}
	}()

	// Based on http.DefaultTransport
	var transport http.RoundTripper = &http.Transport{
		// following are copied from http.DefaultTransport
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// following are my modification
		IdleConnTimeout:     0 * time.Second,
		MaxIdleConnsPerHost: 50,
		DialContext: func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
			separator := strings.LastIndex(addr, ":")
			ips, err := resolver.LookupHost(ctx, addr[:separator])
			if err != nil {
				return nil, err
			}
			dialer := net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}
			for _, ip := range ips {
				conn, err = dialer.DialContext(ctx, network, ip+addr[separator:])
				if err == nil {
					break
				}
			}
			return
		},
	}

	return &http.Client{Transport: transport}
}

func NewRestApi(token string) *RestApi {
	return &RestApi{
		client: NewHttpClient(),
		token:  token,
	}
}

func (api *RestApi) RawRequest(method string, path string, param *string) (*http.Response, error) {
	url := RestEndPoint + path
	var body io.Reader
	if param != nil {
		body = bytes.NewBufferString(*param)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if api.token == "" {
		log.Panic("this client initialized without token")
	}
	req.Header.Add("authorization", api.token)

	trace := &httptrace.ClientTrace{}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	return api.client.Do(req)
}

func (api *RestApi) RequestHelper(method string, path string, param *string,
	reply RestReplier, allow_retry bool) error {
	var err error
	var resp *http.Response
	for retry := 0; retry < MaxRetry && (retry == 0 || allow_retry); retry++ {
		resp, err = api.RawRequest(method, path, param)

		if err != nil {
			log.Println("HTTP failed ", path, " ", err)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		defer resp.Body.Close()

		if 500 <= resp.StatusCode && resp.StatusCode <= 599 {
			log.Println("HTTP failed ", path, " ", resp.StatusCode)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("HTTP Read body fail ", path, err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if err := json.Unmarshal(data, reply); err != nil {
			log.Printf("HTTP failed to decode: %v, json:\n%s", err, data)
			continue
		}
		if !reply.IsSuccess() {
			return errors.New(reply.GetError())
		}
	}
	return err
}

func (api *RestApi) Quote(symbol string, ts int64) (QuoteReply, error) {
	var reply QuoteReply
	str := fmt.Sprintf("/quote/%s", symbol)
	if ts > 0 {
		str += fmt.Sprintf("?ts=%d", ts)
	}
	err := api.RequestHelper(http.MethodGet, str, nil, &reply, false)
	return reply, err
}

func (api *RestApi) QuotePeriod(symbol string, ts1, ts2 int64) (QuotePeriodReply, error) {
	var reply QuotePeriodReply
	str := fmt.Sprintf("/quote_period/%s?ts1=%d&ts2=%d", symbol, ts1, ts2)
	err := api.RequestHelper(http.MethodGet, str, nil, &reply, false)
	for _, candle := range reply.Candles {
		candle.Timestamp = time.Unix(candle.RawTimestamp, 0)
	}
	return reply, err
}

func (api *RestApi) Order(symbol string, volume, price decimal.Decimal, callback string,
	otype int, tag string) (OrderReply, error) {
	var reply OrderReply
	param := fmt.Sprintf(`{"sym":"%s","vol":"%s","pri":"%s","callback":"%s","type":"%d","tag":"%s"}`,
		symbol, volume, price, callback, otype, tag)
	err := api.RequestHelper(http.MethodPost, "/order", &param, &reply, false)
	return reply, err
}

func (api *RestApi) CancelOrder(orderId string) error {
	var reply RestReply
	str := fmt.Sprintf("/cancel?roid=%s", orderId)
	err := api.RequestHelper(http.MethodDelete, str, nil, &reply, false)
	return err
}

func (api *RestApi) OpenOrders() (OrdersReply, error) {
	var reply OrdersReply
	err := api.RequestHelper(http.MethodGet, "/oorder", nil, &reply, false)
	return reply, err
}

func (api *RestApi) HistoricalOrders(page int32) (OrdersReply, error) {
	var reply OrdersReply
	str := fmt.Sprintf("/horder?page=%d", page)
	err := api.RequestHelper(http.MethodGet, str, nil, &reply, false)
	return reply, err
}

func (api *RestApi) HistoricalDeals(page int32) (DealsReply, error) {
	var reply DealsReply
	str := fmt.Sprintf("/hdeal?page=%d", page)
	err := api.RequestHelper(http.MethodGet, str, nil, &reply, false)
	return reply, err
}

func (api *RestApi) Position() (PositionReply, error) {
	var reply PositionReply
	err := api.RequestHelper(http.MethodGet, "/position", nil, &reply, false)
	return reply, err
}

func (api *RestApi) Right() (RightReply, error) {
	var reply RightReply
	err := api.RequestHelper(http.MethodGet, "/right", nil, &reply, false)
	return reply, err
}

func (api *RestApi) DocId() (DocIdReply, error) {
	var reply DocIdReply
	err := api.RequestHelper(http.MethodGet, "/doc", nil, &reply, false)
	return reply, err
}

func (api *RestApi) Watch(symbol string) error {
	var reply RestReply
	param := fmt.Sprintf(`{"sym":"%s"}`, symbol)
	err := api.RequestHelper(http.MethodPost, "/watch", &param, &reply, false)
	return err
}

func (api *RestApi) WatchDel(symbol string) error {
	var reply RestReply
	str := fmt.Sprintf("/watch_del?sym=%s", symbol)
	err := api.RequestHelper(http.MethodDelete, str, nil, &reply, false)
	return err
}

func (api *RestApi) WatchList() (WatchListReply, error) {
	var reply WatchListReply
	err := api.RequestHelper(http.MethodGet, "/watch_list", nil, &reply, false)
	return reply, err
}

func (api *RestApi) Rank() (RanksReply, error) {
	var reply RanksReply
	err := api.RequestHelper(http.MethodGet, "/rank", nil, &reply, false)
	return reply, err
}

func (api *RestApi) Ranks() (RanksReply, error) {
	var reply RanksReply
	err := api.RequestHelper(http.MethodGet, "/ranks", nil, &reply, false)
	return reply, err
}

func (api *RestApi) Sub(hash string) error {
	var reply RestReply
	param := fmt.Sprintf(`{"hash":"%s"}`, hash)
	err := api.RequestHelper(http.MethodPost, "/sub", &param, &reply, false)
	return err
}

func (api *RestApi) SubList() (SubListReply, error) {
	var reply SubListReply
	err := api.RequestHelper(http.MethodGet, "/sub_list", nil, &reply, false)
	return reply, err
}

func (api *RestApi) AllTags() (AllTagsReply, error) {
	var reply AllTagsReply
	err := api.RequestHelper(http.MethodGet, "/alltags", nil, &reply, false)
	return reply, err
}

func (api *RestApi) NetValue() (NetValueReply, error) {
	var reply NetValueReply
	err := api.RequestHelper(http.MethodGet, "/netvalue", nil, &reply, false)
	return reply, err
}

func (api *RestApi) ApiToken() (ApiTokenReply, error) {
	var reply ApiTokenReply
	err := api.RequestHelper(http.MethodGet, "/apitoken", nil, &reply, false)
	return reply, err
}

func (api *RestApi) ReadSymbolData(symbol string, start_time, end_time int64) *SymbolData {
	data := &SymbolData{
		Symbol: symbol,
	}
	if end_time == 0 {
		end_time = time.Now().Unix()
	}
	quoteReply, _ := api.Quote(symbol, 0)
	periodReply, _ := api.QuotePeriod(symbol, start_time, end_time)
	data.CurrentPrice, _ = quoteReply.Price.Float64()
	data.Candles = periodReply.Candles
	return data
}
