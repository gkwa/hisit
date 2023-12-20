package hisit

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Options struct {
	LogFormat string
	LogLevel  string
	BaseDir   string
	Age       string
	Depth     int
}

func Execute() int {
	options := parseArgs()

	logger, err := getLogger(options.LogLevel, options.LogFormat)
	if err != nil {
		slog.Error("getLogger", "error", err)
		return 1
	}

	slog.SetDefault(logger)

	err = run(options)
	if err != nil {
		slog.Error("run failed", "error", err)
		return 1
	}
	return 0
}

func parseArgs() Options {
	options := Options{}

	flag.StringVar(&options.LogLevel, "log-level", "info", "Log level (debug, info, warn, error), default: info")
	flag.StringVar(&options.LogFormat, "log-format", "text", "Log format (text or json)")
	flag.StringVar(&options.BaseDir, "dir", "", "Specify the base path to scan")
	flag.StringVar(&options.Age, "age", "1d", "Specify the age for modification time comparison (e.g., 1d)")
	flag.IntVar(&options.Depth, "depth", 2, "Specify the depth of directory traversal")

	flag.Parse()

	return options
}

func run(options Options) error {
	slog.Debug("test", "test", "Debug")
	slog.Debug("test", "LogLevel", options.LogLevel)
	slog.Info("test", "test", "Info")
	slog.Error("test", "test", "Error")

	ageDuration, err := parseAge(options.Age)
	if err != nil {
		slog.Error("Error parsing age", "error", err)
		os.Exit(1)
	}

	slog.Debug("parseAge", "age", ageDuration)

	basePath, err := expandPath(options.BaseDir)
	if err != nil {
		slog.Error("Error expanding path", "error", err)
		os.Exit(1)
	}

	err = scanDirectories(basePath, ageDuration, options.Depth)
	if err != nil {
		slog.Error("Error scanning directories", "error", err)
		os.Exit(1)
	}

	return nil
}

func parseAge(age string) (time.Duration, error) {
	unit := age[len(age)-1:]
	value, err := strconv.Atoi(age[:len(age)-1])
	if err != nil {
		return 0, err
	}

	slog.Debug("parseAge", "age", age, "unit", unit, "value", value)

	switch unit {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}
}

func expandPath(path string) (string, error) {
	return filepath.Abs(path)
}

func scanDirectories(basePath string, ageDuration time.Duration, depth int) error {
	now := time.Now()

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(basePath, path)
		parts := filepath.SplitList(relPath)
		if len(parts) > depth {
			return filepath.SkipDir
		}

		if path == basePath {
			return nil
		}

		if info.IsDir() {
			modTime := info.ModTime()
			diff := now.Sub(modTime)

			if diff <= ageDuration {
				slog.Debug("Modified directory found", "directory", path)
			}
		}

		return nil
	})

	return err
}
