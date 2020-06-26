package config

import (
	"testing"

	"os"
)

func TestLoad(t *testing.T) {
	os.Setenv("PREST_CONF", "testdata/prest.toml")
	defer os.Unsetenv("PREST_CONF")

	Load()
	if len(PrestConf.AccessConf.Tables) < 2 {
		t.Errorf("expected > 2, got: %d", len(PrestConf.AccessConf.Tables))
	}

	Load()
	if !PrestConf.AccessConf.Restrict {
		t.Error("expected true, but got false")
	}

}

func TestParse(t *testing.T) {
	os.Setenv("PREST_CONF", "testdata/prest.toml")

	viperCfg()
	cfg := &Prest{}
	err := Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.HTTPPort != 6000 {
		t.Errorf("expected port: 6000, got: %d", cfg.HTTPPort)
	}
	if cfg.PGDatabase != "prest" {
		t.Errorf("expected database: prest, got: %s", cfg.PGDatabase)
	}

	os.Setenv("PREST_CONF", "../prest.toml")
	os.Setenv("PREST_HTTP_PORT", "4000")

	viperCfg()
	cfg = &Prest{}
	err = Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.HTTPPort != 4000 {
		t.Errorf("expected port: 4000, got: %d", cfg.HTTPPort)
	}
	if !cfg.EnableDefaultJWT {
		t.Error("EnableDefaultJWT: expected true but got false")
	}

	os.Setenv("PREST_CONF", "")
	os.Setenv("PREST_JWT_DEFAULT", "false")

	viperCfg()
	cfg = &Prest{}
	err = Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.HTTPPort != 4000 {
		t.Errorf("expected port: 4000, got: %d", cfg.HTTPPort)
	}
	if cfg.EnableDefaultJWT {
		t.Error("EnableDefaultJWT: expected false but got true")
	}

	os.Unsetenv("PREST_JWT_DEFAULT")
	os.Setenv("PREST_CONF", "testdata/prest.toml")

	viperCfg()
	cfg = &Prest{}
	err = Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.HTTPPort != 4000 {
		t.Errorf("expected port: 4000, got: %d", cfg.HTTPPort)
	}

	os.Unsetenv("PREST_CONF")
	os.Unsetenv("PREST_HTTP_PORT")
	os.Setenv("PREST_JWT_KEY", "s3cr3t")

	viperCfg()
	cfg = &Prest{}
	err = Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.JWTKey != "s3cr3t" {
		t.Errorf("expected jwt key: s3cr3t, got: %s", cfg.JWTKey)
	}
	if cfg.JWTAlgo != "HS256" {
		t.Errorf("expected (default) jwt algo: HS256, got: %s", cfg.JWTAlgo)
	}

	os.Unsetenv("PREST_JWT_KEY")
	os.Setenv("PREST_JWT_ALGO", "HS512")

	viperCfg()
	cfg = &Prest{}
	err = Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.JWTAlgo != "HS512" {
		t.Errorf("expected jwt algo: HS512, got: %s", cfg.JWTAlgo)
	}

	os.Unsetenv("PREST_JWT_ALGO")
}

func TestGetDefaultPrestConf(t *testing.T) {
	testCases := []struct {
		name        string
		defaultFile string
		prestConf   string
		result      string
	}{
		{"empty config", "./prest.toml", "", ""},
		{"custom config", "./prest.toml", "../prest.toml", "../prest.toml"},
		{"default config", "./testdata/prest.toml", "", "./testdata/prest.toml"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defaultFile = tc.defaultFile
			cfg := getDefaultPrestConf(tc.prestConf)
			if cfg != tc.result {
				t.Errorf("expected %v, but got %v", tc.result, cfg)
			}
		})
	}
}

func TestDatabaseURL(t *testing.T) {
	os.Setenv("PREST_PG_URL", "postgresql://user:pass@localhost:1234/mydatabase/?sslmode=disable")

	viperCfg()
	cfg := &Prest{}
	err := Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.PGDatabase != "mydatabase" {
		t.Errorf("expected database name: mydatabase, got: %s", cfg.PGDatabase)
	}
	if cfg.PGHost != "localhost" {
		t.Errorf("expected database host: localhost, got: %s", cfg.PGHost)
	}
	if cfg.PGPort != 1234 {
		t.Errorf("expected database port: 1234, got: %d", cfg.PGPort)
	}
	if cfg.PGUser != "user" {
		t.Errorf("expected database user: user, got: %s", cfg.PGUser)
	}
	if cfg.PGPass != "pass" {
		t.Errorf("expected database password: pass, got: %s", cfg.PGPass)
	}
	if cfg.SSLMode != "disable" {
		t.Errorf("expected database ssl mode: disable, got: %s", cfg.SSLMode)
	}

	os.Unsetenv("PREST_PG_URL")
	os.Setenv("DATABASE_URL", "postgresql://cloud:cloudPass@localhost:5432/CloudDatabase/?sslmode=disable")

	cfg = &Prest{}
	err = Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.PGPort != 5432 {
		t.Errorf("expected database port: 5432, got: %d", cfg.PGPort)
	}
	if cfg.PGUser != "cloud" {
		t.Errorf("expected database user: cloud, got: %s", cfg.PGUser)
	}
	if cfg.PGPass != "cloudPass" {
		t.Errorf("expected database password: cloudPass, got: %s", cfg.PGPass)
	}
	if cfg.SSLMode != "disable" {
		t.Errorf("expected database SSL mode: disable, got: %s", cfg.SSLMode)
	}

	os.Unsetenv("DATABASE_URL")
}

func TestHTTPPort(t *testing.T) {
	os.Setenv("PORT", "8080")

	viperCfg()
	cfg := &Prest{}
	err := Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("expected http port: 8080, got: %d", cfg.HTTPPort)
	}

	// set env PREST_HTTP_PORT and PORT
	os.Setenv("PREST_HTTP_PORT", "3000")

	cfg = &Prest{}
	err = Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.HTTPPort != 3000 {
		t.Errorf("expected http port: 8080, got: %d", cfg.HTTPPort)
	}

	// unset env PORT and set PREST_HTTP_PORT
	os.Unsetenv("PORT")

	cfg = &Prest{}
	err = Parse(cfg)
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if cfg.HTTPPort != 3000 {
		t.Errorf("expected http port: 3000, got: %d", cfg.HTTPPort)
	}

	os.Unsetenv("PREST_HTTP_PORT")
}

func Test_parseDatabaseURL(t *testing.T) {
	c := &Prest{PGURL: "postgresql://user:pass@localhost:5432/mydatabase/?sslmode=require"}
	if err := parseDatabaseURL(c); err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if c.PGDatabase != "mydatabase" {
		t.Errorf("expected database name: mydatabase, got: %s", c.PGDatabase)
	}
	if c.PGPort != 5432 {
		t.Errorf("expected database port: 5432, got: %d", c.PGPort)
	}
	if c.PGUser != "user" {
		t.Errorf("expected database user: user, got: %s", c.PGUser)
	}
	if c.PGPass != "pass" {
		t.Errorf("expected database password: password, got: %s", c.PGPass)
	}
	if c.SSLMode != "require" {
		t.Errorf("expected database SSL mode: require, got: %s", c.SSLMode)
	}
}
