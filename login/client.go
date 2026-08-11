package login

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	loginNavEndpoint          = "https://api.bilibili.com/x/web-interface/nav"
	loginStatEndpoint         = "https://api.bilibili.com/x/web-interface/nav/stat"
	loginCoinEndpoint         = "https://account.bilibili.com/site/getCoin"
	loginTodayCoinExpEndpoint = "https://api.bilibili.com/x/web-interface/coin/today/exp"
	loginDailyRewardEndpoint  = "https://api.bilibili.com/x/member/web/exp/reward"
	loginAccountInfoEndpoint  = "https://api.bilibili.com/x/member/web/account"
	loginVIPInfoEndpoint      = "https://api.bilibili.com/x/vip/web/user/info"
	loginNoticeEndpoint       = "https://api.bilibili.com/x/safecenter/login_notice"
	loginLogEndpoint          = "https://api.bilibili.com/x/member/web/login/log"
	loginCaptchaEndpoint      = "https://passport.bilibili.com/x/passport-login/captcha"
	loginQRGenerateEndpoint   = "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
	loginQRPollEndpoint       = "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"
)

// Client provides login and authenticated-session operations.
type Client struct{ client *core.Client }

// NewClient binds the login module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// Nav gets the current session's navigation/login state. Anonymous sessions
// normally return an API error for which RequiresLogin reports true.
func (l Client) Nav(ctx context.Context) (Nav, error) {
	return sendLoginPayload[Nav](ctx, l.client, loginNavEndpoint, "login.nav", noLoginQuery)
}

func (l Client) Stat(ctx context.Context) (Stats, error) {
	return sendLoginPayload[Stats](ctx, l.client, loginStatEndpoint, "login.stat", noLoginQuery)
}

func (l Client) Coin(ctx context.Context) (CoinBalance, error) {
	return sendLoginPayload[CoinBalance](ctx, l.client, loginCoinEndpoint, "login.coin", noLoginQuery)
}

func (l Client) TodayCoinExp(ctx context.Context) (TodayCoinExp, error) {
	return sendLoginPayload[TodayCoinExp](ctx, l.client, loginTodayCoinExpEndpoint, "login.today_coin_exp", noLoginQuery)
}

func (l Client) DailyReward(ctx context.Context) (DailyReward, error) {
	return sendLoginPayload[DailyReward](ctx, l.client, loginDailyRewardEndpoint, "login.daily_reward", noLoginQuery)
}

func (l Client) AccountInfo(ctx context.Context) (AccountInfo, error) {
	return sendLoginPayload[AccountInfo](ctx, l.client, loginAccountInfoEndpoint, "login.account_info", noLoginQuery)
}

func (l Client) VIPInfo(ctx context.Context) (VIPInfo, error) {
	return sendLoginPayload[VIPInfo](ctx, l.client, loginVIPInfoEndpoint, "login.vip_info", noLoginQuery)
}

func (l Client) Notice(ctx context.Context, params NoticeParams) (Notice, error) {
	return sendLoginPayload[Notice](ctx, l.client, loginNoticeEndpoint, "login.notice", params.EncodeQuery)
}

func (l Client) Log(ctx context.Context, params LogParams) (Log, error) {
	return sendLoginPayload[Log](ctx, l.client, loginLogEndpoint, "login.log", params.EncodeQuery)
}

func (l Client) GenerateCaptcha(ctx context.Context) (Captcha, error) {
	payload, err := sendLoginPayload[CaptchaPayload](ctx, l.client, loginCaptchaEndpoint, "login.captcha_generate", func() (url.Values, error) {
		return url.Values{"source": {"main_web"}}, nil
	})
	if err != nil {
		return Captcha{}, err
	}
	return Captcha{Token: payload.Token, GT: payload.Geetest.GT, Challenge: payload.Geetest.Challenge}, nil
}

func (l Client) GenerateQR(ctx context.Context) (QRGenerate, error) {
	return sendLoginPayload[QRGenerate](ctx, l.client, loginQRGenerateEndpoint, "login.qr_generate", noLoginQuery)
}

func (l Client) PollQR(ctx context.Context, params QRPollParams) (QRStatus, error) {
	var zero QRStatus
	if l.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "login client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, loginQRPollEndpoint, query)
	if err != nil {
		return zero, err
	}
	response, err := l.client.Do(ctx, request, "login.qr_poll")
	if err != nil {
		return zero, err
	}
	envelope, err := core.DecodeEnvelope[QRStatus](response.Body)
	if err != nil {
		return zero, err
	}
	status, err := envelope.IntoPayload()
	if err != nil {
		return zero, err
	}
	if status.Authenticated() && len(response.Cookies) != 0 {
		status.Cookies = make(map[string]string, len(response.Cookies))
		for _, cookie := range response.Cookies {
			status.Cookies[cookie.Name] = cookie.Value
		}
	}
	return status, nil
}

// QRFlow generates a QR login challenge and immediately performs one poll.
// Callers can continue polling with PollQR and the returned Generate.Key.
func (l Client) QRFlow(ctx context.Context) (QRFlow, error) {
	generated, err := l.GenerateQR(ctx)
	if err != nil {
		return QRFlow{}, err
	}
	params, err := NewQRPollParams(generated.Key)
	if err != nil {
		return QRFlow{}, err
	}
	status, err := l.PollQR(ctx, params)
	if err != nil {
		return QRFlow{}, err
	}
	return QRFlow{Generate: generated, Poll: status}, nil
}

func sendLoginPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "login client is not initialized"}
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

func noLoginQuery() (url.Values, error) { return nil, nil }
