package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey       string
	Model        string
	CommitStyle  string
	CustomFormat string
	Language     string
	MaxDiffLines int
}

const (
	StyleConventional = "conventional"
	StyleGitmoji      = "gitmoji"
	StyleFree         = "free"
	StyleCustom       = "custom"
)

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("model", "claude-sonnet-4-6")
	v.SetDefault("commit_style", StyleConventional)
	v.SetDefault("language", "en")
	v.SetDefault("max_diff_lines", 500)

	v.SetConfigName(".ezgocommit")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(home, ".config", "ezgocommit"))
	}

	_ = v.ReadInConfig()

	v.SetEnvPrefix("")
	v.AutomaticEnv()

	cfg := &Config{
		APIKey:       resolveAPIKey(v),
		Model:        v.GetString("model"),
		CommitStyle:  v.GetString("commit_style"),
		CustomFormat: v.GetString("custom_format"),
		Language:     v.GetString("language"),
		MaxDiffLines: v.GetInt("max_diff_lines"),
	}

	return cfg, nil
}

func LoadWithOverrides(style, model string) (*Config, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	if style != "" {
		cfg.CommitStyle = style
	}
	if model != "" {
		cfg.Model = model
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf(
			"Anthropic API key not found.\n\n" +
				"Set it via environment variable:\n" +
				"  export ANTHROPIC_API_KEY=sk-ant-...\n\n" +
				"Or add it to ~/.config/ezgocommit/config.toml:\n" +
				"  api_key = \"sk-ant-...\"\n",
		)
	}
	return nil
}

func resolveAPIKey(v *viper.Viper) string {
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return key
	}
	return v.GetString("api_key")
}
