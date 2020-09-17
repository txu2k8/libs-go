package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/chenhg5/collection"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
)

// log with level
var (
	logger = logging.MustGetLogger("test")
	debug  bool // debug modle

	TopPath   string // Top path for create files
	DirCount  int    // Sub dir count
	FileCount int    // file count
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "load",
	Short: "Load data",
	Long:  `Load million of files, with TopPath | DirCount | FileCount | FileSize"`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Start Load data: %s/<dirs-%d>/<files-%d>", TopPath, DirCount, FileCount)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	initLogging()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// AddFlagsBase .
func AddFlagsBase(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&TopPath, "top_path", "", "Top path for create files")
	cmd.PersistentFlags().IntVar(&DirCount, "dir_count", 3, "dir count")
	cmd.PersistentFlags().IntVar(&FileCount, "file_count", 10, "file count")
	cmd.MarkPersistentFlagRequired("top_path")
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable console debug log level if true (default false)")
}

func initLogging() {
	dir, _ := os.Getwd()
	fileLogName := "load"
	fileLogPath := path.Join(dir, "log")
	timeStr := time.Now().Format("20060102150405")
	for _, v := range stripArgs() {
		fileLogName = fmt.Sprintf("%s-%s", fileLogName, v)
	}
	fileLogName = fmt.Sprintf("%s-%s.log", fileLogName, timeStr)
	fileLogPathName := path.Join(fileLogPath, fileLogName)
	consoleLoglevel := logging.INFO
	if collection.Collect(os.Args).Contains("--debug") == true {
		consoleLoglevel = logging.DEBUG
	}

	// backend output to log file && Console
	fileStrformat := `%{time:2006-01-02T15:04:05} %{module} %{level:.4s}: (%{shortfile}) %{message}`
	fileFormat := logging.MustStringFormatter(fileStrformat)
	err := os.MkdirAll(path.Dir(fileLogPathName), os.ModePerm)
	file, err := os.OpenFile(fileLogPathName,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	fileBackend := logging.NewLogBackend(io.Writer(file), "", 0)
	fileBackendFormator := logging.NewBackendFormatter(fileBackend, fileFormat)
	fileBackendLeveled := logging.AddModuleLevel(fileBackendFormator)
	fileBackendLeveled.SetLevel(consoleLoglevel, "")

	// Set the backends to be used.
	logging.SetBackend(fileBackendLeveled)

	testCommand := "load " + strings.Join(os.Args[1:], " ")
	logger.Infof("Args: %s", testCommand)
}

func stripArgs() []string {
	commands := []string{}
	args := os.Args[1:]
	ps := ""
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		switch {
		case s == "--":
			// "--" terminates the flags
			break
		case strings.HasPrefix(s, "--") && !strings.Contains(s, "="):
			// If '--flag arg' then
			// delete arg from args.
			fallthrough // (do the same as below)
		case strings.HasPrefix(s, "-") && !strings.Contains(s, "=") && len(s) == 2:
			// If '-f arg' then
			// delete 'arg' from args or break the loop if len(args) <= 1.
			if len(args) <= 1 {
				break
			} else {
				args = args[1:]
				continue
			}
		case s != "" && !strings.HasPrefix(s, "-") && (!strings.HasPrefix(ps, "-") || ps == "--case"):
			commands = append(commands, s)
		}
		ps = s
	}

	return commands
}
