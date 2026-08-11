package user_test

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/user"
)

func TestUserDomainUsesAllPromotedSuccessContracts(t *testing.T) {
	t.Parallel()
	type endpointFixture struct {
		host   string
		batch  string
		folder string
		query  map[string]string
		wbi    bool
	}
	fixtures := map[string]endpointFixture{
		"/link_draw/v1/doc/upload_count":        {"api.vc.bilibili.com", "public-read", "album-count", map[string]string{"uid": "2"}, false},
		"/x/space/bangumi/follow/list":          {"api.bilibili.com", "public-read", "bangumi-follow-list", map[string]string{"vmid": "2", "type": "1", "pn": "1", "ps": "15"}, false},
		"/x/web-interface/card":                 {"api.bilibili.com", "public-read", "card", map[string]string{"mid": "2", "photo": "true"}, false},
		"/account/v1/user/cards":                {"api.vc.bilibili.com", "public-read", "cards", map[string]string{"uids": "2,3"}, false},
		"/x/relation/tags":                      {"api.bilibili.com", "relation-read", "follow-tags", map[string]string{}, false},
		"/x/relation/fans":                      {"api.bilibili.com", "relation-read", "followers", map[string]string{"vmid": "2", "ps": "20", "pn": "1"}, false},
		"/x/relation/followings":                {"api.bilibili.com", "relation-read", "followings", map[string]string{"vmid": "2", "order_type": "attention", "ps": "20", "pn": "1"}, false},
		"/x/im/user_infos":                      {"api.vc.bilibili.com", "public-read", "infos", map[string]string{"uids": "2,3"}, false},
		"/xlive/web-ucenter/user/MedalWall":     {"api.live.bilibili.com", "public-read", "medal-wall", map[string]string{"target_id": "2"}, false},
		"/x/polymer/web-dynamic/v1/name-to-uid": {"api.bilibili.com", "public-read", "name-to-uid", map[string]string{"names": "LexBurner,某科学"}, false},
		"/x/space/navnum":                       {"api.bilibili.com", "public-read", "nav-stat", map[string]string{"mid": "2"}, false},
		"/x/relation/stat":                      {"api.bilibili.com", "public-read", "relation-stat", map[string]string{"vmid": "2"}, false},
		"/x/space/wbi/acc/info":                 {"api.bilibili.com", "public-read", "space-info", map[string]string{"mid": "2"}, true},
		"/x/space/notice":                       {"api.bilibili.com", "public-read", "space-notice", map[string]string{"mid": "2"}, false},
		"/x/space/upstat":                       {"api.bilibili.com", "public-read", "up-stat", map[string]string{"mid": "456664753"}, false},
		"/x/space/wbi/arc/search":               {"api.bilibili.com", "public-read", "uploaded-videos", map[string]string{"mid": "2", "order": "pubdate", "tid": "0", "pn": "1", "ps": "30"}, true},
	}

	var navigationCalls atomic.Int32
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", request.Method)
		}
		if request.URL.Path == "/x/web-interface/nav" {
			navigationCalls.Add(1)
			return testutil.JSONResponse(http.StatusOK, contracttest.WBINavigationBody), nil
		}
		fixture, ok := fixtures[request.URL.Path]
		if !ok {
			t.Fatalf("unexpected request %s", request.URL)
		}
		if request.URL.Host != fixture.host {
			t.Fatalf("host = %q, want %q", request.URL.Host, fixture.host)
		}
		if fixture.wbi {
			contracttest.AssertWBIFields(t, request.URL.Query(), fixture.query)
		} else {
			assertUserQuery(t, request.URL.Query(), fixture.query)
		}
		body := contracttest.Fixture(t, "user", fixture.batch, fixture.folder, "responses", "success.json")
		return testutil.JSONResponse(http.StatusOK, string(body)), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	mid, _ := ids.NewMID(2)
	other, _ := ids.NewMID(3)
	creator, _ := ids.NewMID(456_664_753)
	cardsParams, _ := user.NewCardsParams(mid, other)
	infosParams, _ := user.NewInfosParams(mid, other)
	nameParams, _ := user.NewNameToUIDParams("LexBurner", "某科学")
	followingsParams := user.NewFollowingsParams(mid)
	followingsParams, _ = followingsParams.WithOrderType("attention")
	followingsParams, _ = followingsParams.WithPageSize(20)
	followingsParams, _ = followingsParams.WithPage(1)
	followersParams := user.NewFollowersParams(mid)
	followersParams, _ = followersParams.WithPageSize(20)
	followersParams, _ = followersParams.WithPage(1)
	ctx := context.Background()
	domain := client.User()

	album, err := domain.AlbumCount(ctx, user.NewAlbumCountParams(mid))
	if err != nil || album.All != 0 {
		t.Fatalf("AlbumCount() = %+v, %v", album, err)
	}
	bangumi, err := domain.BangumiFollowList(ctx, user.NewBangumiFollowListParams(mid))
	if err != nil || bangumi.Items == nil || bangumi.Page != 1 {
		t.Fatalf("BangumiFollowList() = %+v, %v", bangumi, err)
	}
	card, err := domain.Card(ctx, user.NewCardParams(mid).WithPhoto(user.CardPhotoInclude))
	if err != nil || card.Card.MID != mid || card.Card.Name == "" {
		t.Fatalf("Card() = %+v, %v", card, err)
	}
	cards, err := domain.Cards(ctx, cardsParams)
	if err != nil || len(cards) != 1 || cards[0].MID != mid {
		t.Fatalf("Cards() = %+v, %v", cards, err)
	}
	tags, err := domain.FollowTags(ctx)
	if err != nil || len(tags) != 2 || tags[0].ID != -10 {
		t.Fatalf("FollowTags() = %+v, %v", tags, err)
	}
	followers, err := domain.Followers(ctx, followersParams)
	if err != nil || len(followers.List) != 1 || followers.List[0].MID != other {
		t.Fatalf("Followers() = %+v, %v", followers, err)
	}
	followings, err := domain.Followings(ctx, followingsParams)
	if err != nil || len(followings.List) != 1 || followings.List[0].MID != mid {
		t.Fatalf("Followings() = %+v, %v", followings, err)
	}
	infos, err := domain.Infos(ctx, infosParams)
	if err != nil || len(infos) != 1 || infos[0].Official == nil {
		t.Fatalf("Infos() = %+v, %v", infos, err)
	}
	wall, err := domain.MedalWall(ctx, user.NewMedalWallParams(mid))
	if err != nil || wall.UID != mid || wall.List == nil {
		t.Fatalf("MedalWall() = %+v, %v", wall, err)
	}
	lookup, err := domain.NameToUID(ctx, nameParams)
	if err != nil || len(lookup.Items) != 1 || lookup.Items[0].MID.Uint64() != 777_536 {
		t.Fatalf("NameToUID() = %+v, %v", lookup, err)
	}
	nav, err := domain.NavStat(ctx, user.NewNavStatParams(mid))
	if err != nil || nav.Channel.Master != 0 {
		t.Fatalf("NavStat() = %+v, %v", nav, err)
	}
	relation, err := domain.RelationStat(ctx, user.NewRelationStatParams(mid))
	if err != nil || relation.MID != mid {
		t.Fatalf("RelationStat() = %+v, %v", relation, err)
	}
	space, err := domain.SpaceInfo(ctx, user.NewSpaceParams(mid))
	if err != nil || space.MID != mid || space.Level != 6 {
		t.Fatalf("SpaceInfo() = %+v, %v", space, err)
	}
	notice, err := domain.SpaceNotice(ctx, user.NewSpaceNoticeParams(mid))
	if err != nil || notice.Content != "sanitized notice" {
		t.Fatalf("SpaceNotice() = %+v, %v", notice, err)
	}
	up, err := domain.UpStat(ctx, user.NewUpStatParams(creator))
	if err != nil || up.Likes != 1 || up.Archive.View != 1 {
		t.Fatalf("UpStat() = %+v, %v", up, err)
	}
	uploaded, err := domain.UploadedVideos(ctx, user.NewUploadedVideosParams(mid))
	if err != nil || uploaded.List.Videos == nil || uploaded.Page.Page != 1 {
		t.Fatalf("UploadedVideos() = %+v, %v", uploaded, err)
	}
	if navigationCalls.Load() != 1 {
		t.Fatalf("navigation calls = %d, want one cached WBI key fetch", navigationCalls.Load())
	}
}

func TestUserDomainDecodesPromotedAnonymousOutcomes(t *testing.T) {
	t.Parallel()
	mid, _ := ids.NewMID(2)
	other, _ := ids.NewMID(3)
	cards, _ := user.NewCardsParams(mid, other)
	infos, _ := user.NewInfosParams(mid, other)
	names, _ := user.NewNameToUIDParams("LexBurner", "某科学")
	followers := user.NewFollowersParams(mid)
	followers, _ = followers.WithPageSize(20)
	followers, _ = followers.WithPage(1)
	followings := user.NewFollowingsParams(mid)
	followings, _ = followings.WithOrderType("attention")
	followings, _ = followings.WithPageSize(20)
	followings, _ = followings.WithPage(1)

	tests := []struct {
		name          string
		path          string
		batch         string
		folder        string
		code          int
		requiresLogin bool
		call          func(context.Context, bpi.UserClient) error
	}{
		{"cards", "/account/v1/user/cards", "public-read", "cards", -101, true, func(ctx context.Context, client bpi.UserClient) error { _, err := client.Cards(ctx, cards); return err }},
		{"infos", "/x/im/user_infos", "public-read", "infos", -101, true, func(ctx context.Context, client bpi.UserClient) error { _, err := client.Infos(ctx, infos); return err }},
		{"medal wall", "/xlive/web-ucenter/user/MedalWall", "public-read", "medal-wall", -101, true, func(ctx context.Context, client bpi.UserClient) error {
			_, err := client.MedalWall(ctx, user.NewMedalWallParams(mid))
			return err
		}},
		{"name to uid", "/x/polymer/web-dynamic/v1/name-to-uid", "public-read", "name-to-uid", -101, true, func(ctx context.Context, client bpi.UserClient) error {
			_, err := client.NameToUID(ctx, names)
			return err
		}},
		{"space info", "/x/space/wbi/acc/info", "public-read", "space-info", -352, false, func(ctx context.Context, client bpi.UserClient) error {
			_, err := client.SpaceInfo(ctx, user.NewSpaceParams(mid))
			return err
		}},
		{"follow tags", "/x/relation/tags", "relation-read", "follow-tags", -101, true, func(ctx context.Context, client bpi.UserClient) error { _, err := client.FollowTags(ctx); return err }},
		{"followers", "/x/relation/fans", "relation-read", "followers", -352, false, func(ctx context.Context, client bpi.UserClient) error {
			_, err := client.Followers(ctx, followers)
			return err
		}},
		{"followings", "/x/relation/followings", "relation-read", "followings", -101, true, func(ctx context.Context, client bpi.UserClient) error {
			_, err := client.Followings(ctx, followings)
			return err
		}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path == "/x/web-interface/nav" {
					return testutil.JSONResponse(http.StatusOK, contracttest.WBINavigationBody), nil
				}
				if request.URL.Path != test.path {
					t.Fatalf("path = %q, want %q", request.URL.Path, test.path)
				}
				body := contracttest.Fixture(t, "user", test.batch, test.folder, "responses", "anonymous.error.json")
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			err = test.call(context.Background(), client.User())
			var apiError *bpi.APIError
			if !errors.As(err, &apiError) || apiError.Code != test.code {
				t.Fatalf("error = %v, want APIError code %d", err, test.code)
			}
			if bpi.RequiresLogin(err) != test.requiresLogin {
				t.Fatalf("RequiresLogin(%v) = %v, want %v", err, bpi.RequiresLogin(err), test.requiresLogin)
			}
		})
	}
}

func TestUserUpStatAcceptsPromotedAnonymousEmptyPayload(t *testing.T) {
	t.Parallel()
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		body := contracttest.Fixture(t, "user", "public-read", "up-stat", "responses", "anonymous.empty.json")
		return testutil.JSONResponse(http.StatusOK, string(body)), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	mid, _ := ids.NewMID(456_664_753)
	result, err := client.User().UpStat(context.Background(), user.NewUpStatParams(mid))
	if err != nil || result != (user.UpStat{}) {
		t.Fatalf("UpStat() = %+v, %v; want zero payload", result, err)
	}
}

func assertUserQuery(t *testing.T, got map[string][]string, want map[string]string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("query = %v, want %v", got, want)
	}
	for key, value := range want {
		if len(got[key]) != 1 || got[key][0] != value {
			t.Fatalf("query[%q] = %v, want %q", key, got[key], value)
		}
	}
}
