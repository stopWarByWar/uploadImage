package process

import "testing"

func TestRunDownloadCCC(t *testing.T) {
	DownloadCCCImages("../files/ccc_nft.json")
}

func TestRunownloadEXT(t *testing.T) {
	DownloadImages("./test_file/nft_image_test.xlsx")
}

func TestCompassGif(t *testing.T) {
	compassGif("../images/4fcza-biaaa-aaaah-abi4q-cai/0.gif")
}
