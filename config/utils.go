package config

import (
	"os"
	"path/filepath"
)

func SqlitePath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(wd, "easymotion.db")
}
