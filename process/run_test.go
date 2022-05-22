package process

import "testing"

func TestRunDownloadCCC(t *testing.T) {
	DownloadCCCImages("./test_file/ccc_nft_ic_test.json")
}

func TestRunownloadEXT(t *testing.T) {
	DownloadImages("../files/nft_image.xlsx")
}
