package lager

import (
	"fmt"
	"runtime"
	"strings"
)

// FileType represents values for including file information of logs
type FileType uint8

const (
	// NoFile does not include the file name or line number of the log
	NoFile FileType = iota
	// ShortFile includes the file name and line number of the log
	ShortFile
	// PackageFile includes the file name, including package, and line number of the log
	PackageFile
	// FullFile includes the absolute file path and line number of the log
	FullFile
)

// Caller returns the appropriate filename and line number of the file type
func (ft FileType) Caller(calldepth int) string {
	if ft == NoFile {
		return ""
	}

	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		return ""
	}

	if ft == PackageFile {
		paths := strings.Split(file, "src/")
		if len(paths) == 2 {
			file = paths[1]
		}
	} else if ft == ShortFile {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}

	return fmt.Sprintf("%s:%d", file, line)
}
