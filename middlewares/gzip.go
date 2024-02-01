package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

func Gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			compressedData := bytes.NewReader(body)

			reader, err := gzip.NewReader(compressedData)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = reader
		}

		h.ServeHTTP(w, r)
	})
}
