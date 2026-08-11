package danmaku

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	danmakuHistoryDatesEndpoint      = "https://api.bilibili.com/x/v2/dm/history/index"
	danmakuSnapshotEndpoint          = "https://api.bilibili.com/x/v2/dm/ajax"
	danmakuThumbupStatsEndpoint      = "https://api.bilibili.com/x/v2/dm/thumbup/stats"
	danmakuAdvStateEndpoint          = "https://api.bilibili.com/x/dm/adv/state"
	danmakuWebSegmentEndpoint        = "https://api.bilibili.com/x/v2/dm/web/seg.so"
	danmakuWebSegmentWBIEndpoint     = "https://api.bilibili.com/x/v2/dm/wbi/web/seg.so"
	danmakuWebViewEndpoint           = "https://api.bilibili.com/x/v2/dm/web/view"
	danmakuMobileSegmentEndpoint     = "https://api.bilibili.com/x/v2/dm/list/seg.so"
	danmakuWebHistorySegmentEndpoint = "https://api.bilibili.com/x/v2/dm/web/history/seg.so"
	danmakuHistoryXMLEndpoint        = "https://api.bilibili.com/x/v2/dm/history"
	danmakuXMLListSOEndpoint         = "https://api.bilibili.com/x/v1/dm/list.so"
)

// Client provides danmaku-domain operations.
type Client struct{ client *core.Client }

// NewClient binds the danmaku module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (d Client) HistoryDates(ctx context.Context, params HistoryDatesParams) ([]string, error) {
	if d.client == nil {
		return nil, &core.ParameterError{Field: "client", Message: "danmaku client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return nil, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, danmakuHistoryDatesEndpoint, query)
	if err != nil {
		return nil, err
	}
	payload, err := core.SendOptionalPayload[[]string](ctx, d.client, request, "danmaku.history.dates")
	if payload == nil || err != nil {
		return nil, err
	}
	return *payload, nil
}

func (d Client) Snapshot(ctx context.Context, params SnapshotParams) ([]string, error) {
	return sendDanmakuPayload[[]string](ctx, d.client, danmakuSnapshotEndpoint, "danmaku.snapshot", params.EncodeQuery)
}

func (d Client) ThumbupStats(ctx context.Context, params ThumbupStatsParams) (ThumbupStats, error) {
	return sendDanmakuPayload[ThumbupStats](ctx, d.client, danmakuThumbupStatsEndpoint, "danmaku.thumbup.stats", params.EncodeQuery)
}

func (d Client) AdvState(ctx context.Context, params AdvStateParams) (AdvState, error) {
	return sendDanmakuPayload[AdvState](ctx, d.client, danmakuAdvStateEndpoint, "danmaku.adv.state", params.EncodeQuery)
}

// WebSegment gets one raw protobuf danmaku segment. It intentionally does not
// attempt JSON envelope decoding; callers may decode it with the official
// DmSegMobileReply protobuf schema.
func (d Client) WebSegment(ctx context.Context, params SegmentParams) ([]byte, error) {
	return sendDanmakuRaw(ctx, d.client, danmakuWebSegmentEndpoint, "danmaku.web.seg", params.EncodeQuery, false)
}

func (d Client) WebSegmentWBI(ctx context.Context, params SegmentParams) ([]byte, error) {
	return sendDanmakuRaw(ctx, d.client, danmakuWebSegmentWBIEndpoint, "danmaku.web.seg_wbi", params.EncodeQuery, true)
}

func (d Client) WebView(ctx context.Context, params WebViewParams) ([]byte, error) {
	return sendDanmakuRaw(ctx, d.client, danmakuWebViewEndpoint, "danmaku.web.view", params.EncodeQuery, false)
}

func (d Client) MobileSegment(ctx context.Context, params SegmentParams) ([]byte, error) {
	return sendDanmakuRaw(ctx, d.client, danmakuMobileSegmentEndpoint, "danmaku.mobile.seg", params.EncodeQuery, false)
}

func (d Client) WebHistorySegment(ctx context.Context, params HistoryBytesParams) ([]byte, error) {
	return sendDanmakuRaw(ctx, d.client, danmakuWebHistorySegmentEndpoint, "danmaku.web.history_seg", params.EncodeQuery, false)
}

func (d Client) HistoryXMLBytes(ctx context.Context, params HistoryBytesParams) ([]byte, error) {
	return sendDanmakuRaw(ctx, d.client, danmakuHistoryXMLEndpoint, "danmaku.history.xml", params.EncodeQuery, false)
}

func (d Client) XMLListSO(ctx context.Context, params XMLListParams) (XML, error) {
	body, err := sendDanmakuRaw(ctx, d.client, danmakuXMLListSOEndpoint, "danmaku.xml.list_so", params.EncodeQuery, false)
	if err != nil {
		return XML{}, err
	}
	return ParseDeflateXML(body)
}

func (d Client) XMLList(ctx context.Context, params XMLListParams) (XML, error) {
	endpoint, err := params.CommentURL()
	if err != nil {
		return XML{}, err
	}
	body, err := sendDanmakuRaw(ctx, d.client, endpoint, "danmaku.xml.comment_xml", func() (url.Values, error) { return nil, nil }, false)
	if err != nil {
		return XML{}, err
	}
	return ParseDeflateXML(body)
}

func sendDanmakuPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "danmaku client is not initialized"}
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

func sendDanmakuRaw(ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error), wbi bool) ([]byte, error) {
	if client == nil {
		return nil, &core.ParameterError{Field: "client", Message: "danmaku client is not initialized"}
	}
	query, err := encode()
	if err != nil {
		return nil, err
	}
	if wbi {
		query, err = client.WBIValues(ctx, query)
		if err != nil {
			return nil, err
		}
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(ctx, request, operation)
	if err != nil {
		return nil, err
	}
	trimmed := bytes.TrimSpace(response.Body)
	if len(trimmed) != 0 && trimmed[0] == '{' {
		envelope, decodeErr := core.DecodeEnvelope[json.RawMessage](response.Body)
		if decodeErr == nil {
			if apiErr := envelope.EnsureSuccess(); apiErr != nil {
				return nil, apiErr
			}
		}
	}
	return response.Body, nil
}
