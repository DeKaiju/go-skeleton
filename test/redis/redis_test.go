package redis_test

import (
	"net"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	redigo "github.com/gomodule/redigo/redis"

	redispkg "github.com/dekaiju/go-skeleton/pkg/redis"
)

func TestRedisPool(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run() error = %v", err)
	}
	defer server.Close()

	host, port, err := net.SplitHostPort(server.Addr())
	if err != nil {
		t.Fatalf("SplitHostPort() error = %v", err)
	}

	previousHost := os.Getenv("CACHE_HOST")
	previousPort := os.Getenv("CACHE_PORT")
	t.Cleanup(func() {
		_ = os.Setenv("CACHE_HOST", previousHost)
		_ = os.Setenv("CACHE_PORT", previousPort)
		_ = redispkg.Close()
	})

	_ = os.Setenv("CACHE_HOST", host)
	_ = os.Setenv("CACHE_PORT", port)

	redispkg.GetCache()
	conn := redispkg.Get()
	defer conn.Close()

	if _, err := conn.Do("SET", "redis_test", "hello"); err != nil {
		t.Fatalf("SET error = %v", err)
	}

	value, err := redigo.String(conn.Do("GET", "redis_test"))
	if err != nil {
		t.Fatalf("GET error = %v", err)
	}

	if value != "hello" {
		t.Fatalf("GET value = %q, want %q", value, "hello")
	}
}
