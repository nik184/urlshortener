package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			wantErr:  "incorrect url was received",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "incorrect url 2",
			body:     "just text",
			wantErr:  "incorrect url was received",
			wantCode: http.StatusBadRequest,
		},
	}
}

type tc interface {
	getResult(t *testing.T, w httptest.ResponseRecorder) *http.Response
	parceResp(res http.Response) string
	testPostResp(t *testing.T, resBodyStr string)
	prepareBody() io.Reader
	getWriter() *httptest.ResponseRecorder
	getHandler() func(http.ResponseWriter, *http.Request)
}

func TestApiShortenRedirectPipeline(t *testing.T) {
	database.ConnectIfNeeded()
	config.DatabaseDSN = "postgres://urlshortener:urlshortener@localhost:5433/urlshortener_test"
	database.DB.Exec("DELETE FROM url;")

	tests := getTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test(t, tt)
		})

		database.DB.Exec("DELETE FROM url;")

		t.Run(tt.name+" api", func(t *testing.T) {
			apiTt := apiTestCase{testCase: tt}
			test(t, apiTt)
		})

		database.DB.Exec("DELETE FROM url;")

		// t.Run(tt.name, func(t *testing.T) {
		// 	cTt := CompressionTestCase{testCase: tt}
		// 	test(t, cTt)
		// })
	}
}

func test(t *testing.T, tt tc) {
	body := tt.prepareBody()
	request := httptest.NewRequest(http.MethodPost, "/", body)
	w := tt.getWriter()
	h := tt.getHandler()
	h(w, request)
	res := tt.getResult(t, *w)

	defer res.Body.Close()
	resBodyStr := tt.parceResp(*res)

	tt.testPostResp(t, resBodyStr)
}

func testSuccessfulGetReq(t *testing.T, tt testCase, path string) {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(RedirectByURLID)
	h(w, request)

	res := w.Result()

	defer res.Body.Close()

	// require.Equal(t, tt.body, res.Header.Get("Location"), "не удалось получить заголовок Location в ответе")
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "статус ответа на get запрос не соответствует ожидаемому")
}

type testCase struct {
	name     string
	body     string
	wantErr  string
	wantCode int
}

func (tt testCase) getWriter() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func (tt testCase) getHandler() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(ShortURL)
}

func (tt testCase) prepareBody() io.Reader {
	return bytes.NewBuffer([]byte(tt.body))
}

func (tt testCase) getResult(t *testing.T, w httptest.ResponseRecorder) *http.Response {
	res := w.Result()
	assert.Equal(t, tt.wantCode, res.StatusCode, "статус ответа на post запрос не соответствует ожидаемому")

	return res
}

func (tt testCase) parceResp(res http.Response) string {
	resBody, _ := io.ReadAll(res.Body)
	resBodyStr := string(resBody)

	return resBodyStr
}

func (tt testCase) testPostResp(t *testing.T, resBodyStr string) {
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

type apiTestCase struct {
	testCase
}

func (tt apiTestCase) prepareBody() io.Reader {
	req := URLReq{URL: tt.body}
	jsonBody, _ := json.Marshal(req)
	body := bytes.NewBuffer(jsonBody)

	return body
}

func (tt apiTestCase) getHandler() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(APIShortURL)
}

func (tt apiTestCase) parceResp(res http.Response) string {
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		resp := URLResp{}
		resBody, _ := io.ReadAll(res.Body)
		json.Unmarshal(resBody, &resp)
		return resp.Result
	} else {
		return tt.testCase.parceResp(res)
	}
}

func TestFailedGetReq(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/wrongID", nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(RedirectByURLID)
	h(w, request)

	res := w.Result()

	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

type CompressionTestCase struct {
	testCase
}

func (tt CompressionTestCase) getWriter() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Accept-Encoding", "gzip")

	return w
}

func (tt CompressionTestCase) prepareBody() io.Reader {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)

	_, err := w.Write([]byte(tt.body))
	if err != nil {
		panic("AAA!!!")
	}

	err = w.Close()
	if err != nil {
		panic("AAA!!!")
	}

	return bytes.NewBuffer(b.Bytes())
}
