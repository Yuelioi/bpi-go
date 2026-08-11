package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"sort"
	"strings"
)

// multipartFile describes one streamed file part for an internal domain
// request. The caller retains ownership of Content.
type multipartFile struct {
	FieldName   string
	FileName    string
	ContentType string
	Content     io.Reader
}

// NewQueryRequest builds a context-bound request and merges query values with
// any query already present in endpoint.
func NewQueryRequest(ctx context.Context, method, endpoint string, query url.Values) (*http.Request, error) {
	if ctx == nil {
		return nil, &ParameterError{Field: "context", Message: "context cannot be nil"}
	}
	request, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, &ParameterError{Field: "endpoint", Message: "endpoint or HTTP method is invalid"}
	}
	request.URL.RawQuery = mergeURLValues(request.URL.Query(), query).Encode()
	return request, nil
}

// NewFormRequest builds a form-encoded POST request.
func NewFormRequest(ctx context.Context, endpoint string, query, form url.Values) (*http.Request, error) {
	request, err := newRequestWithBody(ctx, http.MethodPost, endpoint, query, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}

// NewJSONRequest builds a JSON-encoded POST request.
func NewJSONRequest(ctx context.Context, endpoint string, query url.Values, value any) (*http.Request, error) {
	body, err := json.Marshal(value)
	if err != nil {
		return nil, &ParameterError{Field: "json_body", Message: "request body cannot be encoded as JSON"}
	}
	request, err := newRequestWithBody(ctx, http.MethodPost, endpoint, query, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

// newMultipartRequest streams multipart fields and files through a pipe. It
// is intended for immediate execution by a domain method so uploads are not
// buffered in memory.
func newMultipartRequest(
	ctx context.Context,
	endpoint string,
	query, fields url.Values,
	files []multipartFile,
) (*http.Request, error) {
	if len(fields) == 0 && len(files) == 0 {
		return nil, &ParameterError{Field: "multipart", Message: "multipart body cannot be empty"}
	}
	for _, file := range files {
		if strings.TrimSpace(file.FieldName) == "" {
			return nil, &ParameterError{Field: "multipart_field", Message: "file field name cannot be empty"}
		}
		if strings.TrimSpace(file.FileName) == "" {
			return nil, &ParameterError{Field: "multipart_filename", Message: "file name cannot be empty"}
		}
		if file.Content == nil {
			return nil, &ParameterError{Field: "multipart_content", Message: "file content cannot be nil"}
		}
		if file.ContentType != "" {
			if _, _, err := mime.ParseMediaType(file.ContentType); err != nil {
				return nil, &ParameterError{Field: "multipart_content_type", Message: "content type is invalid"}
			}
		}
	}

	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)
	request, err := newRequestWithBody(ctx, http.MethodPost, endpoint, query, reader)
	if err != nil {
		_ = reader.Close()
		_ = writer.Close()
		return nil, err
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	go func() {
		err := writeMultipartBody(multipartWriter, fields, files)
		if closeErr := multipartWriter.Close(); err == nil {
			err = closeErr
		}
		_ = writer.CloseWithError(err)
	}()
	return request, nil
}

func newRequestWithBody(
	ctx context.Context,
	method, endpoint string,
	query url.Values,
	body io.Reader,
) (*http.Request, error) {
	if ctx == nil {
		return nil, &ParameterError{Field: "context", Message: "context cannot be nil"}
	}
	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, &ParameterError{Field: "endpoint", Message: "endpoint or HTTP method is invalid"}
	}
	request.URL.RawQuery = mergeURLValues(request.URL.Query(), query).Encode()
	return request, nil
}

func writeMultipartBody(writer *multipart.Writer, fields url.Values, files []multipartFile) error {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		for _, value := range fields[key] {
			if err := writer.WriteField(key, value); err != nil {
				return err
			}
		}
	}
	for _, file := range files {
		contentType := file.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		disposition := mime.FormatMediaType("form-data", map[string]string{
			"name":     file.FieldName,
			"filename": file.FileName,
		})
		if disposition == "" || strings.ContainsAny(contentType, "\r\n") {
			return &ParameterError{Field: "multipart_header", Message: "multipart header is invalid"}
		}
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", disposition)
		header.Set("Content-Type", contentType)
		part, err := writer.CreatePart(header)
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, file.Content); err != nil {
			return fmt.Errorf("bpi: stream multipart file: %w", err)
		}
	}
	return nil
}

func (c *Client) csrfValues(values url.Values, fieldNames ...string) (url.Values, error) {
	csrf, err := c.CSRF()
	if err != nil {
		return nil, err
	}
	if len(fieldNames) == 0 {
		fieldNames = []string{"csrf"}
	}
	result := cloneURLValues(values)
	for _, field := range fieldNames {
		if strings.TrimSpace(field) == "" {
			return nil, &ParameterError{Field: "csrf_field", Message: "field name cannot be empty"}
		}
		result.Set(field, csrf)
	}
	return result, nil
}

// WBIValues returns a signed copy of values using the client's cached WBI
// keys and clock. The input values are never mutated.
func (c *Client) WBIValues(ctx context.Context, values url.Values) (url.Values, error) {
	params := make(map[string]string, len(values))
	for key, candidates := range values {
		if len(candidates) != 1 {
			return nil, &ParameterError{Field: key, Message: "WBI parameters require exactly one value per key"}
		}
		params[key] = candidates[0]
	}
	signed, err := c.signWBIParams(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make(url.Values, len(signed))
	for key, value := range signed {
		result.Set(key, value)
	}
	return result, nil
}

func mergeURLValues(base, added url.Values) url.Values {
	result := cloneURLValues(base)
	for key, values := range added {
		for _, value := range values {
			result.Add(key, value)
		}
	}
	return result
}

func cloneURLValues(values url.Values) url.Values {
	result := make(url.Values, len(values))
	for key, source := range values {
		result[key] = append([]string(nil), source...)
	}
	return result
}
