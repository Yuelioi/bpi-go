package video_test

import (
	"encoding/json"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/video"
)

func TestPlayURLParamsMatchPromotedUnsignedQuery(t *testing.T) {
	t.Parallel()

	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	cid, _ := ids.NewCID(62131)
	params := video.PlayURLByBVID(bvid, cid).
		WithQuality(32).
		WithFormatFlags(16).
		WithFormatVersion(0)
	query, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Encode() != "bvid=BV1xx411c7mD&cid=62131&fnval=16&fnver=0&platform=pc&qn=32" {
		t.Fatalf("query = %q, want promoted unsigned query", query.Encode())
	}
}

func TestPlayURLParamsUseAvidAndValidateCID(t *testing.T) {
	t.Parallel()

	aid, _ := ids.NewAID(170001)
	query, err := video.PlayURLByAID(aid, ids.CID(1)).EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Get("avid") != "170001" || query.Get("aid") != "" {
		t.Fatalf("query = %v, want avid", query)
	}
	if _, err := video.PlayURLByAID(aid, ids.CID(0)).EncodeQuery(); err == nil {
		t.Fatal("EncodeQuery(zero CID) error = nil")
	}
}

func TestPlayURLAcceptsNegativeResumeSentinels(t *testing.T) {
	t.Parallel()

	var payload video.PlayURL
	if err := json.Unmarshal([]byte(`{"last_play_time":-1000,"last_play_cid":-1000}`), &payload); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if payload.LastPlayTime != -1000 || payload.LastPlayCID != -1000 {
		t.Fatalf("resume sentinels = %d/%d", payload.LastPlayTime, payload.LastPlayCID)
	}
}

func TestDASHFLACAcceptsNullableAndSingleStreamAudio(t *testing.T) {
	t.Parallel()

	var unavailable video.DASHFLAC
	if err := json.Unmarshal([]byte(`{"display":true,"audio":null}`), &unavailable); err != nil {
		t.Fatalf("Unmarshal(null FLAC audio) error = %v", err)
	}
	if unavailable.Display == nil || !*unavailable.Display || unavailable.Audio != nil {
		t.Fatalf("null FLAC = %+v, want display=true and no audio", unavailable)
	}

	var available video.DASHFLAC
	if err := json.Unmarshal([]byte(`{
		"display":true,
		"audio":{
			"id":30251,
			"baseUrl":"https://example.invalid/flac.m4s",
			"backupUrl":[],
			"bandwidth":1000,
			"mimeType":"audio/mp4",
			"codecs":"fLaC"
		}
	}`), &available); err != nil {
		t.Fatalf("Unmarshal(single FLAC stream) error = %v", err)
	}
	if available.Audio == nil || available.Audio.ID != 30251 {
		t.Fatalf("single FLAC audio = %+v, want stream 30251", available.Audio)
	}
}

func TestPlayURLModelsAcceptNullBackupURLs(t *testing.T) {
	t.Parallel()

	var dash video.DASH
	if err := json.Unmarshal([]byte(`{
		"video":[{
			"id":64,
			"baseUrl":"https://example.invalid/video.m4s",
			"backupUrl":null,
			"bandwidth":1000,
			"mimeType":"video/mp4",
			"codecs":"avc1.640028"
		}],
		"audio":[],
		"dolby":null,
		"flac":null,
		"duration":60
	}`), &dash); err != nil {
		t.Fatalf("Unmarshal(DASH null backupUrl) error = %v", err)
	}
	if len(dash.Video) != 1 || len(dash.Video[0].BackupURLs) != 0 {
		t.Fatalf("DASH backup URLs = %+v, want empty", dash.Video)
	}

	var durl video.DURL
	if err := json.Unmarshal([]byte(`{
		"order":1,
		"length":60000,
		"size":1000000,
		"ahead":"",
		"vhead":"",
		"url":"https://example.invalid/video.mp4",
		"backup_url":null
	}`), &durl); err != nil {
		t.Fatalf("Unmarshal(DURL null backup_url) error = %v", err)
	}
	if len(durl.BackupURLs) != 0 {
		t.Fatalf("DURL backup URLs = %+v, want empty", durl.BackupURLs)
	}
}
