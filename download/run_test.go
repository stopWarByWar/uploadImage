package download

import "testing"

func TestRunDownloadCCC(t *testing.T) {
	DownloadCCCImages("./test_file/ccc_nft_ic_test.json")
}
