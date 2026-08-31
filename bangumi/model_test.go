package bangumi_test

import (
	"encoding/json"
	"testing"

	"github.com/Yuelioi/bpi-go/bangumi"
)

func TestDetailAcceptsNegativeUnknownEpisodeTotal(t *testing.T) {
	t.Parallel()

	var detail bangumi.Detail
	if err := json.Unmarshal([]byte(`{"total":-1}`), &detail); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if detail.Total != -1 {
		t.Fatalf("Total = %d, want -1 sentinel", detail.Total)
	}
}

func TestStatAcceptsCompactSectionAliases(t *testing.T) {
	t.Parallel()

	var stat bangumi.Stat
	if err := json.Unmarshal([]byte(`{"coin":1,"danmakus":2,"likes":3,"play":4,"reply":5,"vt":6}`), &stat); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if stat.Coins != 1 || stat.Views != 4 {
		t.Fatalf("Stat = %+v, want coin/play aliases", stat)
	}
	if stat.Favorite != 0 || stat.FollowText != "" {
		t.Fatalf("compact Stat defaults = %+v, want zero values", stat)
	}
}

func TestBangumiModelsAcceptUnreleasedAndCompactShapes(t *testing.T) {
	t.Parallel()

	var media bangumi.Media
	if err := json.Unmarshal([]byte(`{"media_id":1,"new_ep":{"id":1},"season_id":1}`), &media); err != nil {
		t.Fatalf("Unmarshal(media without rating) error = %v", err)
	}
	if media.Rating.Count != 0 || media.Rating.Score != 0 {
		t.Fatalf("Rating = %+v, want zero value", media.Rating)
	}

	var sections bangumi.Sections
	if err := json.Unmarshal([]byte(`{"section":[]}`), &sections); err != nil {
		t.Fatalf("Unmarshal(sections without main_section) error = %v", err)
	}
	if len(sections.Main.Episodes) != 0 || len(sections.Sections) != 0 {
		t.Fatalf("Sections = %+v, want empty state", sections)
	}

	var episode bangumi.SectionEpisode
	if err := json.Unmarshal([]byte(`{"aid":1,"cid":2,"id":3,"title":"1"}`), &episode); err != nil {
		t.Fatalf("Unmarshal(compact section episode) error = %v", err)
	}
	if episode.AID != 1 || episode.CID != 2 || episode.ID != 3 {
		t.Fatalf("SectionEpisode = %+v, want compact identity", episode)
	}
}

func TestBangumiPlayURLUsesCurrentQualityDURL(t *testing.T) {
	t.Parallel()

	var play bangumi.PlayURL
	if err := json.Unmarshal([]byte(`{
		"durl":[{
			"order":1,
			"length":1000,
			"size":100,
			"ahead":"",
			"vhead":"",
			"url":"https://example.invalid/current.flv",
			"backup_url":[]
		}],
		"durls":[{
			"quality":32,
			"durl":[{"size":999}]
		}]
	}`), &play); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(play.DURLs) != 1 || play.DURLs[0].Size != 100 {
		t.Fatalf("DURLs = %+v, want current-quality durl", play.DURLs)
	}
}
