package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"uploadImage/download"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "uploadImage",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download nft images from ic",
	Long:  `download nft images from ic.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			panic(err)
		}
		download.DownloadImages(filePath)
	},
}

var downloadCCCCmd = &cobra.Command{
	Use:   "downloadCCC",
	Short: "download CCC nft images from ic",
	Long:  `download CCC nft images from ic.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("nft-file")
		if err != nil {
			panic(err)
		}
		download.DownloadCCCImages(filePath)
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload nft images to s3",
	Long:  `upload nft images to s3.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("images-file")
		if err != nil {
			panic(err)
		}
		bucket, err := cmd.Flags().GetString("bucket")
		if err != nil {
			panic(err)
		}
		_ = download.Upload(bucket, filePath)
	},
}

var retryCmd = &cobra.Command{
	Use:   "retry",
	Short: "retry fetch nft images from ic",
	Long:  `retry fetch nft images from ic.`,
	Run: func(cmd *cobra.Command, args []string) {
		errPath, err := cmd.Flags().GetString("err-file")
		if err != nil {
			panic(err)
		}
		download.Retry(errPath)
	},
}

var retryFromLogCmd = &cobra.Command{
	Use:   "retryFromLog",
	Short: "retry fetch nft images from ic by log",
	Long:  `retry fetch nft images from ic by log.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("nft-file")
		if err != nil {
			panic(err)
		}
		errPath, err := cmd.Flags().GetString("log-file")
		fmt.Println(errPath)
		if err != nil {
			panic(err)
		}
		download.RetryFromLogs(filePath, errPath)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("file", "f", "./files/nft_image.xlsx", "nft image info file path")

	rootCmd.AddCommand(retryFromLogCmd)
	retryFromLogCmd.Flags().StringP("nft-file", "n", "./files/nft_image.xlsx", "nft image info file path")
	retryFromLogCmd.Flags().StringP("log-file", "l", "err.log", "error log file path")

	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringP("images-file", "n", "./images", "nft images file path")
	uploadCmd.Flags().StringP("bucket", "b", "pyd-nft", "s3 bucket")

	rootCmd.AddCommand(retryCmd)
	retryCmd.Flags().StringP("err-file", "r", "err.json", "error log file path")

	rootCmd.AddCommand(downloadCCCCmd)
	downloadCCCCmd.Flags().StringP("nft-file", "n", "./files/ccc_nft.json", "nft image info file path")

}
