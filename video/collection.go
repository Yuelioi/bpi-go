package video

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// SeriesID identifies a user-created video series.
type SeriesID uint64

func NewSeriesID(value uint64) (SeriesID, error) {
	if value == 0 {
		return 0, &bpierr.ParameterError{Field: "series_id", Message: "value must be non-zero"}
	}
	return SeriesID(value), nil
}

func (id SeriesID) Validate() error {
	if id == 0 {
		return &bpierr.ParameterError{Field: "series_id", Message: "value must be non-zero"}
	}
	return nil
}

func (id SeriesID) String() string { return strconv.FormatUint(uint64(id), 10) }

type CollectionArchiveSort string

const (
	CollectionArchiveAscending  CollectionArchiveSort = "asc"
	CollectionArchiveDescending CollectionArchiveSort = "desc"
)

func (sort CollectionArchiveSort) validate() error {
	switch sort {
	case CollectionArchiveAscending, CollectionArchiveDescending:
		return nil
	default:
		return &bpierr.ParameterError{Field: "sort", Message: "value must be asc or desc"}
	}
}

type SeasonsArchivesParams struct {
	mid         ids.MID
	seasonID    ids.SeasonID
	sortReverse *bool
	page        uint64
	pageSize    uint64
}

func NewSeasonsArchivesParams(mid ids.MID, seasonID ids.SeasonID) SeasonsArchivesParams {
	return SeasonsArchivesParams{mid: mid, seasonID: seasonID}
}

func (p SeasonsArchivesParams) WithSortReverse(reverse bool) SeasonsArchivesParams {
	p.sortReverse = &reverse
	return p
}

func (p SeasonsArchivesParams) WithPage(page uint64) (SeasonsArchivesParams, error) {
	if page == 0 {
		return SeasonsArchivesParams{}, &bpierr.ParameterError{Field: "page_num", Message: "value must be non-zero"}
	}
	p.page = page
	return p, nil
}

func (p SeasonsArchivesParams) WithPageSize(size uint64) (SeasonsArchivesParams, error) {
	if size == 0 {
		return SeasonsArchivesParams{}, &bpierr.ParameterError{Field: "page_size", Message: "value must be non-zero"}
	}
	p.pageSize = size
	return p, nil
}

func (p SeasonsArchivesParams) EncodeQuery() (url.Values, error) {
	if err := p.mid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "mid", Message: "member ID must be non-zero"}
	}
	if err := p.seasonID.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "season_id", Message: "season ID must be non-zero"}
	}
	page, pageSize := p.page, p.pageSize
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	values := url.Values{
		"mid":       {p.mid.String()},
		"season_id": {p.seasonID.String()},
		"page_num":  {strconv.FormatUint(page, 10)},
		"page_size": {strconv.FormatUint(pageSize, 10)},
	}
	if p.sortReverse != nil {
		values.Set("sort_reverse", strconv.FormatBool(*p.sortReverse))
	}
	return values, nil
}

type HomeSeasonsSeriesParams struct {
	mid      ids.MID
	page     uint64
	pageSize uint64
}

func NewHomeSeasonsSeriesParams(mid ids.MID) HomeSeasonsSeriesParams {
	return HomeSeasonsSeriesParams{mid: mid}
}

func (p HomeSeasonsSeriesParams) WithPage(page uint64) (HomeSeasonsSeriesParams, error) {
	if page == 0 {
		return HomeSeasonsSeriesParams{}, &bpierr.ParameterError{Field: "page_num", Message: "value must be non-zero"}
	}
	p.page = page
	return p, nil
}

func (p HomeSeasonsSeriesParams) WithPageSize(size uint64) (HomeSeasonsSeriesParams, error) {
	if size == 0 {
		return HomeSeasonsSeriesParams{}, &bpierr.ParameterError{Field: "page_size", Message: "value must be non-zero"}
	}
	p.pageSize = size
	return p, nil
}

func (p HomeSeasonsSeriesParams) EncodeQuery() (url.Values, error) {
	if err := p.mid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "mid", Message: "member ID must be non-zero"}
	}
	page, pageSize := p.page, p.pageSize
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return url.Values{
		"mid":       {p.mid.String()},
		"page_num":  {strconv.FormatUint(page, 10)},
		"page_size": {strconv.FormatUint(pageSize, 10)},
	}, nil
}

type SeasonsSeriesParams struct {
	mid      ids.MID
	page     uint64
	pageSize uint64
}

func NewSeasonsSeriesParams(mid ids.MID) SeasonsSeriesParams {
	return SeasonsSeriesParams{mid: mid}
}

func (p SeasonsSeriesParams) WithPage(page uint64) (SeasonsSeriesParams, error) {
	if page == 0 {
		return SeasonsSeriesParams{}, &bpierr.ParameterError{Field: "page_num", Message: "value must be non-zero"}
	}
	p.page = page
	return p, nil
}

func (p SeasonsSeriesParams) WithPageSize(size uint64) (SeasonsSeriesParams, error) {
	if size == 0 {
		return SeasonsSeriesParams{}, &bpierr.ParameterError{Field: "page_size", Message: "value must be non-zero"}
	}
	p.pageSize = size
	return p, nil
}

func (p SeasonsSeriesParams) EncodeQuery() (url.Values, error) {
	if err := p.mid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "mid", Message: "member ID must be non-zero"}
	}
	values := url.Values{"mid": {p.mid.String()}}
	if p.page != 0 {
		values.Set("page_num", strconv.FormatUint(p.page, 10))
	}
	if p.pageSize != 0 {
		values.Set("page_size", strconv.FormatUint(p.pageSize, 10))
	}
	return values, nil
}

type SeriesInfoParams struct{ seriesID SeriesID }

func NewSeriesInfoParams(seriesID SeriesID) SeriesInfoParams {
	return SeriesInfoParams{seriesID: seriesID}
}

func (p SeriesInfoParams) EncodeQuery() (url.Values, error) {
	if err := p.seriesID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"series_id": {p.seriesID.String()}}, nil
}

type SeriesArchivesParams struct {
	mid        ids.MID
	seriesID   SeriesID
	onlyNormal *bool
	sort       CollectionArchiveSort
	page       uint64
	pageSize   uint64
}

func NewSeriesArchivesParams(mid ids.MID, seriesID SeriesID) SeriesArchivesParams {
	return SeriesArchivesParams{mid: mid, seriesID: seriesID}
}

func (p SeriesArchivesParams) WithOnlyNormal(onlyNormal bool) SeriesArchivesParams {
	p.onlyNormal = &onlyNormal
	return p
}

func (p SeriesArchivesParams) WithSort(sort CollectionArchiveSort) (SeriesArchivesParams, error) {
	if err := sort.validate(); err != nil {
		return SeriesArchivesParams{}, err
	}
	p.sort = sort
	return p, nil
}

func (p SeriesArchivesParams) WithPage(page uint64) (SeriesArchivesParams, error) {
	if page == 0 {
		return SeriesArchivesParams{}, &bpierr.ParameterError{Field: "pn", Message: "value must be non-zero"}
	}
	p.page = page
	return p, nil
}

func (p SeriesArchivesParams) WithPageSize(size uint64) (SeriesArchivesParams, error) {
	if size == 0 {
		return SeriesArchivesParams{}, &bpierr.ParameterError{Field: "ps", Message: "value must be non-zero"}
	}
	p.pageSize = size
	return p, nil
}

func (p SeriesArchivesParams) EncodeQuery() (url.Values, error) {
	if err := p.mid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "mid", Message: "member ID must be non-zero"}
	}
	if err := p.seriesID.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{"mid": {p.mid.String()}, "series_id": {p.seriesID.String()}}
	if p.onlyNormal != nil {
		values.Set("only_normal", strconv.FormatBool(*p.onlyNormal))
	}
	if p.sort != "" {
		if err := p.sort.validate(); err != nil {
			return nil, err
		}
		values.Set("sort", string(p.sort))
	}
	if p.page != 0 {
		values.Set("pn", strconv.FormatUint(p.page, 10))
	}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(p.pageSize, 10))
	}
	return values, nil
}

type CollectionArchiveStat struct {
	View uint64  `json:"view"`
	VT   *uint64 `json:"vt"`
}

type CollectionArchive struct {
	AID              ids.AID               `json:"aid"`
	BVID             ids.BVID              `json:"bvid"`
	CreateTime       uint64                `json:"ctime"`
	Duration         uint64                `json:"duration"`
	InteractiveVideo bool                  `json:"interactive_video"`
	Picture          string                `json:"pic"`
	PlaybackPosition uint64                `json:"playback_position"`
	PublishTime      uint64                `json:"pubdate"`
	Stat             CollectionArchiveStat `json:"stat"`
	State            uint64                `json:"state"`
	Title            string                `json:"title"`
	UGCPay           uint64                `json:"ugc_pay"`
	VTDisplay        string                `json:"vt_display"`
	IsLessonVideo    *uint32               `json:"is_lesson_video"`
}

// CollectionPage supports both page_num/page_size and num/size response keys.
type CollectionPage struct {
	Page  uint64
	Size  uint64
	Total uint64
}

func (p *CollectionPage) UnmarshalJSON(data []byte) error {
	var wire struct {
		PageNum  *uint64 `json:"page_num"`
		PageSize *uint64 `json:"page_size"`
		Num      *uint64 `json:"num"`
		Size     *uint64 `json:"size"`
		Total    uint64  `json:"total"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if wire.PageNum != nil {
		p.Page = *wire.PageNum
	} else if wire.Num != nil {
		p.Page = *wire.Num
	}
	if wire.PageSize != nil {
		p.Size = *wire.PageSize
	} else if wire.Size != nil {
		p.Size = *wire.Size
	}
	p.Total = wire.Total
	return nil
}

type SeasonArchivesMeta struct {
	Category    uint64       `json:"category"`
	Cover       string       `json:"cover"`
	Description string       `json:"description"`
	MID         ids.MID      `json:"mid"`
	Name        string       `json:"name"`
	PublishTime uint64       `json:"ptime"`
	SeasonID    ids.SeasonID `json:"season_id"`
	Total       uint64       `json:"total"`
}

type SeasonsArchives struct {
	AIDs     []uint64            `json:"aids"`
	Archives []CollectionArchive `json:"archives"`
	Meta     SeasonArchivesMeta  `json:"meta"`
	Page     CollectionPage      `json:"page"`
}

type SeasonMeta struct {
	Category    uint64       `json:"category"`
	Cover       string       `json:"cover"`
	Description string       `json:"description"`
	MID         ids.MID      `json:"mid"`
	Name        string       `json:"name"`
	PublishTime uint64       `json:"ptime"`
	SeasonID    ids.SeasonID `json:"season_id"`
	Total       uint64       `json:"total"`
}

type SeriesMeta struct {
	Category       uint64   `json:"category"`
	Creator        string   `json:"creator"`
	CreateTime     uint64   `json:"ctime"`
	Description    string   `json:"description"`
	Keywords       []string `json:"keywords"`
	LastUpdateTime uint64   `json:"last_update_ts"`
	MID            ids.MID  `json:"mid"`
	ModifyTime     uint64   `json:"mtime"`
	Name           string   `json:"name"`
	RawKeywords    string   `json:"raw_keywords"`
	SeriesID       SeriesID `json:"series_id"`
	State          uint64   `json:"state"`
	Total          uint64   `json:"total"`
	Cover          *string  `json:"cover"`
}

type SeasonItem struct {
	Archives   []CollectionArchive `json:"archives"`
	Meta       SeasonMeta          `json:"meta"`
	RecentAIDs []uint64            `json:"recent_aids"`
}

type SeriesItem struct {
	Archives   []CollectionArchive `json:"archives"`
	Meta       SeriesMeta          `json:"meta"`
	RecentAIDs []uint64            `json:"recent_aids"`
}

type SeasonsSeriesItems struct {
	Page    CollectionPage `json:"page"`
	Seasons []SeasonItem   `json:"seasons_list"`
	Series  []SeriesItem   `json:"series_list"`
}

type SeasonsSeries struct {
	Items SeasonsSeriesItems `json:"items_lists"`
}

type SeriesInfo struct {
	Meta       SeriesMeta `json:"meta"`
	RecentAIDs []uint64   `json:"recent_aids"`
}

type SeriesArchives struct {
	AIDs     []uint64            `json:"aids"`
	Page     CollectionPage      `json:"page"`
	Archives []CollectionArchive `json:"archives"`
}
