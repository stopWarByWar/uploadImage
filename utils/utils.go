package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	agent "github.com/mix-labs/IC-Go"
	"github.com/mix-labs/IC-Go/utils"
	"github.com/mix-labs/IC-Go/utils/idl"
	"github.com/mix-labs/IC-Go/utils/principal"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const goroutineNum = 100

//const goroutineNum = 3000000000

func TokenId2TokenIdentifier(canisterID string, tokenID uint32) (string, error) {
	var tokenIdentifier string
	canisterBlob, err := principal.Decode(canisterID)
	if err != nil {
		return "", err
	}
	idBlob := make([]byte, 4)
	binary.BigEndian.PutUint32(idBlob, tokenID)

	data := []uint8{10, 116, 105, 100}
	data = append(data, canisterBlob...)
	data = append(data, idBlob...)
	tokenIdentifier = principal.New(data).Encode()
	return tokenIdentifier, err
}

func GetEXTNFTImageInfos(path string) (map[string]EXTNFTImageInfo, error) {
	var nftInfos = make(map[string]EXTNFTImageInfo)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows("Sheet1")

	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < EntrepotNFTInfoRowLen {
			errMsg := fmt.Sprintf("row:%d,the amount of element of nft in entrepot should be %d, but %d", i, EntrepotNFTInfoRowLen, len(rows[i]))
			fmt.Println(errMsg)
			break
		}
		//if rows[i][0] == "" {
		//	fmt.Println(rows[i][1])
		//}
		if rows[i][3] == "N" {
			continue
		}
		supply, _ := strconv.Atoi(rows[i][6])
		nftInfos[rows[i][0]] = EXTNFTImageInfo{
			CanisterID:       rows[i][0],
			Name:             rows[i][1],
			ImageUrlTemplate: rows[i][2],
			Standard:         rows[i][5],
			Supply:           supply,
			FileType:         rows[i][7],
			StorageCanister:  rows[i][8],
		}
	}
	return nftInfos, nil
}

//func PathExists(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//	return false, err
//}

func CompressImageResource(data []byte) []byte {
	imgSrc, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	newImg := image.NewRGBA(imgSrc.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)
	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 40})
	if err != nil {
		return data
	}
	if buf.Len() > len(data) {
		return data
	}
	return buf.Bytes()
}

func GetCCCNFTImageURL(canisterID string, fileType string, imageUrlTemplate string, types string) ([]CCCNFTInfo, error) {
	var infos []CCCNFTInfo
	_agent := agent.New(true, "")
	switch types {

	case "ipfs":
		nftUrls := make(map[uint64]CCCNFTInfo)
		methodName := "getAllNftPhotoLink"
		arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
		_, result, _, err := _agent.Query(canisterID, methodName, arg)
		if err != nil {
			return nil, err
		}
		var myResult []NFTPhotoLink
		utils.Decode(&myResult, result[0])
		for _, token := range myResult {
			nftUrls[token.TokenIndex] = CCCNFTInfo{
				CanisterID:    canisterID,
				TokenID:       token.TokenIndex,
				ImageUrl:      token.NFTLinkInfo,
				ImageFileType: fileType,
			}
		}
		methodName = "getAllNftVideoLink"
		arg, _ = idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
		_, result, _, err = _agent.Query(canisterID, methodName, arg)
		if err != nil {
			return nil, err
		}
		utils.Decode(&myResult, result[0])
		for _, token := range myResult {
			midUrl := nftUrls[token.TokenIndex]
			midUrl.VideoUrl = token.NFTLinkInfo
			nftUrls[token.TokenIndex] = midUrl
		}
		for _, url := range nftUrls {
			infos = append(infos, url)
		}
	case "ipfs-1":
		methodName := "getAllNftLinkInfo"
		arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
		_, result, _, err := _agent.Query(canisterID, methodName, arg)
		if err != nil {
			return nil, err
		}
		var myResult []TokenImageUrl
		utils.Decode(&myResult, result[0])
		for _, token := range myResult {
			imageUrl := fmt.Sprintf(imageUrlTemplate, token.NFTLinkInfo.PhotoLink)
			videoUrl := fmt.Sprintf(imageUrlTemplate, token.NFTLinkInfo.VideoLink)
			infos = append(infos, CCCNFTInfo{
				CanisterID:    canisterID,
				TokenID:       token.TokenIndex,
				ImageUrl:      imageUrl,
				VideoUrl:      videoUrl,
				ImageFileType: fileType,
			})
		}
	case "ipfs-2":
		nftUrls := make(map[uint64]CCCNFTInfo)
		methodName := "getAllNftPhotoLink"
		arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
		_, result, _, err := _agent.Query(canisterID, methodName, arg)
		if err != nil {
			return nil, err
		}
		var myResult []NFTPhotoLink
		utils.Decode(&myResult, result[0])

		for _, token := range myResult {
			imageUrl := fmt.Sprintf(imageUrlTemplate, token.NFTLinkInfo)
			nftUrls[token.TokenIndex] = CCCNFTInfo{
				CanisterID:    canisterID,
				TokenID:       token.TokenIndex,
				ImageUrl:      imageUrl,
				ImageFileType: fileType,
			}
		}
		methodName = "getAllNftVideoLink"
		arg, _ = idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
		_, result, _, err = _agent.Query(canisterID, methodName, arg)
		if err != nil {
			return nil, err
		}
		utils.Decode(&myResult, result[0])
		for _, token := range myResult {
			midUrl := nftUrls[token.TokenIndex]
			midUrl.VideoUrl = fmt.Sprintf(imageUrlTemplate, token.NFTLinkInfo)
			nftUrls[token.TokenIndex] = midUrl
		}
		for _, url := range nftUrls {
			infos = append(infos, url)
		}
	case "ipfs-3":
		methodNameSupply := "getSuppy"
		arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
		_, result, _, err := _agent.Query(canisterID, methodNameSupply, arg)
		if err != nil {
			return nil, err
		}
		var supply uint64
		utils.Decode(&supply, result[0])
		fmt.Println(supply)
		methodName := "getLinkInfoByIndexArr"
		input := idl.NewVec(new(idl.Nat))
		var inputValue []interface{}
		for i := 0; i < int(supply); i++ {
			inputValue = append(inputValue, big.NewInt(int64(i)))
		}

		param, err := idl.Encode([]idl.Type{input}, []interface{}{inputValue})
		if err != nil {
			panic(err)
		}
		_, result, _, err = _agent.Query(canisterID, methodName, param)

		var myResult []IPFS3ImageLink
		utils.Decode(&myResult, result[0])
		for i, token := range myResult {
			//fmt.Println(token.ImageLink.Some)
			imageUrl := fmt.Sprintf(imageUrlTemplate, token.ImageLink.Some)
			imageUrl = strings.ReplaceAll(imageUrl, "."+fileType, "")
			videoUrl := fmt.Sprintf(imageUrlTemplate, token.VideoLink.Some)
			videoUrl = strings.ReplaceAll(videoUrl, "."+fileType, "")
			infos = append(infos, CCCNFTInfo{
				CanisterID:    canisterID,
				TokenID:       uint64(i),
				ImageUrl:      imageUrl,
				VideoUrl:      videoUrl,
				ImageFileType: fileType,
			})
		}
	}
	return infos, nil
}

//func aListFiles(root string) ([]string, error) {
//	var result []string
//	files, err := ioutil.ReadDir(root)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, file := range files {
//		result = append(result, file.Name())
//	}
//
//	return result, nil
//}

func GetDip721TokenMetadata(_agent *agent.Agent, canisterID string, tokenId int) (string, string, error) {
	arg, _ := idl.Encode([]idl.Type{new(idl.Nat)}, []interface{}{big.NewInt(int64(tokenId))})
	_, result, errMsg, err := _agent.Query(canisterID, "tokenMetadata", arg)
	if err != nil || errMsg != "" {
		return "", "", fmt.Errorf("can not get token metadata with ID %s: %d, error: %v, errMsg: %v\n", canisterID, tokenId, err, errMsg)
	}
	var myResult ManualReply_2
	utils.Decode(&myResult, result[0])
	if myResult.EnumIndex != "Ok" {
		fmt.Println("result is not Ok")
		return "", "", nil
	}
	detailMap := map[string]GenericValue{}
	for _, v := range myResult.Ok.Properties {
		detailMap[v.Text] = v.GenericValue
	}
	var imageUrl string
	var videoUrl string
	if url, ok := detailMap["thumbnail"]; ok {
		imageUrl = url.TextContent
	} else {
		panic("no thumbnail")
	}

	if url, ok := detailMap["location"]; ok {
		videoUrl = url.TextContent
	} else {
		panic("no location")
	}
	return imageUrl, videoUrl, nil
}

func GetYumiNFTImageUrl(canisterID string) ([]NFTUrl, error) {
	var infos []NFTUrl
	_agent := agent.New(true, "")
	methodName := "getTokens"
	arg, err := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
	if err != nil {
		return nil, err
	}
	_, result, _, err := _agent.Query(canisterID, methodName, arg)
	if err != nil && err.Error() == "() empty" {
		return nil, nil
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var myResult []MetaData
	utils.Decode(&myResult, result[0])
	for _, token := range myResult {
		metadata := token.MetaData.NonFungible.MetaData.Some
		if metadata == nil {
			//fmt.Printf("metadata of %s#%d is empty\n", canisterID, token.TokenIndex)
			continue
		}

		urlInfo := new(YumUrl)
		if err = json.Unmarshal(metadata, urlInfo); err != nil {
			fmt.Printf("can not unmarshal url of %s#%d with error: %v", canisterID, token.TokenIndex, err)
			continue
		}

		var videoUrl string

		url := urlInfo.Url
		if urlInfo.Thumb != "" {
			url = urlInfo.Thumb
		}

		if urlInfo.MimeType == "video" {
			videoUrl = urlInfo.Url
		}

		infos = append(infos, NFTUrl{
			CanisterID: canisterID,
			TokenID:    token.TokenIndex,
			VideoUrl:   videoUrl,
			ImageUrl:   url,
		})
		fmt.Println(videoUrl, url)
	}
	return infos, nil
}

func SetICPSwapUrls(db *gorm.DB, _agent *agent.Agent, canisterID string) {
	tokenIds, err := getTokens(_agent, canisterID)
	if err != nil {
		panic(err)
	}

	//var idSet = mapset.NewSet()
	//for _, id := range tokenIds {
	//	idSet.Add(id.TokenIndex)
	//}
	//fmt.Println(len(tokenIds), idSet.Cardinality())

	singleGoroutineItemNum := len(tokenIds) / goroutineNum
	if singleGoroutineItemNum < 3 {
		setICPSwapUrls(db, _agent, canisterID, tokenIds)
		return
	}

	var wg sync.WaitGroup
	for i := 0; i < goroutineNum; i++ {
		wg.Add(1)
		if i != goroutineNum-1 {
			go func(i int) {
				setICPSwapUrls(db, _agent, canisterID, tokenIds[i*singleGoroutineItemNum:(i+1)*singleGoroutineItemNum])
				defer wg.Done()
			}(i)
		} else {
			go func(i int) {
				setICPSwapUrls(db, _agent, canisterID, tokenIds[i*singleGoroutineItemNum:])
				defer wg.Done()
			}(i)
		}
	}
	wg.Wait()
}

//type blobOpt struct {
//	Some []uint8 `ic:"some"`
//	None uint8   `ic:"none"`
//}

//type fungible struct {
//	Decimals     uint8   `ic:"decimals"`
//	Metadata     blobOpt `ic:"metadata"`
//	Name         string  `ic:"name"`
//	OwnerAccount string  `ic:"ownerAccount"`
//	Symbol       string  `ic:"symbol"`
//}

type tokens struct {
	TokenIndex uint32 `ic:"0"`
}

type IcsMetadata struct {
	FilePath string `ic:"filePath"`
	FileType string `ic:"fileType"`
	Link     string `ic:"link"`
}

type ResponseResult struct {
	Err   string      `ic:"err"`
	Ok    IcsMetadata `ic:"ok"`
	Index string      `ic:"EnumIndex"`
}

func setICPSwapUrls(db *gorm.DB, _agent *agent.Agent, canisterID string, tokens []tokens) {
	if len(tokens) == 0 {
		return
	}
	var urls []*NFTUrl
	for _, token := range tokens {
		url, err := GetIcsMetadata(_agent, canisterID, token.TokenIndex)
		if err != nil {
			fmt.Printf("\u001B %s:%d can not get metaData with error:%v\n", canisterID, token.TokenIndex, err)
			continue
		}
		urls = append(urls, url)
		fmt.Printf("canister:%s,id:%d,videoUrl:%s,imageUrl:%s\n", canisterID, token.TokenIndex, url.VideoUrl, url.ImageUrl)
	}
	if err := db.Save(&urls).Error; err != nil {
		fmt.Printf("\u001B %s can not set urls into db with error:%v\n", canisterID, err)
	}
}

func getTokens(_agent *agent.Agent, canisterID string) ([]tokens, error) {
	methodName := "getTokens"
	arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})

	_, result, _, err := _agent.Query(canisterID, methodName, arg)

	if err != nil {
		return nil, err
	}

	var myResult []tokens
	utils.Decode(&myResult, result[0])
	return myResult, nil
}

func GetIcsMetadata(_agent *agent.Agent, canisterID string, id uint32) (*NFTUrl, error) {
	methodName := "icsMetadata"
	arg, _ := idl.Encode([]idl.Type{idl.Nat32()}, []interface{}{big.NewInt(int64(id))})
	_, result, _, err := _agent.Query(canisterID, methodName, arg)
	if err != nil {
		return nil, err
	}

	var myResult ResponseResult
	utils.Decode(&myResult, result[0])
	if myResult.Err != "" {
		return nil, fmt.Errorf(myResult.Err)
	}
	switch myResult.Ok.FileType {
	case "image":
		return &NFTUrl{
			CanisterID: canisterID,
			TokenID:    id,
			ImageUrl:   myResult.Ok.FilePath,
		}, nil
	default:
		return &NFTUrl{
			CanisterID: canisterID,
			TokenID:    id,
			VideoUrl:   myResult.Ok.FilePath,
		}, nil
	}
}

func GetEntrepotUrls(db *gorm.DB, c *http.Client, canisterId string, from, to int, types string) {
	var urls []NFTUrl
	for i := from; i < to; i++ {
		url, err := getEntrepotUrl(c, canisterId, i)
		if err != nil {
			fmt.Printf("%s:%d can not get url with error:%v\n", canisterId, i, err)
			//if len(urls) == 0 {
			//	return
			//} else if err := db.Save(&urls).Error; err != nil {
			//	fmt.Printf("%s can not set urls into db with error:%v\n", canisterId, err)
			//}
			//panic(err)
		} else {
			//fmt.Printf("%s:%d urls:%s\n", canisterId, i, url)
			switch types {
			case "image":
				urls = append(urls, NFTUrl{
					CanisterID: canisterId,
					TokenID:    uint32(i),
					ImageUrl:   url,
				})
			case "video":
				urls = append(urls, NFTUrl{
					CanisterID: canisterId,
					TokenID:    uint32(i),
					VideoUrl:   url,
				})

			}
		}
	}
	if len(urls) == 0 {
		return
	} else if err := db.Save(&urls).Error; err != nil {
		fmt.Printf("%s can not set urls into db with error:%v\n", canisterId, err)
	}
}

func getEntrepotUrl(c *http.Client, canisterId string, tokenId int) (string, error) {
	tokenIdentity, err := TokenId2TokenIdentifier(canisterId, uint32(tokenId))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.raw.ic0.app/?tokenid=%s", canisterId, tokenIdentity)
	resp, err := c.Get(url)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content := string(data)
	if strings.Contains(content, "Error") {
		return "", errors.New(content)
	} else if strings.Contains(content, "image href") {
		url := entrepotDealWith1(content)
		return url, nil
	} else if strings.Contains(content, "URL") {
		url := entrepotDealWith2(content)
		return url, nil
	} else if strings.Contains(content, "video/mp4") && strings.Contains(content, "source src") {
		url := entrepotDealWith3(content)
		return url, nil
	}
	return "", fmt.Errorf("not match body:%s,url:%s", string(data), url)
}

func entrepotDealWith1(data string) string {
	datas := strings.Split(data, "<")
	var data1 string
	for _, item := range datas {
		if strings.Contains(item, "image href") {
			data1 = item
		}
	}
	data2 := strings.Split(data1, "\"")
	if len(data2) < 2 {
		panic(fmt.Sprintf("data2:%s,data:%s", data2, data))
	}
	return data2[1]
}

func entrepotDealWith2(data string) string {
	//fmt.Println(data)
	result := strings.Fields(data)
	for _, v := range result {
		if strings.Contains(v, "URL") {
			//fmt.Println(v)
			b, _, _ := strings.Cut(v, "\"")
			//fmt.Printf("before:%s,after:%s,is:%v\n", b, v, is)
			_, a, _ := strings.Cut(b, "URL=")
			//fmt.Printf("before:%s,after:%s,is:%v\n", b, v, is)
			return a
		}
	}
	return ""
}

func entrepotDealWith3(data string) string {
	//fmt.Println(data)
	result := strings.Fields(data)
	for _, v := range result {
		if strings.Contains(v, "src") {
			//fmt.Println(v)
			_, b, _ := strings.Cut(v, "\"")
			//fmt.Printf("before:%s,after:%s,is:%v\n", b, v, is)
			a, b, _ := strings.Cut(b, "\"")
			//fmt.Printf("before:%s,after:%s,is:%v\n", b, v, is)
			return a
		}
	}
	return ""
}

//func getAstroxMetadata(_agent *agent.Agent, canisterID string, tokenId uint32) (string, error) {
//	tokenIdentifier, err := TokenId2TokenIdentifier(canisterID, tokenId)
//	if err != nil {
//		return "", err
//	}
//	methodName := "metadata"
//	arg, _ := idl.Encode([]idl.Type{new(idl.Text)}, []interface{}{tokenIdentifier})
//	_, result, err := _agent.Update(canisterID, methodName, arg, 10)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(result)
//	var myResult result_2
//	utils.Decode(&myResult, result[0])
//	return string(myResult.Ok.Nonfungible.Metadata.Some), nil
//}

func GetYumiNFTImageUrl_1(canisterID string) ([]NFTUrl, error) {
	var infos []NFTUrl
	_agent := agent.New(true, "")
	methodName := "getTokens"
	arg, err := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
	if err != nil {
		return nil, err
	}
	_, result, err := _agent.Update(canisterID, methodName, arg, 10)
	if err != nil && err.Error() == "() empty" {
		return nil, nil
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var myResult []MetaData
	utils.Decode(&myResult, result[0])
	for _, token := range myResult {
		metadata := token.MetaData.NonFungible.MetaData.Some
		if metadata == nil {
			fmt.Printf("metadata of %s#%d is empty\n", canisterID, token.TokenIndex)
			continue
		}

		//urlInfo := new(YumUrl)
		//if err = json.Unmarshal(metadata, urlInfo); err != nil {
		//	fmt.Printf("can not unmarshal url of %s#%d with error: %v", canisterID, token.TokenIndex, err)
		//	continue
		//}
		//
		var videoUrl string
		url := string(metadata)
		if strings.Contains(url, ".map4") {
			videoUrl = url
			url = ""
		}

		infos = append(infos, NFTUrl{
			CanisterID: canisterID,
			TokenID:    token.TokenIndex,
			VideoUrl:   videoUrl,
			ImageUrl:   url,
		})
		fmt.Println(videoUrl, url)
	}
	return infos, nil
}
