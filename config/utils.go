package config

import (
	"os"
	"path/filepath"
)

func SqlitePath() string {
	path := ProjectPath()
	return filepath.Join(path, "easymotion.db")
}

//Path name where is store immudb files
func ImmuDBPath() string {
	path := ProjectPath()
	return filepath.Join(path, "data")
}

func ProjectName() string {
	return "easymotion"
}

//Full, current project path
func ProjectPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

//Converting megabytes to bytes
func ToBytes(m int64) int64 {
	return m * (1024 * 1024)
}
