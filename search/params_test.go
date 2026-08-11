package search

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPromotedQueries(t *testing.T) {
	t.Parallel()

	article, _ := NewArticleParams(" Rust ")
	article = article.WithOrder(OrderPublish).WithCategory(ArticleCategoryTechnology)
	query, err := article.EncodeQuery()
	if err != nil || query.Get("search_type") != "article" || query.Get("keyword") != "Rust" || query.Get("order") != "pubdate" || query.Get("category_id") != "17" || query.Get("page") != "1" {
		t.Fatalf("article query = %v, %v", query, err)
	}
	user, _ := NewUserParams("老番茄")
	query, err = user.WithOrderSort(OrderDescending).EncodeQuery()
	if err != nil || query.Get("search_type") != "bili_user" || query.Get("order_sort") != "0" || query.Get("user_type") != "0" {
		t.Fatalf("user query = %v, %v", query, err)
	}
	video, _ := NewVideoParams("Rust 教程")
	video = video.WithOrder(OrderOnline).WithDuration(Duration10To30).WithTID(171)
	query, err = video.EncodeQuery()
	if err != nil || query.Get("search_type") != "video" || query.Get("duration") != "2" || query.Get("tids") != "171" {
		t.Fatalf("video query = %v, %v", query, err)
	}
	suggest, _ := NewSuggestParams(" rust ")
	query, err = suggest.EncodeQuery()
	if err != nil || query.Get("term") != "rust" {
		t.Fatalf("suggest query = %v, %v", query, err)
	}
}

func TestParamsRejectInvalidValues(t *testing.T) {
	t.Parallel()

	var parameterError *bpierr.ParameterError
	if _, err := NewVideoParams("  "); !errors.As(err, &parameterError) || parameterError.Field != "keyword" {
		t.Fatalf("blank keyword error = %v", err)
	}
	params, _ := NewVideoParams("rust")
	if _, err := params.WithPage(0); !errors.As(err, &parameterError) || parameterError.Field != "page" {
		t.Fatalf("zero page error = %v", err)
	}
	if _, err := params.WithOrder(Order("broken")).EncodeQuery(); !errors.As(err, &parameterError) || parameterError.Field != "order" {
		t.Fatalf("bad order error = %v", err)
	}
	if _, err := params.WithDuration(Duration(99)).EncodeQuery(); !errors.As(err, &parameterError) || parameterError.Field != "duration" {
		t.Fatalf("bad duration error = %v", err)
	}
	if _, err := NewSuggestParams(""); !errors.As(err, &parameterError) || parameterError.Field != "term" {
		t.Fatalf("blank term error = %v", err)
	}
}
