package process

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"uploadImage/utils"
)

func TestDownload(t *testing.T) {
	c := new(http.Client)
	start := time.Now()
	info := utils.EXTNFTImageInfo{
		CanisterID:       "q6hjz-kyaaa-aaaah-qcama-cai",
		ImageUrlTemplate: "https://fp7fo-2aaaa-aaaaf-qaeea-cai.raw.ic0.app/Token/%d",
		Supply:           100,
		Standard:         "id",
		FileType:         "jpg",
	}
	t.Log(Download(c, info))
	t.Log(time.Now().Sub(start))
}

func TestDownloadAll(t *testing.T) {
	start := time.Now()
	c := new(http.Client)
	nftInfos, err := utils.GetEXTNFTImageInfos("./test_file/nft_image_test.xlsx")
	if err != nil {
		t.Fatalf(err.Error())
	}

	for name, info := range nftInfos {
		var url string
		if info.Standard == "identifier" {
			identifier, _ := utils.TokenId2TokenIdentifier(info.CanisterID, 1)
			url = fmt.Sprintf(info.ImageUrlTemplate, identifier)
		} else {
			url = fmt.Sprintf(info.ImageUrlTemplate, 1)
		}
		t.Log(name + ":" + url)
		data, err := DownloadSingle(c, url)
		if err != nil {
			errMsg := fmt.Sprintf("error:%v,url:%v", err, url)
			t.Log(errMsg)
			continue
		}
		path := fmt.Sprintf("images/%s.%s", name, info.FileType)
		t.Log(ioutil.WriteFile(path, data, 0644))
	}
	t.Log(time.Now().Sub(start))
}

func TestDownloadSingle(t *testing.T) {
	c := new(http.Client)
	var url = "https://ep54t-xiaaa-aaaah-qcdza-cai.raw.ic0.app/?type=thumbnail&tokenid=uhfgh-4ykor-uwiaa-aaaaa-b4aq6-iaqca-aaaaa-a"

	data, err := DownloadSingle(c, url)
	if err != nil {
		errMsg := fmt.Sprintf("error:%v,url:%v", err, url)
		t.Log(errMsg)
		return
	}
	t.Log(ioutil.WriteFile("image.gi", data, 0644))
}

//func TestName(t *testing.T) {
//	nftInfos, err := utils.GetEXTNFTImageInfos("./nft_image.xlsx")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	c := new(http.Client)
//	var errUrls []utils.ErrUrl
//	fi, err := os.Open("./log")
//	if err != nil {
//		fmt.Printf("Error: %s\n", err)
//		return
//	}
//	defer fi.Close()
//
//	br := bufio.NewReader(fi)
//	for {
//		a, _, c := br.ReadLine()
//		if c == io.EOF {
//			break
//		}
//
//		result := strings.Split(string(a), ",")
//		canisterID := strings.Split(result[1], ":")[1]
//		tokenID := strings.Split(result[2], ":")[1]
//		url := strings.Split(result[4], ":")[1] + ":" + strings.Split(result[4], ":")[2]
//		errCode := strings.Split(strings.Split(result[0], ":")[1], " ")[1]
//
//		if strings.Contains(errCode, "20") {
//			continue
//		}
//
//		if canisterID == "po6n2-uiaaa-aaaaj-qaiua-cai" {
//			url = "https://po6n2-uiaaa-aaaaj-qaiua-cai.raw.ic0.app/?tokenid="
//			id, _ := strconv.Atoi(tokenID)
//			identifier, _ := utils.TokenId2TokenIdentifier("po6n2-uiaaa-aaaaj-qaiua-cai", uint32(id))
//			url = url + identifier
//		}
//
//		if canisterID == "p5jg7-6aaaa-aaaah-qcolq-cai" {
//			id, _ := strconv.Atoi(tokenID)
//			if id > 9787 {
//				continue
//			}
//		}
//		t.Log(url, tokenID, canisterID, nftInfos[canisterID].FileType, errCode)
//		id, _ := strconv.Atoi(tokenID)
//		errUrls = append(errUrls, utils.ErrUrl{Url: url, TokenID: uint32(id), CanisterID: canisterID, Type: nftInfos[canisterID].FileType})
//	}
//	ReRequestUrlOld(c, errUrls)
//}

func TestDownloadCCCFromIC(t *testing.T) {
	var errUrls []utils.ErrUrl
	data, err := ioutil.ReadFile("./test_file/ccc_nft_ic_test.json")
	if err != nil {
		panic(err)
	}
	var nftInfos map[string]utils.CCCNFTImagesInfo
	if err := json.Unmarshal(data, &nftInfos); err != nil {
		panic(err)
	}

	c := new(http.Client)
	for id, info := range nftInfos {
		errUrl, err := DownloadCCCFromIC(c, info)
		if err != nil {
			fmt.Sprintf("can not download CCC images with error: %v,canisterID:%s", err, id)
			continue
		}
		errUrls = append(errUrls, errUrl...)
	}
	errURLData, err := json.Marshal(errUrls)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(ioutil.WriteFile("err.json", errURLData, 0644))
	fmt.Println("successfully download ")
}

func TestDownloadCCCFromIPFS(t *testing.T) {
	var errUrls []utils.ErrUrl
	data, err := ioutil.ReadFile("./test_file/ccc_nft_ifps_test.json")
	if err != nil {
		panic(err)
	}
	var nftInfos map[string]utils.CCCNFTImagesInfo
	if err := json.Unmarshal(data, &nftInfos); err != nil {
		panic(err)
	}

	c := new(http.Client)
	for _, info := range nftInfos {
		var singleNFTInfos []utils.CCCNFTInfo
		singleNFTInfos, err = utils.GetCCCNFTImageURL(info.CanisterID, info.FileType, info.ImageUrlTemplate, info.Type)
		if err != nil {
			t.Fatalf(err.Error())
		}
		errUrl, err := DownloadCCCFromIPFS(c, singleNFTInfos)
		if err != nil {
			fmt.Sprintf("can not download CCC images with error: %v,canisterID:%s", err, info.CanisterID)
			continue
		}
		errUrls = append(errUrls, errUrl...)
	}
	errURLData, err := json.Marshal(errUrls)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(ioutil.WriteFile("err.json", errURLData, 0644))
}

func TestDownloadCCCFromIPFS1(t *testing.T) {
	var errUrls []utils.ErrUrl
	data, err := ioutil.ReadFile("./test_file/ccc_nft_ifps_1_test.json")
	if err != nil {
		panic(err)
	}
	var nftInfos map[string]utils.CCCNFTImagesInfo
	if err := json.Unmarshal(data, &nftInfos); err != nil {
		panic(err)
	}

	c := new(http.Client)
	for _, info := range nftInfos {
		var singleNFTInfos []utils.CCCNFTInfo
		singleNFTInfos, err = utils.GetCCCNFTImageURL(info.CanisterID, info.FileType, info.ImageUrlTemplate, info.Type)
		if err != nil {
			t.Fatalf(err.Error())
		}
		errUrl, err := DownloadCCCFromIPFS(c, singleNFTInfos)
		if err != nil {
			fmt.Sprintf("can not download CCC images with error: %v,canisterID:%s", err, info.CanisterID)
			continue
		}
		errUrls = append(errUrls, errUrl...)
	}
	errURLData, err := json.Marshal(errUrls)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(ioutil.WriteFile("err.json", errURLData, 0644))
}

func TestReRequestUrl(t *testing.T) {
	var errUrls []utils.ErrUrl
	data, err := ioutil.ReadFile("err.json")
	if err != nil {
		t.Fatalf(err.Error())
	}
	_ = json.Unmarshal(data, &errUrls)
	newErrUrls, _ := ReRequestUrl(&http.Client{}, errUrls)

	_data, _ := json.Marshal(&newErrUrls)
	filePath := fmt.Sprintf("err.json")
	ioutil.WriteFile(filePath, _data, 0644)
}
