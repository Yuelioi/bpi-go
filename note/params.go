// Package note contains validated parameters and stable response models for
// Bilibili video notes.
package note

import (
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type IsForbidParams struct{ aid ids.AID }

func NewIsForbidParams(aid ids.AID) IsForbidParams { return IsForbidParams{aid: aid} }

func (p IsForbidParams) EncodeQuery() (url.Values, error) {
	if err := p.aid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "aid", Message: "video ID is invalid"}
	}
	return url.Values{"aid": {p.aid.String()}}, nil
}

type PrivateInfoParams struct {
	aid    ids.AID
	noteID ids.NoteID
}

func NewPrivateInfoParams(aid ids.AID, noteID ids.NoteID) PrivateInfoParams {
	return PrivateInfoParams{aid: aid, noteID: noteID}
}

func (p PrivateInfoParams) EncodeQuery() (url.Values, error) {
	if err := p.aid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "oid", Message: "video ID is invalid"}
	}
	if err := p.noteID.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "note_id", Message: "note ID is invalid"}
	}
	return url.Values{"oid": {p.aid.String()}, "oid_type": {"0"}, "note_id": {p.noteID.String()}}, nil
}

type PublicInfoParams struct{ cvid ids.CVID }

func NewPublicInfoParams(cvid ids.CVID) PublicInfoParams { return PublicInfoParams{cvid: cvid} }

func (p PublicInfoParams) EncodeQuery() (url.Values, error) {
	if err := p.cvid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "cvid", Message: "public note ID is invalid"}
	}
	return url.Values{"cvid": {p.cvid.String()}}, nil
}

type ArchiveListParams struct{ aid ids.AID }

func NewArchiveListParams(aid ids.AID) ArchiveListParams { return ArchiveListParams{aid: aid} }

func (p ArchiveListParams) EncodeQuery() (url.Values, error) {
	if err := p.aid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "oid", Message: "video ID is invalid"}
	}
	return url.Values{"oid": {p.aid.String()}, "oid_type": {"0"}}, nil
}

type Pagination struct {
	page uint32
	size uint32
}

func NewPagination() Pagination { return Pagination{page: 1, size: 10} }

func (p Pagination) WithPage(page uint32) (Pagination, error) {
	if page == 0 {
		return Pagination{}, &bpierr.ParameterError{Field: "pn", Message: "page must be greater than zero"}
	}
	p.page = page
	return p, nil
}

func (p Pagination) WithPageSize(size uint32) (Pagination, error) {
	if size == 0 {
		return Pagination{}, &bpierr.ParameterError{Field: "ps", Message: "page size must be greater than zero"}
	}
	p.size = size
	return p, nil
}

func (p Pagination) EncodeQuery() (url.Values, error) {
	if p.page == 0 || p.size == 0 {
		return nil, &bpierr.ParameterError{Field: "pagination", Message: "use NewPagination to initialize defaults"}
	}
	return url.Values{"pn": {strconv.FormatUint(uint64(p.page), 10)}, "ps": {strconv.FormatUint(uint64(p.size), 10)}}, nil
}

type PublicArchiveListParams struct {
	aid        ids.AID
	pagination Pagination
}

func NewPublicArchiveListParams(aid ids.AID) PublicArchiveListParams {
	return PublicArchiveListParams{aid: aid, pagination: NewPagination()}
}

func (p PublicArchiveListParams) WithPage(page uint32) (PublicArchiveListParams, error) {
	pagination, err := p.pagination.WithPage(page)
	if err != nil {
		return PublicArchiveListParams{}, err
	}
	p.pagination = pagination
	return p, nil
}

func (p PublicArchiveListParams) WithPageSize(size uint32) (PublicArchiveListParams, error) {
	pagination, err := p.pagination.WithPageSize(size)
	if err != nil {
		return PublicArchiveListParams{}, err
	}
	p.pagination = pagination
	return p, nil
}

func (p PublicArchiveListParams) EncodeQuery() (url.Values, error) {
	if err := p.aid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "oid", Message: "video ID is invalid"}
	}
	values, err := p.pagination.EncodeQuery()
	if err != nil {
		return nil, err
	}
	values.Set("oid", p.aid.String())
	values.Set("oid_type", "0")
	return values, nil
}
