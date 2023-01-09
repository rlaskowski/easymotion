package config

import (
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

//Converting megabytes to bytes
func ToBytes(m int64) int64 {
	return m * (1024 * 1024)
}

func ExecutableName() string {
	executable := "easymotion"

	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s.exe", executable)
	}

	return executable
}

func ExecutablePath() string {
	executable, err := os.Executable()

	if err != nil {
		return "./"
	}

	return executable
}

func WorkingDirectory() string {
	ex := ExecutablePath()

	if filepath.Base(ex) == ExecutableName() {
		return filepath.Dir(ex)
	}

	wd, err := os.Getwd()

	if err != nil {
		return "./"
	}

	return wd
}

// Path to the certificate file, used by TLS
func CertFile() string {
	filename := "easymotion.crt"

	env, ok := os.LookupEnv("EASYMOTION_TLS_CERTFILE")
	if ok {
		filename = env
	}

	return envPath("EASYMOTION_TLS_CERTFILE_TLS_PATH", filename)
}

// Path to the key file, used by TLS
func KeyFile() string {
	filename := "easymotion.key"

	env, ok := os.LookupEnv("EASYMOTION_TLS_CERTFILE_TLS_KEYFILE")
	if ok {
		filename = env
	}

	return envPath("EASYMOTION_TLS_CERTFILE_TLS_PATH", filename)
}

func hubURL() string {
	hubu := "amqp://guest:guest@localhost:5672/"

	env, ok := os.LookupEnv("EASYMOTION_HUB_URL")
	if ok {
		hubu = env
	}

	return hubu
}

func RecordsPath() string {
	path := path.Join(WorkingDirectory(), "videos")

	env, ok := os.LookupEnv("EASYMOTION_RECORDS_PATH")
	if ok {
		path = env
	}

	return path
}

func envPath(env, filename string) string {
	if v, ok := os.LookupEnv(env); ok {
		return filepath.Join(v, filename)
	}
	return filepath.Join(WorkingDirectory(), filename)
}

// Path to log file
func LogPath() string {
	return filepath.Join(WorkingDirectory(), "easymotion.log")
}

// Returns local machine IP address
func ServiceIP() string {
	host, err := os.Hostname()
	if err != nil {
		return ""
	}

	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return ""
	}

	return ip.String()
}

// Returns local machine hostname
func Hostname() string {
	host, err := os.Hostname()
	if err != nil {
		return ""
	}

	return host
}
