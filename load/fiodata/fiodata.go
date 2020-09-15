package fiodata

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"sync"

	"github.com/op/go-logging"
)

var (
	wg     sync.WaitGroup
	logger = logging.MustGetLogger("test")
)

// FIOWrite .
func FIOWrite(filePathName string, fileSize string, numJobs int, ch chan int) {
	fmt.Printf(">> FIO-Write file: %s\n", filePathName)
	filePath, fileName := path.Split(filePathName)
	command := fmt.Sprintf("fio --randrepeat=0 --verify=0 --ioengine=libaio --direct=1 --gtod_reduce=1 --rw=randrw --rwmixread=0 --refill_buffers --group_reporting --norandommap --bs=16K --iodepth=1 --directory=%s --name=%s --size=%s --numjobs=%d", filePath, fileName, fileSize, numJobs)

	cmd := exec.Command("/bin/bash", "-c", command)
	_, err := cmd.Output()
	if err != nil {
		logger.Error(err)
	}
	wg.Done()
	<-ch
}

// MultiFIOWrite Multi FIO Write File
func MultiFIOWrite(files <-chan string, fileSize string, numJobs int) error {
	poolSize := runtime.NumCPU()
	runtime.GOMAXPROCS(poolSize)
	ch := make(chan int, poolSize)

	for filePathName := range files {
		ch <- 1
		wg.Add(1)
		go FIOWrite(filePathName, fileSize, numJobs, ch)
	}
	wg.Wait()
	close(ch)
	return nil
}
