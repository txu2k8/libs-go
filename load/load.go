package main

// Package provide load million of bytes-size-files.

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
)

// log with level
var (
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题

	TopPath   string // Top path for create files
	DirCount  int    // Sub dir count
	FileCount int    // file count
	FileSize  int    // size unit: bytes

	wg sync.WaitGroup
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
	Trace.Printf(">> Create/Write file: %s", filePath)
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

// DirsGenerator Generate dirs based top path
func DirsGenerator(topPath string, dirCount int) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < dirCount; i++ {
			dirPath := path.Join(topPath, "dir_"+strconv.Itoa(i))
			Info.Print(dirPath)
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				Error.Print(err)
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
			// Trace.Print(filePath)
			ch <- filePath
		}
		close(ch)
	}()
	return ch
}

// ============== cmd ==============

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "load",
	Short: "Load million of bytes-size-files",
	Long:  `Load million of files, with topPath | fileSize | fileCount * dirCount"`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		Info.Printf("Start Load data: %s/<dirs-%d>/<files-%d>, Size: %d bytes", TopPath, DirCount, FileCount, FileSize)
		for dirPath := range DirsGenerator(TopPath, DirCount) {
			files := FileNamesGenerator(dirPath, FileCount, ".txt")
			MultiCreateFile(files, FileSize)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Log settings
	file, err := os.OpenFile("load.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(os.Stdout, //ioutil.Discard,
		"TRACE: ",
		log.Ldate|log.Ltime)

	Info = log.New(io.MultiWriter(file, os.Stdout),
		"INFO: ",
		log.Ldate|log.Ltime)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	// cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&TopPath, "top_path", "", "Top path for create files")
	rootCmd.PersistentFlags().IntVar(&DirCount, "dir_count", 3, "dir count")
	rootCmd.PersistentFlags().IntVar(&FileCount, "file_count", 10, "file count")
	rootCmd.PersistentFlags().IntVar(&FileSize, "file_size", 16, "each file size, unit:bytes")
	rootCmd.MarkPersistentFlagRequired("top_path")
}

func main() {
	Execute()
}
