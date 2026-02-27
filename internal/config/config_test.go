package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("ANTHROPIC_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.Model != "claude-sonnet-4-6" {
		t.Errorf("default model = %q, want %q", cfg.Model, "claude-sonnet-4-6")
	}
	if cfg.CommitStyle != StyleConventional {
		t.Errorf("default commit_style = %q, want %q", cfg.CommitStyle, StyleConventional)
	}
	if cfg.Language != "en" {
		t.Errorf("default language = %q, want %q", cfg.Language, "en")
	}
	if cfg.MaxDiffLines != 500 {
		t.Errorf("default max_diff_lines = %d, want 500", cfg.MaxDiffLines)
	}
}

func TestLoad_EnvVarAPIKey(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-test-from-env")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.APIKey != "sk-ant-test-from-env" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "sk-ant-test-from-env")
	}
}

func TestLoad_ConfigFile(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, ".ezgocommit.toml")
	content := []byte(`
model          = "claude-opus-4-6"
commit_style   = "gitmoji"
language       = "pt"
max_diff_lines = 200
api_key        = "sk-ant-from-file"
`)
	if err := os.WriteFile(cfgFile, content, 0600); err != nil {
		t.Fatal(err)
	}

	// Change to that directory so Viper finds the file
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	os.Unsetenv("ANTHROPIC_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.Model != "claude-opus-4-6" {
		t.Errorf("model = %q, want %q", cfg.Model, "claude-opus-4-6")
	}
	if cfg.CommitStyle != StyleGitmoji {
		t.Errorf("commit_style = %q, want %q", cfg.CommitStyle, StyleGitmoji)
	}
	if cfg.Language != "pt" {
		t.Errorf("language = %q, want %q", cfg.Language, "pt")
	}
	if cfg.MaxDiffLines != 200 {
		t.Errorf("max_diff_lines = %d, want 200", cfg.MaxDiffLines)
	}
	if cfg.APIKey != "sk-ant-from-file" {
		t.Errorf("api_key = %q, want %q", cfg.APIKey, "sk-ant-from-file")
	}
}

func TestLoad_EnvVarOverridesFile(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, ".ezgocommit.toml")
	os.WriteFile(cfgFile, []byte(`api_key = "sk-ant-from-file"`), 0600)

	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-from-env")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.APIKey != "sk-ant-from-env" {
		t.Errorf("env var should override file: APIKey = %q, want %q", cfg.APIKey, "sk-ant-from-env")
	}
}

func TestLoadWithOverrides(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	cfg, err := LoadWithOverrides("gitmoji", "claude-opus-4-6", "pt")
	if err != nil {
		t.Fatalf("LoadWithOverrides() error: %v", err)
	}

	if cfg.CommitStyle != StyleGitmoji {
		t.Errorf("CommitStyle = %q, want %q", cfg.CommitStyle, StyleGitmoji)
	}
	if cfg.Model != "claude-opus-4-6" {
		t.Errorf("Model = %q, want %q", cfg.Model, "claude-opus-4-6")
	}
	if cfg.Language != "pt" {
		t.Errorf("Language = %q, want %q", cfg.Language, "pt")
	}
}

func TestLoadWithOverrides_EmptyKeepsDefault(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	cfg, err := LoadWithOverrides("", "", "")
	if err != nil {
		t.Fatalf("LoadWithOverrides() error: %v", err)
	}

	if cfg.CommitStyle != StyleConventional {
		t.Errorf("empty override should keep default: CommitStyle = %q", cfg.CommitStyle)
	}
	if cfg.Model != "claude-sonnet-4-6" {
		t.Errorf("empty override should keep default: Model = %q", cfg.Model)
	}
}

func TestValidate_MissingKey(t *testing.T) {
	cfg := &Config{APIKey: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("Validate() should return error when APIKey is empty")
	}
}

func TestValidate_WithKey(t *testing.T) {
	cfg := &Config{APIKey: "sk-ant-anything"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}
}
