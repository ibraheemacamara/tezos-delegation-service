package config

import (
	"os"
	"strings"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()
	return tmpfile.Name()
}

func TestLoadConfig_Success(t *testing.T) {
	configYAML := `
server:
  host: "localhost"
  port: 8080
  metricsPort: 9090
tzkt:
  url: "http://tzkt.io"
db:
  host: "dbhost"
  port: 5432
  user: "user"
  password: "pass"
  database: "mydb"
`
	path := writeTempConfig(t, configYAML)
	defer os.Remove(path)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if cfg.Server.Host != "localhost" || cfg.Server.Port != 8080 || cfg.Server.MetricsPort != 9090 {
		t.Errorf("server config not loaded correctly: %+v", cfg.Server)
	}
	if cfg.Tzkt.Url != "http://tzkt.io" {
		t.Errorf("tzkt config not loaded correctly: %+v", cfg.Tzkt)
	}
	if cfg.Db.Host != "dbhost" || cfg.Db.Port != 5432 || cfg.Db.User != "user" || cfg.Db.Password != "pass" || cfg.Db.Database != "mydb" {
		t.Errorf("db config not loaded correctly: %+v", cfg.Db)
	}
}

func TestLoadConfig_MissingFields(t *testing.T) {
	cases := []struct {
		name    string
		yaml    string
		errPart string
	}{
		{
			"missing server port",
			`server:
  host: localhost
  metricsPort: 9090
tzkt:
  url: url
db:
  host: h
  port: 1
  user: u
  password: p
  database: d
`,
			"server port is required",
		},
		{
			"missing metrics port",
			`server:
  host: localhost
  port: 8080
tzkt:
  url: url
db:
  host: h
  port: 1
  user: u
  password: p
  database: d
`,
			"server metrics port is required",
		},
		{
			"missing tzkt url",
			`server:
  host: localhost
  port: 8080
  metricsPort: 9090
tzkt:
  url: ""
db:
  host: h
  port: 1
  user: u
  password: p
  database: d
`,
			"tzkt url is required",
		},
		{
			"missing db host",
			`server:
  host: localhost
  port: 8080
  metricsPort: 9090
tzkt:
  url: url
db:
  host: ""
  port: 1
  user: u
  password: p
  database: d
`,
			"db host is required",
		},
		{
			"missing db port",
			`server:
  host: localhost
  port: 8080
  metricsPort: 9090
tzkt:
  url: url
db:
  host: h
  user: u
  password: p
  database: d
`,
			"db port is required",
		},
		{
			"missing db user",
			`server:
  host: localhost
  port: 8080
  metricsPort: 9090
tzkt:
  url: url
db:
  host: h
  port: 1
  user: ""
  password: p
  database: d
`,
			"db user is required",
		},
		{
			"missing db password",
			`server:
  host: localhost
  port: 8080
  metricsPort: 9090
tzkt:
  url: url
db:
  host: h
  port: 1
  user: u
  password: ""
  database: d
`,
			"db password is required",
		},
		{
			"missing db database",
			`server:
  host: localhost
  port: 8080
  metricsPort: 9090
tzkt:
  url: url
db:
  host: h
  port: 1
  user: u
  password: p
  database: ""
`,
			"db database is required",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := writeTempConfig(t, c.yaml)
			defer os.Remove(path)
			_, err := LoadConfig(path)
			if err == nil || !strings.Contains(err.Error(), c.errPart) {
				t.Errorf("expected error containing %q, got %v", c.errPart, err)
			}
		})
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	badYAML := `server: [bad yaml`
	path := writeTempConfig(t, badYAML)
	defer os.Remove(path)
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}
