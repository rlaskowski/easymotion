package config

import (
	"os"
	"path/filepath"
)

func SqlitePath() string {
	path := ProjectPath()

	return filepath.Join(path, "easymotion.db")
}

func ProjectPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}
