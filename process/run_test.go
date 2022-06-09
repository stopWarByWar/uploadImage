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

func TestGetCCCImagesURL(t *testing.T) {
	GetCCCImagesURL("../files/ccc_nft.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/db2?charset=utf8&parseTime=true", true)
}
