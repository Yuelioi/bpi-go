package client

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Yuelioi/bpi-go/internal/sign"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestWBIKeyCacheSingleFlightsConcurrentFetch(t *testing.T) {
	t.Parallel()

	cache := newWBIKeyCache()
	keys, _ := sign.NewWBIKeys("abcdefghijklmnopqrstuvwxyz123456", "ABCDEFGHIJKLMNOPQRSTUVWXYZ654321")
	var calls atomic.Int32
	started := make(chan struct{})
	release := make(chan struct{})
	fetch := func(context.Context) (sign.WBIKeys, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		<-release
		return keys, nil
	}

	const goroutines = 32
	var group sync.WaitGroup
	group.Add(goroutines)
	errors := make(chan error, goroutines)
	for range goroutines {
		go func() {
			defer group.Done()
			got, err := cache.getOrFetch(context.Background(), "2026-08-11T10", fetch)
			if err == nil && got != keys {
				t.Errorf("keys = %#v, want %#v", got, keys)
			}
			errors <- err
		}()
	}
	<-started
	close(release)
	group.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatalf("getOrFetch() error = %v", err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("fetch calls = %d, want 1", calls.Load())
	}
}

func TestClientWBIKeysAreFetchedOnceAndSignDeterministically(t *testing.T) {
	t.Parallel()

	const navigationBody = `{
        "code": -101,
        "message": "not logged in",
        "data": {
            "wbi_img": {
                "img_url": "https://i0.hdslb.com/bfs/wbi/abcdefghijklmnopqrstuvwxyz123456.png",
                "sub_url": "https://i0.hdslb.com/bfs/wbi/ABCDEFGHIJKLMNOPQRSTUVWXYZ654321.png"
            }
        }
    }`
	var calls atomic.Int32
	fixed := time.Date(2026, 8, 11, 10, 30, 0, 0, time.Local)
	client, err := NewClient(
		WithClock(func() time.Time { return fixed }),
		WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
			calls.Add(1)
			return testutil.JSONResponse(http.StatusOK, navigationBody), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	first, err := client.signWBIParams(context.Background(), map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("signWBIParams(first) error = %v", err)
	}
	second, err := client.signWBIParams(context.Background(), map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("signWBIParams(second) error = %v", err)
	}
	if first["w_rid"] == "" || first["w_rid"] != second["w_rid"] || first["wts"] != second["wts"] {
		t.Fatalf("signed params differ: first=%#v second=%#v", first, second)
	}
	if calls.Load() != 1 {
		t.Fatalf("navigation calls = %d, want 1", calls.Load())
	}
}

func TestClientsHaveIsolatedWBICaches(t *testing.T) {
	t.Parallel()

	first, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient(first) error = %v", err)
	}
	second, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient(second) error = %v", err)
	}
	if first.wbiKeys == second.wbiKeys {
		t.Fatal("clients share WBI cache")
	}
}
