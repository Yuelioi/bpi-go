package bpi_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/activity"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/video"
)

func ExampleClient_Activity() {
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: exampleTransport(func(*http.Request) (*http.Response, error) {
		return jsonExampleResponse(`{"code":0,"data":{"id":4017552,"name":"demo activity"}}`), nil
	})}))
	if err != nil {
		panic(err)
	}
	params, err := activity.NewInfoParams(4_017_552)
	if err != nil {
		panic(err)
	}
	info, err := client.Activity().Info(context.Background(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(info.Name)
	// Output: demo activity
}

func ExampleClient_Video() {
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: exampleTransport(func(*http.Request) (*http.Response, error) {
		return jsonExampleResponse(`{"code":0,"data":{"aid":2,"bvid":"BV1xx411c7mD","owner":{"mid":2},"stat":{"aid":2},"cid":62131}}`), nil
	})}))
	if err != nil {
		panic(err)
	}
	bvid, err := ids.NewBVID("BV1xx411c7mD")
	if err != nil {
		panic(err)
	}
	view, err := client.Video().View(context.Background(), video.ViewByBVID(bvid))
	if err != nil {
		panic(err)
	}
	fmt.Println(view.BVID)
	// Output: BV1xx411c7mD
}

func ExampleNewClient_authenticated() {
	client, err := bpi.NewClient(bpi.WithAccount(bpi.Account{
		DedeUserID: "fixture-user",
		SESSDATA:   "fixture-session",
		BiliJCT:    "fixture-csrf",
		Buvid3:     "fixture-device",
	}))
	if err != nil {
		panic(err)
	}
	fmt.Println(client.HasLoginCookies())
	// Output: true
}

func ExampleSendPayload() {
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: exampleTransport(func(*http.Request) (*http.Response, error) {
		return jsonExampleResponse(`{"code":0,"data":{"value":"custom payload"}}`), nil
	})}))
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.bilibili.com/x/example", nil)
	if err != nil {
		panic(err)
	}
	payload, err := bpi.SendPayload[struct {
		Value string `json:"value"`
	}](context.Background(), client, request, "example.custom")
	if err != nil {
		panic(err)
	}
	fmt.Println(payload.Value)
	// Output: custom payload
}

func ExampleResponseDecodeError_Body() {
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: exampleTransport(func(*http.Request) (*http.Response, error) {
		return jsonExampleResponse(`{"code":0,"data":{"count":"unexpected"}}`), nil
	})}))
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.bilibili.com/x/example", nil)
	if err != nil {
		panic(err)
	}
	_, err = bpi.SendPayload[struct {
		Count int `json:"count"`
	}](context.Background(), client, request, "example.decode")
	var decodeError *bpi.ResponseDecodeError
	fmt.Println(errors.As(err, &decodeError), len(decodeError.Body()) > 0)
	// Output: true true
}

type exampleTransport func(*http.Request) (*http.Response, error)

func (transport exampleTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

func jsonExampleResponse(body string) *http.Response {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	return response
}
