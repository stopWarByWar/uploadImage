package process

import "testing"

func TestRunDownloadCCC(t *testing.T) {
	DownloadCCCImages("./test_file/ccc_nft_ic_test.json")
}

func TestRunownloadEXT(t *testing.T) {
	DownloadImages("./test_file/nft_image_test.xlsx")

}
