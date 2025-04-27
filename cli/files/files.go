package files

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"gtithub.com/jgfranco17/lazyfile/cli/logging"
	"gtithub.com/jgfranco17/lazyfile/cli/outputs"
)

func ListDirectoryContents(entries []Entry, asTree bool) {
	logger := logging.NewLogger()

	fileCount := 0
	dirCount := 0
	totalByteSize := 0
	for idx, entry := range entries {
		modTime := outputs.ColorString(color.FgYellow, false, entry.ModTime.Format(time.RFC822))
		totalByteSize += int(entry.Size)

		var name string
		if entry.IsDir {
			dirName := entry.Name + "/"
			name = outputs.ColorString(color.FgBlue, true, dirName)
			dirCount += 1
		} else {
			name = entry.Name
			fileCount += 1
		}
		var prefix string
		if asTree {
			if idx == len(entries)-1 {
				prefix = fileTreeTerminal
			} else {
				prefix = fileTreeNonTerminal
			}
		}
		fileSize := convertBitsToBytesWithUnits(int(entry.Size))
		fileMode := colorModeStringByPermission(entry.Mode)
		fmt.Printf("%s  %8s  %s  %s%s\n", fileMode, fileSize, modTime, prefix, name)
	}
	totalFileSize := convertBitsToBytesWithUnits(totalByteSize)
	logger.Infof("Found %d directories, %d files (%s data)", dirCount, fileCount, totalFileSize)
}

func GetDirectoryContents(path string) ([]Entry, error) {
	var entries []Entry

	// Default to current dir
	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range dirEntries {
		fileInfo, err := entry.Info() // Need to call Info() for detailed metadata
		if err != nil {
			return nil, fmt.Errorf("failed to stat entry %q: %w", entry.Name(), err)
		}

		entry := Entry{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Size:    fileInfo.Size(),
			Mode:    fileInfo.Mode(),
			ModTime: fileInfo.ModTime(),
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func convertBitsToBytesWithUnits(byteSize int) string {
	var units string
	factor := 1
	switch len(strconv.Itoa(byteSize)) {
	case 4, 5, 6:
		factor = 1000
		units = "KB"
	case 7, 8, 9:
		factor = 1000000
		units = "MB"
	case 10, 11, 12:
		factor = 1000000000
		units = "GB"
	case 13, 14, 15:
		factor = 1000000000000
		units = "TB"
	default:
		units = "B"
	}
	convertedNumber := byteSize / factor
	return fmt.Sprintf("%d %s", convertedNumber, units)
}

func colorModeStringByPermission(mode fs.FileMode) string {
	coloredPermissions := ""
	for _, char := range mode.String() {
		stringChar := string(char)
		var colorToUse color.Attribute
		switch stringChar {
		case "w":
			colorToUse = color.FgRed
		case "r":
			colorToUse = color.FgYellow
		case "d":
			colorToUse = color.FgBlue
		case "x":
			colorToUse = color.FgMagenta
		}
		coloredPermissions += outputs.ColorString(colorToUse, true, stringChar)
	}
	return coloredPermissions
}
