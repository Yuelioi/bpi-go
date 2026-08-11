package message

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	messageUnreadCountEndpoint  = "https://api.vc.bilibili.com/x/im/web/msgfeed/unread"
	messageReplyFeedEndpoint    = "https://api.bilibili.com/x/msgfeed/reply"
	messageSingleUnreadEndpoint = "https://api.vc.bilibili.com/session_svr/v1/session_svr/single_unread"
)

// Client provides message-domain operations.
type Client struct{ client *core.Client }

// NewClient binds the message module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (m Client) UnreadCount(ctx context.Context, params UnreadCountParams) (UnreadCount, error) {
	return sendMessagePayload[UnreadCount](ctx, m.client, messageUnreadCountEndpoint, "message.unread_count", params.EncodeQuery)
}

func (m Client) ReplyFeed(ctx context.Context, params ReplyFeedParams) (ReplyFeed, error) {
	return sendMessagePayload[ReplyFeed](ctx, m.client, messageReplyFeedEndpoint, "message.reply_feed", params.EncodeQuery)
}

func (m Client) SingleUnread(ctx context.Context, params SingleUnreadParams) (SingleUnread, error) {
	return sendMessagePayload[SingleUnread](ctx, m.client, messageSingleUnreadEndpoint, "message.single_unread", params.EncodeQuery)
}

func sendMessagePayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "message client is not initialized"}
	}
	query, err := encode()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
