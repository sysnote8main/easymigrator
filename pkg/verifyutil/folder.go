package verifyutil

import (
	"easymigrator/pkg/hashutil"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
)

func VerifyFolderContents(folderA, folderB string, ignoreFileErrors bool) error {
	return filepath.WalkDir(folderA, func(path string, d fs.DirEntry, err error) error {
		slog.Debug("File walking", slog.String("path", path), slog.Bool("isDir", d.IsDir()))

		if d.IsDir() {
			return nil
		}

		// make relative path for folderB
		relPath, err := filepath.Rel(folderA, path)
		if err != nil {
			return err
		}

		// make file in folderB path
		otherFilePath := fmt.Sprintf("%s/%s", folderB, relPath)

		// debug output
		slog.Debug("FilePath created", slog.String("folderASide", path), slog.String("folderBSide", otherFilePath))

		if !hashutil.CheckTwoFileHashSame(path, otherFilePath) {
			slog.Error("Failed to verify hash for this file.")
			if !ignoreFileErrors {
				return fmt.Errorf("Failed to verify file. pathA: %s, pathB: %s", path, otherFilePath)
			}
		}

		slog.Debug("File verified.", slog.String("pathA", path), slog.String("pathB", otherFilePath))

		return nil
	})
}
