package danmaku

import (
	"bytes"
	"compress/flate"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type XML struct {
	XMLName    xml.Name  `xml:"i"`
	ChatServer string    `xml:"chatserver"`
	ChatID     string    `xml:"chatid"`
	Mission    int32     `xml:"mission"`
	MaxLimit   int32     `xml:"maxlimit"`
	State      int32     `xml:"state"`
	RealName   int32     `xml:"real_name"`
	Source     string    `xml:"source"`
	Comments   []Comment `xml:"d"`
}

type Comment struct {
	Parameters string       `xml:"p,attr"`
	Content    string       `xml:",chardata"`
	Meta       *CommentMeta `xml:"-"`
}

type CommentMeta struct {
	Time       float32
	Mode       int32
	FontSize   int32
	Color      int32
	SentAt     int64
	Pool       int32
	UserHash   string
	ID         int64
	BlockLevel int32
}

// ParseDeflateXML decodes a raw-deflate danmaku response and parses each
// comment's compact p attribute into typed metadata.
func ParseDeflateXML(body []byte) (XML, error) {
	reader := flate.NewReader(bytes.NewReader(body))
	defer reader.Close()
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return XML{}, fmt.Errorf("danmaku: decompress XML: %w", err)
	}
	var document XML
	if err := xml.Unmarshal(decoded, &document); err != nil {
		return XML{}, fmt.Errorf("danmaku: decode XML: %w", err)
	}
	for i := range document.Comments {
		meta, err := parseCommentMeta(document.Comments[i].Parameters)
		if err != nil {
			return XML{}, fmt.Errorf("danmaku: decode comment %d metadata: %w", i, err)
		}
		document.Comments[i].Meta = &meta
	}
	return document, nil
}

func parseCommentMeta(value string) (CommentMeta, error) {
	parts := strings.Split(value, ",")
	if len(parts) < 9 {
		return CommentMeta{}, fmt.Errorf("expected 9 fields, got %d", len(parts))
	}
	parseInt := func(index int, bits int) (int64, error) {
		return strconv.ParseInt(parts[index], 10, bits)
	}
	timestamp, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return CommentMeta{}, err
	}
	mode, err := parseInt(1, 32)
	if err != nil {
		return CommentMeta{}, err
	}
	font, err := parseInt(2, 32)
	if err != nil {
		return CommentMeta{}, err
	}
	color, err := parseInt(3, 32)
	if err != nil {
		return CommentMeta{}, err
	}
	sentAt, err := parseInt(4, 64)
	if err != nil {
		return CommentMeta{}, err
	}
	pool, err := parseInt(5, 32)
	if err != nil {
		return CommentMeta{}, err
	}
	id, err := parseInt(7, 64)
	if err != nil {
		return CommentMeta{}, err
	}
	blockLevel, err := parseInt(8, 32)
	if err != nil {
		return CommentMeta{}, err
	}
	return CommentMeta{Time: float32(timestamp), Mode: int32(mode), FontSize: int32(font), Color: int32(color), SentAt: sentAt, Pool: int32(pool), UserHash: parts[6], ID: id, BlockLevel: int32(blockLevel)}, nil
}
