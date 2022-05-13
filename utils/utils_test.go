package utils

import (
	"encoding/json"
	mapSet "github.com/deckarep/golang-set"
	"io/ioutil"
	"testing"
)

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
	result, err := GetCCCNFTImageURL("k3wif-yyaaa-aaaah-qc2ca-cai", "png", "https://gateway.filedrive.io/ipfs/%s", "ipfs-1")
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
	ListFiles("../download/images")
}
