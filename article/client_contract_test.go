package article_test

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/article"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestArticleDomainUsesAllPromotedContractsAndProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			var navigationCalls atomic.Int32
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path == "/x/web-interface/nav" {
					navigationCalls.Add(1)
					return testutil.JSONResponse(http.StatusOK, contracttest.WBINavigationBody), nil
				}
				folder := ""
				wbi := false
				var expected map[string]string
				switch request.URL.Path {
				case "/x/article/viewinfo":
					folder = "info"
					expected = map[string]string{"id": "2"}
				case "/x/article/view":
					folder = "view"
					wbi = true
					expected = map[string]string{"id": "2", "gaia_source": "main_web"}
				case "/x/article/cards":
					folder = "cards"
					wbi = true
					expected = map[string]string{"ids": "av2,cv1,cv2", "web_location": "333.1305"}
				case "/x/article/list/web/articles":
					folder = "articles"
					expected = map[string]string{"id": "207146"}
				default:
					t.Fatalf("unexpected request %s", request.URL)
				}
				if wbi {
					contracttest.AssertWBIFields(t, request.URL.Query(), expected)
				} else {
					contracttest.AssertQuery(t, request.URL.Query(), expected)
				}
				outcome := "success"
				if profile == "anonymous" && (folder == "cards" || folder == "view") {
					outcome = "error"
				}
				body := contracttest.Fixture(t, "article", folder, "responses", profile+"."+outcome+".json")
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			id, _ := ids.NewCVID(2)
			articlesParams, _ := article.NewArticlesInfoParams(207_146)
			cardsParams, _ := article.NewCardsParams("av2,cv1,cv2")
			ctx := context.Background()
			domain := client.Article()

			info, err := domain.Info(ctx, article.NewInfoParams(id))
			if err != nil || info.Title == "" || info.MID.Uint64() == 0 {
				t.Fatalf("Info() = %+v, %v", info, err)
			}
			articles, err := domain.Articles(ctx, articlesParams)
			if err != nil || articles.List.ID != 207_146 || len(articles.Articles) == 0 {
				t.Fatalf("Articles() = list %d items %d, %v", articles.List.ID, len(articles.Articles), err)
			}
			cards, cardsErr := domain.Cards(ctx, cardsParams)
			view, viewErr := domain.View(ctx, article.NewViewParams(id))
			if profile == "anonymous" {
				for name, callErr := range map[string]error{"Cards": cardsErr, "View": viewErr} {
					var apiError *bpi.APIError
					if !errors.As(callErr, &apiError) || apiError.Code != -352 {
						t.Fatalf("%s error = %v, want APIError -352", name, callErr)
					}
				}
			} else {
				if cardsErr != nil || cards["av2"].Kind() != article.CardVideo || cards["cv1"].Kind() != article.CardArticle {
					t.Fatalf("Cards() kinds = %v/%v, %v", cards["av2"].Kind(), cards["cv1"].Kind(), cardsErr)
				}
				if viewErr != nil || view.ID != id || view.Title == "" || view.Content == "" {
					t.Fatalf("View() = id %d title %q content %d, %v", view.ID, view.Title, len(view.Content), viewErr)
				}
			}
			if navigationCalls.Load() != 1 {
				t.Fatalf("navigation calls = %d, want one", navigationCalls.Load())
			}
		})
	}
}
