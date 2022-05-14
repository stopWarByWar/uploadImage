package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	agent "github.com/mix-labs/IC-Go"
	"github.com/mix-labs/IC-Go/utils"
	"github.com/mix-labs/IC-Go/utils/idl"
	"github.com/mix-labs/IC-Go/utils/principal"
	"github.com/xuri/excelize/v2"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strconv"
)

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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

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
		methodName := "getRegistry"
		arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
		_, result, _, err := _agent.Query(canisterID, methodName, arg)
		if err != nil {
			return nil, err
		}
		var myResult []registry
		utils.Decode(&myResult, result[0])
		for _, registry := range myResult {
			for _, token := range registry.Tokens {
				imageUrl := token.PhotoLink.Some
				videoUrl := token.VideoLink.Some
				infos = append(infos, CCCNFTInfo{
					CanisterID:    canisterID,
					TokenID:       token.Index,
					ImageUrl:      imageUrl,
					VideoUrl:      videoUrl,
					ImageFileType: fileType,
				})
			}
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
	}
	return infos, nil
}

func ListFiles(root string) ([]string, error) {
	var result []string
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		result = append(result, file.Name())
	}

	return result, nil
}
