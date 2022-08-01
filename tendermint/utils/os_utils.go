package utils

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

const (
	// DefaultDirMod default unix perms for k9s directory.
	DefaultDirMod os.FileMode = 0755
	// DefaultFileMod default unix perms for k9s files.
	DefaultFileMod os.FileMode = 0600
)

// EnsurePath ensures a directory exist from the given path.
func EnsurePath(path string, mod os.FileMode) {
	dir := filepath.Dir(path)
	EnsureFullPath(dir, mod)
}

// EnsureFullPath ensures a directory exist from the given path.
func EnsureFullPath(path string, mod os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, mod); err != nil {
			log.Fatal().Msgf("Unable to create dir %q %v", path, err)
		}
	}
}

func WriteDataToFile(filePath string, data []byte) int {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Err(err)
		return -1
	}
	defer file.Close()

	bytesWritten, err := file.Write(data)
	if err != nil {
		log.Err(err)
		return -1
	}
	if bytesWritten > 0 {
		return 0
	} else {
		return -2
	}
}
