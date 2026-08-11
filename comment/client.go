package comment

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	commentListEndpoint    = "https://api.bilibili.com/x/v2/reply"
	commentRepliesEndpoint = "https://api.bilibili.com/x/v2/reply/reply"
	commentHotEndpoint     = "https://api.bilibili.com/x/v2/reply/hot"
	commentCountEndpoint   = "https://api.bilibili.com/x/v2/reply/count"
)

// Client provides comment-domain operations.
type Client struct{ client *core.Client }

// NewClient binds the comment module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// List gets the target's main-comment page.
func (c Client) List(ctx context.Context, params ListParams) (List, error) {
	return sendCommentPayload[List](ctx, c.client, commentListEndpoint, "comment.read.list", params.EncodeQuery)
}

// Replies gets one page of replies below a root comment.
func (c Client) Replies(ctx context.Context, params RepliesParams) (List, error) {
	return sendCommentPayload[List](ctx, c.client, commentRepliesEndpoint, "comment.read.replies", params.EncodeQuery)
}

// Count gets the target's total comment count.
func (c Client) Count(ctx context.Context, params CountParams) (Count, error) {
	return sendCommentPayload[Count](ctx, c.client, commentCountEndpoint, "comment.read.count", params.EncodeQuery)
}

// Hot gets hot replies when the API returns a payload. A successful null or
// missing data field returns (nil, nil).
func (c Client) Hot(ctx context.Context, params HotParams) (*Hot, error) {
	if c.client == nil {
		return nil, &core.ParameterError{Field: "client", Message: "comment client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return nil, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, commentHotEndpoint, query)
	if err != nil {
		return nil, err
	}
	return core.SendOptionalPayload[Hot](ctx, c.client, request, "comment.read.hot")
}

func sendCommentPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "comment client is not initialized"}
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
