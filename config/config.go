package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Handle      string `toml:"handle"`
	AppPassword string `toml:"app_password"`
	PDS         string `toml:"pds"`
}

func Load() (*Config, error) {
	cfg := &Config{
		PDS: "https://bsky.social",
	}

	// 環境変数を優先
	if h := os.Getenv("BSKY_HANDLE"); h != "" {
		cfg.Handle = h
	}
	if p := os.Getenv("BSKY_APP_PASSWORD"); p != "" {
		cfg.AppPassword = p
	}

	// 環境変数で揃ってなければ設定ファイルから補完
	if cfg.Handle == "" || cfg.AppPassword == "" {
		path, err := configFilePath()
		if err == nil {
			var fileCfg Config
			if _, err := toml.DecodeFile(path, &fileCfg); err == nil {
				if cfg.Handle == "" {
					cfg.Handle = fileCfg.Handle
				}
				if cfg.AppPassword == "" {
					cfg.AppPassword = fileCfg.AppPassword
				}
				if fileCfg.PDS != "" {
					cfg.PDS = fileCfg.PDS
				}
			}
		}
	}

	return cfg, nil
}

// NeedsLogin は認証情報が不足しているかどうかを返す。
func (c *Config) NeedsLogin() bool {
	return c.Handle == "" || c.AppPassword == ""
}

// Save は認証情報を設定ファイルに保存する。
func (c *Config) Save() error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(c)
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "post2bsky", "config.toml"), nil
}
