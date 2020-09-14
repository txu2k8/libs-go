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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&TopPath, "top_path", "", "Top path for create files")
	rootCmd.PersistentFlags().IntVar(&DirCount, "dir_count", 3, "dir count")
	rootCmd.PersistentFlags().IntVar(&FileCount, "file_count", 10, "file count")
	rootCmd.PersistentFlags().IntVar(&FileSize, "file_size", 16, "each file size, unit:bytes")
	rootCmd.MarkPersistentFlagRequired("top_path")
}