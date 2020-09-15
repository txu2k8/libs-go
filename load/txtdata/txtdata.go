package txtdata

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"sync"

	"github.com/op/go-logging"
)

var (
	wg     sync.WaitGroup
	logger = logging.MustGetLogger("test")
)

// Set of characters to use for generating random strings
const (
	Alphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Numerals     = "1234567890"
	Alphanumeric = Alphabet + Numerals
	ASSIC        = Alphanumeric + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

// String returns a random string n characters long, composed of entities from charset.
func String(n int, charset string) (string, error) {
	randstr := make([]byte, n) // Random string to return
	charlen := big.NewInt(int64(len(charset)))
	for i := 0; i < n; i++ {
		b, err := rand.Int(rand.Reader, charlen)
		if err != nil {
			return "", err
		}
		r := int(b.Int64())
		randstr[i] = charset[r]
	}
	return string(randstr), nil
}

// RandomString returns a random ASSIC string n characters long.
func RandomString(n int) (string, error) {
	return String(n, ASSIC)
}

// CreateFile create original file, each line with line_number, and specified line size
// mode: w 只能写 覆盖整个文件 不存在则创建; a 只能写 从文件底部添加内容 不存在则创建
func CreateFile(filePath string, fileSize int, ch chan int) {
	fmt.Printf(">> Create/Write file: %s\n", filePath)
	flag := os.O_CREATE | os.O_RDWR | os.O_TRUNC
	file, err := os.OpenFile(filePath, flag, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	randomStr, _ := RandomString(fileSize)
	file.WriteString(randomStr)
	wg.Done()
	<-ch
}

// MultiCreateFile Multi Create File
func MultiCreateFile(files <-chan string, fileSize int) error {
	poolSize := runtime.NumCPU()
	runtime.GOMAXPROCS(poolSize)
	ch := make(chan int, poolSize)

	for filePath := range files {
		ch <- 1
		wg.Add(1)
		go CreateFile(filePath, fileSize, ch)
	}
	wg.Wait()
	close(ch)
	return nil
}
