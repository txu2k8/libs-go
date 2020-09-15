package cmd

import (
	"libs-go/load/fiodata"
	"libs-go/load/generator"

	"github.com/spf13/cobra"
)

// FIO args .
var (
	FIOFileSize string // fio file size: 16k,1m,2G
	FIONumJobs  int    // fio numjobs=?
)

// txtCmd represents the "Create txt files"command
var fioCmd = &cobra.Command{
	Use:   "fio",
	Short: "FIO-Write files random ",
	Long:  `For generate million of txt files under top_path by fio tool`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Start Load data: %s/<dirs-%d>/<files-%d>, Size: %d bytes", TopPath, DirCount, FileCount, FIOFileSize)
		for dirPath := range generator.DirsGenerator(TopPath, DirCount) {
			files := generator.FileNamesGenerator(dirPath, FileCount, ".data")
			fiodata.MultiFIOWrite(files, FIOFileSize, FIONumJobs)
		}
	},
}

func init() {
	rootCmd.AddCommand(fioCmd)
	AddFlagsBase(fioCmd)
	fioCmd.PersistentFlags().StringVar(&FIOFileSize, "fio_file_size", "16k", "each fio file size")
	fioCmd.PersistentFlags().IntVar(&FIONumJobs, "fio_numjobs", 1, "fio numjobs=?")
}
