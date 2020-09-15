package generator

import (
	"os"
	"path"
	"strconv"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("test")

// DirsGenerator Generate dirs based top path
func DirsGenerator(topPath string, dirCount int) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < dirCount; i++ {
			dirPath := path.Join(topPath, "dir_"+strconv.Itoa(i))
			logger.Info(dirPath)
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				logger.Error(err)
			}
			ch <- dirPath
		}
		close(ch)
	}()
	return ch
}

// FileNamesGenerator Generate file names based top path
// fileType: ".txt"
func FileNamesGenerator(topPath string, fileCount int, fileType string) <-chan string {
	prefix := "f_"
	ch := make(chan string)
	go func() {
		for i := 0; i < fileCount; i++ {
			fileName := prefix + strconv.Itoa(i) + fileType
			filePath := path.Join(topPath, fileName)
			// logger.Info(filePath)
			ch <- filePath
		}
		close(ch)
	}()
	return ch
}
