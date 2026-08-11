package audio_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/audio"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestAudioDomainUsesAllPromotedContractsAndProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			type endpoint struct {
				host   string
				folder string
				query  map[string]string
			}
			csrf := ""
			if profile != "anonymous" {
				csrf = "csrf-token"
			}
			endpoints := map[string]endpoint{
				"/audio/music-service-c/web/song/info":              {"www.bilibili.com", "info", map[string]string{"sid": "13603"}},
				"/audio/music-service-c/web/tag/song":               {"www.bilibili.com", "tags", map[string]string{"sid": "13603"}},
				"/audio/music-service-c/web/member/song":            {"www.bilibili.com", "members", map[string]string{"sid": "13603"}},
				"/audio/music-service-c/web/song/lyric":             {"www.bilibili.com", "lyric", map[string]string{"sid": "13603"}},
				"/audio/music-service-c/web/stat/song":              {"www.bilibili.com", "status-number", map[string]string{"sid": "13603"}},
				"/audio/music-service-c/web/collections/songs-coll": {"www.bilibili.com", "collection-status", map[string]string{"sid": "13603"}},
				"/audio/music-service-c/web/coin/audio":             {"www.bilibili.com", "coin-count", map[string]string{"sid": "13603"}},
				"/audio/music-service-c/web/url":                    {"www.bilibili.com", "stream-url-web", map[string]string{"sid": "13603", "quality": "2", "privilege": "2"}},
				"/audio/music-service-c/url":                        {"api.bilibili.com", "stream-url", map[string]string{"songid": "15664", "quality": "2", "privilege": "2", "mid": "2", "platform": "android"}},
				"/audio/music-service-c/web/collections/list":       {"www.bilibili.com", "collections-list", map[string]string{"pn": "1", "ps": "2"}},
				"/audio/music-service-c/web/collections/info":       {"www.bilibili.com", "collection-info", map[string]string{"sid": "15967839"}},
				"/audio/music-service-c/web/menu/hit":               {"www.bilibili.com", "hot-menu", map[string]string{"pn": "1", "ps": "3"}},
				"/audio/music-service-c/web/menu/rank":              {"www.bilibili.com", "rank-menu", map[string]string{"pn": "1", "ps": "6"}},
				"/x/copyright-music-publicity/toplist/all_period":   {"api.bilibili.com", "rank-period", map[string]string{"list_type": "2", "csrf": csrf}},
				"/x/copyright-music-publicity/toplist/detail":       {"api.bilibili.com", "rank-detail", map[string]string{"list_id": "76", "csrf": csrf}},
				"/x/copyright-music-publicity/toplist/music_list":   {"api.bilibili.com", "rank-music-list", map[string]string{"list_id": "76", "csrf": csrf}},
			}
			options := []bpi.Option{bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				endpoint, ok := endpoints[request.URL.Path]
				if !ok {
					t.Fatalf("unexpected request %s", request.URL)
				}
				if request.Method != http.MethodGet || request.URL.Host != endpoint.host {
					t.Fatalf("request = %s %s, want GET host %s", request.Method, request.URL, endpoint.host)
				}
				contracttest.AssertQuery(t, request.URL.Query(), endpoint.query)
				cookie := request.Header.Get("Cookie")
				if profile == "anonymous" && cookie != "" {
					t.Fatalf("anonymous Cookie = %q", cookie)
				}
				if profile != "anonymous" && !strings.Contains(cookie, "SESSDATA=session") {
					t.Fatalf("authenticated Cookie = %q", cookie)
				}
				fixture := "success.json"
				switch endpoint.folder {
				case "info", "tags", "members", "lyric", "status-number":
					fixture = profile + ".success.json"
				case "coin-count", "collection-info", "collection-status", "collections-list":
					if profile == "anonymous" {
						fixture = "anonymous.error.json"
					} else {
						fixture = profile + ".success.json"
					}
				case "rank-detail":
					if profile == "vip" {
						fixture = "vip.success.json"
					} else {
						fixture = "public.success.json"
					}
				}
				body := contracttest.Fixture(t, "audio", endpoint.folder, "responses", fixture)
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})})}
			if profile != "anonymous" {
				options = append(options, bpi.WithAccount(bpi.Account{DedeUserID: "2", SESSDATA: "session", BiliJCT: csrf, Buvid3: "buvid"}))
			}
			client, err := bpi.NewClient(options...)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			sid, _ := ids.NewAudioID(13_603)
			streamSID, _ := ids.NewAudioID(15_664)
			collectionsPage, _ := audio.NewPageParams(1, 2)
			hotPage, _ := audio.NewPageParams(1, 3)
			rankPage, _ := audio.NewPageParams(1, 6)
			collectionInfo, _ := audio.NewCollectionInfoParams(15_967_839)
			rankList, _ := audio.NewRankListParams(76)
			ctx := context.Background()
			domain := client.Audio()
			song := audio.NewSongParams(sid)

			info, err := domain.Info(ctx, song)
			if err != nil || info.ID != sid || info.Title == "" {
				t.Fatalf("Info() = %+v, %v", info, err)
			}
			tags, err := domain.Tags(ctx, song)
			if err != nil || len(tags) == 0 {
				t.Fatalf("Tags() = %d, %v", len(tags), err)
			}
			members, err := domain.Members(ctx, song)
			if err != nil || len(members) == 0 {
				t.Fatalf("Members() = %d, %v", len(members), err)
			}
			lyric, err := domain.Lyric(ctx, song)
			if err != nil || lyric == "" {
				t.Fatalf("Lyric() = %q, %v", lyric, err)
			}
			status, err := domain.StatusNumber(ctx, song)
			if err != nil || status.SID != sid {
				t.Fatalf("StatusNumber() = %+v, %v", status, err)
			}
			webStream, err := domain.StreamURLWeb(ctx, audio.NewStreamURLWebParams(sid))
			if err != nil || len(webStream.CDNs) != 1 {
				t.Fatalf("StreamURLWeb() = %+v, %v", webStream, err)
			}
			stream, err := domain.StreamURL(ctx, audio.NewStreamURLParams(streamSID, audio.QualityHigh))
			if err != nil || len(stream.Qualities) == 0 {
				t.Fatalf("StreamURL() = %+v, %v", stream, err)
			}
			hot, err := domain.HotMenu(ctx, hotPage)
			if err != nil || len(hot.Data) == 0 {
				t.Fatalf("HotMenu() = %+v, %v", hot, err)
			}
			rankMenu, err := domain.RankMenu(ctx, rankPage)
			if err != nil || len(rankMenu.Data) == 0 || len(rankMenu.Data[0].Audios) == 0 {
				t.Fatalf("RankMenu() = %+v, %v", rankMenu, err)
			}
			periods, err := domain.RankPeriod(ctx, audio.NewRankPeriodParams(audio.RankOriginal))
			if err != nil || len(periods.List) == 0 {
				t.Fatalf("RankPeriod() = %+v, %v", periods, err)
			}
			detail, err := domain.RankDetail(ctx, rankList)
			if err != nil || detail.ListenFID == 0 || detail.IsSubscribe != (profile == "vip") {
				t.Fatalf("RankDetail() = %+v, %v", detail, err)
			}
			music, err := domain.RankMusicList(ctx, rankList)
			if err != nil || len(music.List) == 0 || music.List[0].Rank != 1 {
				t.Fatalf("RankMusicList() = %+v, %v", music, err)
			}

			coin, coinErr := domain.CoinCount(ctx, song)
			collected, statusErr := domain.CollectionStatus(ctx, song)
			collections, collectionsErr := domain.CollectionsList(ctx, collectionsPage)
			collection, collectionErr := domain.CollectionInfo(ctx, collectionInfo)
			if profile == "anonymous" {
				for name, callErr := range map[string]error{
					"CoinCount": coinErr, "CollectionStatus": statusErr,
					"CollectionsList": collectionsErr, "CollectionInfo": collectionErr,
				} {
					var apiError *bpi.APIError
					if !errors.As(callErr, &apiError) || apiError.Code != 4_511_003 || !bpi.RequiresLogin(callErr) {
						t.Fatalf("%s error = %v, want login APIError 4511003", name, callErr)
					}
				}
			} else {
				if coinErr != nil || statusErr != nil || collectionsErr != nil || collectionErr != nil {
					t.Fatalf("authenticated errors = %v, %v, %v, %v", coinErr, statusErr, collectionsErr, collectionErr)
				}
				if collected || collection != nil || collections.Data == nil {
					t.Fatalf("authenticated payloads = coin %d collected %v collection %+v list %+v", coin, collected, collection, collections)
				}
				if profile == "vip" && coin != 2 {
					t.Fatalf("vip coin count = %d, want 2", coin)
				}
			}
		})
	}
}
