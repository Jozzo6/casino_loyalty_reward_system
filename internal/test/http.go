package test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"
)

type TestRequestFile struct {
	FileName  string
	Data      string
	FieldName string
}

type TestRequest struct {
	Vars      map[string]string
	Body      string
	Context   context.Context
	UrlParams map[string]string
	File      *TestRequestFile
}

func (tr *TestRequest) GetRequest(method string) (*http.Request, error) {
	var r *http.Request
	if tr.Body != "" {
		body := strings.NewReader(tr.Body)
		r = httptest.NewRequest(method, "api://", body)
	} else if tr.File != nil {
		var buf bytes.Buffer
		multipartWriter := multipart.NewWriter(&buf)
		defer multipartWriter.Close()

		// add form field
		filePart, err := multipartWriter.CreateFormFile(tr.File.FieldName, tr.File.FileName)
		if err != nil {
			return nil, err
		}
		_, err = filePart.Write([]byte(tr.File.Data))
		if err != nil {
			return nil, err
		}

		r = httptest.NewRequest(http.MethodPost, "api://", &buf)
		r.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	} else {
		r = httptest.NewRequest(method, "api://", nil)
	}

	if tr.UrlParams != nil {
		q := r.URL.Query()
		for k, v := range tr.UrlParams {
			q.Add(k, v)
		}
		r.URL.RawQuery = q.Encode()
	}
	if tr.Context != nil {
		r = r.WithContext(tr.Context)
	}
	if tr.Vars != nil {
		rctx := chi.NewRouteContext()
		for k, v := range tr.Vars {
			rctx.URLParams.Add(k, v)
		}
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	}

	return r, nil
}
