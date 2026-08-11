package login

import "github.com/Yuelioi/bpi-go/ids"

type Stats struct {
	Following    uint64 `json:"following"`
	Followers    uint64 `json:"follower"`
	DynamicCount uint64 `json:"dynamic_count"`
}

type CoinBalance struct {
	Money float64 `json:"money"`
}

type TodayCoinExp uint32

type DailyReward struct {
	Login        bool   `json:"login"`
	Watch        bool   `json:"watch"`
	Coins        uint32 `json:"coins"`
	Share        bool   `json:"share"`
	Email        bool   `json:"email"`
	Telephone    bool   `json:"tel"`
	SafeQuestion bool   `json:"safe_question"`
	IdentityCard bool   `json:"identify_card"`
}

type AccountInfo struct {
	MID       ids.MID `json:"mid"`
	Username  string  `json:"uname"`
	UserID    string  `json:"userid"`
	Signature string  `json:"sign"`
	Birthday  string  `json:"birthday"`
	Sex       string  `json:"sex"`
	NickFree  bool    `json:"nick_free"`
	Rank      string  `json:"rank"`
}

type VIPInfo struct {
	MID       ids.MID `json:"mid"`
	Type      uint8   `json:"vip_type"`
	Status    uint8   `json:"vip_status"`
	ExpiresAt uint64  `json:"vip_due_date"`
	PayType   uint8   `json:"vip_pay_type"`
	ThemeType uint8   `json:"theme_type"`
}

func (v VIPInfo) Active() bool { return v.Status == 1 && v.ExpiresAt > 0 }

type Notice struct {
	MID        ids.MID `json:"mid"`
	DeviceName string  `json:"device_name"`
	LoginType  string  `json:"login_type"`
	LoginTime  string  `json:"login_time"`
	Location   string  `json:"location"`
	IP         string  `json:"ip"`
}

type LogEntry struct {
	IP        string `json:"ip"`
	Time      uint64 `json:"time"`
	TimeText  string `json:"time_at"`
	Succeeded bool   `json:"status"`
	Type      uint8  `json:"type"`
	Location  string `json:"geo"`
}

type Log struct {
	Count uint32     `json:"count"`
	List  []LogEntry `json:"list"`
}

type Geetest struct {
	Challenge string `json:"challenge"`
	GT        string `json:"gt"`
}

type CaptchaPayload struct {
	Type    string  `json:"type"`
	Token   string  `json:"token"`
	Geetest Geetest `json:"geetest"`
}

type Captcha struct {
	Token     string
	GT        string
	Challenge string
}

type QRGenerate struct {
	URL string `json:"url"`
	Key string `json:"qrcode_key"`
}

const (
	QRSuccess            int32 = 0
	QRExpired            int32 = 86038
	QRScannedUnconfirmed int32 = 86090
	QRWaiting            int32 = 86101
)

type QRStatus struct {
	URL          string            `json:"url"`
	RefreshToken string            `json:"refresh_token"`
	Timestamp    uint64            `json:"timestamp"`
	Code         int32             `json:"code"`
	Message      string            `json:"message"`
	Cookies      map[string]string `json:"-"`
}

func (s QRStatus) Authenticated() bool { return s.Code == QRSuccess }

type QRFlow struct {
	Generate QRGenerate
	Poll     QRStatus
}
