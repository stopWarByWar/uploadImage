package utils

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	mapSet "github.com/deckarep/golang-set"
	agent "github.com/mix-labs/IC-Go"
	"github.com/mix-labs/IC-Go/utils"
	"github.com/mix-labs/IC-Go/utils/idl"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestID(t *testing.T) {
	identifier, _ := TokenId2TokenIdentifier("3sxmo-naaaa-aaaao-aakqa-cai", 12)
	t.Log(identifier)
}

func TestCompressImageResource(t *testing.T) {
	data, err := ioutil.ReadFile("1.png")
	if err != nil {
		panic(err)
	}
	result := CompressImageResource(data)
	ioutil.WriteFile("11.png", result, 0644)

}

func TestSetCCCNftInfo(t *testing.T) {
	info := make(map[string]CCCNFTImagesInfo)
	emptyImageTemplate := CCCImageUrlTemplate{}
	info["empty"] = CCCNFTImagesInfo{
		ImageUrlTemplates: []CCCImageUrlTemplate{emptyImageTemplate},
	}

	data, _ := json.Marshal(info)
	ioutil.WriteFile("../download/ccc_nft.json", data, 0644)
}

func TestGetCCCNFTImageURL(t *testing.T) {
	tokenID := mapSet.NewSet()
	result, err := GetCCCNFTImageURL("o7ehd-5qaaa-aaaah-qc2zq-cai", "jpeg", "https://gateway.filedrive.io/ipfs/%s", "ipfs-3")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(len(result))
	var num int
	for _, info := range result {
		tokenID.Add(info.TokenID)
		t.Log(num, info.ImageUrl, info.TokenID)
		num++
	}
	t.Log(tokenID.Cardinality())

}

func TestListFiles(t *testing.T) {
	_agent := agent.New(true, "")
	canisterID := "ml2cx-yqaaa-aaaah-qc2xq-cai"
	methodNameSupply := "getSuppy"
	arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
	_, result, _, err := _agent.Query(canisterID, methodNameSupply, arg)
	if err != nil {
		t.Fatal(err.Error())
	}
	var supply uint64
	utils.Decode(&supply, result[0])
	t.Log(supply)
}

func Test_getEntrepotUrl(t *testing.T) {
	t.Log(getEntrepotUrl(new(http.Client), "ljt3c-tyaaa-aaaam-qa53a-cai", 3036))
}

func Test_SetICPSwapUrls(t *testing.T) {
	dns := "admin:Gbs1767359487@(database-mysql-instance-1.ccggmi9astti.us-east-1.rds.amazonaws.com:3306)/db2?charset=utf8&parseTime=true"
	_db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: _db}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		panic(err)
	}
	_agent := agent.New(true, "")

	SetICPSwapUrls(db, _agent, "ewuhj-dqaaa-aaaag-qahgq-cai")
}
func TestGetTokens(t *testing.T) {
	_agent := agent.New(true, "")
	canisterID := "wettq-dqaaa-aaaag-qapua-cai"
	tokens, err := getTokens(_agent, canisterID)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(tokens, len(tokens))
}

func TestName(t *testing.T) {
	t.Log(hex.EncodeToString([]byte("Monstea's")))
	//r, _ := hex.DecodeString("6170746f735f636f696e")
	//t.Log(string(r))
}

func Test_GetIcsMetadata(t *testing.T) {
	_agent := agent.New(true, "")
	canisterID := "wettq-dqaaa-aaaag-qapua-cai"
	t.Log(GetIcsMetadata(_agent, canisterID, 22))
}

func TestGetEXTMetaData(t *testing.T) {

}
