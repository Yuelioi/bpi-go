package sign

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
)

var mixinKeyTable = [...]int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35,
	27, 43, 5, 49, 33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13,
	37, 48, 7, 16, 24, 55, 40, 61, 26, 17, 0, 1, 60, 51, 30, 4,
	22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11, 36, 20, 34, 44, 52,
}

// WBIKeys contains the image and sub-resource key stems used for WBI signing.
type WBIKeys struct {
	image string
	sub   string
}

// NewWBIKeys validates key stems.
func NewWBIKeys(image, sub string) (WBIKeys, error) {
	if image == "" {
		return WBIKeys{}, &InvalidError{Field: "img_key", Message: "key cannot be empty"}
	}
	if sub == "" {
		return WBIKeys{}, &InvalidError{Field: "sub_key", Message: "key cannot be empty"}
	}
	return WBIKeys{image: image, sub: sub}, nil
}

// WBIKeysFromURLs extracts key stems from navigation resource URLs.
func WBIKeysFromURLs(imageURL, subURL string) (WBIKeys, error) {
	image, err := keyStem(imageURL)
	if err != nil {
		return WBIKeys{}, &InvalidError{Field: "img_url", Message: err.Error()}
	}
	sub, err := keyStem(subURL)
	if err != nil {
		return WBIKeys{}, &InvalidError{Field: "sub_url", Message: err.Error()}
	}
	return NewWBIKeys(image, sub)
}

// Image returns the image key stem.
func (k WBIKeys) Image() string { return k.image }

// Sub returns the sub-resource key stem.
func (k WBIKeys) Sub() string { return k.sub }

// MixinKey derives the 32-byte WBI mixin key.
func MixinKey(keys WBIKeys) (string, error) {
	if _, err := NewWBIKeys(keys.image, keys.sub); err != nil {
		return "", err
	}
	combined := []byte(keys.image + keys.sub)
	result := make([]byte, 0, 32)
	for _, index := range mixinKeyTable {
		if index < len(combined) {
			result = append(result, combined[index])
			if len(result) == 32 {
				break
			}
		}
	}
	if len(result) != 32 {
		return "", &InvalidError{Field: "wbi_keys", Message: "combined keys must produce a 32-byte mixin key"}
	}
	return string(result), nil
}

// SignWBIAt returns a signed copy of params for the supplied Unix timestamp.
func SignWBIAt(params map[string]string, keys WBIKeys, timestamp uint64) (map[string]string, error) {
	mixin, err := MixinKey(keys)
	if err != nil {
		return nil, err
	}
	signed := make(map[string]string, len(params)+2)
	for key, value := range params {
		signed[key] = filterWBIValue(value)
	}
	signed["wts"] = strconv.FormatUint(timestamp, 10)
	query := EncodeWBIQuery(signed)
	digest := md5.Sum([]byte(query + mixin))
	signed["w_rid"] = hex.EncodeToString(digest[:])
	return signed, nil
}

// EncodeWBIQuery serializes parameters using the RFC 3986 encoding required
// by WBI rather than application/x-www-form-urlencoded space encoding.
func EncodeWBIQuery(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var builder strings.Builder
	for index, key := range keys {
		if index > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(rfc3986Escape(key))
		builder.WriteByte('=')
		builder.WriteString(rfc3986Escape(params[key]))
	}
	return builder.String()
}

func keyStem(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	name := path.Base(parsed.Path)
	if name == "." || name == "/" || name == "" {
		return "", &InvalidError{Field: "wbi_url", Message: "URL has no filename"}
	}
	extension := path.Ext(name)
	stem := strings.TrimSuffix(name, extension)
	if extension == "" || stem == "" {
		return "", &InvalidError{Field: "wbi_url", Message: "filename must contain a stem and extension"}
	}
	return stem, nil
}

func filterWBIValue(value string) string {
	return strings.Map(func(character rune) rune {
		if strings.ContainsRune("!'()*", character) {
			return -1
		}
		return character
	}, value)
}

func rfc3986Escape(value string) string {
	return strings.ReplaceAll(url.QueryEscape(value), "+", "%20")
}
