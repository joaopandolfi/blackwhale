package configurations_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joaopandolfi/blackwhale/v2/configurations"
)

func TestLoadFromMap(t *testing.T) {
	fconf := map[string]string{
		"SERVER_NAME":            "test",
		"MYSQL_USER":             "u",
		"MYSQL_PASSWORD":         "p",
		"MYSQL_HOST":             "h",
		"MYSQL_PORT":             "3306",
		"MYSQL_DB":               "db",
		"MONGO_URL":              "mongodb://x",
		"MONGO_DB":               "m",
		"MONGO_POOL":             "7",
		"SERVER_PORT":            "8080",
		"SERVER_TIMEOUT":         "30",
		"BCRYPT_COST":            "12",
		"TOKEN_VALIDITY_MINUTES": "45",
		"REDIS_USE":              "true",
		"REDIS_DB":               "2",
		"BCRYPT_SECRET":          "bs",
		"RESET_HASH":             "rh",
		"JWT_SECRET":             "js",
		"AES_KEY":                "ak",
	}

	c := configurations.LoadFromMap(fconf)

	if c.Name != "test" {
		t.Errorf("Name: got %q", c.Name)
	}
	if c.MysqlUrl != "u:p@tcp(h:3306)/db" {
		t.Errorf("MysqlUrl: got %q", c.MysqlUrl)
	}
	if c.MongoUrl != "mongodb://x" || c.MongoDb != "m" || c.MongoPool != 7 {
		t.Errorf("Mongo: got %q %q %d", c.MongoUrl, c.MongoDb, c.MongoPool)
	}
	if c.Port != ":8080" {
		t.Errorf("Port: got %q", c.Port)
	}
	if c.Timeout.Read != 30*time.Second || c.Timeout.Write != 30*time.Second {
		t.Errorf("Timeout: got %v", c.Timeout)
	}
	if c.Security.BCryptCost != 12 {
		t.Errorf("BCryptCost: got %d", c.Security.BCryptCost)
	}
	if c.Security.TokenValidity != 45 {
		t.Errorf("TokenValidity: got %d", c.Security.TokenValidity)
	}
	if !c.Redis.Use || c.Redis.DB != 2 {
		t.Errorf("Redis: got %+v", c.Redis)
	}
	if c.Security.JWTSecret != "js" || c.Security.AESKEY != "ak" {
		t.Errorf("Secrets: got %q %q", c.Security.JWTSecret, c.Security.AESKEY)
	}
	if c.BCryptSecret != "bs" || c.ResetHash != "rh" {
		t.Errorf("Top-level secrets: got %q %q", c.BCryptSecret, c.ResetHash)
	}
}

func TestLoadFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"SERVER_NAME": "filetest", "SERVER_PORT": "9090"}`), 0o600); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	c := configurations.LoadFromFile(path)
	if c.Name != "filetest" {
		t.Errorf("Name: got %q", c.Name)
	}
	if c.Port != ":9090" {
		t.Errorf("Port: got %q", c.Port)
	}
}
