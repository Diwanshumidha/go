package store

import (
	"os"
	"path/filepath"
)

const ReadExecPerm = 0o755

func GetWiresharkPath(create bool) (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homePath, "Documents", "Wireshark")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create the directory
		if create {
			err = os.Mkdir(path, ReadExecPerm)
			if err != nil {
				return "", err
			}
		} else {
			return "", os.ErrNotExist
		}
	}
	return path, nil
}
