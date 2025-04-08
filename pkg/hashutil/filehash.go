package hashutil

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func GetFileHashString(targetFilePath string) (*string, error) {
	f, err := os.Open(targetFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))
	slog.Debug("Hash generated", slog.String("hash", hash))
	return &hash, nil
}

func CheckTwoFileHashSame(firstFilePath, secondFilePath string) bool {
	// get first file's hash
	firstHash, err := GetFileHashString(firstFilePath)
	if err != nil {
		slog.Error("Failed to get first file hash", slog.Any("error", err))
		return false
	}

	// get second file's hash
	secondHash, err := GetFileHashString(secondFilePath)
	if err != nil {
		slog.Error("Failed to get second file hash", slog.Any("error", err))
		return false
	}

	// return is hashs same
	return strings.Compare(*firstHash, *secondHash) == 0
}
