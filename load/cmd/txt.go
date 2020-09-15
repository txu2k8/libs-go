package cmd

import (
	"libs-go/load/generator"
	"libs-go/load/txtdata"

	"github.com/spf13/cobra"
)

// TxtFileSize .
var TxtFileSize int // txt file size unit: bytes

// txtCmd represents the "Create txt files"command
var txtCmd = &cobra.Command{
	Use:   "txt",
	Short: "Write txt files with random string",
	Long:  `For create million of txt files under top_path`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Start Load data: %s/<dirs-%d>/<files-%d>, Size: %d bytes", TopPath, DirCount, FileCount, TxtFileSize)
		for dirPath := range generator.DirsGenerator(TopPath, DirCount) {
			files := generator.FileNamesGenerator(dirPath, FileCount, ".txt")
			txtdata.MultiCreateFile(files, TxtFileSize)
		}
	},
}

func init() {
	rootCmd.AddCommand(txtCmd)
	AddFlagsBase(txtCmd)
	txtCmd.PersistentFlags().IntVar(&TxtFileSize, "txt_file_size", 16, "each txt file size, unit:bytes")
}
