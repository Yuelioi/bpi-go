package note_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/note"
)

func TestNotePromotedReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(
				bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					folder := ""
					expectedQuery := ""
					fileName := ""
					switch request.URL.Path {
					case "/x/note/is_forbid":
						folder, expectedQuery, fileName = "is-forbid", "aid=338677252", "success.json"
					case "/x/note/info":
						folder, expectedQuery = "private-info", "note_id=83577722856540160&oid=676931260&oid_type=0"
						switch profile {
						case "anonymous":
							fileName = "anonymous.requires_login.json"
						case "normal":
							fileName = "normal.not_owner.json"
						default:
							fileName = "vip.success.json"
						}
					case "/x/note/publish/info":
						folder, expectedQuery, fileName = "public-info", "cvid=15160286", "success.json"
					case "/x/note/list/archive":
						folder, expectedQuery = "archive-list", "oid=676931260&oid_type=0"
						if profile == "anonymous" {
							fileName = "anonymous.requires_login.json"
						} else {
							fileName = "authenticated.success.json"
						}
					case "/x/note/list":
						folder, expectedQuery = "user-private-list", "pn=1&ps=10"
						if profile == "anonymous" {
							fileName = "anonymous.requires_login.json"
						} else {
							fileName = "authenticated.success.json"
						}
					case "/x/note/publish/list/archive":
						folder, expectedQuery, fileName = "public-archive-list", "oid=338677252&oid_type=0&pn=1&ps=10", "closed.success.json"
					case "/x/note/publish/list/user":
						folder, expectedQuery = "user-public-list", "pn=1&ps=10"
						if profile == "anonymous" {
							var capture struct {
								Body string `json:"body_base64"`
							}
							fixture := contracttest.Fixture(t, "note", "read", folder, "responses", "anonymous.requires_login.binary.json")
							if err := json.Unmarshal(fixture, &capture); err != nil {
								t.Fatalf("decode binary capture: %v", err)
							}
							body, err := base64.StdEncoding.DecodeString(capture.Body)
							if err != nil {
								t.Fatalf("decode captured body: %v", err)
							}
							if request.URL.Query().Encode() != expectedQuery {
								t.Fatalf("request = %s", request.URL)
							}
							return testutil.JSONResponse(http.StatusOK, string(body)), nil
						}
						fileName = "authenticated.success.json"
					default:
						t.Fatalf("unexpected request %s", request.URL)
					}
					if request.Method != http.MethodGet || request.URL.Query().Encode() != expectedQuery {
						t.Fatalf("request = %s %s, want query %q", request.Method, request.URL, expectedQuery)
					}
					body := contracttest.Fixture(t, "note", "read", folder, "responses", fileName)
					return testutil.JSONResponse(http.StatusOK, string(body)), nil
				})}),
				contracttest.ProfileOption(profile),
			)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			aid, _ := ids.NewAID(338_677_252)
			privateAID, _ := ids.NewAID(676_931_260)
			noteID, _ := ids.NewNoteID(83_577_722_856_540_160)
			cvid, _ := ids.NewCVID(15_160_286)
			domain := client.Note()
			ctx := context.Background()

			forbid, err := domain.IsForbid(ctx, note.NewIsForbidParams(aid))
			if err != nil || forbid.ForbidNoteEntrance {
				t.Fatalf("IsForbid() = %+v, %v", forbid, err)
			}
			publicInfo, err := domain.PublicInfo(ctx, note.NewPublicInfoParams(cvid))
			if err != nil || publicInfo.CVID != cvid.Uint64() || publicInfo.Author.Name != "sanitized author" {
				t.Fatalf("PublicInfo() = %+v, %v", publicInfo, err)
			}
			publicArchive, err := domain.PublicArchiveList(ctx, note.NewPublicArchiveListParams(aid))
			if err != nil || publicArchive.ShowPublicNote {
				t.Fatalf("PublicArchiveList() = %+v, %v", publicArchive, err)
			}

			archive, archiveErr := domain.ArchiveList(ctx, note.NewArchiveListParams(privateAID))
			privateList, privateListErr := domain.UserPrivateList(ctx, note.NewPagination())
			privateInfo, privateInfoErr := domain.PrivateInfo(ctx, note.NewPrivateInfoParams(privateAID, noteID))
			publicUser, publicUserErr := domain.UserPublicList(ctx, note.NewPagination())
			if profile == "anonymous" {
				for name, callErr := range map[string]error{"archive": archiveErr, "private-list": privateListErr, "private-info": privateInfoErr, "public-user": publicUserErr} {
					if !bpi.RequiresLogin(callErr) {
						t.Fatalf("%s error = %v, want login error", name, callErr)
					}
				}
				return
			}
			if archiveErr != nil || len(archive.NoteIDs) != 1 || privateListErr != nil || len(privateList.Items) != 1 || publicUserErr != nil || publicUser.Page == nil {
				t.Fatalf("list reads = archive %+v/%v private %+v/%v public %+v/%v", archive, archiveErr, privateList, privateListErr, publicUser, publicUserErr)
			}
			if profile == "normal" {
				var apiError *bpi.APIError
				if !errors.As(privateInfoErr, &apiError) || apiError.Code != 79_511 || !bpi.IsPermissionError(privateInfoErr) {
					t.Fatalf("PrivateInfo() error = %v, want not-owner error", privateInfoErr)
				}
				return
			}
			if privateInfoErr != nil || privateInfo.Title != "sanitized private note title" {
				t.Fatalf("PrivateInfo() = %+v, %v", privateInfo, privateInfoErr)
			}
		})
	}
}
