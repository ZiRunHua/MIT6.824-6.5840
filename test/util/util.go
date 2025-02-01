package util

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func ClickablePath(filePath string) string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error getting absolute path:%v", err))
	}

	switch runtime.GOOS {
	case "windows":
		return absPath
	default:
		return "file://" + filepath.ToSlash(absPath)
	}
}
