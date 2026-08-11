package fav_test

import (
	"context"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/fav"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestFavReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	mediaID, _ := ids.NewMediaID(1052622027)
	mid, _ := ids.NewMID(7792521)
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				folder := map[string]string{
					"/x/v3/fav/folder/info":             "folder-info",
					"/x/v3/fav/folder/created/list-all": "created-list",
					"/x/v3/fav/folder/collected/list":   "collected-list",
					"/x/v3/fav/resource/infos":          "resource-infos",
					"/x/v3/fav/resource/list":           "list-detail",
					"/x/v3/fav/resource/ids":            "resource-ids",
				}[request.URL.Path]
				if request.Method != http.MethodGet || folder == "" {
					t.Fatalf("unexpected request %s %s", request.Method, request.URL)
				}
				body := contracttest.Fixture(t, "fav", "read", folder, "responses", "success.json")
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			domain, ctx := client.Fav(), context.Background()
			folder, err := domain.FolderInfo(ctx, fav.NewFolderInfoParams(mediaID))
			if err != nil || folder.ID != 1052622027 {
				t.Fatalf("FolderInfo() = %+v, %v", folder, err)
			}
			created, err := domain.CreatedList(ctx, fav.NewCreatedListParams(mid))
			if err != nil || len(created.List) != 1 {
				t.Fatalf("CreatedList() = %+v, %v", created, err)
			}
			collected, err := domain.CollectedList(ctx, fav.NewCollectedListParams(mid))
			if err != nil || collected.Count != 0 {
				t.Fatalf("CollectedList() = %+v, %v", collected, err)
			}
			resources, _ := fav.NewResourceInfosParams("371494037:2")
			infos, err := domain.ResourceInfos(ctx, resources)
			if err != nil || len(infos) != 1 {
				t.Fatalf("ResourceInfos() = %+v, %v", infos, err)
			}
			detailParams, _ := fav.NewListDetailParams(mediaID).WithOrder("mtime")
			detailParams = detailParams.WithContentType(0)
			detailParams, _ = detailParams.WithPageSize(5)
			detailParams, _ = detailParams.WithPage(1)
			detail, err := domain.ListDetail(ctx, detailParams)
			if err != nil || len(detail.Medias) != 1 {
				t.Fatalf("ListDetail() = %+v, %v", detail, err)
			}
			resourceIDs, err := domain.ResourceIDs(ctx, fav.NewResourceIDsParams(mediaID))
			if err != nil || len(resourceIDs) != 1 {
				t.Fatalf("ResourceIDs() = %+v, %v", resourceIDs, err)
			}
		})
	}
}
