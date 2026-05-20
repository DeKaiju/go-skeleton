package curl_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	curlpkg "github.com/dekaiju/go-skeleton/pkg/curl"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "ok",
		})
	}))
	defer server.Close()

	content, err := curlpkg.Get(nil, server.URL, nil, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got := content.Get("message").String(); got != "ok" {
		t.Fatalf("Get() message = %q, want %q", got, "ok")
	}
}

func TestPostForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}

		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": r.Form.Get("id"),
		})
	}))
	defer server.Close()

	content, err := curlpkg.PostForm(nil, server.URL, map[string]interface{}{"id": "95"}, map[string]interface{}{})
	if err != nil {
		t.Fatalf("PostForm() error = %v", err)
	}

	if got := content.Get("id").String(); got != "95" {
		t.Fatalf("PostForm() id = %q, want %q", got, "95")
	}
}
