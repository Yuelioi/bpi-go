package note

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	noteIsForbidEndpoint          = "https://api.bilibili.com/x/note/is_forbid"
	notePrivateInfoEndpoint       = "https://api.bilibili.com/x/note/info"
	notePublicInfoEndpoint        = "https://api.bilibili.com/x/note/publish/info"
	noteArchiveListEndpoint       = "https://api.bilibili.com/x/note/list/archive"
	noteUserPrivateListEndpoint   = "https://api.bilibili.com/x/note/list"
	notePublicArchiveListEndpoint = "https://api.bilibili.com/x/note/publish/list/archive"
	noteUserPublicListEndpoint    = "https://api.bilibili.com/x/note/publish/list/user"
)

// Client provides promoted public, authenticated, and private note reads.
type Client struct{ client *core.Client }

// NewClient binds the note module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (n Client) IsForbid(ctx context.Context, params IsForbidParams) (IsForbid, error) {
	return sendNotePayload[IsForbid](ctx, n.client, noteIsForbidEndpoint, "note.is_forbid", params.EncodeQuery)
}

func (n Client) PrivateInfo(ctx context.Context, params PrivateInfoParams) (PrivateInfo, error) {
	return sendNotePayload[PrivateInfo](ctx, n.client, notePrivateInfoEndpoint, "note.private_info", params.EncodeQuery)
}

func (n Client) PublicInfo(ctx context.Context, params PublicInfoParams) (PublicInfo, error) {
	return sendNotePayload[PublicInfo](ctx, n.client, notePublicInfoEndpoint, "note.public_info", params.EncodeQuery)
}

func (n Client) ArchiveList(ctx context.Context, params ArchiveListParams) (ArchiveList, error) {
	return sendNotePayload[ArchiveList](ctx, n.client, noteArchiveListEndpoint, "note.archive_list", params.EncodeQuery)
}

func (n Client) UserPrivateList(ctx context.Context, params Pagination) (PrivateList, error) {
	return sendNotePayload[PrivateList](ctx, n.client, noteUserPrivateListEndpoint, "note.user_private_list", params.EncodeQuery)
}

func (n Client) PublicArchiveList(ctx context.Context, params PublicArchiveListParams) (PublicArchiveList, error) {
	return sendNotePayload[PublicArchiveList](ctx, n.client, notePublicArchiveListEndpoint, "note.public_archive_list", params.EncodeQuery)
}

func (n Client) UserPublicList(ctx context.Context, params Pagination) (PublicUserList, error) {
	var zero PublicUserList
	if n.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "note client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, noteUserPublicListEndpoint, query)
	if err != nil {
		return zero, err
	}
	response, err := n.client.Do(ctx, request, "note.user_public_list")
	if err != nil {
		return zero, err
	}
	envelope, decodeErr := core.DecodeEnvelope[PublicUserList](response.Body)
	if decodeErr == nil {
		return envelope.IntoPayload()
	}

	// One promoted anonymous capture contains the same login-error JSON object
	// twice in a single HTTP body. Recover the first envelope only for error
	// classification; successful malformed bodies still return DecodeError.
	decoder := json.NewDecoder(bytes.NewReader(response.Body))
	var errorEnvelope struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Msg     string `json:"msg"`
	}
	if err := decoder.Decode(&errorEnvelope); err == nil && errorEnvelope.Code != 0 {
		message := errorEnvelope.Message
		if message == "" {
			message = errorEnvelope.Msg
		}
		return zero, &core.APIError{Code: errorEnvelope.Code, Message: message}
	}
	return zero, decodeErr
}

func sendNotePayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "note client is not initialized"}
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
