package flags

import (
	"flag"

	"github.com/rlaskowski/easymotion/config"
)

var (
	VideosPath string
)

func InitFlags() {
	flag.StringVar(&VideosPath, "f", config.ProjectPath(), "path where will be store video files")
}
