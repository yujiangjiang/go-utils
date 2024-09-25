package ufile

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
)

// CountLines count file lines
func CountLines(filePath string) int64 {
	file, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var lineCount int64
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lineCount++
	}
	return lineCount
}

// ReadDirsMap list sub dirs of target dir to a map
func ReadDirsMap(d string) map[string]string {

	var dirs = make(map[string]string)
	filepath.Walk(d, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && (path != d) {
			dirs[info.Name()] = path
		}
		return nil
	})
	return dirs
}
