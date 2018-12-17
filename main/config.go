package main

import (
	"flag"
	"github.com/vitelabs/powclient/log15"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"
)

var (
	mtype         = flag.String("type", "gpu", "processor type, cpu or gpu")
	serverUrl     = flag.String("server", "127.0.0.1:7076", "server url, eg:127.0.0.1:7076")
	metricsEnable = flag.Bool("metrics", false, "enable metrics and export to influxdb")
	influxDBUrl   = flag.String("influxdb", "127.0.0.1:8086", "influxdb url,eg:127.0.0.1:8086")
)

func MakeDefaultLogger(absFilePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   absFilePath,
		MaxSize:    100,
		MaxBackups: 14,
		MaxAge:     14,
		Compress:   true,
		LocalTime:  true,
	}
}

func InitLog(dir, lvl string) {
	logLevel, err := log15.LvlFromString(lvl)
	if err != nil {
		logLevel = log15.LvlInfo
	}
	path := filepath.Join(dir, "pow", time.Now().Format("2006-01-02T15-04"))
	filename := filepath.Join(path, "powclient.log")
	log15.Root().SetHandler(
		log15.LvlFilterHandler(logLevel, log15.StreamHandler(MakeDefaultLogger(filename), log15.LogfmtFormat())),
	)
}

// DefaultDataDir is the default data directory to use for the databases and other persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "GVite")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "GVite")
		} else {
			return filepath.Join(home, ".gvite")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
