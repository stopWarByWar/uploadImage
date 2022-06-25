package utils

import (
	"encoding/json"
	mapSet "github.com/deckarep/golang-set"
	agent "github.com/mix-labs/IC-Go"
	"github.com/mix-labs/IC-Go/utils"
	"github.com/mix-labs/IC-Go/utils/idl"
	"io/ioutil"
	"testing"
)

func TestID(t *testing.T) {
	identifier, _ := TokenId2TokenIdentifier("5movr-diaaa-aaaak-aaftq-cai", 0)
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
