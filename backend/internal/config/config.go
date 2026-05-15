package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
}

type yamlConfig struct {
	Port     string     `yaml:"port"`
	Database dbConfig   `yaml:"database"`
	JWT      jwtConfig  `yaml:"jwt"`
}

type dbConfig struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Name    string `yaml:"name"`
	Charset string `yaml:"charset"`
	Params  string `yaml:"params"`
}

type jwtConfig struct {
	Secret string `yaml:"secret"`
}

func Load() *Config {
	// Load .env from backend/ directory (or cwd)
	loadDotEnv()

	// Load config.yml from cwd or relative to source
	yml := mustLoadYAML(findConfigFile("config.yml"))

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s&charset=%s",
		dbUser, dbPass,
		yml.Database.Host, yml.Database.Port, yml.Database.Name,
		yml.Database.Params, yml.Database.Charset,
	)

	return &Config{
		Port:      yml.Port,
		DBDSN:     dsn,
		JWTSecret: yml.JWT.Secret,
	}
}

func findConfigFile(name string) string {
	// Try cwd first, then source-relative paths
	candidates := []string{
		name,                          // cwd
		filepath.Join("..", name),     // one level up
		filepath.Join("..", "..", "backend", name), // from test/ dir
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return name // let mustLoadYAML report the error
}

func mustLoadYAML(path string) *yamlConfig {
	data, err := os.ReadFile(path)
	if err != nil {
		// Try relative to executable directory
		if execPath, e := os.Executable(); e == nil {
			dir := filepath.Dir(execPath)
			data, err = os.ReadFile(filepath.Join(dir, path))
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: cannot read %s: %v\n", path, err)
		os.Exit(1)
	}

	var y yamlConfig
	if err := yaml.Unmarshal(data, &y); err != nil {
		fmt.Fprintf(os.Stderr, "config: invalid yaml in %s: %v\n", path, err)
		os.Exit(1)
	}

	// Defaults
	if y.Port == "" {
		y.Port = "8080"
	}
	if y.Database.Charset == "" {
		y.Database.Charset = "utf8mb4"
	}
	if y.JWT.Secret == "" {
		y.JWT.Secret = "dev-secret-change-in-production"
	}

	return &y
}

func loadDotEnv() {
	dirs := []string{"", "..", filepath.Join("..", "backend")}
	for _, d := range dirs {
		p := filepath.Join(d, ".env")
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
		break
	}
}
