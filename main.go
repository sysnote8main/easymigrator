package main

import (
	"easymigrator/pkg/verifyutil"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	cp "github.com/otiai10/copy"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func main() {
	slog.Debug("Args", slog.String("args", strings.Join(os.Args, ", ")), slog.Int("len", len(os.Args)))

	// Validate args
	if len(os.Args) != 3 {
		slog.Info("./(exec) (fromPath) (toPath)")
		os.Exit(1)
	}

	// 0. Validate folder A&B pathes
	folderA, err := filepath.Abs(os.Args[1])
	if err != nil {
		slog.Error("Failed to convert first arg to absolute path", slog.Any("error", err))
		os.Exit(1)
	}
	folderB, err := filepath.Abs(os.Args[2])
	if err != nil {
		slog.Error("Failed to convert second arg to absolute path", slog.Any("error", err))
		os.Exit(1)
	}

	// 1. Copy all content in folderA to folderB
	slog.Info("Copying files", slog.String("from", folderA), slog.String("to", folderB))
	err = cp.Copy(folderA, folderB)
	if err != nil {
		slog.Error("Failed to copy", slog.Any("err", err))
		os.Exit(1)
	}
	slog.Info("Files are copied.")

	// folderA := "/Users/sysnote8/testingarea/a"
	// folderB := "/Users/sysnote8/testingarea/b"

	// 2. verify folders
	slog.Info("Verifying folders", slog.String("folderA", folderA), slog.String("folderB", folderB))
	err = verifyutil.VerifyFolderContents(folderA, folderB, false)
	if err != nil {
		slog.Error("Failed to verify folder", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Success to verify folder.")

	// 3. Delete directory (interactive)
	var confirmed bool = false
	if err = huh.NewConfirm().Title("Delete origin folder?").Affirmative("Yes!").Negative("No.").Value(&confirmed).Run(); err != nil {
		slog.Warn("Failed to show confirm ui. Continue with don't delete.")
	}

	// If confirmed, delete origin directory.
	if confirmed {
		slog.Info("Deleting folder & inside contents", slog.String("folderPath", folderA))
		// delete origin folder
		err = os.RemoveAll(folderA)
		if err != nil {
			slog.Error("Failed to remove directory", slog.Any("error", err), slog.String("folder", folderA))
			os.Exit(1)
		}
	} else {
		slog.Info("Delete directory & Create symlink will skip.")
	}

	// 4. create symlink folderA to folderB (if not confirmed, this action will skip.)
	if confirmed {
		slog.Info("Creating symlink", slog.String("from", folderA), slog.String("to", folderB))
		err = os.Symlink(folderB, folderA)
		if err != nil {
			slog.Error("Failed to create symlink", slog.Any("error", err), slog.String("from", folderA), slog.String("to", folderB))
			os.Exit(1)
		}
	}

	slog.Info("All actions completed. see you!")
}
