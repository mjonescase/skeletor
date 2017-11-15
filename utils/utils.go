package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	res := fmt.Sprintf("%x", h.Sum(nil))
	return res
}

func MustEncode(w io.Writer, i interface{}) {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-type", "application/json")
	}

	e := json.NewEncoder(w)
	if err := e.Encode(i); err != nil {
		panic(err)
	}
}
