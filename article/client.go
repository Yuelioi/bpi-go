package article

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	articleInfoEndpoint     = "https://api.bilibili.com/x/article/viewinfo"
	articleViewEndpoint     = "https://api.bilibili.com/x/article/view"
	articleCardsEndpoint    = "https://api.bilibili.com/x/article/cards"
	articleArticlesEndpoint = "https://api.bilibili.com/x/article/list/web/articles"
)

// Client provides article-domain operations.
type Client struct{ client *core.Client }

// NewClient binds the article module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (a Client) Info(ctx context.Context, params InfoParams) (Info, error) {
	return sendArticlePayload[Info](ctx, a.client, articleInfoEndpoint, "article.info", params.EncodeQuery)
}

func (a Client) View(ctx context.Context, params ViewParams) (View, error) {
	return sendArticleWBI[View](ctx, a.client, articleViewEndpoint, "article.view", params.EncodeQuery)
}

func (a Client) Cards(ctx context.Context, params CardsParams) (CardData, error) {
	return sendArticleWBI[CardData](ctx, a.client, articleCardsEndpoint, "article.cards", params.EncodeQuery)
}

func (a Client) Articles(ctx context.Context, params ArticlesInfoParams) (Articles, error) {
	return sendArticlePayload[Articles](ctx, a.client, articleArticlesEndpoint, "article.articles_info", params.EncodeQuery)
}

func sendArticlePayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "article client is not initialized"}
	}
	query, err := encodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}

func sendArticleWBI[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "article client is not initialized"}
	}
	query, err := encodeQuery()
	if err != nil {
		return zero, err
	}
	signed, err := client.WBIValues(ctx, query)
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, signed)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
