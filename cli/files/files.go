package files

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gtithub.com/jgfranco17/lazyfile/cli/logging"
	"gtithub.com/jgfranco17/lazyfile/cli/outputs"
)

func ListDirectoryContents(entries []Entry, asTree bool) {
	logger := logging.NewLogger()

	fileCount := 0
	dirCount := 0
	totalByteSize := 0
	for idx, entry := range entries {
		modTime := entry.ModTime.Format(time.RFC822)
		size := fmt.Sprintf("%10d", entry.Size)
		totalByteSize += int(entry.Size) / 8

		var name string
		if entry.IsDir {
			dirName := entry.Name + "/"
			name = outputs.PrintColoredMessage("green", dirName)
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
		fmt.Printf("%s  %s  %s  %s%s\n", entry.Mode.String(), size, modTime, prefix, name)
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

func convertBitsToBytesWithUnits(bitCount int) string {
	bytes := bitCount / 8
	var convertedNumber int
	var units string
	switch len(strconv.Itoa(bytes)) {
	case 3, 4, 5:
		convertedNumber = bytes / 3
		units = "KB"
	case 6, 7, 8:
		convertedNumber = bytes / 6
		units = "MB"
	case 9, 10, 11:
		convertedNumber = bytes / 9
		units = "GB"
	case 12, 13, 14:
		convertedNumber = bytes / 12
		units = "TB"
	default:
		convertedNumber = bytes
		units = "B"
	}
	return fmt.Sprintf("%d %s", convertedNumber, units)
}
