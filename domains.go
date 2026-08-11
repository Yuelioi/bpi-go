package bpi

import (
	"github.com/Yuelioi/bpi-go/activity"
	"github.com/Yuelioi/bpi-go/article"
	"github.com/Yuelioi/bpi-go/audio"
	"github.com/Yuelioi/bpi-go/bangumi"
	"github.com/Yuelioi/bpi-go/cheese"
	"github.com/Yuelioi/bpi-go/clientinfo"
	"github.com/Yuelioi/bpi-go/comment"
	"github.com/Yuelioi/bpi-go/creativecenter"
	"github.com/Yuelioi/bpi-go/danmaku"
	"github.com/Yuelioi/bpi-go/dynamic"
	"github.com/Yuelioi/bpi-go/electric"
	"github.com/Yuelioi/bpi-go/fav"
	"github.com/Yuelioi/bpi-go/historytoview"
	"github.com/Yuelioi/bpi-go/live"
	"github.com/Yuelioi/bpi-go/login"
	"github.com/Yuelioi/bpi-go/manga"
	"github.com/Yuelioi/bpi-go/message"
	"github.com/Yuelioi/bpi-go/misc"
	"github.com/Yuelioi/bpi-go/note"
	"github.com/Yuelioi/bpi-go/opus"
	"github.com/Yuelioi/bpi-go/search"
	"github.com/Yuelioi/bpi-go/user"
	"github.com/Yuelioi/bpi-go/video"
	"github.com/Yuelioi/bpi-go/videoranking"
	"github.com/Yuelioi/bpi-go/vip"
	"github.com/Yuelioi/bpi-go/wallet"
	"github.com/Yuelioi/bpi-go/webwidget"
)

// Domain-client aliases preserve the original root package names while each
// implementation lives beside its domain parameters and models.
type (
	ActivityClient       = activity.Client
	ArticleClient        = article.Client
	AudioClient          = audio.Client
	BangumiClient        = bangumi.Client
	CheeseClient         = cheese.Client
	ClientInfoClient     = clientinfo.Client
	CommentClient        = comment.Client
	CreativeCenterClient = creativecenter.Client
	DanmakuClient        = danmaku.Client
	DynamicClient        = dynamic.Client
	ElectricClient       = electric.Client
	FavClient            = fav.Client
	HistoryToViewClient  = historytoview.Client
	LiveClient           = live.Client
	LoginClient          = login.Client
	MangaClient          = manga.Client
	MessageClient        = message.Client
	MiscClient           = misc.Client
	NoteClient           = note.Client
	OpusClient           = opus.Client
	SearchClient         = search.Client
	UserClient           = user.Client
	VideoClient          = video.Client
	VideoRankingClient   = videoranking.Client
	VIPClient            = vip.Client
	WalletClient         = wallet.Client
	WebWidgetClient      = webwidget.Client
)

// Activity returns the activity domain client.
func (c *Client) Activity() ActivityClient { return activity.NewClient(c.core()) }

// Article returns the article domain client.
func (c *Client) Article() ArticleClient { return article.NewClient(c.core()) }

// Audio returns the audio domain client.
func (c *Client) Audio() AudioClient { return audio.NewClient(c.core()) }

// Bangumi returns the bangumi domain client.
func (c *Client) Bangumi() BangumiClient { return bangumi.NewClient(c.core()) }

// Cheese returns the PUGV/course domain client.
func (c *Client) Cheese() CheeseClient { return cheese.NewClient(c.core()) }

// ClientInfo returns the client-information domain client.
func (c *Client) ClientInfo() ClientInfoClient { return clientinfo.NewClient(c.core()) }

// Comment returns the comment domain client.
func (c *Client) Comment() CommentClient { return comment.NewClient(c.core()) }

// CreativeCenter returns the creator-center domain client.
func (c *Client) CreativeCenter() CreativeCenterClient { return creativecenter.NewClient(c.core()) }

// Danmaku returns the danmaku domain client.
func (c *Client) Danmaku() DanmakuClient { return danmaku.NewClient(c.core()) }

// Dynamic returns the dynamic-feed domain client.
func (c *Client) Dynamic() DynamicClient { return dynamic.NewClient(c.core()) }

// Electric returns the charging-support domain client.
func (c *Client) Electric() ElectricClient { return electric.NewClient(c.core()) }

// Fav returns the favorites domain client.
func (c *Client) Fav() FavClient { return fav.NewClient(c.core()) }

// HistoryToView returns the history and watch-later domain client.
func (c *Client) HistoryToView() HistoryToViewClient { return historytoview.NewClient(c.core()) }

// Live returns the live-streaming domain client.
func (c *Client) Live() LiveClient { return live.NewClient(c.core()) }

// Login returns the login and authenticated-session domain client.
func (c *Client) Login() LoginClient { return login.NewClient(c.core()) }

// Manga returns the Bilibili Manga domain client.
func (c *Client) Manga() MangaClient { return manga.NewClient(c.core()) }

// Message returns the message domain client.
func (c *Client) Message() MessageClient { return message.NewClient(c.core()) }

// Misc returns the utility and session-bootstrap domain client.
func (c *Client) Misc() MiscClient { return misc.NewClient(c.core()) }

// Note returns the video-note domain client.
func (c *Client) Note() NoteClient { return note.NewClient(c.core()) }

// Opus returns the opus domain client.
func (c *Client) Opus() OpusClient { return opus.NewClient(c.core()) }

// Search returns the search domain client.
func (c *Client) Search() SearchClient { return search.NewClient(c.core()) }

// User returns the user-profile domain client.
func (c *Client) User() UserClient { return user.NewClient(c.core()) }

// Video returns the video domain client.
func (c *Client) Video() VideoClient { return video.NewClient(c.core()) }

// VideoRanking returns the video-ranking domain client.
func (c *Client) VideoRanking() VideoRankingClient { return videoranking.NewClient(c.core()) }

// VIP returns the VIP-center domain client.
func (c *Client) VIP() VIPClient { return vip.NewClient(c.core()) }

// Wallet returns the private wallet domain client.
func (c *Client) Wallet() WalletClient { return wallet.NewClient(c.core()) }

// WebWidget returns the public Web-widget domain client.
func (c *Client) WebWidget() WebWidgetClient { return webwidget.NewClient(c.core()) }
