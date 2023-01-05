package process

import (
	"database/sql"
	"fmt"
	agent "github.com/mix-labs/IC-Go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"testing"
	"uploadImage/utils"
)

func TestRunDownloadCCC(t *testing.T) {
	DownloadCCCImages("../files/ccc_nft.json")
}

func TestDownloadEXT(t *testing.T) {
	DownloadImages("../files/nft_image.xlsx")
}

func TestCompassGif(t *testing.T) {
	_ = compassGif("../images/4fcza-biaaa-aaaah-abi4q-cai/0.gif")
}

func TestGetCCCImagesURL(t *testing.T) {
	//GetCCCImagesURL("../files/ccc_nft.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/db2?charset=utf8&parseTime=true", true)
	//GetCCCImagesURL("../files/ccc_nft.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/nft?charset=utf8&parseTime=true", true)
	//GetCCCImagesURL("../files/ccc_nft.json", "root:xyz12345@(localhost:3306)/xyz?charset=utf8&parseTime=true&parseTime=true", true)
}

func TestGetDip721ImagesURL(t *testing.T) {
	//GetDip721ImageUrl("../files/dip721.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/nft?charset=utf8&parseTime=true", true)
	GetDip721ImageUrl("../files/dip721.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/db2?charset=utf8&parseTime=true", true)
	//GetDip721ImageUrl("../files/dip721.json", "root:xyz12345@(localhost:3306)/xyz?charset=utf8&parseTime=true&parseTime=true", true)
}

func TestGetYumiImagesURL(t *testing.T) {
	GetYumiImagesUrl("../files/yumi.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/nft?charset=utf8&parseTime=true", true)
	//GetYumiImagesUrl("../files/yumi.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/db2?charset=utf8&parseTime=true", true)
}

func TestGetAstroxImagesURL(t *testing.T) {
	//GetAstroXImagesUrl("../files/astrox.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/nft?charset=utf8&parseTime=true", true)
	//GetAstroXImagesUrl("../files/astrox.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/db2?charset=utf8&parseTime=true", true)
}

func TestGetEntrepotsUrls(t *testing.T) {
	//GetEntrepotsUrls("../files/entrepot.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/nft?charset=utf8&parseTime=true", true)
	//GetEntrepotsUrls("../files/entrepot.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/db2?charset=utf8&parseTime=true", false)
	//GetEntrepotsUrls("../files/entrepot.json", "root:xyz12345@(localhost:3306)/xyz?charset=utf8&parseTime=true&parseTime=true", true)
}

func TestGetICPSwapUrls(t *testing.T) {
	//SetICPSwapImageUrl("../files/icpswap.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/nft?charset=utf8&parseTime=true", true)
	//SetICPSwapImageUrl("../files/icpswap.json", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/db2?charset=utf8&parseTime=true", true)
}

func TestName2(t *testing.T) {

	_db, err := sql.Open("mysql", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/nft?charset=utf8&parseTime=true")
	var result []utils.NFTUrl
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: _db}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10000; i++ {
		result = append(result, utils.NFTUrl{
			CanisterID: "kafas-uaaaa-aaaao-aaofq-cai",
			TokenID:    uint32(i),
			ImageUrl:   fmt.Sprintf("https://7gayh-aqaaa-aaaah-abdgq-cai.raw.ic0.app/%d", i),
			VideoUrl:   "",
		})
	}
	t.Log(db.Save(&result).Error)
}

func TestDemail(t *testing.T) {
	_db, err := sql.Open("mysql", "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/nft?charset=utf8&parseTime=true")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: _db}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	_agent := agent.New(true, "")
	canisterID := "qohn4-fiaaa-aaaak-ady6a-cai"
	if err = utils.SaveDMailUrls(db, _agent, canisterID); err != nil {
		t.Log(err)
	}
}
