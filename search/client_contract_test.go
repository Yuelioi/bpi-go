package search_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/internal/testutil"
	searchmodel "github.com/Yuelioi/bpi-go/search"
)

func TestSearchDomainUsesAllPromotedContracts(t *testing.T) {
	t.Parallel()

	var navCalls atomic.Int32
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		var fixture []byte
		switch request.URL.Path {
		case "/x/web-interface/nav":
			navCalls.Add(1)
			return testutil.JSONResponse(http.StatusOK, contracttest.WBINavigationBody), nil
		case "/x/web-interface/wbi/search/type":
			searchType := query.Get("search_type")
			folder := map[string]string{
				"article":       "article",
				"media_bangumi": "bangumi",
				"bili_user":     "bili-user",
				"live":          "live",
				"live_room":     "live-room",
				"live_user":     "live-user",
				"media_ft":      "movie",
				"video":         "video",
			}[searchType]
			if folder == "" {
				t.Fatalf("unexpected search_type %q", searchType)
			}
			if query.Get("keyword") == "" || query.Get("page") != "1" || query.Get("wts") == "" || len(query.Get("w_rid")) != 32 {
				t.Fatalf("typed search query = %v", query)
			}
			fixture = contracttest.Fixture(t, "search", "read", folder, "responses", "success.json")
		case "/x/web-interface/wbi/search/default":
			contracttest.AssertWBIFields(t, query, map[string]string{"foo": "bar"})
			fixture = contracttest.Fixture(t, "search", "read", "default", "responses", "success.json")
		case "/main/suggest":
			if request.URL.Host != "s.search.bilibili.com" || query.Get("term") != "rust" || query.Get("w_rid") != "" {
				t.Fatalf("suggest request = %s", request.URL)
			}
			fixture = contracttest.Fixture(t, "search", "read", "suggest", "responses", "success.json")
		case "/main/hotword":
			if request.URL.Host != "s.search.bilibili.com" || len(query) != 0 {
				t.Fatalf("hotwords request = %s", request.URL)
			}
			fixture = contracttest.Fixture(t, "search", "read", "hotwords", "responses", "success.json")
		default:
			t.Fatalf("unexpected request %s %s", request.Method, request.URL)
		}
		return testutil.JSONResponse(http.StatusOK, string(fixture)), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	ctx := context.Background()

	articleParams, _ := searchmodel.NewArticleParams("Rust")
	articleParams = articleParams.WithOrder(searchmodel.OrderPublish).WithCategory(searchmodel.ArticleCategoryTechnology)
	article, err := client.Search().Article(ctx, articleParams)
	assertSearchSlice(t, "Article", article.Result, err)

	bangumiParams, _ := searchmodel.NewBangumiParams("天气之子")
	bangumi, err := client.Search().Bangumi(ctx, bangumiParams)
	assertSearchSlice(t, "Bangumi", bangumi.Result, err)

	userParams, _ := searchmodel.NewUserParams("老番茄")
	users, err := client.Search().Users(ctx, userParams.WithOrderSort(searchmodel.OrderDescending))
	assertSearchSlice(t, "Users", users.Result, err)

	liveParams, _ := searchmodel.NewLiveParams("游戏")
	live, err := client.Search().Live(ctx, liveParams)
	if err != nil || live.Result == nil || live.Result.Rooms == nil || live.Result.Users == nil || live.Page != 1 {
		t.Fatalf("Live() = %+v, %v", live, err)
	}

	liveRoomParams, _ := searchmodel.NewLiveRoomParams("游戏")
	rooms, err := client.Search().LiveRooms(ctx, liveRoomParams)
	assertSearchSlice(t, "LiveRooms", rooms.Result, err)

	liveUserParams, _ := searchmodel.NewLiveUserParams("散人")
	liveUsers, err := client.Search().LiveUsers(ctx, liveUserParams.WithOrderSort(searchmodel.OrderDescending))
	assertSearchSlice(t, "LiveUsers", liveUsers.Result, err)

	movieParams, _ := searchmodel.NewMovieParams("哈利波特")
	movies, err := client.Search().Movies(ctx, movieParams)
	assertSearchSlice(t, "Movies", movies.Result, err)

	videoParams, _ := searchmodel.NewVideoParams("Rust 教程")
	videoParams = videoParams.WithOrder(searchmodel.OrderOnline).WithDuration(searchmodel.Duration10To30).WithTID(171)
	videos, err := client.Search().Videos(ctx, videoParams)
	assertSearchSlice(t, "Videos", videos.Result, err)

	defaultResult, err := client.Search().Default(ctx)
	if err != nil || defaultResult.ShowName == "" || defaultResult.GotoValue != "1" {
		t.Fatalf("Default() = %+v, %v", defaultResult, err)
	}
	suggestParams, _ := searchmodel.NewSuggestParams("rust")
	suggest, err := client.Search().Suggest(ctx, suggestParams)
	if err != nil || len(suggest.Tags) != 1 || suggest.Tags[0].Value == nil {
		t.Fatalf("Suggest() = %+v, %v", suggest, err)
	}
	hotWords, err := client.Search().HotWords(ctx)
	if err != nil || len(hotWords.Items) != 1 || hotWords.Items[0].Keyword != "sanitized keyword" {
		t.Fatalf("HotWords() = %+v, %v", hotWords, err)
	}
	if navCalls.Load() != 1 {
		t.Fatalf("nav calls = %d, want one for all signed searches", navCalls.Load())
	}
}

func assertSearchSlice[T any](t *testing.T, method string, result *[]T, err error) {
	t.Helper()
	if err != nil || result == nil || *result == nil {
		t.Fatalf("%s() result = %+v, %v", method, result, err)
	}
}
