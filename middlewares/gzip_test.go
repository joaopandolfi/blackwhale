package middlewares_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joaopandolfi/blackwhale/v2/middlewares"
)

func gzipBody(t *testing.T, payload string) *bytes.Buffer {
	t.Helper()
	buf := &bytes.Buffer{}
	zw := gzip.NewWriter(buf)
	if _, err := zw.Write([]byte(payload)); err != nil {
		t.Fatalf("gzipping: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("closing gzip: %v", err)
	}
	return buf
}

func TestGzipDecompressesRequestBody(t *testing.T) {
	var got string
	h := middlewares.Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading body: %v", err)
		}
		got = string(b)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(gzipBody(t, "hello gzip").Bytes()))
	req.Header.Set("Content-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	if got != "hello gzip" {
		t.Fatalf("body: got %q", got)
	}
}

func TestGzipPassthroughWithoutHeader(t *testing.T) {
	var got string
	h := middlewares.Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = string(b)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("plain"))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	if got != "plain" {
		t.Fatalf("body: got %q", got)
	}
}

func TestGzipBadPayload(t *testing.T) {
	h := middlewares.Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not gzip at all"))
	req.Header.Set("Content-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}
