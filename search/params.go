// Package search contains validated parameters and response models for
// Bilibili search endpoints.
package search

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type Order string

const (
	OrderTotalRank Order = "totalrank"
	OrderClick     Order = "click"
	OrderPublish   Order = "pubdate"
	OrderDanmaku   Order = "dm"
	OrderFavorite  Order = "stow"
	OrderScore     Order = "scores"
	OrderAttention Order = "attention"
	OrderOnline    Order = "online"
	OrderLiveTime  Order = "live_time"
	OrderDefault   Order = "0"
	OrderFans      Order = "fans"
	OrderLevel     Order = "level"
)

func (order Order) validate() error {
	switch order {
	case OrderTotalRank, OrderClick, OrderPublish, OrderDanmaku, OrderFavorite, OrderScore,
		OrderAttention, OrderOnline, OrderLiveTime, OrderDefault, OrderFans, OrderLevel:
		return nil
	default:
		return &bpierr.ParameterError{Field: "order", Message: "unsupported search order"}
	}
}

type OrderSort uint8

const (
	OrderDescending OrderSort = 0
	OrderAscending  OrderSort = 1
)

type UserType uint8

const (
	UserAll      UserType = 0
	UserUploader UserType = 1
	UserNormal   UserType = 2
	UserVerified UserType = 3
)

type Duration uint8

const (
	DurationAll     Duration = 0
	DurationUnder10 Duration = 1
	Duration10To30  Duration = 2
	Duration30To60  Duration = 3
	DurationOver60  Duration = 4
)

type ArticleCategory uint8

const (
	ArticleCategoryAll        ArticleCategory = 0
	ArticleCategoryGame       ArticleCategory = 1
	ArticleCategoryAnimation  ArticleCategory = 2
	ArticleCategoryLife       ArticleCategory = 3
	ArticleCategoryLightNovel ArticleCategory = 16
	ArticleCategoryTechnology ArticleCategory = 17
	ArticleCategoryMovie      ArticleCategory = 28
	ArticleCategoryInterest   ArticleCategory = 29
)

type commonParams struct {
	keyword string
	page    uint32
}

func newCommonParams(keyword string) (commonParams, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return commonParams{}, &bpierr.ParameterError{Field: "keyword", Message: "search keyword cannot be blank"}
	}
	return commonParams{keyword: keyword, page: 1}, nil
}

func (p commonParams) withPage(page uint32) (commonParams, error) {
	if page == 0 {
		return commonParams{}, &bpierr.ParameterError{Field: "page", Message: "page number must be at least 1"}
	}
	p.page = page
	return p, nil
}

func (p commonParams) values(searchType string) (url.Values, error) {
	if strings.TrimSpace(p.keyword) == "" {
		return nil, &bpierr.ParameterError{Field: "keyword", Message: "search keyword cannot be blank"}
	}
	if p.page == 0 {
		return nil, &bpierr.ParameterError{Field: "page", Message: "page number must be at least 1"}
	}
	return url.Values{
		"search_type": {searchType},
		"keyword":     {p.keyword},
		"page":        {strconv.FormatUint(uint64(p.page), 10)},
	}, nil
}

type ArticleParams struct {
	common   commonParams
	order    Order
	category ArticleCategory
}

func NewArticleParams(keyword string) (ArticleParams, error) {
	common, err := newCommonParams(keyword)
	return ArticleParams{common: common, order: OrderTotalRank, category: ArticleCategoryAll}, err
}

func (p ArticleParams) WithOrder(order Order) ArticleParams { p.order = order; return p }
func (p ArticleParams) WithCategory(category ArticleCategory) ArticleParams {
	p.category = category
	return p
}
func (p ArticleParams) WithPage(page uint32) (ArticleParams, error) {
	common, err := p.common.withPage(page)
	if err != nil {
		return ArticleParams{}, err
	}
	p.common = common
	return p, nil
}
func (p ArticleParams) EncodeQuery() (url.Values, error) {
	values, err := p.common.values("article")
	if err != nil {
		return nil, err
	}
	if err := p.order.validate(); err != nil {
		return nil, err
	}
	switch p.category {
	case ArticleCategoryAll, ArticleCategoryGame, ArticleCategoryAnimation, ArticleCategoryLife,
		ArticleCategoryLightNovel, ArticleCategoryTechnology, ArticleCategoryMovie, ArticleCategoryInterest:
	default:
		return nil, &bpierr.ParameterError{Field: "category_id", Message: "unsupported article category"}
	}
	values.Set("order", string(p.order))
	values.Set("category_id", strconv.FormatUint(uint64(p.category), 10))
	return values, nil
}

type BangumiParams struct{ common commonParams }

func NewBangumiParams(keyword string) (BangumiParams, error) {
	common, err := newCommonParams(keyword)
	return BangumiParams{common: common}, err
}
func (p BangumiParams) WithPage(page uint32) (BangumiParams, error) {
	common, err := p.common.withPage(page)
	if err != nil {
		return BangumiParams{}, err
	}
	p.common = common
	return p, nil
}
func (p BangumiParams) EncodeQuery() (url.Values, error) { return p.common.values("media_bangumi") }

type UserParams struct {
	common    commonParams
	orderSort OrderSort
	userType  UserType
}

func NewUserParams(keyword string) (UserParams, error) {
	common, err := newCommonParams(keyword)
	return UserParams{common: common, orderSort: OrderAscending, userType: UserAll}, err
}
func (p UserParams) WithOrderSort(sort OrderSort) UserParams   { p.orderSort = sort; return p }
func (p UserParams) WithUserType(userType UserType) UserParams { p.userType = userType; return p }
func (p UserParams) WithPage(page uint32) (UserParams, error) {
	common, err := p.common.withPage(page)
	if err != nil {
		return UserParams{}, err
	}
	p.common = common
	return p, nil
}
func (p UserParams) EncodeQuery() (url.Values, error) {
	values, err := p.common.values("bili_user")
	if err != nil {
		return nil, err
	}
	if p.orderSort > OrderAscending {
		return nil, &bpierr.ParameterError{Field: "order_sort", Message: "value must be descending or ascending"}
	}
	if p.userType > UserVerified {
		return nil, &bpierr.ParameterError{Field: "user_type", Message: "unsupported user type"}
	}
	values.Set("order_sort", strconv.FormatUint(uint64(p.orderSort), 10))
	values.Set("user_type", strconv.FormatUint(uint64(p.userType), 10))
	return values, nil
}

type LiveParams struct{ common commonParams }

func NewLiveParams(keyword string) (LiveParams, error) {
	common, err := newCommonParams(keyword)
	return LiveParams{common: common}, err
}
func (p LiveParams) WithPage(page uint32) (LiveParams, error) {
	common, err := p.common.withPage(page)
	if err != nil {
		return LiveParams{}, err
	}
	p.common = common
	return p, nil
}
func (p LiveParams) EncodeQuery() (url.Values, error) { return p.common.values("live") }

type LiveRoomParams struct {
	common commonParams
	order  Order
}

func NewLiveRoomParams(keyword string) (LiveRoomParams, error) {
	common, err := newCommonParams(keyword)
	return LiveRoomParams{common: common, order: OrderOnline}, err
}
func (p LiveRoomParams) WithOrder(order Order) LiveRoomParams { p.order = order; return p }
func (p LiveRoomParams) WithPage(page uint32) (LiveRoomParams, error) {
	common, err := p.common.withPage(page)
	if err != nil {
		return LiveRoomParams{}, err
	}
	p.common = common
	return p, nil
}
func (p LiveRoomParams) EncodeQuery() (url.Values, error) {
	values, err := p.common.values("live_room")
	if err != nil {
		return nil, err
	}
	if err := p.order.validate(); err != nil {
		return nil, err
	}
	values.Set("order", string(p.order))
	return values, nil
}

type LiveUserParams struct {
	common    commonParams
	orderSort OrderSort
	userType  UserType
}

func NewLiveUserParams(keyword string) (LiveUserParams, error) {
	common, err := newCommonParams(keyword)
	return LiveUserParams{common: common, orderSort: OrderAscending, userType: UserAll}, err
}
func (p LiveUserParams) WithOrderSort(sort OrderSort) LiveUserParams { p.orderSort = sort; return p }
func (p LiveUserParams) WithUserType(userType UserType) LiveUserParams {
	p.userType = userType
	return p
}
func (p LiveUserParams) WithPage(page uint32) (LiveUserParams, error) {
	common, err := p.common.withPage(page)
	if err != nil {
		return LiveUserParams{}, err
	}
	p.common = common
	return p, nil
}
func (p LiveUserParams) EncodeQuery() (url.Values, error) {
	values, err := p.common.values("live_user")
	if err != nil {
		return nil, err
	}
	if p.orderSort > OrderAscending {
		return nil, &bpierr.ParameterError{Field: "order_sort", Message: "value must be descending or ascending"}
	}
	if p.userType > UserVerified {
		return nil, &bpierr.ParameterError{Field: "user_type", Message: "unsupported user type"}
	}
	values.Set("order_sort", strconv.FormatUint(uint64(p.orderSort), 10))
	values.Set("user_type", strconv.FormatUint(uint64(p.userType), 10))
	return values, nil
}

type MovieParams struct{ common commonParams }

func NewMovieParams(keyword string) (MovieParams, error) {
	common, err := newCommonParams(keyword)
	return MovieParams{common: common}, err
}
func (p MovieParams) WithPage(page uint32) (MovieParams, error) {
	common, err := p.common.withPage(page)
	if err != nil {
		return MovieParams{}, err
	}
	p.common = common
	return p, nil
}
func (p MovieParams) EncodeQuery() (url.Values, error) { return p.common.values("media_ft") }

type VideoParams struct {
	common   commonParams
	order    Order
	duration Duration
	tid      uint32
}

func NewVideoParams(keyword string) (VideoParams, error) {
	common, err := newCommonParams(keyword)
	return VideoParams{common: common, order: OrderTotalRank, duration: DurationAll}, err
}
func (p VideoParams) WithOrder(order Order) VideoParams          { p.order = order; return p }
func (p VideoParams) WithDuration(duration Duration) VideoParams { p.duration = duration; return p }
func (p VideoParams) WithTID(tid uint32) VideoParams             { p.tid = tid; return p }
func (p VideoParams) WithPage(page uint32) (VideoParams, error) {
	common, err := p.common.withPage(page)
	if err != nil {
		return VideoParams{}, err
	}
	p.common = common
	return p, nil
}
func (p VideoParams) EncodeQuery() (url.Values, error) {
	values, err := p.common.values("video")
	if err != nil {
		return nil, err
	}
	if err := p.order.validate(); err != nil {
		return nil, err
	}
	if p.duration > DurationOver60 {
		return nil, &bpierr.ParameterError{Field: "duration", Message: "unsupported duration range"}
	}
	values.Set("order", string(p.order))
	values.Set("duration", strconv.FormatUint(uint64(p.duration), 10))
	values.Set("tids", strconv.FormatUint(uint64(p.tid), 10))
	return values, nil
}

type SuggestParams struct{ term string }

func NewSuggestParams(term string) (SuggestParams, error) {
	term = strings.TrimSpace(term)
	if term == "" {
		return SuggestParams{}, &bpierr.ParameterError{Field: "term", Message: "search suggestion term cannot be blank"}
	}
	return SuggestParams{term: term}, nil
}
func (p SuggestParams) EncodeQuery() (url.Values, error) {
	if strings.TrimSpace(p.term) == "" {
		return nil, &bpierr.ParameterError{Field: "term", Message: "search suggestion term cannot be blank"}
	}
	return url.Values{"term": {p.term}}, nil
}
