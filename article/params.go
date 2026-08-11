// Package article contains validated parameters and stable response models
// for Bilibili article endpoints.
package article

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type InfoParams struct{ id ids.CVID }

func NewInfoParams(id ids.CVID) InfoParams { return InfoParams{id: id} }

func (p InfoParams) EncodeQuery() (url.Values, error) {
	if err := p.id.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "id", Message: "article ID is invalid"}
	}
	return url.Values{"id": {p.id.String()}}, nil
}

type ViewParams struct {
	id         ids.CVID
	gaiaSource string
}

func NewViewParams(id ids.CVID) ViewParams {
	return ViewParams{id: id, gaiaSource: "main_web"}
}

func (p ViewParams) WithGaiaSource(source string) (ViewParams, error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return ViewParams{}, &bpierr.ParameterError{Field: "gaia_source", Message: "source cannot be blank"}
	}
	p.gaiaSource = source
	return p, nil
}

func (p ViewParams) EncodeQuery() (url.Values, error) {
	if err := p.id.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "id", Message: "article ID is invalid"}
	}
	if strings.TrimSpace(p.gaiaSource) == "" {
		return nil, &bpierr.ParameterError{Field: "gaia_source", Message: "source cannot be blank"}
	}
	return url.Values{"id": {p.id.String()}, "gaia_source": {p.gaiaSource}}, nil
}

type CardsParams struct {
	ids         string
	webLocation string
}

func NewCardsParams(cardIDs string) (CardsParams, error) {
	cardIDs = strings.TrimSpace(cardIDs)
	if cardIDs == "" {
		return CardsParams{}, &bpierr.ParameterError{Field: "ids", Message: "card IDs cannot be blank"}
	}
	return CardsParams{ids: cardIDs, webLocation: "333.1305"}, nil
}

func (p CardsParams) WithWebLocation(location string) (CardsParams, error) {
	location = strings.TrimSpace(location)
	if location == "" {
		return CardsParams{}, &bpierr.ParameterError{Field: "web_location", Message: "location cannot be blank"}
	}
	p.webLocation = location
	return p, nil
}

func (p CardsParams) EncodeQuery() (url.Values, error) {
	if strings.TrimSpace(p.ids) == "" {
		return nil, &bpierr.ParameterError{Field: "ids", Message: "card IDs cannot be blank"}
	}
	if strings.TrimSpace(p.webLocation) == "" {
		return nil, &bpierr.ParameterError{Field: "web_location", Message: "location cannot be blank"}
	}
	return url.Values{"ids": {p.ids}, "web_location": {p.webLocation}}, nil
}

type ArticlesInfoParams struct{ id uint64 }

func NewArticlesInfoParams(id uint64) (ArticlesInfoParams, error) {
	if id == 0 {
		return ArticlesInfoParams{}, &bpierr.ParameterError{Field: "id", Message: "list ID must be non-zero"}
	}
	return ArticlesInfoParams{id: id}, nil
}

func (p ArticlesInfoParams) EncodeQuery() (url.Values, error) {
	if p.id == 0 {
		return nil, &bpierr.ParameterError{Field: "id", Message: "list ID must be non-zero"}
	}
	return url.Values{"id": {strconv.FormatUint(p.id, 10)}}, nil
}
