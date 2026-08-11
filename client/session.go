package client

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	cookieDedeUserID = "DedeUserID"
	cookieSESSDATA   = "SESSDATA"
	cookieBiliJCT    = "bili_jct"
	cookieBuvid3     = "buvid3"
)

type sessionState struct {
	mu      sync.RWMutex
	cookies map[string]string
	account *Account
}

func newSession() *sessionState {
	return &sessionState{cookies: make(map[string]string)}
}

func (s *sessionState) replaceCookieHeader(header string) error {
	cookies, err := parseCookieHeader(header)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cookies = cookies
	s.account = accountFromCookies(cookies)
	return nil
}

func (s *sessionState) replaceAccount(account Account) error {
	if err := account.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cookies = map[string]string{
		cookieDedeUserID: account.DedeUserID,
		cookieSESSDATA:   account.SESSDATA,
		cookieBiliJCT:    account.BiliJCT,
		cookieBuvid3:     account.Buvid3,
	}
	clone := account
	s.account = &clone
	return nil
}

func (s *sessionState) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.cookies)
	s.account = nil
}

func (s *sessionState) addToRequest(request *http.Request) {
	if request.URL == nil || !isBilibiliHost(request.URL.Hostname()) {
		return
	}
	s.mu.RLock()
	names := make([]string, 0, len(s.cookies))
	for name := range s.cookies {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		request.AddCookie(&http.Cookie{Name: name, Value: s.cookies[name]})
	}
	s.mu.RUnlock()
}

func (s *sessionState) absorbResponse(host string, cookies []*http.Cookie, now time.Time) {
	if !isBilibiliHost(host) || len(cookies) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, cookie := range cookies {
		if cookie == nil || cookie.Name == "" {
			continue
		}
		if cookie.MaxAge < 0 || (!cookie.Expires.IsZero() && !cookie.Expires.After(now)) {
			delete(s.cookies, cookie.Name)
			continue
		}
		s.cookies[cookie.Name] = cookie.Value
	}
	s.account = accountFromCookies(s.cookies)
}

func (s *sessionState) hasLoginCookies() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cookies[cookieSESSDATA] != ""
}

func (s *sessionState) csrf() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	csrf := s.cookies[cookieBiliJCT]
	if csrf == "" {
		return "", ErrAuthenticationRequired
	}
	return csrf, nil
}

func (s *sessionState) getAccount() (Account, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.account == nil {
		return Account{}, false
	}
	return *s.account, true
}

func parseCookieHeader(header string) (map[string]string, error) {
	cookies := make(map[string]string)
	for _, segment := range strings.Split(header, ";") {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}
		name, value, ok := strings.Cut(segment, "=")
		name = strings.TrimSpace(name)
		if !ok || name == "" {
			return nil, &ParameterError{Field: "cookie", Message: "each cookie segment must contain a non-empty name and '='"}
		}
		value = strings.TrimSpace(value)
		if err := (&http.Cookie{Name: name, Value: value}).Valid(); err != nil {
			return nil, &ParameterError{Field: "cookie", Message: fmt.Sprintf("invalid cookie %q: %v", name, err)}
		}
		cookies[name] = value
	}
	if len(cookies) == 0 {
		return nil, &ParameterError{Field: "cookie", Message: "cookie cannot be empty"}
	}
	return cookies, nil
}

func accountFromCookies(cookies map[string]string) *Account {
	account := Account{
		DedeUserID: cookies[cookieDedeUserID],
		SESSDATA:   cookies[cookieSESSDATA],
		BiliJCT:    cookies[cookieBiliJCT],
		Buvid3:     cookies[cookieBuvid3],
	}
	if account.Validate() != nil {
		return nil
	}
	return &account
}

func isBilibiliHost(host string) bool {
	host = strings.TrimSuffix(strings.ToLower(host), ".")
	return host == "bilibili.com" || strings.HasSuffix(host, ".bilibili.com")
}
