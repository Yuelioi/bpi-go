package live_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	bpi "github.com/Yuelioi/bpi-go"
	core "github.com/Yuelioi/bpi-go/client"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/live"
)

type liveContractSpec struct {
	method       string
	query        string
	form         string
	folder       string
	defaultFile  string
	profileFiles map[string]string
	wbi          bool
	roomReferer  bool
}

func TestLivePromotedReadsUseAllProfiles(t *testing.T) {
	specs := map[string]liveContractSpec{
		"/room/v1/Area/getList":                                      livePublicSpec("public-core/area-list", ""),
		"/room/v1/Room/get_info":                                     livePublicSpec("public-core/room-info", "room_id=23174842"),
		"/room/v1/Room/playUrl":                                      livePublicSpec("public-core/stream", "cid=14073662&platform=web&qn=10000"),
		"/xlive/web-interface/v1/webMain/getMoreRecList":             livePublicSpec("public-core/recommend", "platform=web&web_location=333.1007"),
		"/xlive/app-blink/v1/liveVersionInfo/getHomePageLiveVersion": livePublicSpec("public-core/version", "system_version=2"),
		"/gift/v1/master/getGiftTypes": {
			method: http.MethodGet, folder: "gift-read/gift-types", defaultFile: "authenticated.success.json",
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json"},
		},
		"/xlive/web-room/v1/giftPanel/roomGiftList": livePublicSpec("gift-read/room-gift-list", "platform=web&room_id=23174842"),
		"/xlive/general-interface/v1/blindFirstWin/getInfo": {
			method: http.MethodGet, query: "gift_id=32251", folder: "gift-read/blind-gift-info", defaultFile: "authenticated.success.json",
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json"},
		},
		"/xlive/web-room/v1/index/getDanmuInfo": {
			method: http.MethodGet, query: "id=21733448&type=0", folder: "room-interaction-read/danmu-info", defaultFile: "authenticated.success.json", wbi: true,
			profileFiles: map[string]string{"anonymous": "anonymous.error.json"},
		},
		"/xlive/web-ucenter/v2/emoticon/GetEmoticons": {
			method: http.MethodGet, query: "platform=pc&room_id=14047", folder: "room-interaction-read/emoticons", defaultFile: "authenticated.success.json",
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json"},
		},
		"/xlive/lottery-interface/v1/lottery/getLotteryInfoWeb": {
			method: http.MethodGet, query: "roomid=23174842", folder: "room-interaction-read/lottery-info", defaultFile: "authenticated.success.json", wbi: true,
			profileFiles: map[string]string{"anonymous": "anonymous.error.json"},
		},
		"/xlive/app-ucenter/v1/user/GetMyMedals": {
			method: http.MethodGet, query: "page=1&page_size=10", folder: "account-private-read/my-medals", defaultFile: "normal.empty.success.json",
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json", "vip": "vip.sample.success.json"},
		},
		"/xlive/web-ucenter/user/following": {
			method: http.MethodGet, query: "hit_ab=true&ignoreRecord=1&page=1&page_size=2", folder: "account-private-read/follow-up-list", defaultFile: "authenticated.success.json",
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json"},
		},
		"/xlive/web-ucenter/v1/xfetter/GetWebList": {
			method: http.MethodGet, query: "hit_ab=false", folder: "account-private-read/follow-up-web-list", defaultFile: "normal.empty.success.json",
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json", "vip": "vip.sample.success.json"},
		},
		"/xlive/app-blink/v1/anchorVideo/AnchorGetReplayList": {
			method: http.MethodGet, query: "page=1&page_size=2", folder: "account-private-read/replay-list", defaultFile: "authenticated.empty.success.json",
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json"},
		},
		"/xlive/app-room/v2/guardTab/topListNew": livePublicSpec("guard-read/guard-list", "page=1&page_size=20&roomid=23174842&ruid=504140200&typ=5"),
		"/xlive/web-ucenter/v1/banned/GetSilentUserList": {
			method: http.MethodPost, form: "csrf=${csrf}&csrf_token=${csrf}&pn=1&ps=10&room_id=3818081", folder: "moderation-private-read/silent-users", defaultFile: "normal.not_admin.json", roomReferer: true,
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json", "vip": "vip.empty.success.json"},
		},
		"/xlive/app-ucenter/v2/xbanned/banned/GetBlackList": {
			method: http.MethodGet, query: "anchor_id=1000001&csrf=${csrf}&csrf_token=${csrf}&mobi_app=android&platform=android&pn=1&ps=10&spmid=444.8.0.0&visit_id=", folder: "moderation-private-read/banned-users", defaultFile: "normal.empty.success.json", roomReferer: true,
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json", "vip": "vip.sample.success.json"},
		},
		"/xlive/app-ucenter/v1/banned/GetShieldKeywordList": {
			method: http.MethodPost, form: "csrf=${csrf}&csrf_token=${csrf}&mobi_app=android&platform=android&room_id=3818081&spmid=444.8.0.0&visit_id=", folder: "moderation-private-read/shield-keywords", defaultFile: "normal.permission_denied.json", roomReferer: true,
			profileFiles: map[string]string{"anonymous": "anonymous.requires_login.json", "vip": "vip.empty.success.json"},
		},
		"/xlive/rdata-interface/v1/heartbeat/webHeartBeat": livePublicSpec("telemetry-read/heartbeat", "hb=NjB8MjMxNzQ4NDJ8MXww&pf=web"),
	}

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			var navCalls atomic.Int32
			client, err := bpi.NewClient(
				core.WithClock(func() time.Time { return time.Unix(1_700_000_000, 0) }),
				bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					if request.URL.Path == "/x/web-interface/nav" {
						navCalls.Add(1)
						return testutil.JSONResponse(http.StatusOK, contracttest.WBINavigationBody), nil
					}
					spec, ok := specs[request.URL.Path]
					if !ok {
						t.Fatalf("unexpected request %s %s", request.Method, request.URL)
					}
					if request.Method != spec.method {
						t.Fatalf("%s method = %s, want %s", request.URL.Path, request.Method, spec.method)
					}
					query := request.URL.Query()
					if spec.wbi {
						if query.Get("wts") != "1700000000" || len(query.Get("w_rid")) != 32 {
							t.Fatalf("%s WBI query = %v", request.URL.Path, query)
						}
						query.Del("wts")
						query.Del("w_rid")
					}
					csrf := ""
					if profile != "anonymous" {
						csrf = "fixture-csrf"
					}
					wantQuery := strings.ReplaceAll(spec.query, "${csrf}", url.QueryEscape(csrf))
					if got := query.Encode(); got != wantQuery {
						t.Fatalf("%s query = %q, want %q", request.URL.Path, got, wantQuery)
					}
					if spec.form != "" {
						body, readErr := io.ReadAll(request.Body)
						if readErr != nil {
							t.Fatalf("read form: %v", readErr)
						}
						wantForm := strings.ReplaceAll(spec.form, "${csrf}", url.QueryEscape(csrf))
						if string(body) != wantForm {
							t.Fatalf("%s form = %q, want %q", request.URL.Path, body, wantForm)
						}
						if request.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
							t.Fatalf("%s Content-Type = %q", request.URL.Path, request.Header.Get("Content-Type"))
						}
					}
					if spec.roomReferer && request.Header.Get("Referer") != "https://live.bilibili.com/3818081" {
						t.Fatalf("%s Referer = %q", request.URL.Path, request.Header.Get("Referer"))
					}
					if request.Header.Get("User-Agent") == "" || request.Header.Get("Origin") == "" || request.Header.Get("Referer") == "" {
						t.Fatalf("%s missing required Bilibili headers", request.URL.Path)
					}
					if profile == "anonymous" && request.Header.Get("Cookie") != "" {
						t.Fatalf("anonymous Cookie = %q", request.Header.Get("Cookie"))
					}
					if profile != "anonymous" && request.Header.Get("Cookie") == "" {
						t.Fatal("authenticated Cookie is empty")
					}
					fileName := spec.defaultFile
					if candidate, ok := spec.profileFiles[profile]; ok {
						fileName = candidate
					}
					parts := append([]string{"live"}, strings.Split(spec.folder, "/")...)
					parts = append(parts, "responses", fileName)
					return testutil.JSONResponse(http.StatusOK, string(contracttest.Fixture(t, parts...))), nil
				})}),
				contracttest.ProfileOption(profile),
			)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			runLiveProfileContract(t, client, profile)
			if navCalls.Load() != 1 {
				t.Fatalf("WBI navigation calls = %d, want 1", navCalls.Load())
			}
		})
	}
}

func livePublicSpec(folder, query string) liveContractSpec {
	return liveContractSpec{method: http.MethodGet, query: query, folder: folder, defaultFile: "success.json"}
}

func runLiveProfileContract(t *testing.T, client *bpi.Client, profile string) {
	t.Helper()
	roomID, _ := ids.NewRoomID(23_174_842)
	streamRoomID, _ := ids.NewRoomID(14_073_662)
	danmuRoomID, _ := ids.NewRoomID(21_733_448)
	emoticonRoomID, _ := ids.NewRoomID(14_047)
	moderationRoomID, _ := ids.NewRoomID(3_818_081)
	anchorID, _ := ids.NewMID(504_140_200)
	moderationAnchorID, _ := ids.NewMID(1_000_001)
	blindParams, _ := live.NewBlindGiftInfoParams(32_251)
	streamParams, _ := live.NewStreamParams(streamRoomID).WithPlatform("web")
	streamParams, _ = streamParams.WithQN(10_000)
	followParams, _ := live.NewFollowUpListParams().WithPage(1)
	followParams, _ = followParams.WithPageSize(2)
	followParams, _ = followParams.WithIgnoreRecord(1)
	followParams = followParams.WithHitAB(true)
	replayParams, _ := live.NewReplayListParams().WithPage(1)
	replayParams, _ = replayParams.WithPageSize(2)
	domain, ctx := client.Live(), context.Background()

	areas, err := domain.AreaList(ctx)
	if err != nil || len(areas) != 1 || len(areas[0].List) != 1 {
		t.Fatalf("AreaList() = %+v, %v", areas, err)
	}
	room, err := domain.RoomInfo(ctx, live.NewRoomInfoParams(roomID))
	if err != nil || room.RoomID != 23_174_842 || room.LiveStatus != 1 {
		t.Fatalf("RoomInfo() = %+v, %v", room, err)
	}
	stream, err := domain.Stream(ctx, streamParams)
	if err != nil || stream.CurrentQN != 10_000 || len(stream.URLs) != 1 {
		t.Fatalf("Stream() = %+v, %v", stream, err)
	}
	recommend, err := domain.Recommend(ctx)
	if err != nil || len(recommend.Rooms) != 1 || recommend.Rooms[0].RoomID != 1 {
		t.Fatalf("Recommend() = %+v, %v", recommend, err)
	}
	version, err := domain.Version(ctx)
	if err != nil || version.CurrentVersion != "7.61.0.10694" {
		t.Fatalf("Version() = %+v, %v", version, err)
	}
	roomGifts, err := domain.RoomGiftList(ctx, live.NewRoomGiftListParams(roomID))
	if err != nil || len(roomGifts.GiftData) == 0 {
		t.Fatalf("RoomGiftList() = %+v, %v", roomGifts, err)
	}
	guards, err := domain.GuardList(ctx, live.NewGuardListParams(roomID, anchorID))
	if err != nil || len(guards.List) != 1 || len(guards.Top3) != 1 {
		t.Fatalf("GuardList() = %+v, %v", guards, err)
	}
	heartBeat, err := domain.WebHeartBeat(ctx, live.NewWebHeartBeatParams(roomID))
	if err != nil || heartBeat.NextInterval != 60 {
		t.Fatalf("WebHeartBeat() = %+v, %v", heartBeat, err)
	}

	giftTypes, giftTypesErr := domain.GiftTypes(ctx)
	blind, blindErr := domain.BlindGiftInfo(ctx, blindParams)
	danmu, danmuErr := domain.DanmuInfo(ctx, live.NewDanmuInfoParams(danmuRoomID))
	emoticons, emoticonsErr := domain.Emoticons(ctx, live.NewEmoticonsParams(emoticonRoomID))
	lottery, lotteryErr := domain.LotteryInfo(ctx, live.NewLotteryInfoParams(roomID))
	medals, medalsErr := domain.MyMedals(ctx, live.NewMyMedalsParams())
	follow, followErr := domain.FollowUpList(ctx, followParams)
	followWeb, followWebErr := domain.FollowUpWebList(ctx, live.NewFollowUpWebListParams().WithHitAB(false))
	replays, replayErr := domain.ReplayList(ctx, replayParams)
	banned, bannedErr := domain.BannedUsers(ctx, live.NewBannedUsersParams(moderationRoomID, moderationAnchorID))
	silent, silentErr := domain.SilentUsers(ctx, live.NewSilentUsersParams(moderationRoomID))
	shield, shieldErr := domain.ShieldKeywords(ctx, live.NewShieldKeywordsParams(moderationRoomID))

	if profile == "anonymous" {
		for name, callErr := range map[string]error{
			"gift-types": giftTypesErr, "blind": blindErr, "emoticons": emoticonsErr,
			"medals": medalsErr, "follow": followErr, "follow-web": followWebErr,
			"replays": replayErr, "banned": bannedErr, "silent": silentErr, "shield": shieldErr,
		} {
			if !bpi.RequiresLogin(callErr) {
				t.Fatalf("%s error = %v, want login error", name, callErr)
			}
		}
		for name, callErr := range map[string]error{"danmu": danmuErr, "lottery": lotteryErr} {
			if !bpi.IsRiskControl(callErr) {
				t.Fatalf("%s error = %v, want risk-control error", name, callErr)
			}
		}
		return
	}

	if giftTypesErr != nil || len(giftTypes) != 0 || blindErr != nil || len(blind.Gifts) != 1 || danmuErr != nil || danmu.Token != "<redacted>" || emoticonsErr != nil || len(emoticons.Packages) != 1 || lotteryErr != nil || lottery.Extra["activity_box"] == nil {
		t.Fatalf("authenticated live reads failed: gift=%v blind=%v danmu=%v emoticons=%v lottery=%v", giftTypesErr, blindErr, danmuErr, emoticonsErr, lotteryErr)
	}
	if medalsErr != nil || followErr != nil || len(follow.List) != 1 || followWebErr != nil || replayErr != nil || replays.Pagination.Page != 1 || bannedErr != nil {
		t.Fatalf("private live reads failed: medals=%v follow=%v follow-web=%v replay=%v banned=%v", medalsErr, followErr, followWebErr, replayErr, bannedErr)
	}
	if profile == "normal" {
		if len(medals.Items) != 0 || len(followWeb.Rooms) != 0 || banned.Total != 0 {
			t.Fatalf("normal private payloads = medals %d rooms %d banned %d", len(medals.Items), len(followWeb.Rooms), banned.Total)
		}
		if !bpi.IsPermissionError(silentErr) || !bpi.IsPermissionError(shieldErr) {
			t.Fatalf("moderation errors = silent %v shield %v", silentErr, shieldErr)
		}
		return
	}
	if len(medals.Items) != 1 || len(followWeb.Rooms) != 1 || banned.Total != 1 || silentErr != nil || silent.Total != 0 || shieldErr != nil || shield.MaxLimit != 1_000 {
		t.Fatalf("vip private payloads failed: medals=%d rooms=%d banned=%d silent=%v shield=%v", len(medals.Items), len(followWeb.Rooms), banned.Total, silentErr, shieldErr)
	}
}

func TestLiveNilClientAndDecodeRecovery(t *testing.T) {
	var domain bpi.LiveClient
	_, err := domain.AreaList(context.Background())
	var parameterError *bpi.ParameterError
	if !errors.As(err, &parameterError) {
		t.Fatalf("AreaList() error = %v, want ParameterError", err)
	}

	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return testutil.JSONResponse(http.StatusOK, `{"code":0,"data":{"room_id":"wrong-type"}}`), nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	roomID, _ := ids.NewRoomID(1)
	_, err = client.Live().RoomInfo(context.Background(), live.NewRoomInfoParams(roomID))
	var decodeError *bpi.ResponseDecodeError
	if !errors.As(err, &decodeError) || len(decodeError.Body()) == 0 {
		t.Fatalf("RoomInfo() error = %v, want recoverable ResponseDecodeError", err)
	}
}
