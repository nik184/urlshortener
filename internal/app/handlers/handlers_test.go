package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name     string
	body     string
	wantErr  string
	wantCode int
}

func getTestCases() []testCase {
	return []testCase{
		{
			name:     "success shorten 1",
			body:     "http://abcd.net",
			wantCode: http.StatusCreated,
		},
		{
			name:     "success shorten 2",
			body:     "http://textextex.ru/asdasdasd",
			wantCode: http.StatusCreated,
		},
		{
			name:     "success shorten 3",
			body:     "http://textextex.edrftgy/asdasdasd/w3e4rf5tg6yh7uj8ik",
			wantCode: http.StatusCreated,
		},
		{
			name:     "incorrect url 1",
			body:     "http://",
			wantErr:  "incorrect url was received!",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "incorrect url 2",
			body:     "just text",
			wantErr:  "incorrect url was received!",
			wantCode: http.StatusBadRequest,
		},
	}
}

func TestApiShortenRedirectPipeline(t *testing.T) {
	tests := getTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAPIPostReq(t, tt)
		})
	}
}

func TestShortenRedirectPipeline(t *testing.T) {
	tests := getTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPostReq(t, tt)
		})
	}
}

func testPostReq(t *testing.T, tt testCase) {
	body := bytes.NewBuffer([]byte(tt.body))
	request := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(GenerateURL)
	h(w, request)

	res := w.Result()

	assert.Equal(t, tt.wantCode, res.StatusCode, "статус ответа на post запрос не соответствует ожидаемому")

	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	resBodyStr := string(resBody)

	if tt.wantErr != "" {
		assert.Contains(t, resBodyStr, tt.wantErr)
	} else {
		require.Contains(t, resBodyStr, config.BaseURL)
		parsedURL, err := url.ParseRequestURI(resBodyStr)

		if err != nil {
			assert.Fail(t, err.Error())
		}

		testSuccessfulGetReq(t, tt, parsedURL.Path)
	}
}

func testAPIPostReq(t *testing.T, tt testCase) {
	req := Req{URL: tt.body}
	jsonBody, _ := json.Marshal(req)
	body := bytes.NewBuffer(jsonBody)
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(APIGenerateURL)
	h(w, request)

	res := w.Result()

	assert.Equal(t, tt.wantCode, res.StatusCode, "статус ответа на post запрос не соответствует ожидаемому")

	defer res.Body.Close()

	var resBodyStr string

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		resp := Resp{}
		resBody, _ := io.ReadAll(res.Body)
		json.Unmarshal(resBody, &resp)
		resBodyStr = resp.Result
	} else {
		resBody, _ := io.ReadAll(res.Body)
		resBodyStr = string(resBody)
	}

	if tt.wantErr != "" {
		assert.Contains(t, resBodyStr, tt.wantErr)
	} else {
		require.Contains(t, resBodyStr, config.BaseURL)
		parsedURL, err := url.ParseRequestURI(resBodyStr)

		if err != nil {
			assert.Fail(t, err.Error())
		}

		testSuccessfulGetReq(t, tt, parsedURL.Path)
	}
}

func testSuccessfulGetReq(t *testing.T, tt testCase, path string) {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(RedirectByURLID)
	h(w, request)

	res := w.Result()

	defer res.Body.Close()

	require.Equal(t, tt.body, res.Header.Get("Location"), "не удалось получить заголовок Location в ответе")
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "статус ответа на get запрос не соответствует ожидаемому")
}

func TestFailedGetReq(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/wrongID", nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(RedirectByURLID)
	h(w, request)

	res := w.Result()

	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	resBodyStr := string(resBody)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Contains(t, resBodyStr, "wrong id was received!")
}
