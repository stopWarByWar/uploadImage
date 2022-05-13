package download

import "testing"

func TestRunDownloadCCC(t *testing.T) {
	DownloadCCCImages("../files/ccc_nft.json")
}
