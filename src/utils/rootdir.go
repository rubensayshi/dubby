package utils

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

const PACKAGE = "github.com/rubensayshi/duconverter"
const THISFILE = "src/utils/rootdir.go"

var ROOT = mustGetRootDir()

func mustGetRootDir() string {
	rootDir, err := getRootDir()
	if err != nil {
		panic(err)
	}

	return rootDir
}

func getRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrapf(err, "getRootDir")
	}

	// use ROOTDIR env var if set
	envRootdir := os.Getenv("ROOTDIR")
	if envRootdir != "" {
		if path.IsAbs(envRootdir) {
			return envRootdir, nil
		} else {
			return path.Clean(path.Join(cwd, envRootdir)), nil
		}
	}

	// attempt to determine the root dir from the caller filename
	_, filename, _, _ := runtime.Caller(0)
	rootDir, ok := determineRootUpToFilename(filename)
	if ok {
		return rootDir, nil
	}

	// attempt to determine the root dir from the cwd
	rootDir, ok = determineRootUpToPackage(cwd)
	if ok {
		return rootDir, nil
	}

	// attempt to determine the root dir from the executable path
	executable, err := os.Executable()
	if err != nil {
		return "", errors.Wrapf(err, "getRootDir")
	}
	rootDir, ok = determineRootUpToPackage(executable)
	if ok {
		return rootDir, nil
	}

	// fallback to cwd
	return cwd, nil
}

// use the PACKAGE name to determine the part of the path that is the "root"
func determineRootUpToPackage(dir string) (string, bool) {
	// find the last occurrence of our package path
	idx := strings.LastIndex(dir, PACKAGE)
	if idx == -1 {
		return "", false
	}

	// join everything that precedes our last occurrence of our package path with the package path
	rootDir := path.Join(dir[0:idx], PACKAGE)

	return rootDir, true
}

// use the THISFILE name to determine the part of the path that is the "root"
func determineRootUpToFilename(filename string) (string, bool) {
	// find the last occurrence of our filename
	idx := strings.LastIndex(filename, THISFILE)
	if idx == -1 {
		return "", false
	}

	// take everything that precedes our last occurrence of our filename
	rootDir := strings.TrimSuffix(filename[0:idx], "/")

	return rootDir, true
}
