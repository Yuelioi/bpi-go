package client

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/Yuelioi/bpi-go/internal/sign"
)

const wbiNavigationEndpoint = "https://api.bilibili.com/x/web-interface/nav"

type wbiKeyCache struct {
	mu       sync.Mutex
	keys     map[string]sign.WBIKeys
	inFlight map[string]*wbiKeyCall
}

type wbiKeyCall struct {
	done chan struct{}
	keys sign.WBIKeys
	err  error
}

func newWBIKeyCache() *wbiKeyCache {
	return &wbiKeyCache{
		keys:     make(map[string]sign.WBIKeys),
		inFlight: make(map[string]*wbiKeyCall),
	}
}

func (cache *wbiKeyCache) getOrFetch(
	ctx context.Context,
	bucket string,
	fetch func(context.Context) (sign.WBIKeys, error),
) (sign.WBIKeys, error) {
	cache.mu.Lock()
	if keys, ok := cache.keys[bucket]; ok {
		cache.mu.Unlock()
		return keys, nil
	}
	if call, ok := cache.inFlight[bucket]; ok {
		cache.mu.Unlock()
		select {
		case <-ctx.Done():
			return sign.WBIKeys{}, ctx.Err()
		case <-call.done:
			return call.keys, call.err
		}
	}
	call := &wbiKeyCall{done: make(chan struct{})}
	cache.inFlight[bucket] = call
	cache.mu.Unlock()

	call.keys, call.err = fetch(ctx)

	cache.mu.Lock()
	if call.err == nil {
		cache.keys[bucket] = call.keys
	}
	delete(cache.inFlight, bucket)
	close(call.done)
	cache.mu.Unlock()
	return call.keys, call.err
}

func (c *Client) signWBIParams(ctx context.Context, params map[string]string) (map[string]string, error) {
	now := c.now()
	bucket := now.Format("2006-01-02T15")
	keys, err := c.wbiKeys.getOrFetch(ctx, bucket, c.fetchWBIKeys)
	if err != nil {
		return nil, err
	}
	signed, err := sign.SignWBIAt(params, keys, uint64(now.Unix()))
	if err != nil {
		var invalid *sign.InvalidError
		if errors.As(err, &invalid) {
			return nil, &ParameterError{Field: invalid.Field, Message: invalid.Message}
		}
		return nil, err
	}
	return signed, nil
}

func (c *Client) fetchWBIKeys(ctx context.Context) (sign.WBIKeys, error) {
	request, err := http.NewRequest(http.MethodGet, wbiNavigationEndpoint, nil)
	if err != nil {
		return sign.WBIKeys{}, err
	}
	response, err := c.Do(ctx, request, "sign.wbi.navigation")
	if err != nil {
		return sign.WBIKeys{}, err
	}
	envelope, err := DecodeEnvelope[struct {
		WBIImage struct {
			ImageURL string `json:"img_url"`
			SubURL   string `json:"sub_url"`
		} `json:"wbi_img"`
	}](response.Body)
	if err != nil {
		return sign.WBIKeys{}, err
	}
	payload, err := envelope.IntoData()
	if err != nil {
		return sign.WBIKeys{}, err
	}
	keys, err := sign.WBIKeysFromURLs(payload.WBIImage.ImageURL, payload.WBIImage.SubURL)
	if err != nil {
		var invalid *sign.InvalidError
		if errors.As(err, &invalid) {
			return sign.WBIKeys{}, &ParameterError{Field: invalid.Field, Message: invalid.Message}
		}
		return sign.WBIKeys{}, err
	}
	return keys, nil
}
