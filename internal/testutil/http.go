package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ErrorResponse struct {
	Error struct {
		Code    string         `json:"code"`
		Message string         `json:"message"`
		Details map[string]any `json:"details"`
	} `json:"error"`
}

func (a *TestApp) DoJSON(t *testing.T, method string, path string, token string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var payload bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			t.Fatalf("encode request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, &payload)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res := httptest.NewRecorder()
	a.Router.ServeHTTP(res, req)
	return res
}

func DecodeJSON(t *testing.T, res *httptest.ResponseRecorder, out any) {
	t.Helper()

	if err := json.Unmarshal(res.Body.Bytes(), out); err != nil {
		t.Fatalf("decode JSON response: %v\nbody: %s", err, res.Body.String())
	}
}

func AssertStatus(t *testing.T, res *httptest.ResponseRecorder, want int) {
	t.Helper()

	if res.Code != want {
		t.Fatalf("status = %d, want %d\nbody: %s", res.Code, want, res.Body.String())
	}
}

func AssertError(t *testing.T, res *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()

	AssertStatus(t, res, wantStatus)
	var body ErrorResponse
	DecodeJSON(t, res, &body)
	if body.Error.Code != wantCode {
		t.Fatalf("error.code = %q, want %q\nbody: %s", body.Error.Code, wantCode, res.Body.String())
	}
	if body.Error.Details == nil {
		t.Fatalf("error.details is nil, want object\nbody: %s", res.Body.String())
	}
}

func DataMap(t *testing.T, res *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body struct {
		Data map[string]any `json:"data"`
	}
	DecodeJSON(t, res, &body)
	if body.Data == nil {
		t.Fatalf("response data is nil\nbody: %s", res.Body.String())
	}
	return body.Data
}

func DataSlice(t *testing.T, res *httptest.ResponseRecorder) []map[string]any {
	t.Helper()

	var body struct {
		Data []map[string]any `json:"data"`
	}
	DecodeJSON(t, res, &body)
	if body.Data == nil {
		t.Fatalf("response data is nil\nbody: %s", res.Body.String())
	}
	return body.Data
}

func DataString(t *testing.T, data map[string]any, key string) string {
	t.Helper()

	value, ok := data[key].(string)
	if !ok || value == "" {
		t.Fatalf("data[%q] = %#v, want non-empty string", key, data[key])
	}
	return value
}

func (a *TestApp) Login(t *testing.T, email string, password string) (string, string) {
	t.Helper()

	res := a.DoJSON(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email":    email,
		"password": password,
	})
	AssertStatus(t, res, http.StatusOK)
	data := DataMap(t, res)
	return DataString(t, data, "access_token"), DataString(t, data, "refresh_token")
}
